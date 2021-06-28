+++
title = "Bash process substitution failure"
date = 2020-09-12
categories = ["OS"]
tags = ["bash"]
url = "post/bash_process_substitution_failure"
+++

This is just a problem I encountered running my Alpine box, so I publish the fix
for others to not loose time like I did.
I had a software that came with a bash script and got a problem when I tried
to run it. It was failing with the following error: `fopen: No such file or directory`

I added the magical `set -x` line at the beginning to understand what
was going on. I re-ran the script and saw the following output:

```bash
+ wg setconf wg0 /dev/fd/63
fopen: No such file or directory
```

I checked the script to retrieve the line and here is what was written:

```bash
cmd wg setconf "$INTERFACE" <(echo "$WG_CONFIG")
```

The problem is coming from the process substitution part.
I quickly checked that it was really what I was suspecting:

```bash
# ensure I am properly using bash first since ash do not support process substitution
/bin/bash

# this should echo something like 5.0.17(1)-release
echo $BASH_VERSION

# testing process substitution
cat <(echo "foo")
```

I got the same error: `fopen: No such file or directory`.
I then search the Internet for an answer and finally found
[this issue](https://gitlab.alpinelinux.org/alpine/aports/-/issues/1465).
I quickly ran through it and found out that it can be fixed by creating a
symlink: `ln -snf /proc/self/fd /dev/fd`.

__Important: be sure that `/dev/fd` does not exist before running the command,
otherwise, you will end up with `/dev/fd/fd`.__
