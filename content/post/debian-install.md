+++
title = "Debian Install"
date = "2014-10-12"
categories = ["OS"]
tags = ["bash", "admin"]
+++

That's a fact: I can not install a new system without getting into troubles!

So, here is a small tutorial in which I will aggregate the main issues I
encountered and how I solved them.

### Creating a bootable USB key on a Mac

```bash
# plug your USB key, then find it with
diskutil list

# unmount the usb key (it is mandatory), where X in diskX is the number assigned
# to your USB you retrieved with the previous command
diskutil unmountDisk /dev/diskX

# if you are making a bootable usb key for a Mac run this command,
# debian.img will be the output and debian-testing-amd64-CD-1.iso is
# the iso you retrieved from internet https://www.debian.org/CD/http-ftp/
hdiutil convert -format UDRW -o debian.img debian-testing-amd64-CD-1.iso

# if you are making a Mac bootable usb key (X is still the disk number)
sudo dd if=./debian.img.dmg of=/dev/rdiskX bs=1m

# otherwise (X is still the disk number)
sudo dd if=./debian-testing-amd64-CD-1.iso of=/dev/rdiskX bs=1m

# to finish eject
hdiutil eject /dev/diskX
```

### Messing up during the installation

I use to have a lot of issues during the installation, such as:

 * No installable kernel
 * Can not install grub
 * Whatever can happen during the installation

So, it is possible to fix that during the installation process.
At the final step, just before the reboot press ctrl+alt+f3, then when prompted
press enter and you will get into a shell. It is also possible to perform that
with the Advanced Options -> Rescue Mode.

To fix the errors we will chroot into the new system and add the missing part by
hand.

#### Chroot into an other system

First find the partition on which your system has been installed

```bash
# parted will print you the partition of each device
parted /dev/sda print
parted /dev/sdb print
```

Once the partition is found you have to mount it

```bash
# create a directory for the mounted partition
mkdir /mnt/sda1

# mount it
mount /dev/sda1 /mnt/sda1

# bind the main parts
mount -o bind /dev /mnt/sda1/dev
mount -o bind /dev/pts /mnt/sda1/dev/pts
mount -o bind /proc /mnt/sda1/proc
mount -o bind /run /mnt/sda1/run
mount -o bind /sys /mnt/sda1/sys

# chroot in the system
chroot /mnt/sda1 /bin/bash
```

#### Dealing with the system

Actualize the file `/etc/apt/sources.list` with the following
[site](http://debgen.simplylinux.ch/)

```bash
# update
apt-get update

# install the kernel, for me it is an amd64 architecture, to find yours just run
# apt-cache search linux-image and choose the one for your needs
apt-get install linux-image-amd64

# install grub2, when installing it will ask for the device on which you want
# grub to be installed, choose the device not the partition, here it will be
# /dev/sda
apt-get install grub2

# update grub just in case
update-grub

# do not forget to initialize the password
passwd

```

First step is done, and you will now be able to boot on your new system

### Rebooting and first configuration

#### Set up the locales

```bash

# install the package
apt-get install locales

# set the variables up, select the ones you want with space
dpkg-reconfigure locales

```

#### Set up the keyboard

```bash

# install the package
apt-get install console-data

# for me it is a french keyboard
loadkeys fr-latin

```

Test if the configuration works for you, then you can save it by adding to
`/etc/rc.locals`

```bash
# the path to the keymap is displayed when you use the loadkeys command
/usr/bin/loadkeys /usr/share/keymaps/i386/azerty/fr-latin9.kmap.gz
```

#### <a name="users"></a>Users

```bash
# install sudo
apt-get install sudo

# add new user called foo with a home folder (-m),
# users as first group (-g group_name), sudo as additionnal group (-G group_name)
# and bash as login shell
useradd -m -g users -G sudo -s /bin/bash foo

# change user password
passwd foo
```

### Display Manager

Here I will use i3 on top of xorg


#### Init

```bash
# install xorg and i3
apt-get install xorg i3
```

Create a file at the root of the user `~/.bash_profile` if you are using bash
`~/.zprofile` if using zsh, if you have another login shell please refer to the
dedicated doc.

In our configuration it will be bash due to the [user creation](#users).
Add `startx` at any point of the file. This will launch xserver at login.
Then we want to launch i3. To do that add `exec i3` at any point of the file
`~/.xinitrc`

#### Set X keyboard layout

To have the correct layout for X add the above commands in your `.xinitrc` file.

```bash
# reset the options
setxkbmap -option

# I only add the option to quit X by pressing ctrl+alt+backspace
setxkbmap -layout fr -variant latin9 -option terminate:ctrl_alt_bksp
```

To see all available options you can type `localectl list-x11-keymap-options`

To see current configuration type `setxkbmap -query`


Here is the basic configuration for my linux. Trying to allocate the main issues
I had.
