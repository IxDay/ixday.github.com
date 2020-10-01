package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

const ENV_SSH_AUTH_SOCK = "SSH_AUTH_SOCK"

func AuthAgent() (ssh.AuthMethod, error) {
	conn, err := net.Dial("unix", os.Getenv(ENV_SSH_AUTH_SOCK))
	if err != nil {
		return nil, err
	}
	client, err := agent.NewClient(conn), err
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeysCallback(client.Signers), nil
}

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
	clientConn, err := ssh.Dial("tcp", "karadoc.tech:22", clientConfig)
	if err != nil {
		log.Fatalf("failed to connect to the ssh server: %q", err)
	}

	// tunnel traffic between local port 1600 and remote port 1500
	if err := tunnel(clientConn, "localhost:1600", "localhost:1500"); err != nil {
		log.Fatalf("failed to tunnel traffic: %q", err)
	}
}
