# 🏗️ Diretório: `core`

Este é o **coração** da nossa aplicação, o **"Hexágono"** na Arquitetura Hexagonal.

## ⚠️ Regra de Ouro

O código neste diretório **NÃO PODE** ter nenhuma dependência de tecnologia externa. Isso significa:

- ❌ **NÃO** importar pacotes de frameworks web (como `gin`)
- ❌ **NÃO** importar pacotes de ORM/drivers de banco de dados (como `gorm` ou `pgx`)
- ❌ **NÃO** importar pacotes de logging (como `zap`)

> O código aqui deve ser **Go puro**, representando as regras de negócio da forma mais limpa possível.

## 📂 Subdiretórios

### `/entity`
As **entidades** e **value objects** do nosso domínio (ex: `User`, `Patient`, `Email`).

### `/repository`
As **interfaces** (portas de saída) que definem como o domínio interage com a persistência de dados. Ex: `UserRepository`.

### `/usecase`
Os **casos de uso** da aplicação, que orquestram a lógica de negócio.