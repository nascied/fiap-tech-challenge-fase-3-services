🚀 Tech Challenge Fase 3 — CI/CD + DevSecOps
📌 Visão Geral

Este projeto implementa um pipeline de Integração e Entrega Contínua (CI/CD) utilizando GitHub Actions, integrado aos microserviços da Fase 2.

O objetivo é automatizar o processo de build, testes e validações de qualidade, garantindo que apenas código válido seja integrado ao sistema.

🧱 Arquitetura do Projeto
🖌
O sistema é composto por microserviços distribuídos:

auth-service (Go)
evaluation-service (Go)
flag-service (Python)
analytics-service (Python)
targeting-service (Python)
⚙️ Pipeline CI/CD (GitHub Actions)

O pipeline é executado automaticamente a cada push na branch:

feature/devsecops
Etapas do pipeline:
Clone do repositório da Fase 2
Setup do ambiente Go
Build dos microserviços:
auth-service
evaluation-service
Execução de validações:
go vet
go fmt
Execução de testes automatizados (go test)
Validação da estrutura dos serviços Python
🔄 Fluxo do Pipeline

Push na branch → GitHub Actions → Build → Testes → Validações → Sucesso

🖼️ Evidência de Execução

Abaixo o pipeline executando corretamente no GitHub Actions:

<img width="895" height="380" alt="image" src="https://github.com/user-attachments/assets/03099a9f-258f-412a-87d9-eff4622315cb" />


Job aberto
Steps visíveis
Status “Success”
🎯 Objetivo da Implementação

Garantir automação no processo de integração contínua, reduzindo erros em produção e validando o código antes da integração.

🛠️ Tecnologias Utilizadas
GitHub Actions
Go (Golang)
Python
Docker (Fase 2)
Kubernetes (Fase 2)
AWS (infraestrutura da Fase 2)
📌 Status do Projeto
✔ Pipeline CI/CD funcional
✔ Build automatizado dos serviços Go
✔ Testes automatizados
✔ Validações de qualidade
✔ Integração com Fase 2
