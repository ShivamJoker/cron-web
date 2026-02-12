# Simple web app for exiting unix cronjobs

Uses the default [crontab](https://www.man7.org/linux/man-pages/man5/crontab.5.html), make sure it's installed and setup

<img width="1409" height="679" alt="image" src="https://github.com/user-attachments/assets/0d61c07b-4578-4dfa-8279-6e3a766dfbc6" />

## Features
- Single binary ~8 MB (no need to install anything else, apart fron crontab)
- CRON experssion explanation
- CRON validation
- Easy to use web interface

## How to install?

### Get latest build from GitHub release

Download the compatible version for your operating system.
This script will download the compatible binary (`aarch64`, `x86_64`).

```sh
wget https://github.com/ShivamJoker/cron-web/releases/latest/download/cronweb-$(uname -m)
```

### Make sure it's has executable permission

```sh
chmod +x ./cronweb-*
```

### Move the binary to `bin` (optional)

```sh
sudo mv cronweb-$(uname -m) /usr/local/bin/cronweb
```

### Copy paste script

```sh
wget https://github.com/ShivamJoker/cron-web/releases/latest/download/cronweb-$(uname -m)
chmod +x ./cronweb-*
sudo mv cronweb-$(uname -m) /usr/local/bin/cronweb
```

## How to start cron web server

Run the executable, it takes http address and port as an optional argument.

```sh
cronweb localhost:9090
```

## Systemd Service

### Create a systemd unit file with vim

`vim /etc/systemd/system/cronweb.service`

### Edit the service config

```sh
[Unit]
Description=Cronweb

[Service]
ExecStart=/usr/local/bin/cronweb

[Install]
WantedBy=multi-user.target
```
Save and exit from the vim `:wq`

## Enable the service

Tell systemd it has got a new service to run

```sh
sudo systemctl daemon-reload
```

Enable and start the cronweb service

```sh
sudo systemctl enable --now cronweb
```

Now you are all set to keep the service running forever

## Serve at subpath (Caddy example)

If you would like to access the cronweb portal from a subpath like `homelab.local/cronweb` you can setup reverse proxy.

```sh
example.com {
        redir /cronweb /cronweb/ # redirect to path with / (won't work without)
        handle_path /cronweb/* {
                reverse_proxy localhost:8816 # Cron web
        }
}
```
