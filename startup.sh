#!/bin/bash

cd /home/ec2-user/word-search-go-web

# Start Tailwind CLI in watch mode
nohup /usr/local/bin/tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch > tailwind.log 2>&1 &

# Start the Go server (adjust ./main if needed)
nohup ./main > server.log 2>&1 &
