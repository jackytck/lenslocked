#!/bin/bash
sudo cp lenslocked.service /etc/systemd/system
sudo systemctl enable lenslocked.service
