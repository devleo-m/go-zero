# 📁 Diretório: `internal`

Este diretório contém todo o **código-fonte privado** da nossa aplicação.

> **Convenção Go:** De acordo com as convenções da comunidade Go, o código dentro de `/internal` não pode ser importado por outros projetos. Isso garante que toda a nossa lógica de negócio e de infraestrutura esteja encapsulada e não seja exposta acidentalmente.

## 📂 Subdiretórios

### `/core`
O **núcleo** da nossa aplicação. Contém a lógica de negócio pura, sem dependências de frameworks ou tecnologias externas. É o nosso **"Hexágono"**.

### `/infra`
A **camada de infraestrutura**. Contém os "Adaptadores" que implementam as interfaces definidas no Core e lidam com o mundo exterior (banco de dados, servidor web, etc.).