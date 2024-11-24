# Stage 1: Base builder - configuração para baixar dependências
FROM golang:1.22.5-alpine3.19 AS base_builder

WORKDIR /myapp/

COPY ["go.mod", "go.sum", "./"]

# Baixa todas as dependências do Go
RUN go mod download

# Stage 2: Builder - compila o aplicativo Go
FROM base_builder AS builder

WORKDIR /myapp/

COPY . .

# Argumentos para versão e commit
ARG PROJECT_VERSION=1
ARG CI_COMMIT_SHORT_SHA=1

# Compila o binário com flags de otimização (-s e -w removem símbolos de depuração)
RUN go build -ldflags="-s -w -X 'main.VERSION=$PROJECT_VERSION' -X main.COMMIT=$CI_COMMIT_SHORT_SHA" -o app cmd/worker/main.go

# Stage 3: Minimal Docker Image - cria uma imagem final com o binário compilado
FROM alpine:3.19

WORKDIR /app/

# Copia o binário compilado da etapa anterior para a imagem final
COPY --from=builder /myapp/app .

# Define o ponto de entrada da aplicação
ENTRYPOINT ["./app"]
