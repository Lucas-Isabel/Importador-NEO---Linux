@echo off

REM Definir as variáveis de ambiente para comrootlação cruzada
set GOOS=linux
set GOARCH=arm
set GOARM=7
set CGO_ENABLED=1
set CC=arm-linux-gnueabihf-gcc
REM Garantir que as dependências Go estejam atualizadas

REM Definir flags de compilação e linkagem sem -pthread
set CGO_CFLAGS=-mcpu=cortex-a7 -mfloat-abi=hard -mfpu=neon-vfpv4
set CGO_LDFLAGS=-mcpu=cortex-a7 -mfloat-abi=hard -mfpu=neon-vfpv4

echo Atualizando as dependências do Go...
go mod tidy
go get -d ./...



REM Comrootlar o programa Go
REM Compilar o programa Go com CGO
go build -o import

if errorlevel 1 (
    echo Erro durante a comrootlação do programa.
    pause
    exit /b 1
)

echo Comrootlação concluída com sucesso!

REM Transferir o arquivo para a Raspberry root via SCP
echo Transferindo o arquivo para a Raspberry root...
scp import root@132.144.55.28:
scp import.db root@132.144.55.28:

if errorlevel 1 (
    echo Erro ao transferir o arquivo para a Raspberry root.
    pause
    exit /b 1
)

REM Conectar à Raspberry root e executar o arquivo
echo Conectando à Raspberry root...
ssh root@132.144.55.28 "chmod +x import && ./import"

if errorlevel 1 (
    echo Erro ao conectar ou executar o arquivo na Raspberry root.
    pause
    exit /b 1
)

echo Processo concluído!
pause
