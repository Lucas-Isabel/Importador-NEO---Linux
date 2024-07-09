# Use a imagem oficial do Golang como base
FROM golang:1.22

# Atualize o sistema e instale as dependências necessárias para compilação para ARM7
RUN dpkg --add-architecture armhf \
    && apt-get update \
    && apt-get install -y \
        gcc-arm-linux-gnueabihf \
        libc6-dev-armhf-cross \
        linux-libc-dev-armhf-cross \
        --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# Configure as variáveis de ambiente para compilação cruzada para ARM7
ENV GOOS=linux \
    GOARCH=arm \
    GOARM=7 \
    CC=arm-linux-gnueabihf-gcc \
    CGO_ENABLED=1

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /usr/src/myapp

# Copie o código-fonte para o contêiner
COPY . .

# Compile o programa Go
RUN go build -o import-embbed

# Comando padrão a ser executado quando o contêiner for iniciado
CMD ["./import-embbed"]
