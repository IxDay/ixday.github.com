---
title:      "Socat, Telnet and Unix sockets"
date:       2021-07-25
categories: ["Snippet"]
tags:       ["cli", "admin", "bash"]
url:        "post/socat_telnet"
---

Once in a while I use `telnet`, mostly to check if a port is open
(the infamous `telnet localhost 22`) and sometimes to send a random http request.

However, `telnet` has a few caveats which are:
- not able to read from stdin
- not able to deal with https
- not able to handle unix sockets


## The escape sequence

On the other hand `socat` does not provide the `^]` escape sequence so if you
connect to an interactive shell you will close the connection when pressing `^C`
or `^D` (the `^` notation represents the `Ctrl` button, it is equivalent to `Ctrl+]`, `Ctrl+C`, `Ctrl+D`).

## One last thing

During my investigation around `socat` I noticed a weird
potential bug. It seems that https connection is broken when connecting
to a website with a wildcard certificate. You can check if that's the case
by running:

```sh
echo | openssl s_client -connect google.com:443 2> /dev/null | grep -i subject
```

When running socat against this domain using https you will get the following
error:

```text
socat[26639] E SSL_connect(): error:1416F086:SSL routines:tls_process_server_certificate:certificate verify failed
```

I managed to make it work by piping `socat` to `cat`, I have no idea why this
happen or the reason why the piping fixes it (I discovered this out of pure luck). So the final working command will be:

```sh
echo -e "GET / HTTP/1.1\nHost:google.com\n" | socat - openssl:google.com:443  | cat
```

http://technostuff.blogspot.com/2008/10/some-useful-socat-commands.html
