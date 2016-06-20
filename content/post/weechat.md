+++
title = "weechat"
date = "2015-07-16"
categories = ["tuto"]
tags = ["admin"]
+++

[Weechat](https://weechat.org/) is an ncurse client for IRC,
which I use and have tweaked to fit my needs.

But Weechat have one major issue: IT IS NOT USER FRIENDLY.
The documentation is poor, there is a lot of plugins,
which documentation is even poorer, and the configuration is a hell.

WELCOME!

First of all the installation, `apt-get install weechat`
will be sufficient on a real OS. Then, just type `weechat` in
order to launch the client.

[![First screen]({filename}/images/weechat_1.png){:.image-process-article-image}]({filename}/images/weechat_1.png){:target="_blank"}

## Install plugins
At this point all the configuration can be done inside weechat,
the command `script` will install all the plugins you want, just like `apt`.

First script to install: `iset.pl`, just type `script search iset` and a
selection of available plugins will appears.
To leave it type `q` then press enter, if you want to install the script
type `i` then press enter (an `i` will appears in front of the name package,
type `q` and press enter to exit).

[![script search screen]({filename}/images/weechat_script.png){:.image-process-article-image}]({filename}/images/weechat_script.png){:target="_blank"}

If you didn't install the package through `script search` you can install it
with the following command: `script install iset.pl`.

## Configure
We will now be able to configure, just type `/iset` to enter the iset screen.
You will now see a list of all the parameters which can be modified.
If you type something in the input bar, it will look for the pattern in
the list of variables. If you want to search through the values, put an
`=` before the pattern.

To change the value, press `Alt + Enter` then enter the new value
(it is possible to navigate through values depending on variable type by
pressing the `Tab` key).

For example, I don't like the background color of iset selector.
On this screenshot I replace the value by `darkgray`,
changed values appear in magenta.

[![iset screen]({filename}/images/weechat_iset.png){:.image-process-article-image}]({filename}/images/weechat_iset.png){:target="_blank"}

## The buffer list
When using weechat I like to know on which buffer I am.
A buffer is the way weechat display informations, for example if I type:
`/iset` weechat will open a new buffer to display the informations.
If I want to close a buffer I just have to type: `/close`, If I want
to navigate through buffers I type `/buffer +1`. The command `/buffer list`
will display the buffer list in the first buffer, type: `/buffer 1` to see it.

There is a convenient plugin to display that: `chanlist.rb`.
Sadly this one is not supported and you will need to download it.
Here is the command to download the script:

`wget https://weechat.org/files/scripts/unofficial/chanlist.rb ~/.weechat/ruby/autoload/`

Relauch weechat or type `/script load chanlist.rb`.
We now have the buffer list on the left, but there is some commands to run in order to
have something good.

```bash
# It is mandatory and appears in the chanlist documentation.
/set irc.look.server_buffer independent

# Add UTF8 (for the delimiter) only if your terminal is UTF8 compatible (I hope so).
/set plugins.var.ruby.chanlist.utf8 on

# Fix chanlist size
/set weechat.bar.chanlist.size 30
/set weechat.bar.chanlist.size_max 30
```

The two last config fixed the weechat chanlist bar size.

Finally, I do not really like the color which is used for displaying the current
buffer. I want it magenta, chanlist is not very well developed and I will have to
modify the source code. Run the following line for changing this:

`sed -i -- 's/white,red/magenta,default/g' ~/.weechat/ruby/autoload/chanlist.rb`

Quit and relaunch weechat to see the changes. Here is my current result:

[![chanlist screen]({filename}/images/weechat_chanlist.png){:.image-process-article-image}]({filename}/images/weechat_chanlist.png){:target="_blank"}

## Connect to a server and a channel
Now we can visualize the buffers we will connect to a server.

```bash
# Setup nickname and connect
/set irc.server_default.nicks "MaxV, MaxV_, MaxV__"
/connect freenode
/join #freenode

# Customize nicklist
/set irc.look.color_nicks_in_nicklist on
/set weechat.bar.nicklist.size_max 12
/set weechat.bar.nicklist.size 12
```

Here is the result:

[![nicklist screen]({filename}/images/weechat_nicklist.png){:.image-process-article-image}]({filename}/images/weechat_nicklist.png){:target="_blank"}

## Bars

I will now customize the bars (because I can :p), this is inspired by
[Pascal Poitras blog entry](http://pascalpoitras.com/my-weechat-configuration/).

```bash
# Change original separator (need UTF8 support)
/set weechat.look.separator_horizontal "—"

# Create and customize activetitle bar
/bar add activetitle window top 1 0 buffer_title
/set weechat.bar.activetitle.priority 500
/set weechat.bar.activetitle.conditions "${active}"
/set weechat.bar.activetitle.color_fg red
/set weechat.bar.activetitle.color_bg default
/set weechat.bar.activetitle.separator on

# Customize the title bar
/set weechat.bar.title.conditions "${inactive}"
/set weechat.bar.title.color_fg white
/set weechat.bar.title.color_bg default
/set weechat.bar.title.separator on

# Remove status bar
/bar del status

# Create and customize the rootinput bar
/bar add rootinput root bottom 1 0 [buffer_name]+[input_prompt]+(away),\
[input_search],[input_paste],input_text
/set weechat.bar.rootinput.separator on

# Remove the input bar
/bar del input
```

Here is the result:

[![bars screen]({filename}/images/weechat_bars.png){:.image-process-article-image}]({filename}/images/weechat_bars.png){:target="_blank"}

## Layout

I want to display 4 buffers at a time, I will change the layout to fit my needs.

```bash
/window splith 50
/window splitv 50
/window 1
/window splitv 50

/layout store default
```

[![layout screen]({filename}/images/weechat_layout.png){:.image-process-article-image}]({filename}/images/weechat_layout.png){:target="_blank"}


## Key bindings

I have a few key bindings, at this point it fits my needs,
maybe I will add some in the future and edit this post accordingly.

```bash
# Activate the mouse on Alt + x, desactivate by pressing a second time
/key bind meta-x /mouse toggle

# Navigate through windows by pressing Alt + right, Alt + left
/key bind meta-meta2-C /window +1
/key bind meta-meta2-D /window -1
```
