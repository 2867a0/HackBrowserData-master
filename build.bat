@echo off

% windows%
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/windows_amd64.exe .\cmd\hack-browser-data\main.go

% linux%
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/linux_amd64 .\cmd\hack-browser-data\main.go

% MacOS%
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-s -w" -o bin/macos_amd64 .\cmd\hack-browser-data\main.go