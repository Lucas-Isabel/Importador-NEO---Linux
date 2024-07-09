@echo off
go install github.com/tc-hib/go-winres@latest
go-winres simply --icon youricon.png
go build -o IMPORTADOR_SYSTEL_console.exe
echo Compilation finished!
pause
