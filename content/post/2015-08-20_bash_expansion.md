---
title:      "Bash Expansion"
date:       2015-08-20
categories: ["Snippet"]
tags:       ["bash", "cli"]
url:        "post/bash_expansion"
---

A small post just to share a useful bashism the
[Brace Expansion](http://wiki.bash-hackers.org/syntax/expansion/brace).
It is really simple to make it works:

```bash
for i in {1..50}
do
	echo "Hello World $i"
done
```

It will print fifty "Hello World". Ok it seems cool but not amazing?
Ok, now the second feature
```bash
echo something/{foo,bar}
> something/foo something/bar
```

Still not amazed, ok now type this one:
```bash
cp some_file{,.old}
```
It will copy your file
adding a `.old` extension. I do a lot of things like this and it saves me a lot
of time, so think about it.
