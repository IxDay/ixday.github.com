---
title:      "My_linux_config"
date:       2020-11-12
draft:      true
categories: []
tags:       []
---

Here are all the packages I installed for my Archlinux config


- base
- base-devel
- linux
- linux-firmware
- lvm2
- intel-ucode
- nmap
- man-db
- gimp
- gnome-screenshot

- zsh
- git
- tig
- tmux
- htop
- neovim
- code
- hugo
- go
- go-tools
- firefox
- gopass
- docker
- podman #
- wireguard-dkms
- wireguard-tools
- systemd-resolvconf
- jq
- strace
- bind
- gnome-keyring
- evince

- lightdm-gtk-greeter
- xclip
- xdg-utils
- xorg-server
- xorg-xrdb
- ttf-droid
- redshift
- pkgfile
- openssh
- network
- cinnamon
- file-roller
- arc-gtk-theme
- arc-icon-theme
- libva-intel-driver

- network-manager-applet
- wmctrl
- youtube-dl
- blueberry
- spotify
- gnome-mplayer

- retroarch
- retroarch-assets-xmb
- retroarch-assets-ozone
- libretro-desmume
- termite
- gtk-recordmydesktop



aur
- yay
- direnv
- sublime-text-3
- franz-bin
- google-chrome
- age
- epson-inkjet-printer-escpr
- skypeforlinux-stable-bin
- networkmanager-wireguard
- scrcpy
- android-tools
- android-udev # https://github.com/M0Rf30/android-udev-rules

vscode
- base16 theme generator
- go
- vim
- editorconfig
- shellcheck



mkdir -p ~/.local/share/gnupg
chmod 700 ~/.local/share/gnupg
gpg --fingerprint
gpg --keyserver pool.sks-keyservers.net --recv-keys 8FD3D9A8D3800305A9FFF259D1742AD60D811D58

modprobe -r pcspkr
http://www.thinkwiki.org/wiki/How_to_disable_the_pc_speaker_(beep!)
https://wiki.archlinux.org/index.php/Kernel_module#Blacklisting
echo -e "\07"
