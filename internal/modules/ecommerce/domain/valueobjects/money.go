package valueobjects

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
)

// Money representa um valor monetário com precisão decimal
type Money struct {
	amount   *big.Int // Valor em centavos (para evitar problemas de ponto flutuante)
	currency string   // Moeda (BRL, USD, EUR, etc.)
}

// NewMoney cria um novo valor monetário
func NewMoney(amount float64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, errors.ErrInvalidMoney
	}

	if currency == "" {
		return nil, errors.ErrInvalidMoney
	}

	// Converter para centavos (multiplicar por 100)
	// Usar big.Float para precisão
	amountFloat := big.NewFloat(amount)
	hundred := big.NewFloat(100)
	amountFloat.Mul(amountFloat, hundred)

	// Converter para big.Int (centavos)
	amountInt, _ := amountFloat.Int(nil)

	return &Money{
		amount:   amountInt,
		currency: strings.ToUpper(currency),
	}, nil
}

// NewMoneyFromCents cria um valor monetário a partir de centavos
func NewMoneyFromCents(cents int64, currency string) (*Money, error) {
	if cents < 0 {
		return nil, errors.ErrInvalidMoney
	}

	if currency == "" {
		return nil, errors.ErrInvalidMoney
	}

	return &Money{
		amount:   big.NewInt(cents),
		currency: strings.ToUpper(currency),
	}, nil
}

// Amount retorna o valor em centavos
func (m Money) Amount() *big.Int {
	return new(big.Int).Set(m.amount)
}

// AmountFloat retorna o valor como float64
func (m Money) AmountFloat() float64 {
	amountFloat := new(big.Float).SetInt(m.amount)
	hundred := big.NewFloat(100)
	amountFloat.Quo(amountFloat, hundred)
	result, _ := amountFloat.Float64()
	return result
}

// Currency retorna a moeda
func (m Money) Currency() string {
	return m.currency
}

// String retorna a representação string do valor
func (m Money) String() string {
	amount := m.AmountFloat()
	return fmt.Sprintf("%.2f %s", amount, m.currency)
}

// Equals verifica se dois valores monetários são iguais
func (m Money) Equals(other Money) bool {
	return m.amount.Cmp(other.amount) == 0 && m.currency == other.currency
}

// IsZero verifica se o valor é zero
func (m Money) IsZero() bool {
	return m.amount.Cmp(big.NewInt(0)) == 0
}

// IsPositive verifica se o valor é positivo
func (m Money) IsPositive() bool {
	return m.amount.Cmp(big.NewInt(0)) > 0
}

// IsNegative verifica se o valor é negativo
func (m Money) IsNegative() bool {
	return m.amount.Cmp(big.NewInt(0)) < 0
}

// Add soma dois valores monetários
func (m Money) Add(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("não é possível somar moedas diferentes: %s e %s", m.currency, other.currency))
	}

	result := new(big.Int).Add(m.amount, other.amount)
	return Money{
		amount:   result,
		currency: m.currency,
	}
}

// Subtract subtrai dois valores monetários
func (m Money) Subtract(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("não é possível subtrair moedas diferentes: %s e %s", m.currency, other.currency))
	}

	result := new(big.Int).Sub(m.amount, other.amount)
	return Money{
		amount:   result,
		currency: m.currency,
	}
}

// Multiply multiplica o valor por um número
func (m Money) Multiply(factor float64) Money {
	if factor < 0 {
		panic("fator de multiplicação não pode ser negativo")
	}

	factorFloat := big.NewFloat(factor)
	amountFloat := new(big.Float).SetInt(m.amount)
	resultFloat := new(big.Float).Mul(amountFloat, factorFloat)
	resultInt, _ := resultFloat.Int(nil)

	return Money{
		amount:   resultInt,
		currency: m.currency,
	}
}

// Divide divide o valor por um número
func (m Money) Divide(divisor float64) Money {
	if divisor <= 0 {
		panic("divisor deve ser maior que zero")
	}

	divisorFloat := big.NewFloat(divisor)
	amountFloat := new(big.Float).SetInt(m.amount)
	resultFloat := new(big.Float).Quo(amountFloat, divisorFloat)
	resultInt, _ := resultFloat.Int(nil)

	return Money{
		amount:   resultInt,
		currency: m.currency,
	}
}

// Compare compara dois valores monetários
// Retorna: -1 se m < other, 0 se m == other, 1 se m > other
func (m Money) Compare(other Money) (int, error) {
	if m.currency != other.currency {
		return 0, fmt.Errorf("não é possível comparar moedas diferentes: %s e %s", m.currency, other.currency)
	}

	return m.amount.Cmp(other.amount), nil
}

// GreaterThan verifica se m > other
func (m Money) GreaterThan(other Money) (bool, error) {
	compare, err := m.Compare(other)
	if err != nil {
		return false, err
	}
	return compare > 0, nil
}

// LessThan verifica se m < other
func (m Money) LessThan(other Money) (bool, error) {
	compare, err := m.Compare(other)
	if err != nil {
		return false, err
	}
	return compare < 0, nil
}

// GreaterThanOrEqual verifica se m >= other
func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	compare, err := m.Compare(other)
	if err != nil {
		return false, err
	}
	return compare >= 0, nil
}

// LessThanOrEqual verifica se m <= other
func (m Money) LessThanOrEqual(other Money) (bool, error) {
	compare, err := m.Compare(other)
	if err != nil {
		return false, err
	}
	return compare <= 0, nil
}

// Format retorna o valor formatado para exibição
func (m Money) Format() string {
	amount := m.AmountFloat()

	// Formatar com separadores de milhares
	amountStr := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(amountStr, ".")

	// Adicionar separadores de milhares
	if len(parts[0]) > 3 {
		formatted := ""
		for i, digit := range parts[0] {
			if i > 0 && (len(parts[0])-i)%3 == 0 {
				formatted += "."
			}
			formatted += string(digit)
		}
		parts[0] = formatted
	}

	return fmt.Sprintf("%s,%s %s", parts[0], parts[1], m.currency)
}

// ParseMoney cria um Money a partir de uma string
func ParseMoney(valueStr, currency string) (*Money, error) {
	// Remover espaços e caracteres não numéricos exceto vírgula e ponto
	valueStr = strings.ReplaceAll(valueStr, " ", "")
	valueStr = strings.ReplaceAll(valueStr, ".", "")
	valueStr = strings.ReplaceAll(valueStr, ",", ".")

	amount, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf("valor inválido: %s", valueStr)
	}

	return NewMoney(amount, currency)
}

// MustNewMoney cria um Money sem retornar erro (para casos onde sabemos que é válido)
func MustNewMoney(amount float64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return *m
}
