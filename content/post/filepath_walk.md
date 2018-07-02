+++
title = "Filepath Walk bug?"
date = 2018-05-09T10:22:33+02:00
categories = ["Snippet"]
tags = ["golang"]
+++

During the development of a small project I noticed something strange with
golang's function [filepath.Walk](https://golang.org/pkg/path/filepath/#Walk).

If you are going over a directory, and return `filepath.SkipDir` error on
a file entry. The walker will stop.


Setup directory
---------------

First, set up a simple directory to walk

```bash
mkdir -p foo/baz
touch foo/bar
touch foo/baz/qux
```

Run the script
--------------

Now, let's run this example script, with the previous directory as the entry
argument. I will reuse the example from golang documentation website with small
changes.

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := os.Args[1]   // this is not solid, but will work for the example
	fileToSkip := "bar" // we skip on a file to highlight the bug

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == fileToSkip { // no check on entry type info.IsDir()
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file: %q\n", path)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
}
```

This will output:

```text
visited file: "foo"
skipping a dir without errors: bar
```

As you can notice, there's no trace of subdir `baz` and file `baz/qux`.

Fix the script
--------------

To make this work _"as expected"_ you will have to change line 17. With:
```go
if info.IsDir() && info.Name() == fileToSkip { // now test on entry type
	//unchanged
	fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
	return filepath.SkipDir
}
```

Output is now:

```text
visited file: "foo"
visited file: "foo/bar"
visited file: "foo/baz"
visited file: "foo/baz/qux"
```

__So, next time using `filepath.Walk` be careful on this one!__
