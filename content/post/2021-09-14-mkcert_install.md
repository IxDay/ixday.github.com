---
title:      "Install custom certificate with mkcert"
date:       2021-09-14
categories: ["Snippet"]
tags:       ["bash", "admin"]
url:        "post/mkcert_install"
---

Following my [previous post][previous_post] about generating a custom root certificate and a
PKI. I will share today a small snippet to install your self made certificate on
your laptop (or fleet) to be used by the system or the browsers.
It works on OSX, Linux and Windows by using a nice hidden feature of [`mkcert`][mkcert_repo].

Following [the documentation][mkcert_install], we can see that we can install a certificate
from a custom location. The idea here will be to rename our custom certificate
to `rootCA.pem` and pass the folder to the environment variable.

If I was starting from [my previous post][previous_post] terraform directory,
I would run the following commands to install the certificate on my system:

```sh
cp {ca_certificate,rootCA}.pem
CAROOT="$(pwd)" mkcert -install
```

And that's it! It should work out of the box, you can check by running `curl` without
the need of the `--cacert` option from now on.

[previous_post]: /post/terraform_pki
[mkcert_repo]: https://github.com/FiloSottile/mkcert
[mkcert_install]: https://github.com/FiloSottile/mkcert#installing-the-ca-on-other-systems
