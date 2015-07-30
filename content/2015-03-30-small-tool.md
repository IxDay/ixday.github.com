Title: Small docker project
Date: 2015-03-30
Category: Project
Tags: admin, bash, dev

I really like docker (even if I will look at
[systemd-nspawn](http://www.freedesktop.org/software/systemd/man/
systemd-nspawn.html)), and also [gulp](http://gulpjs.com/).
So I decided to create a small tool for serving a directory with a livereload.

The repo is available [here](https://github.com/IxDay/docker-html5-boilerplate)

## What I have learned
* Docker, especially with boot2docker (I am on MacOSX shame on me), is not
  really flexible:
    * no evaluation for environment variables
    * you can not store a variable through multiple run, you will need to do
      a oneliner e.g:

			#!bash
            # you will need to write this in the Dockerfile
            RUN TMPFILE=$(tempfile) && \
              echo "Hello World" > $TMPFILE && \
              rm $TMPFILE

    * you cannot build remotely from a custom branch:
    [see this post](http://stackoverflow.com/questions/25509828/
    can-a-docker-build-use-the-url-of-a-git-branch)
    * accessibility to the container through boot2docker is pretty hard:
    [see this post](http://stackoverflow.com/questions/28047809/
    docker0-interface-missing-on-osx/)

* awk is an awesome tool for manipulating strings and console outputs.
  Best one hour investment so far. If you have to write bash scripts, awk is
  definitely a best to know tool.
* nodejs is definitely a hell for developers but is the only platform for
  frontend dev. Here are some issues I have gone through:
    * bug in collecting interfaces which forced me to use a syscall:
      [bug report](https://github.com/joyent/node/issues/9029)
    * npm ecosystem is messy. You will have to test many plugin, which does
      quite the same thing, to find the one with *THE* option you need.
      And finally find that the lib is 10 lines long and wraps another lib.
    * no good packaging, which leads to a custom installation path,
      thanks [this blog post](https://nodesource.com/blog/
      nodejs-v012-iojs-and-the-nodesource-linux-repositories) for the help.
* gulpjs is way better than grunt and must be ported to other languages,
really excited about the 4.0 version coming soon.

## How I can improve

The tool is currently working, I am thinking of adding a markdown compiler in
the chain because I really use markdown all the time.

I still have to test it a bit to be sure it is okay to make a release
(first thing I release yay \o/)
