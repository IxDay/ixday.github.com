+++
title = "Alpine on a Raspberry"
date = 2020-01-05
categories = ["OS"]
tags = ["admin"]
url = "post/alpine_raspberry"
+++

Lately, I decided to re-install my old [Raspberry Pi](https://www.raspberrypi.org/)
(version 1, yes I'm that old) to create a small home server.
Since I am using [Alpine](https://alpinelinux.org/) and liking it quite a lot,
I wanted to install it on my new little toy project.

I followed documentations and tutorials (
[here](https://wiki.alpinelinux.org/wiki/Raspberry_Pi) and
[here](https://www.rigon.tk/documentation/alpine-raspberry-pi))
and finally succeeded.
I am now writing it down so I remember what has been achieved.

Setup SD Card
-------------

This one was quite tricky because if anything is not exactly what is expected
your raspberry will never boot! I am copy-pasting here the commands executed
to properly setup the SD card partitions
(I took the fdisk automation from [this stack overflow answer][1]):

```bash
# replace /dev/...  with your sdcard entry
sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk /dev/mmcblk0
  o # clear the in memory partition table
  n # new partition
  p # primary partition
  1 # partition number 1
    # default - start at beginning of disk
  +256M # 256 MB boot parttion
  n # new partition
  p # primary partition
  2 # partion number 2
    # default, start immediately after preceding partition
    # default, extend partition to end of disk
  t # change partition type
  1 # for partition 1
  c # partition type is W95 FAT32 (LBA)
  p # print the in-memory partition table
  w # write the partition table
  q # and we're done
EOF

# here as well replace with proper path (lsblk will help you)
mkfs.vfat /dev/mmcblk0p1
mkfs.ext4 /dev/mmcblk0p2
```

Now, we have to run some commands because the raspberry I am using need so.

```bash
cd $(mktemp -d) && mount /dev/mmcblk0p1 . && cd .
curl -o - http://dl-cdn.alpinelinux.org/alpine/v3.11/releases/armhf/alpine-rpi-3.11.2-armhf.tar.gz | \
	tar xzf -
rm boot/*-rpi2
mv boot/* .

cat > config.txt <<-EOF
	disable_splash=0
	boot_delay=0
	gpu_mem=256
	gpu_mem_256=32
	kernel=vmlinuz-rpi
	initramfs initramfs-rpi
EOF
```

Once this is done, the SD card is ready to be plugged and the Raspberry to be booted.

First boot and in-memory setup
------------------------------

When the booting is done, you will get a prompt, user is `root` and no password
should be required.

```bash
# run the setup and answer the various questions
setup-alpine

# be sure to have properly installed ssh as it ease the administration,
# we now create a user to be able to connect in
adduser -G wheel foo

# the clock may display some errors, better fix it
ntpd -q -p ptbtime1.ptb.de

# we can now commit the changes
lbu commit -d
# if you want reboot in order to ensure that changes are properly saved
reboot
```

Doing the persistent install
----------------------------

The `sys` install is the next phase, it will allow your installation to survive
reboots (as you would expect from a regular machine).

```bash
# we mount the ext4 partition
cd $(mktemp -d) && mount /dev/mmcblk0p2 . && cd .

# we copy from the previous commited image (you can ignore the syslinux/extlinux errors)
setup-disk -o /media/mmcblk0p1/MYHOSTNAME.apkovl.tar.gz $(pwd)

# now we update the fstab
echo "/dev/mmcblk0p1 /media/mmcblk0p1 vfat defaults 0 0" >> ./etc/fstab

```

We now need to change the boot partition a bit in order to switch to our newly
installed system.

```bash
mount -o remount,rw /media/mmcblk0p1 && cd media/mmcblk0p1
sed -i '$ s/^/root=\/dev\/mmcblk0p2 /' /media/mmcblk0p1/cmdline.txt

mkdir kernel-installer
mv System.map-rpi config-rpi initramfs-rpi vmlinuz-rpi kernel-installer
cd - && cp -v ./boot/* /media/mmcblk0p1
```

__Everything should now works on next boot! You have successfully installed Alpine
on a Raspberry!__

[1]: https://superuser.com/questions/332252/how-to-create-and-format-a-partition-using-a-bash-script#answer-984637

