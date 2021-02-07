+++
title = "Debug Terraform"
date = 2021-02-07T08:25:05+01:00
categories = ["Tuto"]
tags = ["cli"]
+++

How to easily debug Terraform? This was one of my biggest problem when dealing with the tool.
I followed instructions to use some `output` ressources or browse the _"debug"_ logs.
Disclaimer: none of it was working properly.

So here is what I am actually doing and I am perfectly happy with.
Oh! And just before starting to enable `debug` logs you have to pass the
environment variable `TF_LOG=debug` (because there is no man page and it is not
written in the `--help` content).

Before starting
---------------

The way Terraform is working is by applying some changes to your infrastructure
and writting it to a state file (I am simplifying here). Terraform will not
display any information until this state has been written down. Namely: generated
ressource ids, names or anything related to resource creation.

My way of debugging will then be pretty simple. I will apply my terraform changes
until the point I am blocked (so all the remote resources will be created and
saved in the state file), then ask Terraform to display this information nicely :).


The actual way
--------------

If you have read the previous section and knows a bit how terraform is working
you will not be surprised here.

So for the sake of simplicity let's say that I am working with a file looking
like this:

```tf
resource "a" "this" {
	name = "foo"
}

resource "b" "this" {
	some_config = a.this.id
	another_config = a.this.name
}
```

Let's say that when running my `terraform apply` I get an error saying that `a.this.id` does not exist.
It is possible that I made a mistake and want to know what is actually held by
the resource `a` to fix the instructions.

Here is what I will run:

```sh
# this will generate the state up until my error which is at the next step
terraform apply -target a.this

# and now I can display what is held by the resource a
terraform show
```

As simple as that, but it took me quite some time to figure this out. Hope it
can help someone out there.

Bonus
-----

Sometimes it is obvious but sometimes it may be complicated to understand the
dependency path of a resource. Usually, if I want to visualize this I run the
following command: `terraform graph | dot -Tsvg > dependencies.svg`
