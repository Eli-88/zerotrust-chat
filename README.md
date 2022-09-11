# <Network Monitoring>

## Description

A very simple command line interface for chat service that requires the client to generate and shared its secret key to be use for encrypting chat messages. Hence it is known as zero trust as it does not require any server or external services to generate the keys.

## Usage

Simply just run: `go run src/cmd/main.go <own port number>` to start the chat service

Commands
- `?` to display all active client
- `> <active client>` to switch to active client for chat
- `just type any messages` send messages to current active client 
- `$ <ip addr>` connect to other chat client and your it will be switch to your current active client

## Design
![Alt text](images/secret_key_exchange.png?raw=true "Network Monitor Design")