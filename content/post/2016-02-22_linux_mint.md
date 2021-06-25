+++
title = "Linux Mint"
date = 2016-02-22
categories = ["OS"]
tags = ["admin"]
+++

After a long journey, I finally found the distribution which fits
all my need: **Linux Mint Debian Edition**

## Why ?

* It is based on debian a distribution which I really like (and not Ubuntu).
* The default tilling (similar to W7) is really comfortable.
* Lot of configuration can be done for shortcuts.
* Support for `systemd` (Yay \o/).
* Beautiful UI out of the box (which I didn't have with TWM).
* And a lot more

## Installation


As I am on Macbook (once upon a time, I was young and dumb), I need a
particular configuration of the partitions for the installation.

Here is how my HDD is partitionned:

- First, the `/boot/efi` partition with the boot flag is about 200MB
- Then the Macintosh partition, in hfs+
- The `/boot`, partition, as I will use LVM I need to have an external
  partition to hold the initram images, with an ext4 filesystem.
- Finally, the LVM partition which is subdivided in two partitions:
    - lvm_root: which will hold the root filesystem
    - lvm_swap: my system swap

Mount the partitions in the right order into `/mnt/target`, follow the
expert disk partioning in the installation wizard.

Install LVM into the new system (with chroot), then fillup the fstab file
to look similar to something like this:

```bash
/dev/mapper/vg_ssd-lv_root  /	    ext4	rw,relatime,discard,data=ordered  0 1
/dev/mapper/vg_ssd-lv_swap  none	swap	defaults                          0	0

proc          /proc	    proc  defaults        0	0
/dev/sda3     /boot	    ext4  defaults        0 2
/dev/sda2     /boot/efi vfat  defaults        0 2
```

## Configuration

Now that we have a running distribution there is some small details to fix.
Here are some piece of configuration I needed to fix to have a satisfying
install.

### Chromium as default

First, I want to use chromium as the default browser, simply install it with
the package manager, then just run the `Preferred Application` program to
setup the default.


### Microphone issue
As I use a Macbook I have an issue with my mic by default. To fix that, just
create the `/etc/modprobe.d/alsa-base.conf` file, and just past the following:

```text
options snd-hda-intel model=mbp101 index=1
```

Reboot, and tadaaa!

### Systemd
By default, on LMDE Betsy (my current installation), the init system is
still sysvinit.
I really like systemd and good news the skeleton of it is already present,
so we just have to explain the system to change.
And... it is simpler as I firstly thought, because there is a package for
that: ``systemd-sysv``.

Just run the installation and you will be good

### Powertop
As systemd is already installed, here is the service file:
`/etc/systemd/system/powertop.service`

```ini
[Unit]
Description=Powertop tunings

[Service]
Type=oneshot
ExecStart=/usr/sbin/powertop --auto-tune
Environment="TERM=xterm"

[Install]
WantedBy=multi-user.target
```

**Note**: I had to install `xterm` in order to make the service work, because
powertop needs a shell at runtime to perform the auto-tune statement.

It is a systemd service which will be loaded at startup, so just enable it:
`systemctl enable powertop`

### Skype
This one is easy, but I also want my beautiful cinnamon skin on it ;)

Just download and install the .deb file from
[here](http://www.skype.com/en/download-skype/skype-for-computer/)

Enable multiple architecture if you are running x86_64 distribution:
`dpkg --add-architecture i386`

And now install some complementary packages:
`apt-get install gtk2-engines-murrine:i386 gtk2-engines-pixbuf:i386`

## Last word

This article will evolve to fit the latest change on my system, stay tuned!
