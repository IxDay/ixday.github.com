+++
title = "HTML5 Boilerplate"
date = "2014-09-30"
categories =  ["Snippet"]
tags = ["bash", "web"]
+++

Sometimes we only need to have a boilerplate quickly to test it through a
browser. We only want to have the basis, and having it working fine.

You can find [here](http://www.initializr.com/) a good generator for what you
want. But sometimes, having a snippet in the bash prompt can be needed.

So here is an example:

```bash
# stop script if something bad happen
set -e

# unzip need to have a tempfile to extract properly
TMPFILE="/tmp/tempfile.$(date +%s)"

PWD=$(pwd)
DEST="$1"

# here is the link of my configuration, but you can easily create yours
URL="http://www.initializr.com/builder?h5bp-content&html5shiv&"
URL="$URL""h5bp-css&h5bp-csshelpers&h5bp-mediaqueryprint&h5bp-mediaqueries&"
URL="$URL""simplehtmltag&izr-emptyscript" # string concatenation example

curl -o "$TMPFILE" "$URL" 2> /dev/null

# if argument is provided, move to the specified directory
if [[  -z "$DEST" ]]
then
    # do not display what has been inflated
    unzip -qq -d "$PWD" "$TMPFILE"
    rm "$TMPFILE"
else

    # if the path provided is a file warn the user and exit
    if [[ -f "$DEST" ]]
    then
        echo "$DEST already exists and is a file" && exit 1
    fi

    # if the directory is not empty warn the user and exit
    if [[ -d "$DEST" && ! -z "$(ls -A $DIR)" ]]
    then
        echo "$DEST is not empty" && exit 1
    fi

    # if the directory does not exist, we create it
    if [[ ! -a "$DEST" ]]
    then
        mkdir -p "$DEST"
    fi

    unzip -qq -d "$DEST" "$TMPFILE"

    #then we move the content to the directory
    for f in "$DEST"/initializr/*
    do
        mv "$f" "$DEST"
    done

    # delete the other files
    rm "$TMPFILE"
    rmdir "$DEST/initializr"
fi
```

it is also possible to create a sample on github or whatever, do not forget that
automation will lead to time saving (maybe, ...sometimes).
