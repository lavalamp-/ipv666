#!/bin/bash

SERVER_IP=142.93.68.181

# Build the binary
# cd ..
# go build
# cd deploy

# Copy it over to the server
# ssh root@$SERVER_IP mkdir -p /root/ipv666
# scp ../ipv666 root@$SERVER_IP:/root/ipv666/
# scp ../config.json root@$SERVER_IP:/root/ipv666/config.json

# Copy over the systemd script
scp systemd.conf root@$SERVER_IP:/etc/systemd/system/ipv666.service

# Start the service
ssh root@$SERVER_IP systemctl daemon-reload
ssh root@$SERVER_IP systemctl start ipv666
ssh root@$SERVER_IP systemctl enable ipv666
