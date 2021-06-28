+++
title = "ssh, rsync and fswatch"
date = 2015-03-02
categories = ["Tuto"]
tags = ["admin", "bash"]
url = "post/rsync_fswatch"
+++

Sometimes I just can't work on a local environment (particular architecture,
particular services, local configuration too complex, etc...).
So, I have to synchronize my local directory with a remote one and test the
web interface on my local machine.


## ssh
First, ssh! For this I need a ssh connection to the remote server, here I use
a particular ssh key.

```bash
ssh -i ~/.ssh/my_ssh.key mylogin@192.168.0.1

# urls also work
ssh -i ~/.ssh/my_ssh.key mylogin@my.url.com
```

Okay, at this moment we can use a config file for ssh: `$HOME/.ssh/config` :)

```text
Host 192.168.0.1 my.url.com
   user mylogin
   IdentityFile ~/.ssh/my_ssh.key
```

The CLI is now:

```bash
# with IP
ssh 192.168.0.1

# with url
ssh my.url.com
```

I have a connection but now I want to share my dev server on the remote host to
my local web browser. For that I will use the `-L` of ssh, this will forward
the local port with the remote one.

```bash
ssh -L 5000:localhost:9999 192.168.0.1 # same as before it works with url
```

This command will connect my local port 5000 to the port 9999 of the remote
host local interface (so I do not have to open another port on the remote
except the ssh one)


## rsync
Then, we will use rsync to send files to our server


```bash
rsync -avz -e "ssh" local/folder my.url.com:/home/username/remote_folder
```
* -a the archive option (recursive, preserve links, times, permissions,
group, owner, devices files
* -v the verbose option
* -z the compress option
* -e remote shell to use (basically we specify here all the ssh configuration
needed. If we didn't had the .ssh/config file the command line should have been
`rsync -avz -e "ssh -i ~/.ssh/my_ssh.key" local/folder
mylogin@my.url.com:/home/username/remote_folder`

It is possible to not want to synchronize all files
(.git folder, generated files, etc...), so we will use the `--exclude-from`
option. In the folder we want to synchronize, we create a file `exclude.txt`
(the name is not important), then we fill it with the needed files or folder:

```text
.git
/static
*.pyc
```

Take care that the `/` apply at the point where the rsync command is launched.
So it will not have the same effect if we change the working directory.

The command will look like:

```bash
# move to working dir
cd local/folder

remote_loc="my.url.com:/home/username/remote_folder"

rsync -avz -e "ssh" . $(remote_loc) --exclude-from 'exclude.txt'
```

## fswatch (or inotify)

I also want my folder to synchronize automatically with the remote one when a
file change. For this purpose I will use `fswatch` because I use MacOS
(shame on me), `inotify` can be use on linux platforms.

First, check the changes on my working directory:

```bash
fswatch -e .git/ -e .pyc -e $(pwd)/static .
```

Here I use `$(pwd)` in order to not catch the `/static` folder at the root of
the folder, but keep the nested one included. This catch the same files as the
exclude file from rsync. At this point I haven't found any solution to unify
those two commands.

## xargs

The last piece needed is xargs, this will read stdin and execute a command on
each entry.

## All together

Here is my final command:

```bash
cd local/folder

remote_loc="my.url.com:/home/username/remote_folder"

fswatch -0 -o -e .git/ -e .pyc -e $(pwd)/static . | \
xargs -0 -I {} rsync -avz -e "ssh" . $remote_loc --exclude-from 'exclude.txt'
```

* The option `-0` indicates that fswatch will use `\0` as a line separator.
* The option `-o` will only indicates how many files have been modified,
has long has I do not need the filename to perform the command.

* The command xargs take the same option `-0` so it will accept `\0` as the
separator between each command.
* The `-I {}` option will tell xargs that the
caught at first will be injected in the command at the place of `{}`
(this is the same as the -exec command in find). We do not use it because
rsync will take care to check which file has changed, this is a trick to avoid
xargs to complain.

