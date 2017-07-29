+++
date = "2017-07-27T10:53:10+02:00"
title = "direnv"
tags = ["bash", "cli", "dev"]
categories = ["Snippet"]

+++

I recently discovered the [direnv](https://direnv.net/) project. Which helped
me a lot for setting up my development environments. I will share a bit of
things I use in my daily basis.

## Python projects

This is well documented but here is what I use in python projects to set up
a default virtualenv.

```bash
layout python
```

That's it! It sets up a virtualenv in your `.direnv` directory,
and load the updated `PATH`.

## Golang projects

This one is a bit trickier.

```bash
export GOPATH="$(pwd)/.go"
export GOBIN="$GOPATH/bin"
PATH_add "$GOBIN"
```

I usually put my `$GOPATH` in an hidden directory, here it's `.go` and link
the project to the top directory. But you can also set up the `vendor` directory,
depending on your needs.

## Misc

This is a trick I use with a bunch of project to move quickly around. I set up
this in `.envrc` file:

```bash
export PROJECT="$(pwd)"
```

And then in the `.zshrc` (works also with `.bashrc`), I put this:

```bash
alias cd='HOME=${PROJECT:-$HOME} cd'
```

This will make my `cd` command to go back to the root of my project (where the
`.envrc` file lands) instead of going to `$HOME`. But only inside my project
and without overriding my `$HOME` variable globally. Which is something you
really don't want, seriously don't do it!
