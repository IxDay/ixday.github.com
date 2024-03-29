---
title:      "Pop-up with xmobar"
date:       2015-08-24
categories: []
tags:       []
url:        "post/pop_up_with_xmobar"
---

Sometimes I just want to click, even for simple things.
So I created a little popup to shutdown, restart, suspend my computer,
which looks like this:

[![Xmobar popup]({filename}/images/xmobar_popup.png){:.image-process-article-image}]({filename}/images/xmobar_popup.png){:target="_blank"}


## Preparation

First of all install [xmobar](http://projects.haskell.org/xmobar/),
`apt-get install xmobar/testing`. Here I will use testing because
I need the multiple fonts support,
see [here](http://projects.haskell.org/xmobar/releases.html#version-0.23-mar-8-2015)
for the release note.

The second thing I will need is an iconic font for the power button.
[Font Awesome](https://fortawesome.github.io/Font-Awesome/) fits my need.
I downloaded the `.ttf` file [here](https://github.com/FortAwesome/Font-Awesome/blob/master/fonts/fontawesome-webfont.ttf?raw=true)
and placed it under `~/.local/share/fonts`. I ran `fc-cache -vf` in order to
update the font cache of my machine.

And that's all, now let's configure


## Configuration
```haskell
Config {
	font = "xft:Droid Sans Mono-10"
	, additionalFonts = ["xft:FontAwesome-10"]
	, border = FullB
	, borderColor = "#8a8a8a"
	, borderWidth = 2
	, bgColor = "#3c3b37"
	, fgColor = "#ffffff"
	, lowerOnStart = False
	, position = Static { xpos = 660 , ypos = 490, width = 600, height = 100 }
	, template = "}<fn=1></fn>   \
				  \<action=`systemctl poweroff`> Shutdown </action>    \
				  \<action=`systemctl reboot`>Reboot</action>    \
				  \<action=`ps -o pid | sed '3q;d' |xargs kill && systemctl suspend`>\
				  \Suspend</action>    \
				  \<action=`ps -o pid | sed '3q;d' |xargs kill`>Exit</action>{"
}
```

Explanation line by line:

* line 2: Define the font which will be use, this is my default font for my
machine, which is available in the package `fonts-droid`.
* line 3: I declare an additional font which will be used for the icon.
* line 4: The border style, this will tell xmonad that my border applies on full perimeter.
* line 5, 6: The border color and width.
* line 7, 8: The background color and the font color.
* line 9: Tells xmonad to launch the popup above all windows.
* line 10: Here is the position of the popup, my screen is 1920x1080,
I want a popup 600x100, so `xpos = (1920 - 600) / 2` and `ypos = (1080 - 100) / 2`
* line 11 to 16: The template to display, see [the template section](#template)
for a full description

### The template <a name="template"></a>

* `}{` symbol describe the text alignment, in a template like this `foo}bar{qux`,
foo will be left align, bar will be centered and qux will be left align.
Here I just want all the text of the popup centered
* `<fn=1>` indicates the use of the font number one in additionals fonts,
here it will be the "Font Awesome" font.
* The curious symbol will appear like [this](http://fortawesome.github.io/Font-Awesome/icon/power-off/) if Font Awesome is installed on yout system, like this: Ϳ if not.
To retrieve a character from font awesome and paste it in a file,
you can launch the following bash command `awk 'BEGIN{printf "%c", 0x00<font_awesome_code>}'`.
Where `font_awesome_code` is the code provided on the icon page of Font Awesome
(for the power button it will be f011).
* `</fn>` indicates that we finished using the Font Awesome font.
* `\\` is the way to perform multiline string in Haskell.
* `<action=\`systemctl poweroff\`>` will perform a poweroff when the user
clicks on the text enclosed in the tag.
* `Shutdown` the displayed text.
* `</action>` close the clickable area.
* same thing for the reboot section
* last specificity `ps -o pid | sed '3q;d' |xargs kill`, this command line close
the popup, by killing its process.




