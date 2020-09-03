# Quick start

Krane is a self-hosted container management solution. It manages containers on remote servers so you dont have to.



##### 1. Install

Linux

One a remote server install krane

```bash
curl -sL linux.krane.sh | tar xz && chmod +x krane
```

##### 2. Run 

Run Krane as a background process

```bash
krane &
```

##### 3. Authentication

On your local machine

```bash
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

This creates a `private` and `public` key. 

The private key stays on your local machine, the public key gets appended to `~/.ssh/authorized_keys` on the server where Krane is running.

> Refer to [Authentication](authentication.md) for more details

##### 4. Download the Krane CLI

Mac

```bash
curl -fL linux-cli.krane.sh -o krane-cli && chmod +x krane-cli
```
