@echo off
cd "%~dp0src\nw\"
go generate
cd "%~dp0\..\..\"
set GOPATH="%~dp0"
go build -o "%~dp0build\nw_debug.exe" nw