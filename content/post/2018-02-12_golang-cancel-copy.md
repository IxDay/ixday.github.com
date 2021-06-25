+++
title = "Cancel copy of huge file in Go"
date = 2018-02-12
categories = ["Snippet"]
tags = ["dev", "golang"]
+++

I recently came across [this video](https://www.youtube.com/watch?v=xyDkyFjzFVc) on Golang programming.
I think this was the moment I finally fully understood the power of Go.
It is smart, simple and elegant, I love it.

Then, a few days later, I was coding on a [toy project](https://github.com/IxDay/antfarm)
and I was doing some stuff around the [io package](https://golang.org/pkg/io/) to
copy huge files. I wanted to achieve copy cancelation during the processing, basically, being able to interrupt.

The first idea was to replicate the source code of golang [io.Copy](https://golang.org/pkg/io/#Copy)
and insert a `context` argument, then find a place to interrupt the copy
when looping over the file, when context got canceled.

This brings some issues:

- I had to copy/paste code from a standard lib, which will force me to maintain it.
- The code was not relevant to what I was doing. Basically, I copied low level
code in order to insert hooks.
- I had to test code which was already tested.

But in Go, extends is done through interfacing, and the copy function uses
interfaces as arguments. I can sneak in this! Here is how I achieved this:

```go
import (
	"io"
	"context"
)

// here is some syntaxic sugar inspired by the Tomas Senart's video,
// it allows me to inline the Reader interface
type readerFunc func(p []byte) (n int, err error)
func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// slightly modified function signature:
// - context has been added in order to propagate cancelation
// - I do not return the number of bytes written, has it is not useful in my use case
func Copy(ctx context.Context, dst io.Writer, src io.Reader) error {

	// Copy will call the Reader and Writer interface multiple time, in order
	// to copy by chunk (avoiding loading the whole file in memory).
	// I insert the ability to cancel before read time as it is the earliest
	// possible in the call process.
	_, err := io.Copy(out, readerFunc(func(p []byte) (int, error) {

		// golang non-blocking channel: https://gobyexample.com/non-blocking-channel-operations
		select {

		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return in.Read(p)
		}
	}))
	return err
}
```

There is a huge topic in Go community around the use of context. Should all the
standard lib use contexts in order to be able to cancel some calls. I think
interfaces can be the start of a solution. Allowing you to extend what is needed
and leave the stdlib as simple as possible. This is a good example, hope it
can help some people to better understand how Golang interfaces can be used.
