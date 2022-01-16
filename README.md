My personal blog
================

Some technical stuff I want to keep around.

Install
-------

```bash
# To install run a clone command with recursive enabled
git clone --recursive git@github.com:IxDay/ixday.github.com.git

# Then clone the master branch of the repo in the "public" directory
git clone -b master git@github.com:IxDay/ixday.github.com.git public
```

Publish
-------

```bash
# Default command build the website inside the public directory
hugo
# Add changes
cd public
git add -A
git commit -m "whatever"
git push
```

Create a new post
-----------------

You can run: `hugo new <name_of_the_post>`, to create a new post file with the
proper path and file name. This is due to the `script/wrapper.sh` file in
conjunction with `.envrc` which allows wrapping of the Hugo CLI.

__Do not forget:__ you have to remove the `draft=true` in order for the post
to be properly published.
