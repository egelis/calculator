package core

import (
	"fmt"
)

var operatorPrecedence = map[string]int{
	"(":   10,
	"OR":  20,
	"AND": 30,
	"=":   40,
	"!=":  40,
	">":   50,
	"<":   50,
	">=":  50,
	"<=":  50,
	"+":   120,
	"-":   120,
	"/":   130,
	"*":   130,
}

type UnknownTokenTypeError struct {
	TokenType TokenType
}

func (e *UnknownTokenTypeError) Error() string {
	return fmt.Sprintf("unknown token type: %s", e.TokenType)
}

// ToPostfixExp implements the "Shunting yard" algorithm (https://en.wikipedia.org/wiki/Shunting_yard_algorithm)
func ToPostfixExp(infixExp []Token) ([]Token, error) {
	var (
		operationStack TokenStack
		output         []Token
	)

	for _, token := range infixExp {
		switch token.Type {
		case NUMBER, BOOL, IDENT:
			output = append(output, token)
		case LBR:
			operationStack.Push(token)
		case RBR:
			operation, ok := operationStack.Pop()
			for ok {
				if operation.Value == "(" {
					break
				}
				output = append(output, operation)
				operation, ok = operationStack.Pop()
			}
		case LOG_OP, COMP_OP, ARITH_OP:
			if weight, ok := operatorPrecedence[token.Value]; ok {
				stackOperation, ok := operationStack.Peek()

				for ok && operatorPrecedence[stackOperation.Value] > weight {
					operationStack.Pop()
					output = append(output, stackOperation)
					stackOperation, ok = operationStack.Peek()
				}

				operationStack.Push(token)
			}
		default:
			return nil, &UnknownTokenTypeError{TokenType: token.Type}
		}
	}

	for token, ok := operationStack.Pop(); ok; token, ok = operationStack.Pop() {
		output = append(output, token)
	}

	return output, nil
}
