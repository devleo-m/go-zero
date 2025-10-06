package main

import "fmt"

// 🎯 O que você vai aprender

// Declaração com := (a forma que você vai usar 90% do tempo)
// Tipos básicos e quando usar cada um
// Zero values (adeus undefined/null!)
// Type inference
// Comparação direta: Go vs TypeScript

// 1. Declaração de Variáveis - A Forma Certa
// TypeScript - 3 formas

// const name: string = "João"  // imutável
// let age: number = 25         // mutável
// var old = true               // evitar! (function scope)

// Go - Use := (Short Declaration)

// name := "João"    // string (inferido)
// age := 25         // int (inferido)
// active := true    // bool (inferido)
// price := 19.99    // float64 (inferido)

// Pronto! Isso é 90% do que você vai usar.

// 2. Formas de Declarar Variáveis
// Forma 1: := (Short Declaration) - USE ESTA! ✅

// Dentro de funções
func exemploVariaveis() {
	name := "João"
	age := 25

	// Múltiplas variáveis
	x, y := 10, 20

	// Exemplo de função que retorna múltiplos valores
	// user, err := getUser(1)  // Comentado pois getUser não existe

	// Usar as variáveis para evitar erro de "declared and not used"
	fmt.Printf("Nome: %s, Idade: %d\n", name, age)
	fmt.Printf("x: %d, y: %d\n", x, y)
}

// Quando usar:

// ✅ Dentro de funções (99% dos casos)
// ✅ Quando o tipo é óbvio
// ✅ Variáveis locais

// Limitações:

// ❌ NÃO funciona fora de funções (escopo de pacote)

// Forma 2: var com tipo - Use quando necessário

// Fora de funções (escopo de pacote)
var AppName string = "MeuApp"
var Port int = 3000

// func main() {
//     // Ou quando precisa declarar sem inicializar
//     var user User
//     var count int

//     // Múltiplas variáveis do mesmo tipo
//     var x, y, z int
// }

// Quando usar:

// ✅ Variáveis de pacote (globais)
// ✅ Quando precisa declarar sem valor inicial
// ✅ Quando quer ser MUITO explícito sobre o tipo

// Forma 3: var com inferência
// var name = "João"    // string inferido
// var age = 25         // int inferido
// Quando usar:

// ⚠️ Raramente! Prefira := dentro de funções
// ✅ Apenas para variáveis de pacote quando o tipo é óbvio

// ❌ O que NÃO fazer
// ❌ Redundante
// var name string = "João"  // tipo óbvio, use :=

// ❌ Fora de função
// name := "João"  // ERRO! := só funciona dentro de funções

// ❌ Redeclaração
// name := "João"
// name := "Maria"  // ERRO! name já foi declarado
