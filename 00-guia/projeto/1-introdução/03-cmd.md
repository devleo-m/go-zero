# 📁 Diretório `cmd`

Este diretório contém os **pontos de entrada (entrypoints)** da nossa aplicação. Cada subdiretório aqui representa um executável que pode ser compilado.

## 🎯 Propósito

A principal responsabilidade do código dentro de `/cmd` é:

- **Carregar configurações**
- **Inicializar dependências** (Logger, Conexão com DB, etc.)
- **"Montar" a aplicação** injetando as dependências (adaptadores) nas camadas de caso de uso e domínio
- **Iniciar a aplicação** (ex: subir um servidor HTTP)

> **Regra de Ouro:** O código aqui deve ser mínimo. Toda a lógica complexa deve residir nas camadas `internal/core` e `internal/infra`.