#!/bin/bash
mkdir -p bin

GOOS=windows GOARCH=amd64 go build -o bin/maketap.exe   maketap.go
GOOS=windows OOARCH=i386  go build -o bin/maketap.win32.exe maketap.go
GOOS=linux   GOARCH=amd64 go build -o bin/maketap.linux maketap.go
GOOS=darwin  GOARCH=arm64 go build -o bin/maketap.mac   maketap.go

GOOS=windows GOARCH=amd64 go build -o bin/loadtap.exe   loadtap.go
GOOS=windows GOARCH=i386  go build -o bin/loadtap.win32.exe loadtap.go
GOOS=linux   GOARCH=amd64 go build -o bin/loadtap.linux loadtap.go
GOOS=darwin  GOARCH=arm64 go build -o bin/loadtap.mac   loadtap.go
