---
title:      "Golang SSH tunneling"
date:       2020-10-04
categories: ["Snippet"]
tags:       ["dev", "golang"]
url:        "post/golang_ssh_tunneling"
---

I already did [an article around Golang and SSH previously][post_certs],
since it is a huge part of my current work. Today, I will share a small
snippet which can be quite handy when automating stuff. This is how to
perform a [SSH tunnel][tunnel_doc] using Golang. I will not explain what it is
or what it is useful for, if you need more details please check the [doc][tunnel_doc].
In this tutorial, I will use the official [SSH library][ssh_lib] from Google.

[post_certs]: {{< ref "/post/2020-04-10_golang_ssh_certs">}}
[tunnel_doc]: https://help.ubuntu.com/community/SSH/OpenSSH/PortForwarding
[ssh_lib]: https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200930160638-afb6bcd081ae/ssh

Setup remote
------------

To be sure that our tunneling is working we will need a small service to reach.
Here is a small shell snippet to display the current time over HTTP:

```sh
while true; do
  echo -e "HTTP/1.1 200 OK\n\n $(date)" | nc -l -p 1500
done
```

Just run this on your remote host. You can then check if it works using curl:
`curl <remote_host>:1500`. You will also be sure to have a properly configured
SSH server. Check this by using: `ssh <remote_host>`.


Connect to local agent
----------------------

If you have some keys in your local agent you may want to use them to connect.
This has to be done in the code, and I will use the [agent module][agent_module]
from the Golang SSH library.

Here is a quick way to instantiate an SSH-agent client connection, it re-uses
the code from the [client example][client_example], I am also directly casting it as an `AuthMethod`:

```go
import (
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"
)

const EnvSSHAuthSock = "SSH_AUTH_SOCK"

func AuthAgent() (ssh.AuthMethod, error) {
	conn, err := net.Dial("unix", os.Getenv(EnvSSHAuthSock))
	if err != nil {
		return nil, err
	}
	client, err := agent.NewClient(conn), err
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeysCallback(client.Signers), nil
}
```

[agent_module]: https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200930160638-afb6bcd081ae/ssh/agent
[client_example]: https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200930160638-afb6bcd081ae/ssh/agent#NewClient

Keyboard interactive
--------------------

Another way of authenticating to a server is by using a keyboard-interactive challenge.
This may be necessary for additional authentication methods like 2FA.
For this snippet, we want to let the user enter its answers and pass it to the server.
Since code may be a bit complicated to understand I will comment as best as I
can to help understanding the flow. Also, check the [documentation][challenge_doc]
of the challenge function to understand the arguments.

```go
import (
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func AuthInteractive() ssh.AuthMethod {
	return ssh.KeyboardInteractive(
		func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			// this part is not really clear, as far as I get it we should print this
			// only if no question is provided
			if len(questions) == 0 {
				fmt.Printf("%s %s\n", user, instruction)
			}
			// instanciate the answers slice
			answers := make([]string, len(questions))

			// we iterate over each question and print it to the console
			for i, question := range questions {
				fmt.Print(question)

				// here is the trick, if the echo is true, we want to display user input
				// otherwise we want to hide it and use the terminal module to perform this
				if echos[i] {

					// simple scan over the console
					if _, err := fmt.Scan(&answers[i]); err != nil {
						return nil, err
					}
				} else {
					// here we use the ReadPassword function to hide user input
					answer, err := terminal.ReadPassword(syscall.Stdin)
					if err != nil {
						return nil, err
					}
					answers[i] = string(answer)
				}
			}
			return answers, nil
		})
}
```

__NOTE:__ If you want to implement the password authentication method you
should also use the `terminal.ReadPassword` function to hide the user input.

[challenge_doc]: https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200930160638-afb6bcd081ae/ssh#KeyboardInteractiveChallenge

Initiate ssh connection
-----------------------

Now we are going to create our SSH connection which we will use to
tunnel our traffic. We need a raw connection to perform byte copies. Looking
at the [SSH library][ssh_lib] we will use the `Dial` function. It takes
three arguments, the network (which is TCP), an address (which is your remote SSH server),
and a configuration. Let's take a look at this structure.

```go
import (
	"os/user"

	"golang.org/x/crypto/ssh"
)

func config(methods ...ssh.AuthMethod) (*ssh.ClientConfig, error) {
	// here I am retrieving user from current execution,
	// you can pass it as argument if you want
	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: current.Username,
		Auth: methods,
		// you should not pass this option, but for the sake of simplicity
		// we use it here
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}
```

Once config is ready we can create our connection instance. Take a look at the
last section to see how this is performed.

Tunnel
------

This is where the magic happens. We want to tunnel traffic from local
to remote endpoint. The first step is to initialize a TCP server, it will listen
on a specific port. Once we receive a connection to our local port, we
open a connection over the SSH connection to the remote address. We now have
two network connections, we just need to pipe each one to the other.
__NOTE:__ Most of the code has to be asynchronous to not block on
one side or another, however, in this example I am taking shortcuts and this
is not solid at all! Use with caution!

```go

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func tunnel(conn *ssh.Client, local, remote string) error {
	pipe := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("failed to copy: %s", err)
		}
	}
	listener, err := net.Listen("tcp", local)
	if err != nil {
		return err
	}
	for {
		here, err := listener.Accept()
		if err != nil {
			return err
		}
		go func(here net.Conn) {
			there, err := conn.Dial("tcp", remote)
			if err != nil {
				log.Fatalf("failed to dial to remote: %q", err)
			}
			go pipe(there, here)
			go pipe(here, there)
		}(here)
	}
}
```

Plugging everything together
----------------------------

It's now time to write our main function to bring everything together.
We set up our authentications methods, then initiate the SSH connection
to the remote host. Once the link has been established we tunnel our traffic
between our local and remote hosts.

```go
import (
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	// initiate auths methods
	authInteractive := AuthInteractive()
	authAgent, err := AuthAgent()
	if err != nil {
		log.Fatalf("failed to connect to the ssh agent: %q", err)
	}

	// initialize SSH connection
	clientConfig, err := config(authAgent, authInteractive)
	if err != nil {
		log.Fatalf("failed to create ssh config: %q", err)
	}
	clientConn, err := ssh.Dial("tcp", "<remote_host>:22", clientConfig)
	if err != nil {
		log.Fatalf("failed to connect to the ssh server: %q", err)
	}

	// tunnel traffic between local port 1600 and remote port 1500
	if err := tunnel(clientConn, "localhost:1600", "localhost:1500"); err != nil {
		log.Fatalf("failed to tunnel traffic: %q", err)
	}
}
```

You can now test the connection. Just run `curl localhost:1600` and you should
see the same thing as you would have by requesting the remote host
with `curl <remote_host>:1500`.

As always this code is not production-ready by any means. It falls short on
a lot of things, like error handling, or connection pool management.
However, this is the basic you will need to iterate onto. As always you
can find [the code in the blog Github repository][code_url]. Hope this will help.

[code_url]: /code/golang_ssh_certs/ssh_tunneling.go
