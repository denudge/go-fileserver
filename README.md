# Go Fileserver

A simple streaming file server demonstration

## Scope

This package contains an absolute simple streaming file server written in Golang.
It can store files from clients to disk (if not already there), and delete (existing) files from disk again.

It's purpose is first and foremost to demonstrate streaming some data.
Second, the `client` package can be embedded into other Go programs, thus making it easy to store files
to a remote server, e.g. backup files on a regular basis or dump processed data.
There's no authentication, though, so it's production usage is quite limited.

## Prerequisites

- Golang 1.22 or later
- make (optional)

## Installation

With `make`, you can just build client and server with this:

```bash
$ make
```

## Usage

Upload a file by using the client:
```bash
$ ./bin/client upload <folder> <filename> <local-file-path>
```

Delete a file on the server by using the client:
```bash
$ ./bin/client delete <folder> <filename>
```

## Server Options

- Port
- Root Directory
- Debug flag
- Buffer size

