+++
title = "Golang SSH, marshalling, unmarshalling"
date = 2020-04-10T15:55:12+02:00
categories = ["Snippet"]
tags = ["golang", "dev"]
+++

Golang is a wonderful language to deal with the SSH protocol. It's mostly due
to the [SSH library](https://pkg.go.dev/golang.org/x/crypto/ssh?tab=doc) which
is pretty exhaustive.

However, when I had to deal with external requirements like SSH Agent or
OpenSSH I experienced a lack of example and struggled a bit interfacing.

In this article we will see how to exchange keys between a program written
in Go and those tools using files. This means to output our keys in a proper
format (marshaling) and being able to read them (unmarshalling) from the
external tool format.

Also, we will not only cover the RSA keypairs but also the shiny SSH certificates.
If you are dealing with an infrastructure using SSH a lot, you should definitely
take a look at those. Here is a [good blog post](https://smallstep.com/blog/use-ssh-certificates/),
which explains why you should consider SSH certificates.


RSA keys
--------

Asymmetric keys are at the foundation of authentication for the SSH protocol.
I will cover the RSA format here, but there are also others, however, logic should be the same.

### Create
The first thing to do is to create a key pair. It is pretty straightforward in
Go and there is plenty of examples on the internet.
Will still put a snippet here, in order to save a search:

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"

	"golang.org/x/crypto/ssh"
)

const keySize = 2048

func generateKey() (*rsa.PrivateKey, ssh.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}
```

I also converted the public part of the key to the SSH library structure.
This is not mandatory, but it might be useful for developers.

### Marshal

We now want to marshal this in order to save it to a file.
We also want this file to be usable by the command line tooling.
For instance, we want to be able to load the key in the SSH agent.
To perform this we will make our output string to be in the right format.

```go
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const typePrivateKey = "RSA PRIVATE KEY"

func marshalRSAPrivate(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: typePrivateKey, Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
}
```

If you save the output to a file, you can load this key in your agent with the
following command: `ssh-add <output_file>`

You may want to also save the public part in a format readable by OpenSSH
to grant access to a user. It is usually the format you can find in the
`~/.ssh/authorized_keys` file.
Here is a quick snippet on how to generate this string.

```go
package main

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

func marshalRSAPublic(pub ssh.PublicKey) []byte {
	return bytes.TrimSuffix(ssh.MarshalAuthorizedKey(pub), []byte{'\n'})
}
```

### Unmarshal

Same as the previous section but now we are loading from a file. To simplify
the snippets I will use a slice of bytes, and remove the file handling logic.

Here is the code for the private key:

```go
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func unmarshalRSAPrivate(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
        return nil, fmt.Errorf("failed to parse PEM block containing the key")
    }

    return x509.ParsePKCS1PrivateKey(block.Bytes)
}
```

And now the code for the public part.
You can use a [line scanner](https://golang.org/pkg/bufio/#example_Scanner_lines)
in order to parse `authorized_keys` files.
This also the format you will see when listing keys from your agent with this command: `ssh-add -L`

```go
package main

import (
	"golang.org/x/crypto/ssh"
)

func unmarshalRSAPublic(bytes []byte) (ssh.PublicKey, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	return pub, err
}
```

### Extra

The `~/.ssh/known_hosts` file, contains public keys in a format similar to the one we have seen previously.
It appears that the encoding is the same, only a prefix is added.
We can check this out [in the source code](https://github.com/golang/crypto/blob/056763e48d71/ssh/keys.go#L123).


Certificates
------------

Certificates are supported by Golang `x/crypto` package. However, I had a hard
time finding how to perform the same tasks on the Internet. The interface
provided by the library is not aimed for certificates making functions a bit
harder to locate. Here are a few snippets that can be useful.

### Create

The creation of a certificate is pretty straightforward.
When generating a certificate you will need a pair of keys associated with it,
here we will use the one generated previously.
The certificate structure contains way more fields than what I am showing here,
I strongly advise you to add a time validity window and a serial number for proper
security and tracking.
Check out [the documentation](https://pkg.go.dev/golang.org/x/crypto/ssh?tab=doc#Certificate) for more information.

```go
package main

import (
	"golang.org/x/crypto/ssh"
)

func generateCert(pub ssh.PublicKey) *ssh.Certificate {
	permissions := ssh.Permissions{
		CriticalOptions: map[string]string{},
		Extensions: map[string]string{ "permit-agent-forwarding": ""},
	}
	return &ssh.Certificate{
		CertType: ssh.UserCert, Permissions: permissions, Key: pub,
	}
}
```

### Marshal

In order to marshal a certificate to a valid string you will need to sign it
first. Here I will self sign the certificate. Self-signing is using the private
part of the certificate key to sign it.

We first need to create a signer interface from our private key:

```go
package main

import (
	"crypto/rsa"

	"golang.org/x/crypto/ssh"
)

func generateSignerFromKey(priv *rsa.PrivateKey) (ssh.Signer, error) {
	return ssh.NewSignerFromKey(priv)
}

func generateSignerFromBytes(bytes []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(bytes)
}
```

The `generateSignerFromBytes` function can take the output of the previous
`marshalRSAPrivate` function. Since you already have a `*rsa.PrivateKey`
structure you do not want to marshal/unmarshal again, this would be a bit
overkill. However, the purpose of this post is to show you how
all those structures, interfaces and types plug together.
This is in case you load the key from a file and thus from a slice of bytes.

Once we have the signer we will use it to sign the certificate:

```go
package main

import (
	"crypto/rsa"
	"crypto/rand"

	"golang.org/x/crypto/ssh"
)

func generateAndSign() (*rsa.PrivateKey, *ssh.Certificate, error) {
	priv, pub, err := generateKey()
	if err != nil {
		return nil, nil, err
	}
	signer, err := generateSignerFromKey(priv)
	if err != nil {
		return nil, nil, err
	}
	cert := generateCert(pub)
	return priv, cert, cert.SignCert(rand.Reader, signer)
}
```

Now that we have a self-signed certificate we can properly marshal it. Note that
if we did not do the previous steps the certificate would not be complete.
Therefore it would not have been possible to marshal it.

Marshaling is actually quite simple, but the function name does not make it obvious:

```go
package main

import (
	"golang.org/x/crypto/ssh"
)

func marshalCert(cert *ssh.Certificate) []byte {
	return ssh.MarshalAuthorizedKey(cert)
}
```

We do not use the `cert.Marshal()` function here, we will see later what is its
purpose.

### Unmarshal

Unmarshaling is the opposite operation, we just need to additionally cast to
the structure we want to have:

```go
package main

import (
	"golang.org/x/crypto/ssh"
)

func unmarshalCert(bytes []byte) (*ssh.Certificate, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return nil, err
	}
	cert, ok := pub.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("failed to cast to certificate")
	}
	return cert, nil
}
```


SSH Agent
---------

Now that we have created all those slices of bytes, we can dump them into files.
Those files can be loaded directly in your SSH agent using the command line interface.

To load a certificate in your agent you will need:

- To dump the private key in a file without extension (ex: "foo") using the
	`marshalRSAPrivate` function [from the section above](#marshal).
- Then, to dump the certificate in a file with the same name suffixed with `-cert.pub`
	(ex: "foo-cert.pub") using the `marshalCert` function
	[from the certificate marshaling section](#marshal-1).
- Finally, to load the files in the agent by issuing the following command: `ssh-add foo`.

You can now connect to an OpenSSH server using the certificate
(you will have a bit of configuration to do though).

### The "wire" format

In the certificate marshaling section I talked about the `cert.Marshal()` function.
This function does not marshal to text it encodes to binary format. This is
the format used by the SSH Agent when communicating through its socket.
The `x/crypto` library actually supports the agent protocol and you can find
[the API on godoc](https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200414173820-0848c9571904/ssh/agent?tab=doc).

So, given this "binary" format we may be interested in one last snippet.
If we take a close look at this agent we see a `List` function. This function
actually returns an `agent.Key` pointer. We can transform this structure into
an `ssh.PublicKey` interface and potentially cast it to an `ssh.Certificate` pointer.
Here is an example on how to do this:

```go
package main

import (
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	formatCert = "ssh-rsa-cert-v01@openssh.com"
)

func listAndCast(keys []*agent.Key) error {
	for _, key := range keys {
		pub, err := ssh.ParsePublicKey(key.Blob)
		if err != nil {
			return err
		}
		if key.Format == formatCert {
			cert, ok := pub.(*ssh.Certificate)
			if !ok {
				return fmt.Errorf("failed to cast key to certificate: %q", err)
			}
			// ... do whatever you want with the certificate
		}
	}

}
```

Conclusion
----------

I hope this code will help other people. I am also publishing a simple
[go file](/code/golang_ssh_certs/ssh_certs.go) containing all the snippets
except for the agent. The file plugs everything together creating structures,
marshaling them, then unmarshaling.

Enjoy coding!
