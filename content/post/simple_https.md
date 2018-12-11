+++
title = "Simple HTTPs server"
date = 2018-12-11T15:46:28+01:00
categories = ["Snippet"]
tags = ["cli", "admin"]
+++

Sometimes we need to create a simple server saying hello through https.
Here is a simple snippet to achieve this in a shell.

```bash
# first we generate a self signed certificate for domain foo
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes \
	-out foo.crt -keyout foo.key -subj "/CN=foo.com"

# concatenate key and certificate to create a pem file
cat foo.key foo.crt > foo.pem

# then we start a server using socat
sudo socat "ssl-l:443,cafile=foo.crt,cert=foo.pem,verify=0,fork,reuseaddr" \
	SYSTEM:"echo HTTP/1.0 200; echo Content-Type\: text/plain; echo; echo Hello World\!;"
```

You can now request using curl

```bash
# using ip address
curl -k "https://127.0.0.1:443"

# providing a resolve entry matching our domain name
curl -k --resolve "foo.com:443:127.0.0.1" https://foo.com
```

If you got time, take a look at [`socat`](https://linux.die.net/man/1/socat),
it's a super powerful tool.
