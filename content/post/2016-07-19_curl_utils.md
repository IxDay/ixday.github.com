+++
categories = ["Snippet"]
date = 2016-07-19
tags = ["bash"]
title = "Curl utils"
url = "post/curl_utils"
+++

Here are some options and command I use with `curl` when dealing with stuff
I have to develop.

```bash
curl -si <ip> # -s is the silent flag, it removes the progress
              # -i displays the headers

curl -X POST -H "Content-Type: application/json" -u "admin:admin" -d '{}' <ip>
# -X set up the http method (here POST)
# -H set up an header, format is: "header_name: value"
# -u support for Basic Auth, format is: "user:password"
# -d set up data to send to the server
```

I mostly use those options, the `-s` is really interesting when you want to
grep the content.

So, when testing availability of an http service, I use this snippet:

```bash
while true; do curl -si <ip>|awk 'NR==1||NR==3'; sleep 1; done
```

This will do kind of a ping for an http service, displaying this:

```text
HTTP/1.1 200 OK
Date: Tue, 19 Jul 2016 12:23:43 GMT
HTTP/1.1 200 OK
Date: Tue, 19 Jul 2016 12:23:45 GMT
HTTP/1.1 200 OK
Date: Tue, 19 Jul 2016 12:23:46 GMT
```

The status of the http call, and the date of the call has been done.
