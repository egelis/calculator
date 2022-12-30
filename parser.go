// nolint:nestif
package calculator

import (
	"errors"
	"fmt"

	"github.com/egelis/calculator/core"
	"github.com/egelis/jparser"
)

const (
	errSyntax = "found a syntax error"
	errCalc   = "calculation failed"
)

type ParseError struct {
	Reason string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("error: %s", e.Reason)
}

type parser struct {
	rawSet jparser.RawMessageSet
	tokens []core.Token

	tokensSize        int
	it                int
	calculationTokens []core.Token
}

func newParser(
	tokens []core.Token,
	rawSet jparser.RawMessageSet,
) *parser {
	return &parser{
		rawSet:            rawSet,
		tokens:            tokens,
		tokensSize:        len(tokens),
		it:                -1,
		calculationTokens: make([]core.Token, 0, len(tokens)),
	}
}

// START: LOG_EXP

// LOG_EXP: LOG_TERM => {LOG_OP | COMP_OP => LOG_TERM}
// LOG_TERM: BOOL | EXISTS | ARITH_EXP | ( "(" => LOG_EXP => ")" )

// ARITH_EXP: ARITH_TERM => {ARITH_OP => ARITH_TERM}
// ARITH_TERM: NUM | IDENT ( "(" => ARITH_EXP => ")" )

// EXISTS: 'exists' => '(' => IDENT => ')'

// Конечные:
// BOOL: true, false
// LOG_OP: AND, OR
// COMP_OP: > < != = >= <=
// NUM: 2.45, 2
// IDENT: param_123, denmt123

// START: LOGIC_EXP
func (p *parser) start() (core.Token, error) {
	if !p.checkNext(p.LogicExp) {
		// TODO: уточнить ошибку
		return core.Token{}, &ParseError{Reason: errSyntax}
	}

	// Если остались неразобранные токены, то они не подошли под правила
	if p.it+1 != p.tokensSize {
		// TODO: уточнить ошибку
		return core.Token{}, &ParseError{Reason: errSyntax}
	}

	// Result calculation
	res, err := core.Calculate(p.calculationTokens, p.rawSet)
	if err != nil {
		var paramErr *core.UnknownParameterError
		if errors.As(err, &paramErr) {
			return core.Token{
				Type:      core.BOOL,
				Value:     "false",
				ValueType: core.BOOL_TYPE,
			}, nil
		}

		return core.Token{}, &ParseError{Reason: fmt.Sprintf("%s: %s", errCalc, err)}
	}

	return res, nil
}

// LOGIC_EXP: LOGIC_TERM => {LOG_OP | COMP_OP => LOGIC_TERM}
func (p *parser) LogicExp() bool {
	if !p.checkNext(p.LogicTerm) {
		return false
	}

	for {
		savedIt := p.it

		if !p.checkNext(p.LogicOperator) {
			p.it = savedIt

			if !p.checkNext(p.CompOperator) {
				p.it = savedIt
				break
			}
		}

		// Add logic or comparison operator
		p.calculationTokens = append(p.calculationTokens, p.tokens[p.it])

		if !p.checkNext(p.LogicTerm) {
			p.it = savedIt
			break
		}
	}

	return true
}

// LOGIC_TERM: BOOL | EXISTS | ARITH_EXP | ( "(" => LOGIC_EXP => ")" )
func (p *parser) LogicTerm() bool {
	savedIt := p.it

	if !p.checkNext(p.Bool) {
		p.it = savedIt

		if !p.checkNext(p.ExistsFunc) {
			p.it = savedIt

			startIt := p.it

			if !p.checkNext(p.ArithmeticExp) {
				p.it = savedIt

				if !p.checkNext(p.LBracket) {
					return false
				}

				// add bracket
				p.calculationTokens = append(p.calculationTokens, p.tokens[p.it])

				if !p.checkNext(p.LogicExp) {
					return false
				}

				if !p.checkNext(p.RBracket) {
					return false
				}

				// add bracket
				p.calculationTokens = append(p.calculationTokens, p.tokens[p.it])
			} else {
				// Add arithmetic expression
				p.calculationTokens = append(p.calculationTokens, p.tokens[startIt+1:p.it+1]...)
			}
		}
	} else {
		// Add bool
		p.calculationTokens = append(p.calculationTokens, p.tokens[p.it])
	}

	return true
}

// ARITH_EXP: ARITH_TERM => {ARITH_OP => ARITH_TERM}
func (p *parser) ArithmeticExp() bool {
	if !p.checkNext(p.ArithmeticTerm) {
		return false
	}

	for {
		savedIt := p.it

		if !p.checkNext(p.ArithmeticOp) {
			p.it = savedIt
			break
		}

		if !p.checkNext(p.ArithmeticTerm) {
			p.it = savedIt
			break
		}
	}

	return true
}

// ARITH_TERM: NUM | IDENT | ( "(" => ARITH_EXP => ")" )
func (p *parser) ArithmeticTerm() bool {
	savedIt := p.it
	if !p.checkNext(p.Num) {
		p.it = savedIt

		if !p.checkNext(p.Ident) {
			p.it = savedIt

			if !p.checkNext(p.LBracket) {
				return false
			}

			if !p.checkNext(p.ArithmeticExp) {
				return false
			}

			if !p.checkNext(p.RBracket) {
				return false
			}
		}
	}

	return true
}

// EXISTS: 'exists' -> '(' -> IDENT -> ')'
func (p *parser) ExistsFunc() bool {
	if !p.checkNext(p.IdentExists) {
		return false
	}

	if !p.checkNext(p.LBracket) {
		return false
	}

	if !p.checkNext(p.Ident) {
		return false
	}

	field := p.tokens[p.it].Value

	if !p.checkNext(p.RBracket) {
		return false
	}

	_, ok := p.rawSet[field]
	p.calculationTokens = append(p.calculationTokens, core.Token{
		Type:      core.BOOL,
		Value:     fmt.Sprintf("%t", ok),
		ValueType: core.BOOL_TYPE,
	})

	return true
}

// Нетерминалы

func (p *parser) IdentExists() bool {
	p.it++

	return p.tokens[p.it].Type == core.EXISTS_FUNC
}

func (p *parser) CompOperator() bool {
	p.it++

	return p.tokens[p.it].Type == core.COMP_OP
}

func (p *parser) LogicOperator() bool {
	p.it++

	return p.tokens[p.it].Type == core.LOG_OP
}

func (p *parser) LBracket() bool {
	p.it++

	return p.tokens[p.it].Type == core.LBR
}

func (p *parser) RBracket() bool {
	p.it++

	return p.tokens[p.it].Type == core.RBR
}

func (p *parser) Bool() bool {
	p.it++

	return p.tokens[p.it].Type == core.BOOL
}

func (p *parser) Num() bool {
	p.it++

	return p.tokens[p.it].Type == core.NUMBER
}

func (p *parser) Ident() bool {
	p.it++

	return p.tokens[p.it].Type == core.IDENT
}

func (p *parser) ArithmeticOp() bool {
	p.it++

	return p.tokens[p.it].Type == core.ARITH_OP
}

func (p *parser) checkNext(f func() bool) bool {
	if p.it+1 < p.tokensSize && f() {
		return true
	}

	return false
}
