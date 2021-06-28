+++
date = 2016-06-03
title = "Libvirt through vagrant"
tags = ["admin"]
categories = ["Tuto"]
url = "post/libvirt_through_vagrant"
+++

I will not install VirtualBox! That's all (nor VMWare, don't be ridiculous).
But I really like [Vagrant](https://www.vagrantup.com/), and use it every time
I need something closer to a running machine. So, I dig up the internet and
found that there is an unofficial support of libvirt.

## Installation

I do not remember where I found the documentation to do this or if I did it by
myself, so no link here, just what I do in order to make this work. I use a
debian jessie distribution, so packages name may vary.

```bash
# Install packages
apt-get install qemu-kvm libvirt-daemon-system libvirt-dev zlib1g-dev vagrant

# Install vagrant plugin
vagrant plugin install vagrant-libvirt

# Add user to group
gpasswd -a user_name libvirt

# Reboot in order to load kvm and group change
reboot
```

Now everything is installed and ready to work with the libvirt provider.
In order to verify that everything works, run:
`vagrant init debian/jessie64; vagrant up --provider libvirt`.
It must download a debian/jessie image then run it through vagrant.

## In action

For this example, I will just run a Coreos image, and show you how you can
migrate an existing Vagrantfile to the libvirt provider.

Clone this repo https://github.com/coreos/coreos-vagrant, and follow the
instructions.

Then according to this PR, I made a modification to the Vagrantfile
https://github.com/coreos/coreos-vagrant/pull/290/files, I also added
a valid `box_url`

```ruby

  # if the provider asked is libvirt
  config.vm.provider :libvirt do |v|

    # change the box to a libvirt compatible one
    config.vm.box = "dongsupark/coreos-%s" % $update_channel

    # this line indicates where the box can be found
    config.vm.box_url = "https://atlas.hashicorp.com/dongsupark/boxes/coreos-%s" % $update_channel

    # change the driver name and pass the parameters needed to comply with
    # specifications
    v.driver = "kvm"
    v.memory = $vm_memory
    v.cpus = $vm_cpus

  end
```

Now you can do a `vagrant up` and it will work (it works on my machine :p).

## Bonus

You can migrate Vagrant images to libvirt with another plugin:
`vagrant-mutate`

```bash
# Install the plugin
vagrant plugin install vagrant-mutate

# Mutate an image
vagrant mutate https://atlas.hashicorp.com/debian/boxes/jessie64/versions/8.2.2/providers/virtualbox.box libvirt

# Now the image is available, you can rename it if needed
mv ~/.vagrant.d/boxes/virtualbox ~/.vagrant.d/boxes/debian-VAGRANTSLASH-jessie64

# Then run it
vagrant init debian/jessie64; vagrant up --provider libvirt
```

This is no longer really needed as there is more and more images with libvirt
provider.
