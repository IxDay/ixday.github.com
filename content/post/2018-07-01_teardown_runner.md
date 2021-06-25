---
title: "Parallel runners with teardown in go"
date: 2018-07-01
tags: ["golang", "dev"]
categories: ["Snippet"]
---

From time to time in go I have to start multiple small services acting in
parallel. For example, a ssh server tied to an administration console
(one is running on port 80 the other on port 22), or a kafka consumer pushing
to a database and a website to serve those informations. You can decouple this
in multiple programs, or run them through some kind of a manager and handle
everything at the same place. Both solutions have trades off, but we will look
at the last one today, because we will write a Golang wrapper for this.

What do we want to achieve
--------------------------

The basis is running multiple services in parallel, handle the potential errors
and bubble them to stop the program if needed. Furthermore, we'd like to have
support for teardown, and interrupt.

In go there's multiple way to do this, one of them is `channel` the other one
is the `sync` package. During this development I found it easier to deal with
`sync` (you can still go for `channel`).

We will also take a deeper look at Golang interfaces, and design patterns.
If you still haven't seen this excellent
[talk](https://www.youtube.com/watch?v=xyDkyFjzFVc), just do it,
we will use some of its points.

I will also publish the code in
[my snippet repo](https://github.com/IxDay/snippets/tree/master/golang/src/runner).

First draft, build base interface
---------------------------------

Let's define what we will provide to the outside world as a simple interface.

```go
import (
	"context"
)
type (

	// This interface is simple enough, it gets a context which will be used
	// to propagate some event down the pipe. We also want to handle potential
	// errors.
	Runner interface {
		Run(ctx context.Context) error
	}
	// This is more boilerplate, it will help us write runners without needing
	// structs. This is related to Tomas Senart video.
	RunnerFunc    func(context.Context) error
)

func (rf RunnerFunc) Run(ctx context.Context) error { return rf(ctx) }
```

Okay, let's check what a simple implementation can look like. An example is
worth a thousand words.

```go
import (
	"context"
	"fmt"
	"time"
)
// Miscellaneous runner because why not.
func Misc() Runner {

	// Here we use a bit of functionnal programming, to return a function implementing
	// our interface.
	return RunnerFunc(func(ctx context.Context) error {
		fmt.Println("Starting doing something...")

		// We will now block on two different events, either our work finish
		// normally, or we receive a done event from the context.
		select {

		// The normal workflow need to use a channel, to allow event base handling
		// here is a good example: https://gobyexample.com/select.
		// We just use the time package to wait a fixed duration.
		case <-time.After(20 * time.Second):
			fmt.Println("Done doing something...")
			return nil // nothing went wrong so no need to return an error

		// Here is the event received from the context, it asks to stop the current
		// work. We perform a teardown operation which takes a given time (just sleeping here).
		case <-ctx.Done():
			fmt.Println("Tearing down...")
			time.Sleep(5 * time.Second)
			fmt.Println("Tearing down, finished...")
			// Everything went ok, so no need to return any error, if teardown fails
			// we can propagate error
			return nil
		}
	})
}
```

We can now implement the interrupt logic:

```go
import (
	"context"
	"os"
	"os/signal"
)
func Interrupt() Runner {
	// Plug the interrupt call to a channel, so we can easily get the event.
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	// Build our runner.
	return RunnerFunc(func(ctx context.Context) error {

		// Here we block and return on either event is the first one to fire.
		select {
		case <-interrupt:
			return nil
		case <-ctx.Done():
			return nil
		}
	})
}
```

And finally, the service consumming this runner thing and handling the logic.
As I said we will use the `sync` package here and rely on `go` keyword for
the parallelization.

```go
import (
	"context"
	"sync"
)

type (
	// Our manager is just a slice of tasks which we want to run.
	RunnerManager []Runner
)

// Only one function required, we block the program until all execution has been done.
// We use a context as argument, so user can provide more constraints, like
// timeouts or external cancellation
// (see: https://golang.org/pkg/context/#WithCancel or https://golang.org/pkg/context/#WithTimeout for ideas).
// We also return an error, for the moment it will be the first caught one if any,
// we will improve this later.
func (rm RunnerManager) Wait(ctx context.Context) error {

	// Required variables:
	wg := sync.WaitGroup{}                  // WaitGroup keep tracks of running goroutines.
	ctx, cancel := context.WithCancel(ctx)  // Enhance the base context with a cancel function.
	errors := []error{}                     // We need to store errors, as multiple can come.

	// We loop over all the tasks we want to run
	for _, runner := range rm {
		// Add it to our group, this has to be done before going parallel as
		// go routines may not have started before the end of the function
		// and main function, causing program to exit before the work starts
		wg.Add(1)

		// Start a goroutine, collect errors, decrement group when task is done
		go func(runner Runner) {
			if err := runner.Run(ctx); err != nil {
				errors = append(errors, err)
			}

			// Task ended, decrement WaitGroup counter
			wg.Done()

			// Once task is done, we want to cancel all the remaining ones.
			// This is a design choice, nothing stops you from continuing execution
			// and only cancel when an error is caught.
			cancel()
		}(runner)
	}

	// Use the WaitGroup counter to block execution until all tasks finished.
	wg.Wait()

	// Return first caught error.
	if len(errors) != 0 {
		return errors[0]
	} else {
		return nil
	}
}
```

How to use it
-------------

Now, let's write our main function to see how all of this perform.

```go
import (
	"fmt"
)

func main() {
	// We queue our two tasks, and use a default context
	RunnerManager{Interrupt(), Misc()}.Wait(context.Background())
}
```

Just play around, try to interrupt or let go, things should work as expected
(I hope so ;)).

Going further
-------------

### Better error handling

We may need to record all errors and return them.
To do so we will define a new kind of error based on an array.

```go
type (
	// Simple type, just a list of errors
	Errors []error
)

// We only print the first one, yet another choice, feel free to customize according to your needs.
// We let it panic, cause we should have at least one error here.
func (e Errors) Error() string { return e[0].Error() }

// Adapt the function accordingly
func (rm RunnerManager) Wait(ctx context.Context) error {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	errors := []error{}

	for _, runner := range rm {
		wg.Add(1)
		go func(runner Runner) {
			if err := runner.Run(ctx); err != nil {
				errors = append(errors, err)
			}
			wg.Done()
			cancel()
		}(runner)
	}

	wg.Wait()
	// Here is our change, we check if we gathered an error
	if len(errors) != 0 {
		// We cast our array to the new type
		return Errors(errors)
	} else {
		// Otherwise we just return nil, as no error occured.
		return nil
	}
}
```

### Better API

We can improve the API by providing a way to deal with runners without using
complex context usage. Basically, our API allows you to run stuff in a blocking
way, and then to teardown resources. Let's create a new interface with those
attributes, and a function in case someone does not need a `struct`.

```go
import (
	"context"
)

// Make a synchronous function asynchronous using a channel and a goroutine.
// We will use this in our next function.
func Async(cb func() error) <-chan error {
	out := make(chan error)
	go func() { out <- cb() }()
	return out
}

// We use here two function to create a runner. Those functions do not rely on
// any context logic, and are both blocking.
func TeardownRunnerFunc(run func() error, teardown func() error) Runner {

	return RunnerFunc(func(ctx context.Context) error {
		// We make the run function asynchronous and we wait both on return and context
		select {
		case err := <-Async(run):
			return err
		case <-ctx.Done():
			return teardown()
		}
	})
}
```

And here the interface with a function to cast it to a runner:

```go
// Define an interface with the two function needed.
type TeardownRunner interface {
	Run() error
	Teardown() error
}

// Define a function to use our new interface as a runner.
func RunnerWithTeardown(tr TeardownRunner) Runner {
	return TeardownRunnerFunc(tr.Run, tr.Teardown)
}
```

That's it, I tried to put the code in the simpler way I could, and brings some
way to extend it. Do not hesitate to checkout the
[repo](https://github.com/IxDay/snippets/tree/master/golang/src/runner)
in order to see it complete with tests.
