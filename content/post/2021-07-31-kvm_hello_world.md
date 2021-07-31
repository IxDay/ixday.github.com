---
title:      "Kvm Hello World"
date:       2021-07-31
categories: ["Tuto"]
tags:       ["kvm", "cli"]
url:        "post/kvm_hello_world"
---

Let's start a small serie of post on KVM. It is a tool I sometime need to use
when Docker is too simple for the problem at hand, or if I really need virtualization
(to emulate another arch or a totally different system).

I will start by giving a series of snippets to quick start your KVM usage. I am
doing this because I often struggle to retrieve simple information when I want
to come back on this kind of technology. Most of the articles are deep diving
on a specific option and I always wonder what are all those options passed to
the CLI.

For this serie of articles I will use the [Qemu][qemu_website] CLI since it is
the simplest and most used tool out there.

## The hard disk drive

Qemu mostly rely on the [qcow][qcow_wiki] format for the data persistence.
It offers a variety of capabilities I will not dig into in this article.
Let's keep things simple, and create a simple disk for our example:

```sh
qemu-img create -f qcow2 hdd.qcow2 10G
```

This will generate a new file which will grow up to 10G. It will be of qcow2 format
as indicated by the option `-f qcow2`.

Another useful option is to create a new file backed by another one. This new
file will only record the difference from the previous one. This is handy when
you have finished the installation of the base image and start doing some configuration.
It will allow you to go back to the initial state easily without the need of
reinstalling everything:

```sh
# let's create our base image
qemu-img create -f qcow2 base.qcow2 10G
# we do what has to be done to have a clean install
qemu-system-x86_64 ...
# now we create a new image based on the previous one
qemu-img create -F qcow2 -f qcow2 -b base.qcow2 hdd.qcow2
```

We will now use `hdd.qcow2` for our VM, keeping `base.qcow2` safe from subsequent
modifications. It will ease our development because we can now come back to
a clean install by running the same command (beware it will erase `hdd.qcow2`):

```sh
# reset the hdd.qcow2 image
qemu-img create -F qcow2 -f qcow2 -b base.qcow2 hdd.qcow2
```

## The options

Qemu comes with a lot of options, making it highly configurable. However, it
may be overwhelming at first. I will share here the default values I use when
starting a VM:

```sh
qemu-system-x86_64 -m 4G -smp 4 -cpu host -accel kvm \
	-monitor "unix:monitor.sock,server=on,wait=off" \
	-serial "chardev:serial0" \
	-chardev "socket,id=serial0,path=console.sock,server=on,wait=off" \
	-netdev "user,id=user.1" \
	-device "virtio-net,netdev=user.1" \
	-hda hdd.qcow2 -display none
```

The options are:

- `-m 4G`: amount of RAM, here 4G
- `-smp 4`: number of CPUs, here 4
- `-cpu host`: CPU model, use the same as the host
- `-accel kvm`: virtualization accelerator, here we use KVM
- `-monitor unix:monitor.sock,server=on,wait=off`: we want a qemu console monitor,
  this helps managing the VM and allow the user to send commands
  (to shutwdown the system for example).
  This option contains a few instructions, because you will need a two way communication link.
  For this example we will use a unix socket, but it can also be a regular socket
  which can be reachable over the network (I will show an example of this a bit later).
  - `unix:monitor.sock`: This is the kind of socket we want (unix),
    with the path for the file to create.
  - `server=on` indicates to qemu that it should be a listening socket. You want to connect to it, not the other way around (not sure this sentence makes any sense).
  - `wait=off` tells qemu to not block waiting for the client to connect. If this
    option is set to true, you will have to connect to the socket before qemu can
    go further (it will not boot the VM until you connect).
- `-serial "chardev:serial0"` redirect a serial port to a specified device.
  Here we use a named device called `serial0` and define it as a unix socket with
  the next option. We could have directly defined it as a `unix:console.sock,...`
  and it would have been equivalent. I am just showing a different notation here.
- `-chardev "socket,id=serial0,path=console.sock,server=on,wait=off"` describe
  a character device (this is a synonym for a two way communication link). It will
  be a unix `socket` with id `serial0`, the path for the socket will be `console.sock`. This device has the same options as the monitoring one.
- `-netdev "user,id=user.1"` define a network device with id `user.1` in user mode.
  User mode is the simplest connectivity setup and requires no additional privilege to run. On the downside performance is poor and the guest is not accessible
  from the host, you can check the [official doc][qemu_wiki] for more info.
  You can still define some port forwarding (we will do this a bit later in the post).
- `-device "virtio-net,netdev=user.1"` attach the network device `user.1` using
  the driver `virtio-net`.
- `-hda hdd.qcow2` use the `hdd.qcow2` as a block device
  (`-hda`, `-hdb`, ... should map to `/dev/vda`, `/dev/vdb`, ...).
- `-display none` we do not want to attach something here since we already have
  a unix socket `console.sock`.

__Disclaimer__: `server,nowait` is equivalent to `server=on,wait=off` and can be found
out there on the internet. This has to be confirmed but it should be the legacy
way of declaring those options.

### Connecting to the control sockets

In this example I defined two unix sockets:
- `console.sock`: this is the serial console of the emulated machine, it will
  behave as a regular shell once you properly connect to it.
- `monitor.sock`: this is the Qemu controller interface, you can use it to control
  the VM (shutdown, inpect, connect stuff).

I wrote [a small post][socat_post] explaining the options of the next command,
check it out if you need details. Here is an example I am using personally to
properly connect:

```sh
socat -,rawer,escape=0x1d unix-connect:console.sock
```

### Some variation

If you do not want to use a unix socket you can directly publish the serial
console to a local tcp port instead. Here is the proper option:

```sh
qemu-system-x86_64 ... \
	-chardev "socket,id=serial0,port=4444,host=localhost,telnet=on,server=on,wait=off" \
	-serial "chardev:serial0"
```

The connection string will become:

```sh
socat -,rawer,escape=0x1d tcp:localhost:4444
```


## The first boot

Unless you are getting an image with a system already installed you will have
to first boot from an [ISO file][iso_wiki]. Since we are using an ISO the
bootable device will be an emulated CD-ROM drive. You need to pass 2 additional
options to Qemu. The boot order with `-boot` option, here we set it to `d` for
drive. The second option is `-cdrom` to pass an ISO which will be mounted in
the guest VM at `/dev/cdrom`.

```sh
qemu-system-x86_64 ... \
	-boot "d" -cdrom "<your_distro>.iso"
```

You will need to remove those options once the installation is done, otherwise,
you will keep booting using the CD-ROM drive.

__Alternatively:__ you can keep the option in the command line, to always run the
same command even after the first installation by prefixing `d` with the `once:`
keyword. It will become:

```sh
qemu-system-x86_64 ... \
	-boot "once:d" -cdrom "<your_distro>.iso"
```

You will need to be sure that installation is properly done because this will
only boot once using the CD-ROM drive.

## Publishing services and make them reachable from the host

Last part of this post will be about reaching the guest VM from the host.
Given the current networking configuration, it will not be possible to connect
through `ssh` or expose a web service.

The solution will be to perform a port forward from the guest to the host.
This can be achieved when defining the network interface by passing an additional
`hostfwd` option. Here is an example using the `-netdev` configuration I
presented with the default options I set in one of the previous section:

```sh
qemu-system-x86_64 ... \
	-netdev "user,id=user.1,hostfwd=tcp::2222-:22" \
	-device "virtio-net,netdev=user.1"
```

This will publish the port 22 from the guest to the port 2222 of the host.
You can combine as much `hostfwd` options as you want to publish all the services
you need.

## That's all folks

This conclude this post with the options I would have loved to know when I started
using Qemu. I hope it could help someone out there and save a lot of time when
dealing with something as complicated as virtualization.


[qemu_website]: https://www.qemu.org/
[qcow_wiki]: https://en.wikipedia.org/wiki/Qcow
[iso_wiki]: https://en.wikipedia.org/wiki/Optical_disc_image
[qemu_wiki]: https://wiki.qemu.org/Documentation/Networking#User_Networking_.28SLIRP.29
[socat_post]: /post/socat_telnet
