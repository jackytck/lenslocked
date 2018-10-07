#!/bin/bash
curl https://getcaddy.com | bash -s personal
sudo cp caddy.service /etc/systemd/system
sudo systemctl enable caddy.service
mkdir $HOME/app
sudo systemctl start caddy