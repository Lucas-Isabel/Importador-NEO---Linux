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
