+++
title = "Introducing Janus (SSH Agent written in Go)"
date = 2020-11-02
categories = ["Project"]
tags = ["dev", "golang"]
+++

Why this project
----------------

I am currently using [gopass][gopass_url] to store and share my passwords.
I was relying on GPG to handle the encryption side of the process. Then a colleague of
mine introduced me to [age][age_url]. This encryption specification allows the
use of SSH keys and specifically the ED25519 ones. I decided to make the switch
and moved all my stores to age encryption. 

I am protecting my private key with
a password and only loading it inside my agent. I needed an easy way to decrypt
files and found out [this project][sagent_url] when browsing the age repository 
issues. I decided to "fork" it and implements all the stuff I needed to make it
usable.

[gopass_url]: https://www.gopass.pw/ 
[age_url]: https://github.com/FiloSottile/age
[sagent_url]: https://github.com/42wim/sagent

How to use it
-------------

From the current repo, I just installed it under `/usr/local/bin`

```sh
git clone "https://github.com/IxDay/janus"
cd janus
make
PREFIX=/usr/local make install
```

My current laptop is using Archlinux I just replaced my SystemD unit file with
my new binary (previous version from [Archlinux wiki][arch_wiki]):

```ini
[Unit]
Description=SSH key agent (Janus)

[Service]
Type=simple
Environment=SSH_AUTH_SOCK=%t/ssh-agent.socket
# Display required for ssh-askpass to work
Environment=DISPLAY=:0
ExecStart=/usr/local/bin/janus

[Install]
WantedBy=default.target
```
Here a few useful commands:
- You can activate this using SystemD: `systemctl --user enable ssh-agent`.
- You can check the logs using the journal command: `journalctl --user -fu ssh-agent`.
- You can add keys using the usual `ssh-add` command.

Also, Janus provides an `ssh-decrypt` command to perform decryption using
a key in the agent. Here is a quick example to show how this works:

The tool has a lot of limitations. It does not handle stdin, this
is why I am using a process substitution here. In this line, I am encrypting the
string "foo" using my public key. I pass the result to the `ssh-decrypt` binary
and it outputs back my string.

Further work
------------

The tool is working as a beta for my current use case. However, I will need to
add a few capabilities to the tool to make it properly ready. Here is a quick 
list, I hope I will bring it down shortly:

- stdin input.
- armor format support.
- better logs (for the agent process and the decrypt tool).
- options and documentation for both CLIs.
- do the `gopass` integration, which is the reason why I started this project. 
    You can track the progress on [my fork][gopass_fork].


[arch_wiki]: https://wiki.archlinux.org/index.php/SSH_keys#Start_ssh-agent_with_systemd_user
[gopass_fork]: https://github.com/IxDay/gopass