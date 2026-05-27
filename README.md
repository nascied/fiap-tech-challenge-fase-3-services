🚀 Tech Challenge Fase 3 — CI/CD + DevSecOps

<img width="929" height="191" alt="image" src="https://github.com/user-attachments/assets/40af84c8-317a-4cdb-8f78-2b1d58415289" />
📋 Sobre o Projeto

Este projeto representa a implementação da pipeline DevSecOps da plataforma ToggleMaster, utilizando GitHub Actions para automação de integração contínua, validação de qualidade, análise estática de código, segurança e atualização automatizada do ambiente GitOps.

A solução foi construída para garantir rastreabilidade, padronização e segurança no fluxo de entrega dos microserviços da aplicação.

🧱 Arquitetura da Pipeline

A pipeline foi estruturada em workflows independentes por microserviço:

CI - Auth Service
CI - Evaluation Service
CI - Flag Service
CI - Targeting Service
CI - Analytics Service

Cada workflow executa validações específicas de build, testes, qualidade e segurança.

🔄 Fluxo DevSecOps
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
⚙️ Pipeline CI/CD

A pipeline é executada automaticamente em:

Push na branch principal
Pull Requests
Atualizações dos microserviços

O fluxo foi dividido em múltiplas etapas para garantir validação progressiva antes do deploy.

✅ Build e Testes Automatizados

A etapa inicial realiza:

preparação do ambiente
instalação de dependências
execução de testes automatizados
validação estrutural dos serviços
Serviços Python
pip install -r requirements.txt
pytest
Serviços Go
go mod tidy
go vet
go test

O objetivo desta etapa é impedir que código inválido avance para as próximas fases da pipeline.

🔎 Quality Gates e Static Analysis

A pipeline implementa quality gates automatizados para validação da qualidade do código.

Ferramentas utilizadas
Python
flake8
Go
go vet
golangci-lint
Validações executadas

As validações verificam:

padronização do código
linhas excessivamente longas
trailing spaces
comentários inválidos
erros estruturais
organização do código

Exemplos detectados durante a análise:

E501 line too long
W293 blank line contains whitespace
E302 expected 2 blank lines

As inconsistências identificadas foram corrigidas antes da aprovação final da pipeline.

🔐 DevSecOps — Segurança Integrada

A pipeline foi construída seguindo práticas DevSecOps, incorporando verificações de segurança diretamente no processo de integração contínua.

Etapas de Segurança
Etapa	Objetivo
Static Analysis	Detectar falhas estruturais
SAST	Identificar vulnerabilidades no código
SCA	Validar dependências e bibliotecas
Docker Scan	Validar segurança das imagens
Security Scan (SAST & SCA)

A esteira executa verificações automáticas de segurança durante o pipeline:

Security Scan (SAST & SCA)

Objetivos:

reduzir vulnerabilidades
validar dependências
aumentar confiabilidade
aplicar segurança desde o CI/CD
🐳 Docker Build e Publicação no Amazon ECR

Após aprovação das validações:

as imagens Docker são construídas
ocorre integração com Amazon ECR
as imagens ficam disponíveis para deploy no Kubernetes

Etapa executada:

Docker Build & Push to ECR
🔄 Integração GitOps

Após o push das imagens:

os manifests Kubernetes são atualizados automaticamente
o repositório GitOps recebe a nova tag da imagem
o ArgoCD sincroniza automaticamente o cluster EKS

Etapa executada:

Update GitOps Manifests
☁️ Tecnologias Utilizadas
Tecnologia	Finalidade
GitHub Actions	Pipeline CI/CD
Docker	Containerização
Go	Microserviços
Python	Microserviços
flake8	Static Analysis
GitOps	Deploy contínuo
ArgoCD	Sincronização Kubernetes
Amazon ECR	Registry Docker
Amazon EKS	Orquestração Kubernetes
📊 Evidências da Pipeline
Workflows separados por microserviço

--

Tela principal GitHub Actions

Descrição:
Cada microserviço possui sua própria pipeline independente, permitindo rastreabilidade, isolamento e validação individual.

Execução da Pipeline

--

Pipeline completa executando

Descrição:
A pipeline foi estruturada em múltiplos estágios sequenciais com quality gates e validações automatizadas.

Build & Unit Test

--

Stage de testes automatizados

Descrição:
Execução automatizada de build e testes unitários antes da continuidade da esteira.

Linter / Static Analysis

--

Execução do flake8

Descrição:
Validação automática de qualidade do código utilizando ferramentas de análise estática.

Security Scan

--

Security Scan (SAST & SCA)

Descrição:
Execução de validações de segurança integradas ao pipeline DevSecOps.

Docker Build & Push

--

Push das imagens Docker

Descrição:
Construção automatizada das imagens e publicação no Amazon ECR.

GitOps Update

--

Atualização automática dos manifests

Descrição:
Atualização automatizada do repositório GitOps para sincronização via ArgoCD.

📄 Exemplo Simplificado do Workflow
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
🎯 Objetivos da Implementação

A implementação teve como foco:

automatizar integração contínua
aumentar confiabilidade das entregas
reduzir falhas manuais
aplicar quality gates
implementar práticas DevSecOps
automatizar o fluxo GitOps
melhorar rastreabilidade das entregas
📌 Status do Projeto

✅ CI/CD implementado
✅ Workflows independentes
✅ Build automatizado
✅ Testes automatizados
✅ Quality Gates
✅ Static Analysis
✅ DevSecOps
✅ Docker Build
✅ Push para Amazon ECR
✅ Integração GitOps
✅ Atualização automatizada dos manifests
✅ Integração com Kubernetes/EKS
