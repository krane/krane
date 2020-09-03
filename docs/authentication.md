# Authentication

## Creating authentication keys

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'

-t type
-b bytes
-C comments
-m key format
```

This will generate 2 different keys, a `private` & `public (.pub)` key. The `private key` is kept on the user's machine, the `public key` is stored where Krane is running and appended to `~/.ssh/authorized_keys`
