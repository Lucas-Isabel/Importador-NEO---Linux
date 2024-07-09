@echo off
REM Definir as variáveis de ambiente para compilação cruzada
set GOOS=linux
set GOARCH=arm
set GOARM=7




REM Compilar o programa Go para o Raspberry Pi 3
go build -o IMPORTADOR_SYSTEL

echo Compilação finalizada!
pause
