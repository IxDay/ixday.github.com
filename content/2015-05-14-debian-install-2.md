Title: Debian install 2
Date: 2015-05-14
Category: Tuto
Tags: bash, adminsys
Status: draft

After many days of configuration, I finally complete the installation of my
perfect "workstation" ;).


## source.list
I have to directory for this, first the 
[preferences](https://github.com/IxDay/config_files/tree/new_conf/preferences.d),
then the
[sources.list](https://github.com/IxDay/config_files/tree/new_conf/source.list.d).

Just run an `aptitude update`, and you will some missing gpg keys, now run
```
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com <key_number>
```
To have your keys installed.


## sudo 

A sudo file is important in order to correctly manage a computer, a simple
`aptitude install sudo` will give you the tool. 

Then uncomment the following line
```
# Allow members of group sudo to execute any command
%sudo   ALL=(ALL:ALL) ALL

```

Then add your user to the group sudo `gpasswd -a <your_user> sudo`


## install all the packages o/

Here is the list of packages I install by default, you can click on them to go
to the proper section
* [rxvt-unicode-256color](#shell)
* [zsh](#shell)
* [vim-gtk](#vim)
* [xserver-xorg](#xserver-i3-gtk), [i3](#xserver-i3-gtk)
* [fonts-droid](#xserver-i3-gtk), [feh](#theming)
* [gtk2-engine-murrine](#xserver-i3-gtk), 
[gtk2-engine-pixbuf](#xserver-i3-gtk), [libgtk2.0-0](#xserver-i3-gtk),
[libgtk-3-0](#xserver-i3-gtk)
* [libnotify-bin](#notifications), [dunst](#notifications)
* xclip, jq
* [git](#git), [keychain](#git)
* [google-chrome-stable](#google-chrome)
* [docker.io](#docker)
* 


## <a name="shell"></a> Shell

First thing first, a computer is not a real computer without a good shell.
So, I install zsh, then I change my default shell `chsh -s $(which zsh)`, and
to finish [Oh My Zsh](rxvt-unicode-256color). I have a custom 
[prompt](https://github.com/IxDay/config_files/blob/new_conf/max.zsh-theme) 
which is placed in `.oh-my-zsh/themes/`. The promp is loaded in my 
[zshrc file](https://github.com/IxDay/config_files/blob/new_conf/zshrc)
and some [aliases](https://github.com/IxDay/config_files/blob/new_conf/zshrc) 
are sourced for convenience.


## <a name="vim"></a> Vim

I use vim for quite everything so I have some customization here.
I directly install the gtk version in order to have support for the clipboard
see [here](http://stackoverflow.com/questions/11489428/how-to-make-vim-paste-from-and-copy-to-systems-clipboard)
for the explanation.

First a [.vimrc](https://github.com/IxDay/config_files/blob/new_conf/vimrc), 
the plugin system is based on 
[Vundle](https://github.com/gmarik/Vundle.vim). To install it just run
`git clone https://github.com/gmarik/Vundle.vim.git ~/.vim/bundle/Vundle.vim`
and `vim +PluginInstall +qall`

I also use a global [editorconfig file](https://github.com/IxDay/config_files/blob/new_conf/editorconfig) 
to keep tidy the maximum of files.


## <a name="xserver-i3-gtk"></a> Xserver, i3 and GTK

I use i3 for my window manager directly on top of X, all the configuration can
be found easily on my github, here are the files used to configure my desktop.
* The i3 config directory: [here](https://github.com/IxDay/config_files/tree/new_conf/i3) 
* The i3status config file: [here](https://github.com/IxDay/config_files/blob/new_conf/i3status.conf)
* The gtkrc2 file: [here](https://github.com/IxDay/config_files/blob/new_conf/gtkrc-2.0)
* The settings.ini file for gtk3: [here](https://github.com/IxDay/config_files/blob/new_conf/gtkrc-2.0)
* The Xresource file: [here](https://github.com/IxDay/config_files/blob/new_conf/Xresources)

In order to launch X with i3 the following files are required in your home 
directory:
* [.zprofile](https://github.com/IxDay/config_files/blob/new_conf/zprofile) is
the first file automatically sourced at the login, it launch the Xserver for
the session.
* [.xinitrc](https://github.com/IxDay/config_files/blob/new_conf/xinitrc)
contains the configuration for Xserver, it loads the keyboard layout, 
initialize some properties from the `.Xresource` file (colors, fonts, etc...),
and finally load the i3 window manager.

There is some issues with the dmenu provided by suckless-tools. 
In order to support xft font I have reinstalled it from the minos repository, 
mentionned [here](https://wiki.archlinux.org/index.php/Dmenu#Fonts)


### <a name="notification"></a>Notification

For the notifications I use dunst, which is a notification service, it is
started at i3 startup and will display in a configurable way the notifications
from the system. It needs `libnotify-bin` to run, and has a 
[configuration file](https://github.com/IxDay/config_files/blob/new_conf/dunstrc)


### <a name="theming"></a>Theming
I have installed the [Vertex theme](https://github.com/horst3180/Vertex-theme),
and the [Awoken White icon pack](https://github.com/IxDay/config_files/blob/new_conf/AwOken-2.5.zip).

For the last one some packages are required in order to configure it correctly:
* imagemagick
* zenity

I also use `feh` for managing the desktop wallpaper, it is launched at the 
Xserver startup from the `.xinitrc` file.


## <a name="git"></a>git
The versionning tool!

My configuration require here the two following files:
* [gitconfig](https://github.com/IxDay/config_files/blob/new_conf/gitconfig)
* [gitignore_global](https://github.com/IxDay/config_files/blob/new_conf/gitconfig)

After copying those files you can create a rsa key:
```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

And to manage the keys passphrase I use `keychain` in an alias
```
alias keychain_default="eval $(keychain --eval --agents ssh -Q --quiet id_rsa)"
```
It is possible to configure multiple keys, in order to separate some services.
The archlinux wiki has a good article about the SSH keys management: 
[here](https://wiki.archlinux.org/index.php/SSH_keys)


## <a name="google-chrome"></a>Google Chrome

Why Chrome, and not Chromium? I had a lot of issues with Chromium, on flash,
with google-talk and so on. It was too complicated and never works, 
I wasn't able to figure out why so I give up on this


## <a name="docker"></a>Docker
Docker is now vastly known by the community of developpers. For the moment,
I do not think it is a good production tool. But, for the development it is 
awesome. 

Simply install the package and add your user to the docker group 
`sudo gpasswd -a <your_user> docker`


