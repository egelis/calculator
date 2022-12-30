// nolint:nlreturn,wsl,gomnd,varnamelen
package calculator

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/egelis/calculator/core"
)

type InvalidTokenError struct {
	Position int
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("invalid token at position: %d", e.Position)
}

func tokenize(input string, paramTypes map[string]core.ValueType) ([]core.Token, error) {
	chars := []rune(input)
	inputLen := len(chars)

	var tokens []core.Token
	for i := 0; i < inputLen; {
		char := chars[i]

		if unicode.IsSpace(char) {
			i++
			continue
		}

		if isAlpha(char) {
			start := i

			i++
			for i < inputLen && (isAlpha(chars[i]) || isDigit(chars[i]) || chars[i] == '_') {
				i++
			}

			var (
				tokenType core.TokenType
				valueType core.ValueType
			)
			switch {
			case isBool(chars[start:i]):
				tokenType = core.BOOL
				valueType = core.BOOL_TYPE
			case isOrAnd(chars[start:i]):
				tokenType = core.LOG_OP
			case isExistsFunc(chars[start:i]):
				tokenType = core.EXISTS_FUNC
			default:
				tokenType = core.IDENT

				var ok bool
				if valueType, ok = paramTypes[string(chars[start:i])]; !ok {
					valueType = core.UNKNOWN_TYPE
				}
			}

			tokens = append(tokens, core.Token{
				Type:      tokenType,
				Value:     string(chars[start:i]),
				ValueType: valueType,
			})

			continue
		}

		if isArithmeticOp(char) {
			i++
			tokens = append(tokens, core.Token{Type: core.ARITH_OP, Value: string(char)})
			continue
		}

		if isLeftBracket(char) {
			i++
			tokens = append(tokens, core.Token{Type: core.LBR, Value: string(char)})
			continue
		}

		if isRightBracket(char) {
			i++
			tokens = append(tokens, core.Token{Type: core.RBR, Value: string(char)})
			continue
		}

		start := i
		if isLogicOp(chars, &i, inputLen) {
			i++
			tokens = append(tokens, core.Token{Type: core.COMP_OP, Value: string(chars[start:i])})
			continue
		}

		start = i
		if isNumber(chars, &i, inputLen) {
			tokens = append(tokens, core.Token{
				Type:      core.NUMBER,
				Value:     string(chars[start:i]),
				ValueType: core.NUMBER_TYPE,
			})
			continue
		}

		return nil, &InvalidTokenError{Position: i}
	}

	return tokens, nil
}

func isExistsFunc(chars []rune) bool {
	return string(chars) == "exists"
}

var boolWords = map[string]struct{}{
	"true":  {},
	"false": {},
}

func isBool(chars []rune) bool {
	_, ok := boolWords[string(chars)]
	return ok
}

var logicOperators = map[string]struct{}{
	"OR":  {},
	"AND": {},
}

func isOrAnd(chars []rune) bool {
	_, ok := logicOperators[string(chars)]
	return ok
}

func isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isLeftBracket(char rune) bool {
	return char == '('
}

func isRightBracket(char rune) bool {
	return char == ')'
}

var arithmeticOp = map[string]struct{}{
	"+": {},
	"-": {},
	"*": {},
	"/": {},
}

func isArithmeticOp(char rune) bool {
	_, ok := arithmeticOp[string(char)]
	return ok
}

func isLogicOp(chars []rune, i *int, length int) bool {
	for *i < length {
		switch chars[*i] {
		case '>', '<':
			if *i+1 < length && chars[*i+1] == '=' {
				*i++
			}
			return true
		case '=':
			return true
		case '!':
			if *i+1 < length && chars[*i+1] == '=' {
				*i++
				return true
			}
			return false
		default:
			return false
		}
	}
	return false
}

func isNumber(chars []rune, i *int, inputLen int) bool {
	c := 0
	for *i < inputLen {
		_, err := strconv.ParseFloat(string(chars[*i-c:*i+1]), 64)
		if err != nil {
			break
		}
		c++
		*i++
	}
	return c > 0
}
