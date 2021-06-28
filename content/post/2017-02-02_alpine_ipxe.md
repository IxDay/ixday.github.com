+++
categories = ["tuto"]
date = 2017-02-02
tags = ["admin"]
title = "Alpine iPXE"
url = "post/alpine_ipxe"
+++

For a personal project I want to iPXE boot alpine. I did not found any step by
step guide, so I ended up testing multiple solutions until it works. This
post is a "copy" of the answer I made [here](https://github.com/antonym/netboot.xyz/issues/30#issuecomment-276722892) for the [netboot.xyz project](https://netboot.xyz).

To fix this issue I used a bunch of thread and resources but here are the three
main entry used:

<!--more-->
- Alpine documentation to create a custom ISO image:
[https://wiki.alpinelinux.org/wiki/How_to_make_a_custom_ISO_image](https://wiki.alpinelinux.org/wiki/How_to_make_a_custom_ISO_image)
- Alpine documentation for the network boot:
[https://wiki.alpinelinux.org/wiki/PXE_boot](https://wiki.alpinelinux.org/wiki/PXE_boot)
- This post on the Alpine forum from a guy who managed to successfully boot from PXE:
[https://forum.alpinelinux.org/forum/installation/boot-pxe](https://forum.alpinelinux.org/forum/installation/boot-pxe)

Now here is the process to make all of this work:

- clone the [alpine-iso git repo](http://git.alpinelinux.org/cgit/alpine-iso/)
- cd to the directory and create the following two files:
  - `alpine-pxe.packages` which will be empty (those packages are installed in the iso not in the initrd)
  - `alpine-pxe.conf.mk` with the following content:

         ```text
         ALPINE_NAME     := alpine-pxe
         KERNEL_FLAVOR   := grsec
         INITFS_FEATURES := ata base bootchart squashfs ext4 usb virtio network dhcp
         MODLOOP_EXTRA   :=
         ```
    some options may not be needed, I did not had the time to check this correctly,
	`virtio` and `network` are needed according to
	[the wiki](https://wiki.alpinelinux.org/wiki/PXE_boot#Using_pxelinux_instead_of_gPXE)

- create the image following the instructions [here](https://wiki.alpinelinux.org/wiki/How_to_make_a_custom_ISO_image) and pass your new profile: `make PROFILE=alpine-pxe`
- start a simple http server in `isotmp.alpine-pxe/isofs/boot/` the ipxe boot will need the two following files:
  - `initramfs-grsec`
  - `vmlinuz-grsec`
- create your ipxe script with an adaptation of the following to your own url

    ```text
    #!ipxe

    dhcp

    set base-url http://192.168.122.0:5050
    set kernel-params ip=dhcp modules=loop,squashfs,usb-storage nomodeset

    kernel ${base-url}/vmlinuz-grsec ${kernel-params}
    initrd ${base-url}/initramfs-grsec

    boot
    ```
  Here again I am not sure that all the options are needed, still have to perform some tests.
- start your network boot!
