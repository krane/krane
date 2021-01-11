# Vultr

Run Krane on a $5 [Vultr](https://vultr.com) instance

> The following guides assumes you have a [Vultr](https://my.vultr.com/) account and can access the dashboard

### $5 / month

Under **Products** click the **+** sign or find _Deploy New Server_

Select **Cloud Compute** and the location closest to you.

<span class="img-wrapper">![Select compute](./assets/vultr/vultr-01.png)</span>

For the server type, under the **Application** tab select **Docker (Ubuntu 18.04)**

<span class="img-wrapper">![Select compute](./assets/vultr/vultr-02.png)</span>

Under **SSH Keys** add your public key.

This is usually under `~/.ssh/id_rsa`. **David's Personal Mac** being my id_rsa.

You can ignore the key labeled **Krane**.

<span class="img-wrapper">![Select compute](./assets/vultr/vultr-03.png)</span>

Click **Deploy Now** to create the server.

...

Once the server has been created, you'll want configure the [Vultr dns](https://www.vultr.com/docs/introduction-to-vultr-dns) with your domain provider. This will setup the domain you'll be using to point to your new server.

> tldr; DNS → (+) Add Domain → `Domain` example.com and `Default IP Address` server ip

---

Once you're able to `ssh` into your new server, follow the [getting started](docs/getting-started) to install Krane.
