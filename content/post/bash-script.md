+++
title = "Bash Script"
date = "2014-10-19"
categories = ["Tuto"]
tags = ["bash"]
+++

Writting a bash script is kind of a complicated task, there is a strict
syntax, multiple external tools, and some tricks which depends on the version
and in some cases on the distribution itself (for example, grep is not the same
either you are on BSD or Debian).

In this article I will talk about four code habits which can improve
the maintainability of your shell scripts

### Double quotes

Everyone who already used a bash script will tell you to mark every variable
reading with double quotes:

```bash
# like this
echo "$a"
```

The question is why do we have to mark them, to show that just a little example,
run the following script

```bash
#!/bin/bash
LS="ls *"

echo "Not Quoted"
echo $LS
echo
echo "Quoted"
echo "$LS"
```

As you can see there is a huge difference between the two versions, in fact,
if there is no double quotes, the variable is subject to word splitting and
file globbing see [here](http://mywiki.wooledge.org/BashPitfalls#echo%5f.24foo)
for further informations. This can lead to major errors and
failures.

### The braquets notation

To display a variable there is two notations

```bash
#!/bin/bash
FOO="foo"

echo "$FOO"
echo "${FOO}"

```

These two notations displays the same thing, so why do we need the second one?
This notation is here to clarify an ambiguity, because variables can be used
in a string interpolation it can lead to a miscomprehension.


```bash
#!/bin/bash
FO="fo"

echo "because $FOo you "
echo "because ${FO}o you"
```

As you can see the result of the previous script the first line does not write
what is contained into the variable, because bash interpreter try to display
the content of the variable `FOo` which does not exist.

### The set command

At the any level in your script you can use the `set` command,
this command is powerful because it allows some extra behaviour in your shell
script:

* `set -e` Will stop the script if an error occurs

```bash
#!/bin/bash
set -e

echo "this command exist"
ls
echo
echo "hope this one not"
foo
echo
echo "this will not be displayed"
ls
```

Try removing the second line and see what happened. Sometime this behaviour is
needed, sometime not.

* `set -x` Will display the line running and evaluate the variables

```bash
#!/bin/bash
set -x

FOO="foo"
BAR="$FOO"
```

This will display

```bash
+ FOO=foo
+ BAR=foo
```

### Command return code testing

There is a simple way to test some command return in bash if we do not need the
result

Here is the common way to do that

```bash
#!/bin/bash

RESULT=$(grep "toto" /dev/null)

if [[ $? -eq 1 ]]
then
  echo "Command failed"
fi

RESULT=$(ls /dev/null)

if [[ ! $? -eq 1 ]]
then
  echo "Command succeed"
fi
```

But if we do not want the result, and just want to test some command return,
here is the simple way to perform that

```bash
#!/bin/bash

if grep "toto" /dev/null
then
  echo "Command succeed"
else
  echo "Command failed"
fi
```
