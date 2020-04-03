#!/bin/sh



env GOOS=linux GOARCH=amd64 go build -o $1ipcamemapi-linux-amd64
env GOOS=linux GOARCH=arm GOARM=5 go build -o $1ipcamemapi-linux-arm5-raspi
env GOOS=linux GOARCH=arm GOARM=6 go build -o $1ipcamemapi-linux-arm6-raspi
env GOOS=linux GOARCH=arm GOARM=7 go build -o $1ipcamemapi-linux-arm7-raspi
env GOOS=linux GOARCH=arm64 go build -o $1ipcamemapi-linux-arm8-raspi4
env GOOS=windows GOARCH=amd64 go build -o $1ipcamemapi-windows-amd64.exe
env GOOS=darwin GOARCH=amd64 go build -o $1ipcamemapi-darwin-amd64
