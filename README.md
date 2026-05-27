# 🚀 Tech Challenge Fase 3 — Pipeline DevSecOps

\[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-CI/CD-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)](https://github.com/)
[![Docker](https://img.shields.io/badge/Docker-Container-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![AWS](https://img.shields.io/badge/AWS-ECR/EKS-232F3E?style=for-the-badge&logo=amazonaws&logoColor=white)](https://aws.amazon.com/)
[![ArgoCD](https://img.shields.io/badge/ArgoCD-GitOps-EF7B4D?style=for-the-badge&logo=argo&logoColor=white)](https://argo-cd.readthedocs.io/)

> ⚠️ Projeto desenvolvido como parte do Tech Challenge Fase 3 da Pós-Tech FIAP em Arquitetura Cloud e DevOps.

---

# 📋 Sobre o Projeto

Este projeto representa a implementação da pipeline DevSecOps da plataforma ToggleMaster, utilizando GitHub Actions para automação de integração contínua, validação de qualidade, análise estática de código, segurança e atualização automatizada do ambiente GitOps.

A solução foi construída para garantir:

- rastreabilidade
- automação
- padronização
- segurança
- integração contínua
- atualização automatizada do ambiente Kubernetes

---

# 🧱 Arquitetura da Pipeline

A pipeline foi estruturada em workflows independentes por microserviço:

```text
.github/workflows/
├── ci-auth-service.yaml
├── ci-evaluation-service.yaml
├── ci-flag-service.yaml
├── ci-targeting-service.yaml
└── ci-analytics-service.yaml
```

Cada workflow executa validações específicas de:

- build
- testes
- quality gates
- análise estática
- segurança
- Docker build
- integração GitOps

---

# 🔄 Fluxo DevSecOps

```text
Push / Pull Request
        ↓
 Build & Unit Test
        ↓
 Linter / Static Analysis
        ↓
 Security Scan (SAST & SCA)
        ↓
 Docker Build & Push to ECR
        ↓
 Update GitOps Manifests
        ↓
 ArgoCD Sync no EKS
```

---

# ⚙️ Pipeline CI/CD

A pipeline é executada automaticamente em:

- Push na branch principal
- Pull Requests
- Atualizações dos microserviços

O fluxo foi dividido em múltiplos estágios sequenciais para garantir validação progressiva antes do deploy.

---

# ✅ Build e Testes Automatizados

A etapa inicial realiza:

- preparação do ambiente
- instalação de dependências
- execução de testes automatizados
- validação estrutural dos serviços

## Serviços Python

```bash
pip install -r requirements.txt
pytest
```

## Serviços Go

```bash
go mod tidy
go vet
go test
```

O objetivo desta etapa é impedir que código inválido avance para as próximas fases da pipeline.

---

# 🔎 Quality Gates e Static Analysis

A pipeline implementa quality gates automatizados para validação da qualidade do código.

## Ferramentas utilizadas

### Python
- flake8

### Go
- go vet
- golangci-lint

---

## Execução do Static Analysis

```bash
flake8 . --max-line-length=120
```

---

## Validações executadas

As validações verificam:

- padronização do código
- linhas excessivamente longas
- trailing spaces
- comentários inválidos
- erros estruturais
- organização do código

Exemplos detectados durante a análise:

```text
E501 line too long
W293 blank line contains whitespace
E302 expected 2 blank lines
```

As inconsistências identificadas foram corrigidas antes da aprovação final da pipeline.

---

# 🔐 DevSecOps — Segurança Integrada

A pipeline foi construída seguindo práticas DevSecOps, incorporando verificações de segurança diretamente no processo de integração contínua.

---

# 🛡️ Etapas de Segurança

| Etapa | Objetivo |
|---|---|
| Static Analysis | Detectar falhas estruturais |
| SAST | Identificar vulnerabilidades no código |
| SCA | Validar dependências |
| Docker Scan | Validar segurança das imagens |

---

# 🔍 Security Scan (SAST & SCA)

A esteira executa verificações automáticas de segurança durante o pipeline:

```text
Security Scan (SAST & SCA)

✔ Static Analysis
✔ Dependency Validation
✔ Docker Image Scan
✔ Vulnerability Verification
```

Objetivos:

- reduzir vulnerabilidades
- validar dependências
- aumentar confiabilidade
- aplicar segurança desde o CI/CD

---

# 🐳 Docker Build e Publicação no Amazon ECR

Após aprovação das validações:

- as imagens Docker são construídas
- ocorre integração com Amazon ECR
- as imagens ficam disponíveis para deploy no Kubernetes

---

## Docker Build

```bash
docker build -t analytics-service .
```

---

## Push para Amazon ECR

```bash
docker tag analytics-service:latest \
123456789012.dkr.ecr.us-east-1.amazonaws.com/analytics-service:latest

docker push \
123456789012.dkr.ecr.us-east-1.amazonaws.com/analytics-service:latest
```

---

# ☁️ Integração AWS

## Configuração AWS CLI

```bash
aws configure
```

---

## Login no Amazon ECR

```bash
aws ecr get-login-password --region us-east-1 \
| docker login --username AWS --password-stdin <ecr-registry>
```

---

# 🔄 Integração GitOps

Após o push das imagens:

- os manifests Kubernetes são atualizados automaticamente
- o repositório GitOps recebe a nova tag da imagem
- o ArgoCD sincroniza automaticamente o cluster EKS

---

# 📄 Exemplo de Atualização do Manifest

```yaml
containers:
  - name: analytics-service
    image: 123456789012.dkr.ecr.us-east-1.amazonaws.com/analytics-service:latest
```

---

# ☁️ Tecnologias Utilizadas

| Tecnologia | Finalidade |
|---|---|
| GitHub Actions | Pipeline CI/CD |
| Docker | Containerização |
| Go | Microserviços |
| Python | Microserviços |
| flake8 | Static Analysis |
| GitOps | Deploy contínuo |
| ArgoCD | Sincronização Kubernetes |
| Amazon ECR | Registry Docker |
| Amazon EKS | Orquestração Kubernetes |

---

# 📊 Evidências da Pipeline

## 1. Workflows separados por microserviço

Adicionar screenshot:
- Tela principal GitHub Actions

Descrição:
Cada microserviço possui sua própria pipeline independente, permitindo rastreabilidade, isolamento e validação individual.

---

## 2. Execução da Pipeline

Adicionar screenshot:
- Pipeline completa executando

Descrição:
A pipeline foi estruturada em múltiplos estágios sequenciais com quality gates e validações automatizadas.

---

## 3. Build & Unit Test

Adicionar screenshot:
- Stage de testes automatizados

Descrição:
Execução automatizada de build e testes unitários antes da continuidade da esteira.

---

## 4. Linter / Static Analysis

Adicionar screenshot:
- Execução do flake8

Descrição:
Validação automática de qualidade do código utilizando ferramentas de análise estática.

---

## 5. Security Scan

Adicionar screenshot:
- Security Scan (SAST & SCA)

Descrição:
Execução de validações de segurança integradas ao pipeline DevSecOps.

---

## 6. Docker Build & Push

Adicionar screenshot:
- Push das imagens Docker

Descrição:
Construção automatizada das imagens e publicação no Amazon ECR.

---

## 7. GitOps Update

Adicionar screenshot:
- Atualização automática dos manifests

Descrição:
Atualização automatizada do repositório GitOps para sincronização via ArgoCD.

---

# 📄 Exemplo Simplificado do Workflow

```yaml
name: CI - Analytics Service

on:
  push:
    branches:
      - main

  pull_request:
    branches:
      - main

jobs:
  build-test:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v5

      - name: Install dependencies
        run: pip install -r requirements.txt

      - name: Run Tests
        run: pytest

      - name: Static Analysis
        run: flake8 .

      - name: Security Scan
        run: echo "SAST/SCA"

      - name: Docker Build
        run: docker build -t analytics-service .

      - name: Push ECR
        run: echo "Push image"

      - name: Update GitOps
        run: echo "Update manifests"
```

---

# 🎯 Objetivos da Implementação

A implementação teve como foco:

- automatizar integração contínua
- aumentar confiabilidade das entregas
- reduzir falhas manuais
- aplicar quality gates
- implementar práticas DevSecOps
- automatizar o fluxo GitOps
- melhorar rastreabilidade das entregas

---

# 📌 Status do Projeto

| Status | Implementação |
|---|---|
| ✅ | CI/CD implementado |
| ✅ | Workflows independentes |
| ✅ | Build automatizado |
| ✅ | Testes automatizados |
| ✅ | Quality Gates |
| ✅ | Static Analysis |
| ✅ | DevSecOps |
| ✅ | Docker Build |
| ✅ | Push para Amazon ECR |
| ✅ | Integração GitOps |
| ✅ | Atualização automática dos manifests |
| ✅ | Integração Kubernetes/EKS |

---

# 👨‍💻 Autor

**Sandro Moraes de Souza**  
Pós-Tech FIAP — Arquitetura Cloud e DevOps  
Tech Challenge Fase 3

---

# 📄 Licença

Projeto desenvolvido para fins educacionais como parte da Pós-Tech FIAP.
