@echo off
cd "%~dp0src\nw\"
go generate
cd "%~dp0\..\..\"
set GOPATH="%~dp0"
go build -ldflags "-H windowsgui" -o "%~dp0build\nw.exe" nw