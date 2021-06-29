---
date:       2017-05-01
title:      "Start libvirt VM as unprivileged user"
tags:       ["admin", "bash"]
categories: ["os"]
url:        "post/unprivileged_libvirt"
---

Quick post for starting a VM inside libvirt as a non-root user. Also contains
some useful snippets.

I want to start an alpine virt iso ([from here](https://alpinelinux.org/downloads/))
inside kvm through `libvirt`. But I am sick to run all my virsh commands prefixed with sudo.

**DISCLAIMER**: This will contain some of my conclusions with my partial understanding
of those tools. I am quite sure that all of this can be improved, but I don't
have time to invest on this for the moment.


## Packages and doc

To do this I have to install all the libvirt stuff and dependencies. I recently
changed my distrib for an Archlinux, relevant articles can be found here:

- [KVM](https://wiki.archlinux.org/index.php/KVM)
- [QEMU](https://wiki.archlinux.org/index.php/QEMU)
- [libvirt](https://wiki.archlinux.org/index.php/libvirt)

Installed packages are:

- `libvirt`
- `qemu`
- `virt-install`

Libvirt is a wrapper around virtualization solutions, you can use it to
start containers and VMs through different tools and hypervisors. The `virt-install`
package is just a plugin of `libvirt` to manage VMs installations. You can
perform all those operations directly through `QEMU` if you want.

## Setup

*Note: libvirt command line interface is `virsh`, because fuck logic, so
you don't get lost in the following.*

To make all this work I had to add myself to those groups: `kvm`, `libvirt`.

Libvirt comes with a definition of a default network, you can see it by typing:
`sudo virsh net-dumpxml default`. This will create and start a virtual bridge
named `virbr0`, for my case (check the name it will be important later on), when
the network is started.

What we will do is reuse this bridge for the VMs we will start in our user space.
To perform that we have to ensure that this network is started by libvirt at
startup:

```bash
# start libvirt at computer startup
sudo systemctl enable libvirtd
# ensure network is started by default
sudo virsh net-autostart default
```

Now, just one last config, QEMU provides some ACL which will allow users to
interact with parts of it, we will modificate it to add access to the virtual
bridge. We also have to set suid on the qemu bridge command so a simple user
can get permissions for setting up bridge configuration.

```bash
# create ACL directory and set the value
sudo mkdir -p /etc/qemu && echo "allow virbr0" | sudo tee /etc/qemu/bridge.conf
# enable suid on qemu bridge tool
locate qemu-bridge-helper | sudo xargs chmod +s
```

**BEWARE:** The default network bridge name has been used here (`virbr0`),
check if it is the one from your default network.

## Shoot!

We just have to start our new VM now, be sure to have correctly logged out and
in again before running the following commands (it makes the group changes effective).

```bash
# download the ISO
wget https://nl.alpinelinux.org/alpine/v3.5/releases/x86_64/alpine-virt-3.5.2-x86_64.iso
# start it
virt-install			# create a new VM through libvirt \
	--virt-type kvm		# use the kvm hypervisor \
	--name alpine		# VM name  \
	--memory 1024		# allocate 1024MB of RAM \
	--disk size=10		# allocate a volume of 10G for disk storage of the VM \
	--noautoconsole		# I don't want to start a gui on top of my VM (this part has to be investigated) \
	--cdrom alpine-virt-3.5.2-x86_64.iso	# the disk iso for the installation \
	--network bridge=virbr0,model=virtio	# the tricky part now, I define the bridge to attach my VM to \
						# the model=virtio is something which makes it work properly, \
						# also have to investigate it

virsh console --domain alpine # attach our console to the VM console
```

You should have a running VM and be able to control it without any `sudo` command.
Yay \o/

## Some snippets

I use some shell snippets to control the destruction of my VM.

```bash
vm=alpine
# stop a dedicated vm
virsh list | awk '$2 ~ /'"${vm}"'/ {system("virsh destroy " $2)}'
# delete it
virsh list --all | awk '$2 ~ /'"${vm}"'/ {system("virsh undefine " $2)}'
# delete all the volumes
virsh vol-list default | awk 'NR > 2 && NF > 0 {system("xargs virsh vol-delete --pool default " $1)}'
```
