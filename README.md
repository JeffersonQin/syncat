<p>
    <a href="https://goreportcard.com/report/github.com/JeffersonQin/syncat"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/JeffersonQin/syncat"/></a>
    <img src="https://img.shields.io/github/go-mod/go-version/JeffersonQin/syncat" alt="GitHub go.mod Go version"/>
    <a href="https://pkg.go.dev/github.com/JeffersonQin/syncat"><img src="https://godoc.org/github.com/JeffersonQin/syncat?status.svg" alt="Go Reference"/></a>
    <img src="https://img.shields.io/github/license/JeffersonQin/syncat" alt="GitHub"/>
    <img src="https://img.shields.io/github/last-commit/JeffersonQin/syncat" alt="GitHub last commit"/>
    <a href="https://github.com/JeffersonQin/syncat"><img src="https://img.shields.io/github/stars/JeffersonQin/syncat?logo=github" alt="Github Star"></a>
</p>

# syncat [WIP]

> Files should never have been synced by OneDrive, Qsync, iCloud automatically ...
> 
> -- one almost lost all his files once after a sync

## Introduction

A server-client tool with a centralized server hosting the files, and clients syncing files cross platforms.

## Features

* [ ] Sync files between server and multiple clients
* [ ] Files on server can be modified **in place**, because it is intended to be deployed on NAS
* [ ] Auto conflict detection. Conflict will be notified to the client and then notify the user to resolve manually
* [ ] Custom network protocol based on TCP
* [ ] GUI for clients (mainly used for resolving conflicts)
* [ ] Special support for NTFS, use win32api to listen to file changes instead of polling

## Project Structure

```
.
├── bin                     # Executables
├── cmd                     # Entry points
│   ├── client
│   │   └── main.go
│   └── server
│       └── main.go
├── config                  # Configuration files
│   ├── config.yml
│   ├── client_config.yml
│   └── server_config.yml
├── data                    # Database
├── internal                # Internal packages
│   ├── client
│   └── server
├── pkg                     # Reused code for server and client
│   ├── config              # Configuration
│   ├── database            # Database
│   ├── proto               # Protobuf
│   ├── sync                # Sync
│   └── syncnet             # Network Protocol
└── go.mod
```

## Resources

* http://go-database-sql.org/index.html
