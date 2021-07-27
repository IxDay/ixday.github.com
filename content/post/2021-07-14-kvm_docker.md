---
title:      "KVM with Docker bridge"
date:       2021-07-14
categories: ["Tuto"]
tags:       ["admin", "cli"]
draft:      true
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
qemu-system-x86_64 -m 4G -smp 4 -cpu host -accel kvm \
	-monitor unix:$@,server,nowait \
	-netdev "bridge,id=user.1,br=docker0" \
	-device "virtio-net,netdev=user.1" \
	-hda <your_qcow_file>
```


[wiki_nat]: https://en.wikipedia.org/wiki/Network_address_translation
[previous_post]: /post/kvm_
