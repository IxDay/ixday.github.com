---
title:      "Curl for REST"
date:       2015-01-07
categories: ["tuto"]
tags:       ["admin", "cli", "bash"]
url:        "post/curl_for_rest"
---

It has been a long time since the last post. But today, I will just show two
tools I use for debugging my REST APIs.

First one is the well known [curl](http://curl.haxx.se/docs/manpage.html) and
the second one is [jq](http://stedolan.github.io/jq/manual/).

One important feature of curl is its hability to load external files for datas
with `@` before file name:
```bash
curl -X POST -H "Content-Type: application/json" -d @filepath
```

Then you can remove the progress bar by adding `-s` in the options

Finally, you can use jq for parsing the output with a request syntax, here is
what the final line looks like:

```bash
curl -X POST -H "Content-Type: application/json" -s -d @filepath | jq '.'
```

It is a small article but I just wanted to show that those tools are really
great, we do not need complex software for this. I prefer using simple cli
because it gives a better understanding on what we are doing, and on what we
are relying on.

