+++
categories = ["snippet"]
date = 2016-11-23
tags = []
title = "Docker Nspawn"
url = "url/docker_nspawn"
+++

Don't want to use Docker? Still want to start containers for tests or whatever?
Don't want to install yet another software to perform this? Want to understand
a bit of how all those things work? Great! I will show you how to boot a
container from the internet only through `systemd-nspawn`

Thanks to the CoreOS team (love those guys) a new hub for storing container now exists:
[quay](https://quay.io/). The other good news is the ACI,
the container image format defined in the
[App Container (appc) spec](https://github.com/appc/spec). It basically define
what a container image should be: Basically, a filesystem under a `rootfs`
directory with a `manifest` file. The `manifest` file is in a JSON format and
furnish metadatas for the container (env variables, run command at start, ...).
It is a good news because this specification does not include the **sh\*\*y**
layer system of Docker (will publish a rant someday about this).

Okay, so... what the point? Good question! If I simplify all those informations,
it means that there is a point in the internet where I can download filesystem
archives and directly boot them through systemd-nspawn.

Find the url
------------

Just go to [https://quay.io](https://quay.io) and type the kind of image you
want in the search bar. For this example, I will retrieve an
[Alpine linux](https://alpinelinux.org/) image. I take the first one available,
which is a clone of the one in the docker hub
[https://quay.io/repository/aptible/alpine](https://quay.io/repository/aptible/alpine).
In the page open the web debugger and in the console enter the following:

```javascript
$("meta[name='ac-discovery']").eq(0) // simple jquery request

// this should be the output
[<meta name="ac-discovery" content="quay.io https://quay.io/c1/aci/{name}/{version}/{ext}/{os}/{arch}/">]
```

The part between `{}` are variables, and should be replaced by what we want to
retrieve, spec [here](https://github.com/appc/spec/blob/master/spec/discovery.md)

- **name**: the name of the image, here it is `alpine`, but there is a trick
  the name is a full qualified one, and is the one you usually pull with the
  docker command. Here it is `quay.io/aptible/alpine`
- **version**: can be retrieve under the tags tab (add `?tab=tags` at the end
  of the url). We will take `latest` here. **BEWARE**: using latest is not
  recommended as it is not a fixed version and can change from day to day. This
  can cause non-reproductible builds (same advice when using Docker as well).
- **ext**: `aci` to get a tarball. The other possibility is `aci.asc` which is
  the signature.
- **os**: `linux` of course
- **arch**: `amd64` for my part

Resulting URL will be:
`https://quay.io/c1/aci/quay.io/aptible/alpine/latest/aci/linux/amd64`

Retrieve - Deflate - Start
--------------------------

```bash
# Download the image inside a alpine.tgz file
wget -O alpine.tgz https://quay.io/c1/aci/quay.io/aptible/alpine/latest/aci/linux/amd64

# Untar it
tar xvf alpine.tgz

# Boot it, with systemd-nspawn
sudo systemd-nspawn -M alpine -D rootfs

# TADAAA! You are now inside a container 
```

You can also download, untar, rename and place the directory in your tree with
one command:

```bash
wget -O - "https://quay.io/c1/aci/quay.io/aptible/alpine/latest/aci/linux/amd64" | \
		tar -C "/tmp/alpine" --transform="s|rootfs/|/|" -xzf -

```

Conclusion
----------

Easy duh!?
