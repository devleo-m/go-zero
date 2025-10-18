# 🚀 Projeto Go-Zero (API Profissional)

Este repositório contém o desenvolvimento de uma API backend completa em Go, seguindo as melhores práticas de mercado como **Arquitetura Hexagonal**, **Domain-Driven Design (DDD)**, **SOLID** e testes robustos.

O objetivo é servir como um guia prático e um projeto de referência para a construção de aplicações escaláveis, testáveis e de fácil manutenção.

## 📁 Estrutura do Projeto

A estrutura de pastas foi projetada para refletir a separação de responsabilidades da Arquitetura Hexagonal.

- **`/cmd`:** Pontos de entrada da aplicação (entrypoints)
- **`/internal`:** Todo o código-fonte principal da nossa aplicação, não destinado a ser importado por outros projetos
- **`/core`:** O coração da aplicação, contendo a lógica de negócio pura (o "Hexágono")
- **`/infra`:** Adaptadores e implementações concretas de tecnologias externas (banco de dados, frameworks web, etc.)
- **`/pkg`:** Código compartilhado que poderia, teoricamente, ser importado por outros projetos (neste projeto, será pouco utilizado)

## 🚀 Como Começar

### Ambiente de Desenvolvimento
As dependências (PostgreSQL, Redis) são gerenciadas via Docker. Para iniciar o ambiente, execute:

```bash
docker-compose up -d
```

### Migrations
O schema do banco de dados é versionado. Para aplicar a última versão, use o Makefile:

```bash
make migrate-up
```
