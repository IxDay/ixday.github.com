+++
title = "Unable to find LVM volume... with LVM on top of Luks"
date = 2015-08-19
categories = ["OS"]
tags = ["admin"]
url = "post/unable_to_find_lvm_volume_with_lvm_on_top_of_luks"
+++


Following [this article]({filename}/2015-08-18.debian-install-2.md)
I have LVM on top of Luks for my system.
When I boot I encounter the following message:

```text
  Volume group "vg_ssd" not found
  Skipping volume group vg_ssd
Unable to find LVM volume vg_ssd/lv_root
```

It appears that LVM is started before I open the crypted partition and display this error.
To fix this we will manipulate the initramfs \o/.

The issue is in the file `/usr/share/initramfs-tools/scripts/local-top/cryptroot`
<!--more-->
which starts like this:
```bash
#!/bin/sh

PREREQ="cryptroot-prepare"

#
# Standard initramfs preamble
#
prereqs()
{
	# Make sure that cryptroot is run last in local-top
	for req in $(dirname $0)/*; do
		script=${req##*/}
		if [ $script != cryptroot ]; then
			echo $script
		fi
	done
}

case $1 in
prereqs)
	prereqs
	exit 0
	;;
esac

# source for log_*_msg() functions, see LP: #272301
. /scripts/functions
```

Line 10 to 16 it says that cryptroot is run last, in the same directory there is a `lvm2`
script. So, what is happening is that cryptroot is launched after lvm2 which is not what we want.
To fix this remove lines 10 to 16 and replace them with `echo "$PREREQ"`

```bash
#!/bin/sh

PREREQ="cryptroot-prepare"

#
# Standard initramfs preamble
#
prereqs()
{
	echo "$PREREQ"
}

case $1 in
prereqs)
	prereqs
	exit 0
	;;
esac

# source for log_*_msg() functions, see LP: #272301
. /scripts/functions
```

At the init step files are taken in alphabetical order so `cryptroot` will come before
`lvm2` (c < l).

Now we just regenerate the initrd files: `update-initramfs -u -k all` reboot and voila, it works :D


