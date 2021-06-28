+++
title = "Gogs + Drone"
date = 2016-05-12
categories = ["Snippet", "Tuto"]
tags = ["admin"]
url = "post/gogs_drone_compose"
+++

Jenkins is everywhere now, but I really don't like it. So I am looking at
a replacement from day to day. I discovered [Gogs](https://gogs.io/) an
I though that a CI is also a good use case for the Golang language.

And I finally found [Drone](https://drone.io/) (which was not really difficult
as it is mentionned in
[a ticket on gogs github](https://github.com/gogits/gogs/issues/1232)).

So I decided to make them work together in order to test that.

## Environment

We want to test them in a temporary location with a simple "Hello world" test.
Here is the architecture we will use:
```text
  /tmp/drone_gogs_test      # just a custom directory for our test
  |-- test                  # the git repository with the file to test
  |   |-- hello.py          # simple python file with a doctest test
  |-- gogs                  # directory which will be used by gogs to store
  |                         # datas (sqlite, git, ssh)
  |-- drone
      |-- dronerc           # config file
      |-- var               # directory to store drone data (mostly sqlite)
```

## Test repo
We just use the doctest feature of python here, this allow us to perform a
simple test without bootstraping a bunch of code

```python
"""Simple hello world function testing with doctest"""

def hello():
    """Simple hello world function

    Here the test we want to perform
    >>> hello()
    'Hello World!'
    """
    return 'Hello World!'


if __name__ == '__main__':
    import doctest
    doctest.testmod()
```

And this is the run

```bash
$ python hello.py -v
Trying:
    hello()
Expecting:
    'Hello World!'
ok
1 items had no tests:
    __main__
1 items passed all tests:
   1 tests in __main__.hello
1 tests in 2 items.
1 passed and 0 failed.
Test passed.
```

## Gogs

This one is really basic, the tutorial is really simple and it works out
of the box with the docker container provided. I haven't tested it through
ssh, but there is some disclaimer so I will test this later.

docs: https://github.com/gogits/gogs/tree/master/docker

I choose the sqlite backend, which is easier to configure because there is
no configuration to do.

I just run the docker container with the following line:
`docker run --name gogs -p 10022:22 -p 10080:3000 -v $(pwd)/gogs:/data gogs/gogs`

I connect it through the **localhost:10080** url and use the default config.
The first user created (through signup) is the administrator of the application.
Just be sure to replace all the localhost mentions with the address of
the container (which can be accessed with `docker inspect gogs`)

Add a new repository (like you do with github), and push the test project
described in the previous section to it.

You now have a working Gogs! (Be careful this is a temporary build, do not
use this as your day to day git repo).

## Drone

This one is a bit more complicated, first we have to create a dronerc file
with the following content:

```bash
REMOTE_DRIVER=gogs
REMOTE_CONFIG=http://172.17.0.3:3000
DEBUG=true
```

Here I put the application in debug mode, this is not mandatory, but I like
to have outputs if something goes bad.

Now we can start the container:
```text
docker run \
        --volume $(pwd)/drone/var:/var/lib/drone \
        --volume /var/run/docker.sock:/var/run/docker.sock \
        --env-file $(pwd)/drone/dronerc \
        --restart=always \
        --publish=80:8000 \
        --detach=true \
        --name=drone \
        drone/drone:0.4
```

