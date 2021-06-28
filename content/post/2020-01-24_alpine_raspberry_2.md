+++
title = "Alpine on a Raspberry (part 2)"
date = 2020-01-24
categories = ["OS"]
tags = ["admin"]
url = "post/alpine_raspberry_2"
+++

I had to buy a new raspberry (the old one did not survive :p), a version 3 A+.
Now, I will go over again the installation process. Since I encountered a few
new problems.

Setup SD Card
-------------

This is exactly the same as the previous version. So I let you check out the other
article.

First boot and in-memory setup
------------------------------

Here a few changes, using a version 3A+ requires to setup the wifi.
Following this [wiki](https://wiki.alpinelinux.org/wiki/Connecting_to_a_wireless_access_point),
here are the steps I followed:

```bash
# first install the necessary packages
apk add wpa_supplicant

# setup the config file
wpa_passphrase 'ExampleWifi' 'ExampleWifiPassword' > /etc/wpa_supplicant/wpa_supplicant.conf

# start the service
wpa_supplicant -B -i wlan0 -c /etc/wpa_supplicant/wpa_supplicant.conf

# request a dhcp lease
udhcpc -i wlan0

# enforce the service startup at boot
rc-update add wpa_supplicant boot
```

__Fix the clock:__ I found a better solution to the clock problem in alpine.
This is a mix of [this article](https://wiki.alpinelinux.org/wiki/Classic_install_or_sys_mode_on_Raspberry_Pi#Installation)
and [this one](https://gitlab.alpinelinux.org/alpine/aports/issues/8093).

```bash
# install package
apk add chrony

# add the service to the boot
rc-update add chronyd boot

# create a file to signal usage of this hack
touch /etc/init.d/.use-swclock

# we need to insert a hack in the file
# https://github.com/OpenRC/openrc/blob/master/sh/init.sh.Linux.in#L53
sed -i '53iif [ -e /etc/init.d/.use-swclock ]; then\n' /lib/rc/sh/init.sh
sed -i '54i\t“$RC_LIBEXECDIR/sbin/swclock” /etc/init.d\n' /lib/rc/sh/init.sh
sed -i '55ifi' /lib/rc/sh/init.sh
```

Once this done, you can stick to the previous install instructions.

Doing the persistent install
----------------------------

The steps are the same as the previous one except for the boot part.
I had to use those instructions [from the wiki](https://wiki.alpinelinux.org/wiki/Classic_install_or_sys_mode_on_Raspberry_Pi#Installation).

```bash
# fix the boot directories
rm -f /media/mmcblk0p1/boot/*

cd /mnt
rm boot/boot
mv boot/* /media/mmcblk0p1/boot/
rm -Rf boot
mkdir media/mmcblk0p1
ln -s media/mmcblk0p1/boot boot

# fix the fstab
cat > etc/fstab <<- EOF
	/dev/mmcblk0p2 /                ext4 defaults 0 0
	/dev/mmcblk0p1 /media/mmcblk0p1 vfat defaults 0 0
EOF

# fix the kernel boot flags
sed -i 's/^/root=\/dev\/mmcblk0p2 /' /media/mmcblk0p1/cmdline.txt
```

