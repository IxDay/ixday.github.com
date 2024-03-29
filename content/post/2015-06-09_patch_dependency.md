---
title:      "Patch dependency"
date:       2015-06-09
categories: ["Tuto"]
tags:       ["dev", "python"]
url:        "post/patch_dependency"
---

When developing on a project it is possible that a dependency can have an issue.
First you want to be able to debug it (pdb, ipdb), then modify it if you find
a bug.
To do that there is a naive way in python, which consist in editing directly
the sources of the module. But there is a cleaner way based on `pip`.

The `-e` option allows you to pass a path (git, http, file) for a given module
and link it to your environment. Then you just have to modify the files in
your filesystem and they will be provided to your project.

Here are the command lines needed (two line is enough).

```bash
# first I download the sources from wherever you want
git clone <address of the project>

# then I install and link the lib in my current env
pip install -e <path of the downloaded lib>

# now I can modify the lib file in order to debug it
```

The documentation for the feature can be found [here](https://pip.pypa.io/en/latest/reference/pip_install.html#editable-installs)

This is really useful when developing with third part libraries, and it is not
so well explained when you use `pip`, even if this feature is really awesome.

