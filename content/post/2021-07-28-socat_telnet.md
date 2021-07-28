---
title:      "Socat, Telnet and Unix sockets"
date:       2021-07-28
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

Those limitations are pushing me towards the use of the `socat` utility.
Here I will show a few situations in which I am now using `socat`. As always
I recommend to read the [man page][documentation] to get a good idea of the
various options.

Now let's get started!

## The escape sequence

First of all `socat` does not provide the `^]` escape sequence by default.
This means that if you connect to an interactive shell you will close the
connection when pressing `^C` or `^D` (the `^` notation represents the `Ctrl` button,
it is equivalent to `Ctrl+]`, `Ctrl+C`, `Ctrl+D`).
However, you can mimic this feature with the following
flags: `,rawer,escape=0x1d`.
- The `rawer` option will avoid closing the `socat` process by pressing `^C` or `^D`,
  and will pass it through the connection allowing
  you to control remote processes.
- The `escape=0x1d` option is here to actually close the process by pressing
  `^]`.

## Some examples

All those examples use `socat` in a client configuration. However, `socat`
provides a lot of additional use cases, notably as a server
(I wrote a [small post][previous_post] a while ago).
Now that all is clear, let's introduce the examples:

__Connect to remote port:__

```sh
echo | socat - tcp:foo.bar:22
```

__Connect to remote virtual terminal over TCP:__

```sh
socat -,rawer,escape=0x1d tcp:foo.bar:23
```

__Connect to a local virtual terminal over a Unix socket:__

```sh
socat file:`tty`,rawer,escape=0x1d unix-connect:console.sock
```

The option `` file:`tty` `` is equivalent to `-` here. The following command will
have the same behavior:

```sh
socat -,rawer,escape=0x1d unix-connect:console.sock
```

__Issue an http request:__

```sh
echo -e "GET / HTTP/1.1\nHost: www.example.com\n" | socat - tcp:www.example.com:80
```

__Issue an https request:__

```sh
echo -e "GET / HTTP/1.1\nHost:github.com\n" | socat - openssl:github.com:443
```

__Moar examples:__ follow [this link][moar_examples] for more ideas on how to
use `socat`.

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

[documentation]: https://linux.die.net/man/1/socat
[moar_examples]: http://technostuff.blogspot.com/2008/10/some-useful-socat-commands.html
[previous_post]: /post/simple_https/
