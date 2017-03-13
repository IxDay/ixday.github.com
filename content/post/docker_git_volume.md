+++
date = "2017-03-13T14:52:21+01:00"
title = "Git clone inside a mounted volume with Docker"
tags = ["dev"]
categories = ["Snippet"]

+++

**DISCLAIMER:** This is now fixed in git new releases and does not need to be done
anymore. I use an old version of alpine in order to have an unpached version of
git.

I ran into an interesting issue lastly. I wanted to mount a volume inside
a container and clone a repo in it. I also wanted to avoid messing with the
permissions and pass my user to the container as well.

## Set test env

Here is my setup, a Dockerfile with a container and git installed

```dockerfile
FROM alpine:3.2
RUN apk --update add git
```

Then simply run the git clone command inside a volume mounted with correct uid
and gid.

```console
docker build -t foo --rm .
docker run --rm -u 1000:1000 -v $(pwd):/mnt -w /mnt test \
	git clone https://github.com/octocat/Spoon-Knife
```

And I get the following error:
`fatal: unable to look up current user in the passwd file: no such user`. This
is actually right, my user only exists on the host and not into my container.

## Googling

I checked my best friend for technical questions: Google. And found
[this thread](http://www.spinics.net/lists/git/msg263682.html)
which describe exactly what I am facing right now. The thread is really
interesting and I invite you to read it entirely.

A solution is suggested [here](http://www.spinics.net/lists/git/msg263958.html).
What does it say? Just disable the `reflog` at clone time. I still have an issue
here because I do not have the hand on the cloning command in my real life issue.

Googling again, nothing really obvious came here. So better check at the git
documentation, and more precisely at the
[git config man page](https://git-scm.com/docs/git-config). Looking for `reflog`
and ended up [here](https://git-scm.com/docs/git-config#git-config-corelogAllRefUpdates).
Great, there is a trigger!

## Fixing

We now have to set this option inside the Dockerfile in order to make the
process able to clone the repo correctly. There is an option inside the
`git config` we will use: `--system` described
[here](https://git-scm.com/docs/git-config#git-config---system).
This will set the configuration system wide, because the cloning will occur with a random uid,
we need a generic setup, whereas `--global` as root will only set the option
for the root user.

## Final setup

We modify our Dockerfile according to the last results.

```dockerfile
FROM alpine:3.2

RUN apk --update add git
RUN git config --system core.logallrefupdates false
```
The git config command is case insensitive, so all lowercase will work here.
Now just run the previous build and run again your git clone.

**It works!**
