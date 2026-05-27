# =========================
# 1 - Estágio de build 
# =========================
FROM golang:1.24-alpine AS builder
WORKDIR /app
#Criar dependências do módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia arquivos para build
COPY . .

# limpa 
RUN go mod tidy
# Compila binário estático e otimizado dando o nome do binãrio de auth-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /app/evaluation-service .

#RUN go mod tidy
RUN go build -o /app/evaluation-service . 
# =========================
# Stage 2 - Runtime
# =========================
FROM golang:1.24-alpine AS runner

WORKDIR /app

# Cria usuário app + instala CA + wget (para healthcheck)
RUN addgroup -S app && adduser -S app -G app \
    && apk add --no-cache ca-certificates wget \
    && chown -R app:app /app 


# Copia o binário para execucao da aplicacao
COPY --from=builder  --chown=app:app /app/evaluation-service /app/evaluation-service

# # Permissões do binario para executar com o usuario app
# RUN chown app:app /app/auth-service

# muda o contexto para o usu'ario
USER app

# Porta alta para non-root
ARG PORT=8004
EXPOSE 8004

# Healthcheck (HTTP 200 em /health)
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
  CMD wget -qO- "http://127.0.0.1:${PORT}/health" >/dev/null 2>&1 || exit 1

ENTRYPOINT ["/app/evaluation-service"]
