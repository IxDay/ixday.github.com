---
date:       2017-07-25
title:      "Golang err shadowing"
tags:       ["golang", "dev"]
categories: ["Snippet"]
url:        "post/golang-err-shadowing"
---

A feature I like in golang is the hability to declare a variable at the
assignation time. Something like this:

```go
foo := "bar"
```

Here the variable foo will automatically set up as a string with the value "bar".
One more feature is to be able to allocate on same line as doing a comparison.
Like this:

```go
if foo := "bar"; foo == "baz" {
	// do something
} else {
	// do something else
}
```

This is really handy when it comes to catch errors from an other function:

```go

if err := potentialyFailingFn(); err != nil {
	// if there is an error we return and stop execution
	return err
}
// here is executed if there is no error caught
```

This assignation however is scoped to the `if` block, which means that out of
it `err` does not exist.

```go
if res, err := potentialyFailingFn(); err != nil {
	// error caught aborting
	return err
} else {
	// here res exist, we can do whatever we want with the variable
	whatever(res)
}
// here nor res or err exist and we can't access them
```

Still with me? Ok, one last thing, the named returned value
```go
import (
	"errors"
)

func foo() (err error) {
	err = errors.New("bar")
	return	// no need to specify the err variable here,
	        // this is a golang feature, see: https://tour.golang.org/basics/7
}
```

So now, what if I do this:

```go
import (
	"errors"
)

func foo() (err error) {
	if err := errors.New("bar"); err != nil {
		return
	}
	return
}
```

Compiler says: `err is shadowed during return`, this says, we have a named
variable `err` defined in the function declaration. Then, in the `if` we
redeclare for the block a `err` variable with an error value. Compiler spot that
there can be a misunderstanding on which one to use then, and warns you.
If type had been different, compiler should just have errored, but here I can
reuse a variable name.

**BTW**, this compile and will return a `"bar"` error.
```go
import (
	"errors"
)

func foo() (err error) {
	if err := errors.New("bar"); err != nil {
		return err // note that I do not use the named variable but the scoped one.
	}
	return
}
```

Now, how can I use this? To perform this for example:

```go
import (
	"log"
)

func foo() (err error) {
	// note the difference between err = and err :=

	if err = bar(); err != nil { // I want to keep this error and return it

		// but I have one more call to do and want to log the error
		if err := baz(); err != nil {
			// here err is the scoped one and does not override the one returned
			log.Printf("%s", err)
		}
	}
	return
}
```

Interesting isn't it?

