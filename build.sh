#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -o flagship-sync-agent_linux_amd64 .
sha512sum flagship-sync-agent_linux_amd64 > flagship-sync-agent_linux_amd64_checksums.txt
tar -czvf flagship-sync-agent_linux_amd64.tar.gz flagship-sync-agent_linux_amd64 flagship-sync-agent_linux_amd64_checksums.txt

env GOOS=linux GOARCH=386 go build -o flagship-sync-agent_linux_386 .
sha512sum flagship-sync-agent_linux_386 > flagship-sync-agent_linux_386_checksums.txt
tar -czvf flagship-sync-agent_linux_386.tar.gz flagship-sync-agent_linux_386 flagship-sync-agent_linux_386_checksums.txt