+++
title = "Wireguard on a Linux Alpine with Docker"
date = 2020-09-22T16:51:22+02:00
categories = ["Tuto"]
tags = ["alpine", "admin"]
+++

For most of my infrastructure, I am now using [Alpine Linux](https://alpinelinux.org/).
I like it because it only has a small number of moving parts. It's easy to
know and master them, it is making my life easier :).

So, I decided to install one on my VPS. Like my distro I wanted it to
be simple and small. For all those reasons I went for [Wireguard](https://www.wireguard.com/).
The fact that it is the new cool kid, may also have helped.

Last but not least, I want this VPN to run inside a docker container which is
running Alpine as well!

__My cloud provider:__ Like everything here I wanted my cloud provider to be
simple to use. Also, I needed to be able to launch an Alpine Linux VPS. I tried
OVH (would not recommend), then Digital Ocean, but it was too expensive even
for the most simple setup. Finally, I discovered [Scaleway](https://www.scaleway.com/en/).
I am using the smallest instance, plus an IPv4 address plus a 50GB SSD volume.
This cost me less than 5€ per month. Needless to say that I am pretty happy!

Setup of the host
-----------------

To run Wireguard in a container we need to configure the underlying host.
Since containers share the host kernel you have to do some changes to make it work.

First, you have to install the kernel module:

```bash
# first check your kernel version
uname -r

# install wireguard kernel module
apk add wireguard-${your_kernel_version}
```

Once this is done, we also need to configure kernel parameters to allow IP forwarding.
I want the VPN to tunnel all traffic. I will forward IPv4 and IPv6.

```bash
# temporary set up
sysctl -w net.ipv4.ip_forward=1
sysctl -w net.ipv6.conf.all.forwarding=1
```

I want those changes to be permanent. I first looked at writing this to
`/etc/sysctl.conf`. However, the change was not picked up at boot. This is
mostly due, as far as I get it, to the fact that the IPv6 support is a kernel
module that is loaded after the kernel config is called. I found some
reports [here](https://www.raspberrypi.org/forums/viewtopic.php?t=113560)
and [here](https://bugs.launchpad.net/ubuntu/+source/procps/+bug/50093).

I decided to use a local startup script from OpenRC. You only need to create
a shell script in `/etc/local.d` and make it executable
(it also needs to have a `.start` suffix).

```bash
#!/bin/sh
# 60-forward.start script, set up the IP forward kernel parameter

sysctl -w net.ipv4.ip_forward=1
sysctl -w net.ipv6.conf.all.forwarding=1
```

The last step here is to activate the local scripts at boot with this command:
`rc-update add local default`.
Now that we have IP forwarding, it's time to set up iptables.
Here the configuration is easy, you need to let iptables know about the forwarding
of packets for both IPv4 and IPv6:

```sh
iptables -P FORWARD ACCEPT
rc-update add iptables
/etc/init.d/iptables save

ip6tables -P FORWARD ACCEPT
rc-update add ip6tables
/etc/init.d/ip6tables save
```

__Important:__ The saving of the forwarding rules may be disabled in your
configuration, ensure that `IPFORWARD="yes"` is properly set in the files:
`/etc/conf.d/iptables` and `/etc/conf.d/ip6tables`

Setup the container
-------------------

This will be a bit more straightforward. We first need to create the container
Dockerfile:

```dockerfile
FROM alpine:3.12

RUN apk add --no-cache wireguard-tools ip6tables
COPY server.sh /usr/local/bin/wireguard
EXPOSE 5555
CMD ["wireguard"]
```

For the container, I am using the same version as my host system. Since Wireguard
is using a kernel module on the host system I would like to avoid any incompatibility.
I am also installing iptables because the Wireguard script needs to add a
few more rules at startup and shutdown.
You would have noticed that I am using a custom script as the executable, this
is because I need extra configuration to be generated at startup. Here is the content:

```sh
#!/bin/sh

#  10.0.0.x, fd47:d1a9:8d26:c99b:xxxx:xxxx:xxxx:xxxx
cat > /etc/wireguard/wg0.conf << EOF
[Interface]
PrivateKey = $(wg genkey)
Address = 10.0.0.1/24, fd47:d1a9:8d26:c99b::/64
ListenPort = 5555
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE; ip6tables -A FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE; ip6tables -D FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
SaveConfig = true

EOF

wg-quick up wg0
watch wg
```

Let's explain what I am doing:
- I am first putting in the comment the network blocks I am using for the VPN.
I am doing IPv4 and IPv6 so I describe both of them.
- Next, I am generating the config file. I want to create the `wg0` interface like the default one
from documentation. You can of course replace it with another naming convention.
- I am generating a private key for my peer. This is why I am using a script,
the key is generated at startup and is unique to every restart. You can change
the logic if it does not fit your need.
- Then, I configure the addresses blocks for IPv4 and IPv6. I am also putting the
port to listen to. It is the same as the Dockerfile.
- The `PostUp` and `PostDown` sections are just copy-pasted from the [Linode documentation][linode]
(most of my configuration is coming from there).
- Last line will save configuration to file, I am keeping this here for debug
purpose, and __it's not required for this to work__. Also, notice that any change
to the file when Wireguard is running will be erased at shutdown since the program
is writing the current running state to the file.

Once the container is built (with a command like `docker build -t wireguard .`)
we can start it with the following command:

```sh
docker run --rm --detach --name wireguard \
	--cap-add=NET_ADMIN --cap-add=SYS_MODULE \
	--network=host --volume /lib/modules:/lib/modules \
	wireguard
```

This command may need a bit of clarification. First, I am running the container
with `--rm`, `--detach` and `--name`. This is a kind of convention I am using
for most of my _"production"_ containers. This allows me to always start with
a clean state. The `--rm` and `--name` are a good flags combination, I can always
connect to the same containers using the same commands which speed
up my operations and it properly delete my container at the end of the use so
I can easily re-use the name. The `--detach` option is mostly here to not block
the shell and copy-paste those commands to a script file without effort.

Then, the capabilities, here I want to pass the bare minimum to the container.
I usually run with `--privileged` when I don't know the tooling and need to experiment.
However, in my _"production"_ environment, containers run with the smallest
amount of privileges possible. You can check Docker documentation for a
short introduction on the [various capabilities][capabilities].


[linode]: https://www.linode.com/docs/networking/vpn/set-up-wireguard-vpn-on-ubuntu/#configure-wireguard-server
[capabilities]: https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities

Synchronize configuration
-------------------------

Now our server is running, we need to configure the client and synchronize both
endpoints to make them communicate properly.

To do that I run the following command:

```sh
ssh <remote host> docker exec -i wireguard wg show wg0 public-key | \
	sudo xargs ./setup.sh | \
	xargs -rI{} ssh <remote host> docker exec -i wireguard \
		wg set wg0 peer {} allowed-ips 10.0.0.2/32 allowed-ips fd47:d1a9:8d26:c99b::1/128
```

This will use the following script:

```sh
#!/bin/sh

public="${1}"
private="$(wg genkey)"

cat > /etc/wireguard/wg0.conf << EOF
[Interface]
PrivateKey = ${private}
Address = 10.0.0.2/15, fd47:d1a9:8d26:c99b::1/64
DNS = 8.8.8.8

[Peer]
PublicKey = ${public}
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = <remote host>:5555
EOF

echo "${private}" | wg pubkey
```

Now let's do a bit of explanation. In the first line of the shell command,
I retrieve the public part of the server key over the network using ssh.
It is then piped to the script. The script captures the input and store the
public key in the `public` variable. Once this is done, it generates its own
key-pair for the client side and store it in the `private` variable.

Now that we have all the keys, we generate the client config file. We set up
the IPv4 and IPv6 address of our current client endpoint as well as a DNS. Later
I will use my internal DNS with [Hashicorp Consul][consul], but for now
I am using Google's one.

The second part of the configuration is the connection with the server. First,
we populate the public key of our remote host. Next line, we define the network
traffic we want to direct through Wireguard. Here, I redirect absolutely everything,
IPv4 and IPv6 traffic will go through our VPN. Last line, we indicate the endpoint
of our server. It is using port `5555` and you will have to indicate the IP address
or DNS name.

Once the configuration file is properly set up, we echo the public part of our client.
This is passed up to the shell process and pipe to our last command (the `xargs` part).
It is a bit convoluted but I will try to explain what it does.

I use [xargs][xargs] to capture the script output, here it will get the public key
of my client peer. I assign this value to the `{}` symbol, it will allow me to
inject the value in the middle of my command later on. The next part of the line
is the `ssh` + `docker exec` combo. Since I am running my VPN on a remote host
inside a Docker container, this should make sense, otherwise check the previous
sections. The last part is the interesting one. I am calling the `wg` binary
to add a peer to my `wg0` interface. This interface will use the injected public
key and have an IPv4 and IPv6 address.

If on the remote host you run `docker logs -f wireguard`, you should see something like that:

```text
interface: wg0
  public key: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
  private key: (hidden)
  listening port: 5555

peer: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
  endpoint: <your_ip_address>
  allowed ips: 10.0.0.2/32, fd47:d1a9:8d26:c99b::1/128
```

Everything is ready! Just running `wg-quick up wg0` should bring the VPN up and
you can browse the Internet from wherever you are ;)

[consul]: https://www.hashicorp.com/products/consul
[xargs]: https://man7.org/linux/man-pages/man1/xargs.1.html

Miscellaneous
-------------

Sometimes I want to echo the private key of a Wireguard
configuration. I want to read directly from the file to pipe to a potential
following command. Thus to avoid leaking the key to the bash history for example
(a shell recorder may also be present). Here is the more portable one-liner I found:

```sh
sed -n -e 's|PrivateKey = ||p' /etc/wireguard/wg0.conf
```
