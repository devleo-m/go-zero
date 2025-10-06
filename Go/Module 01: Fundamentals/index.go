package main

import "fmt"

// 1. Variáveis
// Declaração de Variáveis
// Go oferece várias formas de declarar variáveis:
// Forma 1: var com tipo explícito

var nome1 string
var idade1 int
var altura1 float64
var isEstudante1 bool

// Forma 2: var com inicialização (tipo inferido)

var nome2 = "John"
var idade2 = 25
var altura2 = 1.75
var isEstudante2 = true

// Forma 3: := (declaração e inicialização)

// name := "John"
// age := 25
// height := 1.75
// isStudent := true

// ⚠️ IMPORTANTE: O operador := só funciona dentro de funções!
func main() {
    // ✅ CORRETO - := funciona dentro de funções
	nome3 := "John"
	idade3 := 25
	altura3 := 1.75
	isEstudante3 := true
	
	// Usar as variáveis para evitar erro de "declared and not used"
	fmt.Printf("Nome: %s, Idade: %d, Altura: %.2f, É estudante: %t\n", 
		nome3, idade3, altura3, isEstudante3)
}

// Forma 4: Múltiplas variáveis
// Mesmo tipo
// var x, y, z int

// Tipos diferentes
var (
    // nome   string
    // idade  int
    // altura float64
)

// Com inicialização
// var a, b, c = 1, 2, 3

// Declaração curta múltipla (só funciona dentro de funções!)
// nome2, idade := "Maria", 30  // ❌ ERRO fora de função

// Regras de Nomenclatura
// ✅ Permitido:
// nome         // minúscula
// Nome         // Maiúscula (exportado/público)
// nome123      // com números
// _nome        // com underscore
// nomeCompleto // camelCase

// 2. Tipos de Dados
// Tipos Numéricos Inteiros
// Inteiros com sinal (podem ser negativos)

/*
--------------------------------------------------------|
Tipo  |  Tamanho       | Intervalo                      |
--------------------------------------------------------|
int8  |  8 bits        | -128 a 127                     |
int16 |  16 bits       | -32,768 a 32,767               |
int32 |  32 bits       | -2 bilhões a 2 bilhões         |
int64 |  64 bits       | -9 quintilhões a 9 quintilhões |
int   |  32 ou 64 bits | Depende da arquitetura         |
--------------------------------------------------------|
*/

//Inteiros sem sinal (apenas positivos)
/*
--------------------------------------------------------|
Tipo   |  Tamanho       | Intervalo                      |
--------------------------------------------------------|
uint8  |  8 bits        | 0 a 255                        |
uint16 |  16 bits       | 0 a 65,535                    |
uint32 |  32 bits       | 0 a 4,294,967,295             |
uint64 |  64 bits       | 0 a 18 quintilhões            |
uint   |  32 ou 64 bits | Depende da arquitetura        |
--------------------------------------------------------|
*/

// Tipos especiais
/*
--------------------------------------------------------|
Tipo  |  Tamanho       | Elias                          |
--------------------------------------------------------|
byte  |  8 bits        | uint8                          |
--------------------------------------------------------|
rune  |  32 bits       | int32                          |
--------------------------------------------------------|
*/

// exemplos

var idade int = 25
var temperatura int8 = -10
var populacao uint64 = 7_800_000_000  // underscores para legibilidade
var letra rune = 'A'
var dados byte = 255

// 💡 Dica: Use int na maioria dos casos. Use tamanhos específicos quando necessário (ex: protocolos de rede, otimização de memória).
// Tipos Numéricos de Ponto Flutuante

/*
--------------------------------------------------------|
Tipo    |  Tamanho       | Precisão                     |
--------------------------------------------------------|
float32 |  32 bits       | 7 dígitos decimais           |
float64 |  64 bits       | 15 dígitos decimais          |
--------------------------------------------------------|
*/

// exemplos
var altura float64 = 1.75
var peso float32 = 68.5
var pi float64 = 3.14159265359

// Notação científica
var planck float64 = 6.62607015e-34
var avogadro float64 = 6.022e23

// 💡 Dica: Prefira float64 (mais preciso e é o padrão).
// Tipos Numéricos Complexos
var z1 complex64 = 3 + 4i
var z2 complex128 = complex(5, 6)  // 5 + 6i

// real := real(z1)  // 3
// imag := imag(z1)  // 4

//Tipo Boolean

/*
--------------------------------------------------------|
Tipo  |  Tamanho       | Valor                          |
--------------------------------------------------------|
bool  |  8 bits        | true ou false                  |
--------------------------------------------------------|
*/

var isStudent bool = true
var isAdmin bool = false

// Tipo String

/*
--------------------------------------------------------|
Tipo  |  Tamanho       | Valor                          |
--------------------------------------------------------|
string |  16 bits       | 0 a 65,535                    |
--------------------------------------------------------|
*/

var nome string = "João Silva"
var mensagem = "Olá, mundo!"

// String multilinha (raw string)
var texto = `
    Isso é uma string
    com múltiplas linhas
    e preserva formatação
`

// String vazia
var vazia string  // "" (zero value)

// Operações com Strings
// Concatenação
// nome := "João"
// sobrenome := "Silva"
// nomeCompleto := nome + " " + sobrenome  // "João Silva"

// Tamanho (bytes, não caracteres!)
// tamanho := len("Olá")  // 4 bytes (á ocupa 2 bytes)

// Acessar caractere por índice (retorna byte)
// primeira := nome[0]  // 'J' (tipo byte/uint8)

// Substring
// sub := nome[0:3]  // "Joã"

// String é imutável!
// nome[0] = 'X'  // ❌ ERRO - não pode modificar

// Strings e Unicode (Runes)

/*
--------------------------------------------------------|
Tipo  |  Tamanho       | Valor                          |
--------------------------------------------------------|
rune  |  32 bits       | int32                          |
--------------------------------------------------------|
*/

// texto := "Olá, 世界"

// len() retorna bytes, não caracteres
// fmt.Println(len(texto))  // 12 bytes

// Para contar caracteres, use runes
// runes := []rune(texto)
// fmt.Println(len(runes))  // 6 caracteres

// Iterar sobre runes
// for i, r := range "Olá" {
//     fmt.Printf("Índice: %d, Rune: %c\n", i, r)
// }
// Saída:
// Índice: 0, Rune: O
// Índice: 1, Rune: l
// Índice: 2, Rune: á  (pula bytes intermediários)

// 3. Constantes
// Constantes são valores imutáveis definidos em tempo de compilação.
// Declaração de Constantes
const Pi = 3.14159
const MaxUsuarios = 100
const NomeApp = "MeuApp"

// Múltiplas constantes
const (
    Domingo1 = 0
    Segunda2 = 1
    Terca3   = 2
    Quarta4 = 3
    Quinta5  = 4
    Sexta6   = 5
    Sabado7  = 6
)

// Com tipo explícito
const (
    StatusOK    int = 200
    StatusError int = 500
)

// iota - Enumerador Automático

const (
    Segunda = iota  // 0
    Terca           // 1
    Quarta          // 2
    Quinta          // 3
    Sexta           // 4
    Sabado          // 5
    Domingo         // 6
)

// Começar de 1
const (
    Janeiro = iota + 1  // 1
    Fevereiro           // 2
    Marco               // 3
    // ...
)

// Pular valores
const (
    _ = iota  // pula 0
    KB = 1 << (10 * iota)  // 1024
    MB                      // 1048576
    GB                      // 1073741824
    TB                      // 1099511627776
)

// Expressões complexas
const (
    Flag1 = 1 << iota  // 1 (binário: 0001)
    Flag2              // 2 (binário: 0010)
    Flag3              // 4 (binário: 0100)
    Flag4              // 8 (binário: 1000)
)

// Constantes Tipadas vs Não-Tipadas
// Não-tipada (mais flexível)
const x1 = 42
var a2 int = x1
var b3 float64 = x1
var c4 byte = x1
// Todas funcionam!

// Tipada (mais restrita)
const y1 int = 42
var d int = y1       // OK
// var e float64 = y // ❌ ERRO - tipos incompatíveis

// 4. Operadores
// Operadores Aritméticos
// a := 10
// b := 3

// soma := a + b        // 13
// sub := a - b         // 7
// mult := a * b        // 30
// div := a / b         // 3 (divisão inteira!)
// resto := a % b       // 1 (módulo)

// Divisão com float
// divFloat := float64(a) / float64(b)  // 3.333...

// Incremento e decremento
// a++  // a = a + 1
// b--  // b = b - 1

// ❌ NÃO existe ++a ou --b em Go
// ❌ NÃO existe a = b++ em Go (++ não retorna valor)

// Operadores de Atribuição
// a = 10
// b = 3
// a += b  // a = a + b = 13
// a -= b  // a = a - b = 10
// a *= b  // a = a * b = 30
// a /= b  // a = a / b = 10

// Operadores de Comparação
// a := 10
// b := 5

// a == b   // false (igual)
// a != b   // true  (diferente)
// a > b    // true  (maior)
// a < b    // false (menor)
// a >= b   // true  (maior ou igual)
// a <= b   // false (menor ou igual)

// Operadores Lógicos
// t := true
// f := false

// t && f   // false (E lógico - AND)
// t || f   // true  (OU lógico - OR)
// !t       // false (NÃO lógico - NOT)

// Short-circuit evaluation
// func custoso() bool {
//     fmt.Println("Executado!")
//     return true
// }

// false && custoso()  // custoso() NÃO é executado
// true || custoso()   // custoso() NÃO é executado

// Operadores Bitwise
// a := 60  // 0011 1100
// b := 13  // 0000 1101

// a & b    // 12   (0000 1100) - AND
// a | b    // 61   (0011 1101) - OR
// a ^ b    // 49   (0011 0001) - XOR
// a &^ b   // 48   (0011 0000) - AND NOT (bit clear)
// ^a       // -61  (1100 0011) - NOT

// Shifts
// a << 2   // 240  (1111 0000) - shift left
// a >> 2   // 15   (0000 1111) - shift right

// Precedência de Operadores
// Prioridade (maior para menor):
// 1. *, /, %, <<, >>, &, &^
// 2. +, -, |, ^
// 3. ==, !=, <, <=, >, >=
// 4. &&
// 5. ||

// resultado := 2 + 3 * 4    // 14 (não 20)
// resultado = (2 + 3) * 4   // 20 (com parênteses)

// 5. Conversão de Tipos
// Go é fortemente tipado - conversões precisam ser explícitas!
// Conversão entre Numéricos
var i int = 42
var f float64 = float64(i)  // int → float64
var u7 uint = uint(f)         // float64 → uint

// Cuidado com perda de precisão!
var x float64 = 3.99
var y7 int = int(x)  // 3 (trunca, não arredonda!)

// Overflow
var a7 int8 = 127
// var b int8 = int8(128)  // Overflow! (wrap around)



