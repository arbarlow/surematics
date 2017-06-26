# Install

To install all the examples here to play with..

```
$ go install github.com/arbarlow/surematics/...
```

## Postcode

Simple postcode distance calculator, using haversine formula. This is "as the crow flies" and so not accurate enough for say, a maps service or driving directions.

To use...

```
$ postcode 'SW1A 1AA' 'E8 4AA'
```

## Chat Server/Client

Super simple and stupid chat client and server.

The server is a simple tcp socket that sends broadcasts all things to all clients and terminates responses with a `\r` so that respones with multiple lines can be handled. A list of current users is simply send with a prefix of `u:` which the chat client handles differently.

First run the server..

```
$ chatserver
```

then in a seperate tmux pane or similar, run the interactive gui

```
$ chatclient
```
