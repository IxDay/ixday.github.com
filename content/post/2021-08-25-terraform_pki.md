---
title:      "Terraform PKI"
date:       2021-08-25
categories: ["Snippet"]
tags:       ["terraform"]
url:        "post/terraform_pki"
---

Once upon a time, I tried to setup an infrastructure using Terraform and Docker.
Spoiler alert: I reevaluated my setup to go for a Kubernetes approach which
is easier and more robust.
However, this experiment gave me the opportunity to better understand a really
important part of a self hosted infrastructure: how to set up a PKI.

__Disclaimer:__ I may not use the proper terms or be sloppy on some concepts.
I am not fully mastering this domain and want to share what I discovered.
Take this article with a grain of salt. Now bear with me and let's start.

## The goal

The basic idea of a KPI is to use a root certificate (most of the time self signed)
which will be trusted by your infrastructure components. Then sign child
certificates for every endpoint for a small period of time.

A good example of this is the [documentation of vault][vault_pki]. It consists
of a long living root certificate (valid for a few years for example). Which will
be stored in the most secure manner possible, ideally on a physical device not
connected to the internet. This root certificate is your last defense line and
is used to invalidate all or part of the infrastructure in case a breach occurs
(and some certificates got leaked).

Once you have a this base certificate you will sign an intermediate one with
a shorter TTL (Time To Live). Usually, we sign it for one year, at least in
the various companies I worked for, but I guess it may vary. However, the basic
idea is to have it last long enough to avoid a complicated rotation, but not too
long as well.

Now that we have this intermediate certificate we will use it to sign all the
child certificates used accross the infrastructure. Those child certificates
will need to be rotated on a regular basis (usually a few days).

This last point may be hard to reach and it is not uncommon to have child certificates
signed for a few months as well.

## The implementation

In my example I will use Terraform because it makes it pretty explicit and avoid
a lot of command, we are using the [official terraform provider][tls_provider].
You can perform all of this using the `openssl` CLI,
I may write a small article in the future to better explain this.

__Last but not least:__ In this example I will use an elliptic curve encryption
algorithm instead of the well known RSA one. Both are totally fine and can be
used here, I am just using ECDSA for personal reasons and because there is still
a lack of examples using it on the internet.

```tf

# define a few values for my use case, replace it with your naming
locals {
  cn       = "rulz.xyz"
  org      = "Rulz Corp."
  validity = 87600 # This is approximatively 10 years
}

# To build a certificate you actually need a private key, so we are creating it here
resource "tls_private_key" "rulz_ca" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

# Here is our root certificate, using our private key previously created
resource "tls_self_signed_cert" "rulz_ca" {
  key_algorithm     = "ECDSA"
  private_key_pem   = tls_private_key.rulz_ca.private_key_pem  # we use our private key
  is_ca_certificate = true                                     # this is a certificate authority

  subject {
    common_name  = local.cn
    organization = local.org
  }

  validity_period_hours = local.validity

  allowed_uses = [
    "cert_signing",  # this is what will make this certificate able to sign child certificates
  ]
}
```

Now we want to create our intermediate certificate. Since we made the root one
valid for 10 years, we want this one to be valid for only one year:

```tf

# We need to create a key for the intermediate cert
resource "tls_private_key" "rulz_int" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

# Once we have the key we can create a request
resource "tls_cert_request" "rulz_int" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.rulz_int.private_key_pem

  # the wildcard here does not have an impact but will help remind that this
  # is an intermediate signing all *.rulz.xyz certificates
  subject {
    common_name  = "*.${local.cn}"
  }
}

# And sign the request with our root ca
resource "tls_locally_signed_cert" "rulz_int" {
  cert_request_pem   = tls_cert_request.rulz_int.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.rulz_ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.rulz_ca.cert_pem
  is_ca_certificate  = true

  validity_period_hours = local.validity / 10 # one year of validity

  allowed_uses = [
    "cert_signing",
  ]
}
```

We are now ready to create the child certificate. It will not have the signing
capability and will only be able to authenticate. The creation is the same as
our intermediate certificate, only the capabilities and the fact that
it is not a certificate authority vary.

```tf
resource "tls_private_key" "rulz_child" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "rulz_child" {

  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.rulz_child.private_key_pem
  dns_names       = [ "child.${local.cn}" ]

  subject {
    common_name = "child.${local.cn}"
  }
}

resource "tls_locally_signed_cert" "rulz_child" {
  cert_request_pem   = tls_cert_request.rulz_child.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.rulz_int.private_key_pem
  ca_cert_pem        = tls_locally_signed_cert.rulz_int.cert_pem

  validity_period_hours = local.validity / 40 # valid for ~ 3 months

  # this will only allow to be used as a server
  allowed_uses = [
    "server_auth",
  ]
}
```

Deploy and Test
---------------

To properly use what we generated we need to output it to files to use it with
other softwares.  Here I will use the `local_file` resource to output
the certificates and keys to my filesystem:

First we need to output our root CA to a file which will be distributed across
the whole infrastructure.

```tf
locals {
  path = "${path.module}/ca_certificate.pem"
}

resource "local_file" "rulz_ca" {
  filename        = local.path
  content         = tls_self_signed_cert.rulz_ca.cert_pem
  file_permission = "0644"

  provisioner "local-exec" {  # I want to remove the certificate when I destroy the infra
    when    = destroy
    command = "rm ${self.filename}"
  }
}
```

Now we need to output our child certificate to a file as well. This child
certificate will need to contain the full chain of trust.

__Why should I need the full chain of trust:__ This come down to the way certificates
are implemented. On a network infrastructure you may have as many intermediate
signers as you want. The only thing is each intermediate signer certificate is
signed by the previous signer private key. However, your client only has the
root certificate. So, when checking the child it must ensure that the whole hierarchy
matches and goes to the root certificate which is trusted. This is why the certificate
handed over by your child server must contain this full chain to let the client
check it.

The order in which each certificate appear is important, the root is at the bottom,
the child at the top. Then each intermediate must appear in order between those two.
The extension does not have any impact you can use `.crt`, `.key`, `.pem` without
any repercussion.

```tf
resource "local_file" "rulz_child_key" {
  filename        = "${path.module}/child.key.pem"
  content         = tls_private_key.rulz_child.private_key_pem
  file_permission = "0600" # we want to protect the private key from others reading

  provisioner "local-exec" {
    when    = destroy
    command = "rm ${self.filename}"
  }
}

resource "local_file" "rulz_child_cert" {
  filename        = "${path.module}/child.crt.pem"
  file_permission = "0644"

  content = join("", [
    tls_locally_signed_cert.rulz_child.cert_pem,
    tls_locally_signed_cert.rulz_int.cert_pem,
    tls_self_signed_cert.rulz_ca.cert_pem,
  ])

  provisioner "local-exec" {
    when    = destroy
    command = "rm ${self.filename}"
  }
}
```

We are now ready to start our server, I will use `socat` following one of
[my previous post][socat_post] setup:

```sh
socat "ssl-l:8443,cert=child.crt.pem,key=child.key.pem,verify=0,fork,reuseaddr" \
        SYSTEM:"echo HTTP/1.0 200; echo Content-Type\: text/plain; echo; echo Hello World\!;"
```

And now we can use `curl` in another terminal to check if our setup is valid:

```sh
curl --cacert ca_certificate.pem --resolve "child.rulz.xyz:8443:127.0.0.1" https://child.rulz.xyz:8443
```

That's it! It should be a good introduction to KPI and also explain just enough
to have a basic understanding of certificate chain trust.
You can find the code of this post on the [github repository][post_code], and
launch it with `terraform init && terraform apply`.
Happy deployment!

[vault_pki]: https://www.vaultproject.io/docs/secrets/pki
[tls_provider]: https://registry.terraform.io/providers/hashicorp/tls/latest/docs
[socat_post]: /post/simple_https/
[post_code]: https://raw.githubusercontent.com/IxDay/ixday.github.com/source/content/code/terraform_pki/pki.tf
