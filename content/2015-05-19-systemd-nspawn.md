Title: Systemd-nspawn
Category: Tuto
Tags: systemd, admin

I am a huge fan of docker for my dev environments, it helps me keeping things
clear and understanding what are the ressources needed for a project.
A few month ago a friend told me that there already is a similar feature on
Linux, and this feature is systemd-nspawn.

## Creating your first container

So like docker I wanted to first start a container. Nspawn has no environment
so everything has to be done "by hand". First, as long as we do not have any
registery we need to retrieve an "image" of the distribution we need.

(all the commands are ran on behalf of the root user)
```
#!bash
apt-get install debootstrap
debootstrap --arch=amd64 stable /tmp/my-debian-machine
```

This is a quite long to retrieve, so you can go and take a coffee during the
download ;)

Now, we have a folder containing the debian basic configuration, we can "spawn"
it.

```
#!bash
systemd-nspawn -D /tmp/my-debian-machine
```

And that's it! You now have your own container up and running.
The main difference is that, it is not running on AUFS, and every change will
be definitive.

Your machine is now visible if you run `machinectl`.

## Booting a container like a real VM

Now, we will see a feature that I didn't find in docker, the ability to boot
your system.

```
#!bash
systemd-nspawn -D /tmp/my-debian-machine -b
```

And a prompt will appear, so to log in just 'type the root password'.
I'm sure you forget to set that before ;). Just stop your container,
relaunch it without boot and set the password.

```
#!bash
# stop the container
machinectl poweroff my-debian-machine
# launch without boot
systemd-nspawn -D /tmp/my-debian-machine

# when you have the prompt type a new password
passwd
```

Now, you can reboot and loggin normally. And the most interesting thing about
the normal boot, is that now you can profit of all the power of *systemd*.

## Autologin at startup

I have a normal boot, so I have to loggin, but I am in a container and I want
to automatically be logged when I start it.

So, here is how we can do,
**note: I have tested it on a debian image, the path may not been the same**:

When we connect into our container we are using the system console, so
we will override the systemd service handling the console to force the
autologin. In linux, the console manager is *getty* and the service dedicated
to the console is *console-getty@.service*

To override a service, you just have to create a *service_name.service.d*
directory in `/etc/systemd/system` and a file *whatever_name.conf*.

So, I created the dedicated file and put this inside:

```
#!bash
# /etc/systemd/system/console-getty.service.d/autologin.conf
[Service]
ExecStart=
ExecStart=-/sbin/agetty --noclear --autologin root --keep-baud console 115200,38400,9600 $TERM

```

The first `ExecStart=` clean the old call so that we can override it in the
following line.

## First conclusion

Systemd-nspawn comes with great features and is directly shipped in the system,
it is a good tool for stuff you want to do in a closed environment.
I wanted to ship *Steam* in a Docker container after reading
[this](http://fabiorehm.com/blog/2014/09/11/running-gui-apps-with-docker/),
but I think I will use systemd instead.

For development, I will continue to stick with Docker because of the
ecosystem and the community. I have started to use
[docker-compose](https://maci0.wordpress.com/2014/05/02/run-any-applications-on-rhel7-containerized-with-3d-acceleration-and-pulseaudio-steam-pidgin-vlc/) and it allows a lot
of interesting suffs.
