# syntax=docker/dockerfile:1

############################
# 1) Builder (deps + venv)
############################
FROM python:3.11-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /build

# Dependências mínimas para instalar pacotes via pip (mantém builder limpo)
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Cria venv e instala dependências
COPY requirements.txt .
RUN python -m venv /opt/venv \
 && /opt/venv/bin/pip install --upgrade pip \
 && /opt/venv/bin/pip install --no-cache-dir -r requirements.txt


############################
# 2) Runtime (mais enxuto)
############################
FROM python:3.11-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PATH="/opt/venv/bin:$PATH" \
    # default (pode sobrescrever no docker run / k8s)
    PORT=8002

# user não-root
RUN useradd -m -u 10001 app

WORKDIR /app

# Copia venv pronto do builder
COPY --from=builder /opt/venv /opt/venv

# Copia app
COPY app.py /app/app.py

# (Opcional) se você tiver outras pastas/arquivos no repositório, inclua aqui:
# COPY . /app

EXPOSE 8002

# Healthcheck via endpoint /health (sua app possui /health) :contentReference[oaicite:3]{index=3}
RUN apt-get update && apt-get install -y --no-install-recommends curl \
 && rm -rf /var/lib/apt/lists/*
HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --retries=3 \
  CMD curl -fsS "http://127.0.0.1:${PORT}/health" || exit 1

USER app

# Gunicorn servindo a variável app do app.py (app:app) :contentReference[oaicite:4]{index=4}
CMD ["gunicorn", "--bind", "0.0.0.0:8002", "--workers", "2", "--threads", "4", "--timeout", "60", "app:app"]
