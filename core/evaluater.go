package core

import (
	"fmt"

	"github.com/egelis/jparser"
)

const (
	errUnknownToken           = "unknown token"
	errDifferentType          = "operands must be of the same type"
	errInvalidOperatorForType = "operator not defined for types"
	errTypeCast               = "typecast error"
	errDivisionByZero         = "division by zero"
)

type UnknownParameterError struct {
	Param string
}

func (e *UnknownParameterError) Error() string {
	return fmt.Sprintf("parameter not found: %s", e.Param)
}

type CalculationError struct {
	Reason string
	Value  string
}

func (e *CalculationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Reason, e.Value)
}

func Evaluate(tokens []Token, knownParams jparser.RawMessageSet) (Token, error) {
	resStack := make(TokenStack, 0, len(tokens))

	for _, token := range tokens {
		switch token.Type {
		case LOG_OP, COMP_OP, ARITH_OP:
			opFunc, ok := operatorFuncs[token.Value]
			if !ok {
				return Token{}, &CalculationError{Reason: errUnknownToken, Value: token.Value}
			}

			y, _ := resStack.Pop()
			x, _ := resStack.Pop()

			res, err := opFunc(x, y)
			if err != nil {
				return Token{}, err
			}

			resStack.Push(*res)
		case NUMBER, BOOL:
			resStack.Push(token)
		case IDENT:
			rawValue, ok := knownParams[token.Value]
			if !ok {
				return Token{}, &UnknownParameterError{Param: token.Value}
			}

			resStack.Push(Token{
				Type:      token.Type,
				Value:     string(rawValue),
				ValueType: token.ValueType,
			})
		// TODO: case DATE:
		default:
			return Token{}, &UnknownTokenTypeError{TokenType: token.Type}
		}
	}

	res, _ := resStack.Pop()

	switch res.ValueType {
	case BOOL_TYPE:
		return Token{
			Type:      BOOL,
			Value:     res.Value,
			ValueType: BOOL_TYPE,
		}, nil
	case NUMBER_TYPE:
		return Token{
			Type:      NUMBER,
			Value:     res.Value,
			ValueType: NUMBER_TYPE,
		}, nil
	default:
		return Token{}, &CalculationError{Reason: errUnknownToken, Value: res.Value}
	}
}
