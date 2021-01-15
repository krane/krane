# Authentication

Krane uses [private and public key authentication](https://en.wikipedia.org/wiki/Public-key_cryptography). Both keys are used for ensuring authenticity of incoming requests.

Start by creating a public and private key

> This command will generate 2 different keys, a `private` & `public (.pub)` key.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM' -f $HOME/.ssh/krane
```

Place the `public key` on the server where Krane is running, appended to `~/.ssh/authorized_keys`

The `private key` is kept on the user's machine.

Now try authenticating. The CLI will prompt you to select the `private key` you just created. This will be used for authenticating with the `public key` located on the Krane server.

```
krane login
```