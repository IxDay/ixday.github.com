+++
title = "Pdbpp"
date = "2016-02-28"
categories = ["Tuto"]
tags = ["dev", "python"]
+++

I am a huge fan of ipython and its debugger ipdb (I have also done a
patch on it). Then I discovered [pdbpp](https://pypi.python.org/pypi/pdbpp/)
and I found it so great that I no more use ipdb, here is why.

## More features

### Sticky mode

Pdbpp comes with a lot of additional features which are really convenient.
The first and more well known is the sticky mode:

{{< figure src="/img/pdbpp_sticky.png" >}}

This will display the code currently, executed and shows you with an arrow at
which exact line you are. Just type `sticky` to enable this mode and "voila" it
works

### Pdb disable

This one is quite usefull when debugging in tests. For example, you have
a pdb statement in a fixture which is called by a bunch of tests.
Maybe there is something I do not really understand but when I want
to exit, I am stuck pressing Ctrl+d Ctrl+c hoping it will stop the tests
before another fixture call happen.

Pdbpp extends the pdb command with a disable statement
(`pdb.disable()` in the console), which will ignore subsequent calls to pdb.

### No variable overrides

When using ipdb, sometimes in the code I use some reserved variables:
"c", "rv", "args",.. when entering the debugger they are overriden by
keywords, and I have to change their name and start again to see what they
contain.

With pdbpp, this no longer happens, as explained in the documentation, pdbpp
can infers which variable is currently used and avoid the override.

### And many more

Just read the documentation and looks at the other features, this library
is pure gold.

## Just a simple wrapper

This is just a wrapper around pdb, which means that you can simply call it
with `import pdb; pdb.set_trace()`. And more importantly, when a tool provides
a postmortem feature (scrapy, pytest) it brings you directly to pdbpp without
having to dig in the code or documentation (when you wanted to activate ipdb).

## Simple config

Here is my config (a real simple one):

```python
import pdb

class Config(pdb.DefaultConfig):
    sticky_by_default = True # start in sticky mode
    current_line_color = 40  # black
```
