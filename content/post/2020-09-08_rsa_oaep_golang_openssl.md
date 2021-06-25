+++
title = "RSA OAEP in Go and Openssl equivalent"
date = 2020-09-08
categories = ["Snippet"]
tags = ["golang", "bash"]
+++

From time to time you write some code to deal with data that will be stored
somewhere on a drive. When debugging it may be nice to have _"shell"_ commands
to mimic code and help understanding what is going on.

When it comes to encryption the swiss army knife you will most probably have
in your shell environment is `openssl`. But reading and understanding the documentation
can sometimes be "challenging".

Let's cook
----------

First, I generated a RSA key in order to encrypt our cipher:

```bash
openssl genrsa -out rsa.priv 2048
```

You will use the following Golang code to load this key:

```go
import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func DecodeRSAKey(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
```

We will use [the encrypt function from the standard lib](https://golang.org/pkg/crypto/rsa/#EncryptOAEP).
Usage will be as follow:

```go
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func EncryptRSAOAEP(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, []byte(""))
}
```

You now want decrypt this using `openssl`. If you saved the bytes to a file
you can use the following snippet:

```bash
cat <file> | openssl pkeyutl -decrypt -inkey key.rsa  \
	-pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 \
	-pkeyopt rsa_mgf1_md:sha256
```

Testing it out!
---------------

I created [a directory in the blog repo](https://github.com/IxDay/ixday.github.com/tree/source/content/code/golang_openssl_rsa_oaep)
with a simple golang file reading from stdin and encoding using the private key
present in the directory.

Here is how to use it:

```bash
# encrypt stdin
echo "foo" | go run rsa_oaep.go

# encrypt from stdin and decrypt right away
echo "foo" | go run rsa_oaep.go | openssl pkeyutl -decrypt -inkey rsa.priv \
	-pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 \
	-pkeyopt rsa_mgf1_md:sha256
```

