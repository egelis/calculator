package core

import (
	"github.com/egelis/jparser"
)

func Calculate(tokens []Token, knownParams jparser.RawMessageSet) (Token, error) {
	exp, err := ToPostfixExp(tokens)
	if err != nil {
		return Token{}, err
	}

	res, err := Evaluate(exp, knownParams)
	if err != nil {
		return Token{}, err
	}

	return res, nil
}
