+++
title = "Simple HTTPs server"
date = 2018-12-11
categories = ["Snippet"]
tags = ["cli", "admin"]
+++

Sometimes we need to create a simple server saying hello through https.
Here is a simple snippet to achieve this in a shell.

```bash
# first we generate a self signed certificate for domain foo
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes \
	-out foo.crt -keyout foo.key -subj "/CN=foo.com"

# then we start a server using socat
sudo socat "ssl-l:443,cert=foo.crt,key=foo.key,verify=0,fork,reuseaddr" \
	SYSTEM:"echo HTTP/1.0 200; echo Content-Type\: text/plain; echo; echo Hello World\!;"
```

You can now request using curl

```bash
# using ip address
curl -k "https://127.0.0.1:443"

# providing a resolve entry matching our domain name
curl -k --resolve "foo.com:443:127.0.0.1" https://foo.com
```

If you don't want to use the `-k/--insecure` option you can install the root
certificate on your machine by running the following commands (may vary on distributions):

```bash
cp foo.crt /usr/local/share/ca-certificates/
update-ca-certificates
```

Last but not list, you can use socat as an ssl termination proxy. It is pretty
straightforward:

```bash
# start a simple http server, here darkhttpd to serve some directory
darkhttpd /tmp --port 8080 --daemon

# now start your socat proxy
sudo socat "ssl-l:443,cert=foo.crt,key=foo.key,verify=0,fork,reuseaddr" \
	'tcp4:0.0.0.0:8080'
```

Et voilà! Simplest way I ever found to perform some testing or run ad-hoc
simple https servers.

If you got time, take a look at [`socat`](https://linux.die.net/man/1/socat),
it's a super powerful tool.
