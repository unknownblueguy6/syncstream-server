# syncstream-server

A Golang webserver to sync video stream state among connected clients.

## Firefox extension
+ [Firefox Extension Repo](https://github.com/unknownblueguy6/syncstream-firefox)

## Purpose 
The primary goal of SyncStream is to implement a product that enables virtual watch parties for any website with HTML5 video playback, allowing people to watch movies, TV shows, and sports events together online.

## Features
- Browser extension that clients can connect to
- Clients can join a room and synch their video playback state with other users
- Supports any HTML5 Video Player
- Create Watch parties and join using a provided code
- Following events are supported: play/pause & video seeking
- Chat during the watch party session

## Documentation
+ [System Design & Documentation](/docs/system_design.md)

## Prerequisites
Go 1.22.1

## How to run
1. Clone the repository
2. In the root of the repo, run:
``` shell
go mod download
go run main.go
go run main.go --debug #to run in debug mode
``` 

## Demo
Located at `docs/demo.mp4`
<video src='docs/demo.mp4'>

## Contributing
+ [pull_request_template.md](/docs/pull_request_template.md)
+ [issue_template.md](/docs/issue_template.md)

## Team Roles 
+ [TEAM.md](/docs/TEAM.md)