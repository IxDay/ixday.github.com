+++
title = "Debian Install 2"
date = "2015-08-18"
categories = []
tags = []
+++

Today I want to install debian on my personal server.
And I want to crypt the FS using Luks, and add LVM on top.

I want it to look like that (like this scheme, a lot of reference
come from the [arch wiki](https://wiki.archlinux.org/index.php/Dm-crypt/Encrypting_an_entire_system)
<!--more-->
```text
+-----------------------------------------------+ +----------------+
|Logical volume1        | Logical volume2       | |                |
|/dev/vg_ssd/lv_swap   | /dev/vg_ssd/lv_root  | | Boot partition |
|_ _ _ _ _ _ _ _ _ _ _ _|_ _ _ _ _ _ _ _ _ _ _ _| |	               |
|                                               | |                |
|		LUKS encrypted partition                | |                |
|               /dev/sda2                       | |   /dev/sda1    |
+-----------------------------------------------+ +----------------+
```

## Prepare disks

I use parted for this configuration, I left 1MB at the start of the device
to be compliant with `parted <disk> align-check optimal <partition>`

```bash
apt install cryptsetup lvm2

parted /dev/sda mkpart primary 1MB 201MB # boot partition
parted /dev/sda mkpart primary 201MB 100% # lvm partition
parted /dev/sda toggle 1 boot
parted /dev/sda toggle 2 lvm

cryptsetup -c aes-xts-plain -y -s 512 luksFormat /dev/sda2
# enter passphrase

# this command will open the crypted /dev/sda2 partition and link it to ssd in /dev/mapper
cryptsetup luksOpen /dev/sda2 ssd

lvm pvcreate /dev/mapper/ssd
lvm vgcreate vg_ssd /dev/mapper/ssd

lvm lvcreate -L 4GB -n swap vg_ssd
lvm lvcreate -l 100%FREE -n root vg_ssd

# now format the partitions
mkfs.ext4 /dev/sda1 # boot partition
mkfs.ext4 /dev/mapper/vg_ssd-lv_root
mkswap /dev/mapper/vg_ssd-lv_swap
```

## Prepare base system

In order to prepare the base system I will use the `debootstrap` tool,
which will take care of creating base files. Then I will `chroot` in the new
system in order to setup packages and boot initialisation

```bash
# mount new filesystem and debootstrap
mount /dev/mapper/<vg_name>-root /mnt
mkdir /mnt/boot
mount /dev/sda1 /mnt/boot

debootstrap --arch=amd64 stable /mnt

# copy somefiles
cp /etc/resolv.conf /mnt/etc/
cp /etc/network/interfaces /mnt/etc/network

# bind and chroot
mount -o bind /dev /mnt/dev
mount -o bind /dev/pts /mnt/dev/pts
mount -o bind /proc /mnt/proc
mount -o bind /run /mnt/run
mount -o bind /sys /mnt/sys

chroot /mnt/ /bin/bash
```

## Configuration

First install vim! `apt-get install vim`

Then create the fstab file, I use the UUID notation for this part.
Basically, I run two commands in order to get the needed informations

```bash
ls -l /dev/disks/by-uuid

total 0
lrwxrwxrwx 1 root root 10 Aug 19 12:52 2dfb2c7a-e99b-4073-9267-8e517bc0ce82 -> ../../sda1
lrwxrwxrwx 1 root root 10 Aug 19 13:09 59da986c-ef82-485e-bc8b-f36cc440273c -> ../../sda2
lrwxrwxrwx 1 root root 10 Aug 19 12:53 de3c4d79-beb9-427e-9dea-d25726a5f492 -> ../../dm-2
lrwxrwxrwx 1 root root 10 Aug 19 12:52 f7c9d9f4-671e-45be-97a9-be775196545e -> ../../dm-1
lrwxrwxrwx 1 root root 10 Aug 19 12:40 fcc7394f-6865-4ac5-989a-f6c58dc129d5 -> ../../dm-0

ls -l /dev/mapper
total 0
crw------- 1 root root 10, 236 Aug 19 12:37 control
lrwxrwxrwx 1 root root       7 Aug 19 12:51 ssd -> ../dm-0
lrwxrwxrwx 1 root root       7 Aug 19 12:52 vg_ssd-lv_root -> ../dm-2
lrwxrwxrwx 1 root root       7 Aug 19 12:53 vg_ssd-lv_swap -> ../dm-1

```

Which will give me the following /etc/fstab.
```bash
UUID=de3c4d79-beb9-427e-9dea-d25726a5f492   /       ext4    defaults    0  1
UUID=2dfb2c7a-e99b-4073-9267-8e517bc0ce82   /boot   ext4    defaults    0  2
UUID=f7c9d9f4-671e-45be-97a9-be775196545e   none    swap    defaults    0  0
```

Edit the /etc/apt/sources.list file to set up basic packages and update package list: `apt-get update`
```bash
deb http://ftp.us.debian.org/debian stable main contrib non-free
```


Install needed packages

```bash
# install the package
apt-get install locales console-data keyboard-configuration

# set the variables up, select the ones you want with space
dpkg-reconfigure locales

# setup password
passwd

# change hostname
echo "suika" > /etc/hostname

apt-get install lvm2 cryptsetup

# this one is specific to my mothercard
apt-get install firmware-realtek
```

## Grub

Now we want to make the new system bootable, to perform this the system needs a kernel
`apt-get install linux-image-amd64`

Set up the /etc/crypttab file, which will tell how to map the crypted partition
```bash
# <target name> <source device>         <key file>      <options>
ssd             /dev/sda2               none            luks
```

Install grub on /dev/sda (on disk, not on partition) `apt-get install grub2`

Load keyboard on initramfs in order to avoid keyboard layout collision when opening luks volume.

In file `/etc/initramfs-tools/initramfs.conf`
```bash
#
# KEYMAP: [ y | n ]
#
# Charger une configuration de clavier à l'étape d'initramfs.
#

KEYMAP=y

```

Update initramfs `update-initramfs -u` to apply change on initrd.

## The end!

Exit the `chroot` unmount partitions and reboot on your new fresh installed system.

If at startup you get an error message like:
```text
  Volume group "vg_ssd" not found
  Skipping volume group vg_ssd
Unable to find LVM volume vg_ssd/lv_root
```
You can look at [this article]({filename}/2015-08-19.unable-to-find-lvm-volume-with-lvm-on-top-of-luks.md) which explain how to fix this issue


