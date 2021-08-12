---
title:      "KVM with Docker bridge"
date:       2021-08-10
categories: ["Snippet"]
tags:       ["admin", "cli"]
url:        "post/kvm_docker"
---

This post will explain how to use the docker bridge as a KVM bridge. In this
post I will use the Qemu command line to manage my VMs.

There is a lot of ways to connect a VM to the internet. The most common one
is via network address translation ([NAT][wiki_nat]). This method has a few
down side, the main one being that you need to explicitely configure port
forwarding for your VM services to be reachable from the host.

I wasn't a network person when I first started playing with docker. So when I
did the first installation I discovered its bridge networking. I found the solution
to better fit my needs. Every container has a specific IP and I can reach them
from my laptop. I want the same for my VM but this is a bit less user friendly
since you have to create and maintain the bridge by yourself. But on the other
hand the bridge is not managed by Qemu making it possible to reuse the one from
Docker.

## The command

```sh
qemu-system-x86_64 ... \
	-netdev "bridge,id=user.1,br=docker0" \
	-device "virtio-net,netdev=user.1"
```

If you want a proper explanation of the options I use, I will let you check
my [previous post][previous_post] where I explain this in details.
For now the options we are looking into are:
- `-netdev "bridge,id=user.1,br=docker0"` this define a network device with
  `id`: `user.1` of type bridge called `docker0`. Which is the bridge created by the
  docker daemon. Here, we will reuse that bridge to connect our VM.
- `-device "virtio-net,netdev=user.1"` second option is here to attach the device
  to our VM. We define the network driver and the device id. As previously shown
  in the other post you need two options to fully describe and attach a network device.

## The configuration

To make this properly work you will need a configuration step. Since the docker
bridge do not provide you with dynamics IPs you will have to set up a static one.
This is easily done using the `/etc/network/interfaces` file:

```txt
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet static
    address 172.17.1.1/16
    gateway 172.17.0.1
```

In this example I am setting up a static IP of `172.17.1.1` using the default
docker bridge network (`172.17.0.0/16`). This configuration will allow you to reach the internet
through the bridge, but also the containers you might start in this network as well.

__Disclaimer:__ As far as I know the docker daemon will not be aware of this setup
and might start a container with the same IP. I am using the address `172.17.1.1`
because I will need to have 254 containers registered in this network before
it conflicts. Keep this in mind when you are setting this up!

[wiki_nat]: https://en.wikipedia.org/wiki/Network_address_translation
[previous_post]: /post/kvm_hello_world#the-options
