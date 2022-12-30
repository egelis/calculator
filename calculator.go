package calculator

import (
	"strconv"

	"github.com/egelis/calculator/core"
	"github.com/egelis/jparser"
)

type Color string

const (
	BlackColor  = Color("black")
	GreyColor   = Color("grey")
	GreenColor  = Color("green")
	YellowColor = Color("yellow")
	RedColor    = Color("red")
)

type Formula struct {
	Name       string
	Expression string
	Color      Color
	Version    int64
	IsEnable   bool
}

type (
	FormulaResult map[string]Value

	Value struct {
		Version int64 `json:"version"`
		Color   Color `json:"color"`
		Result  bool  `json:"result"`
	}
)

// Calculate calculates each formula from 'formulas' for each set of parameters from 'rawSets'
func Calculate(formulas []Formula, rawSets []jparser.RawMessageSet, paramTypes map[string]core.ValueType,
) (Color, []FormulaResult, error) {
	tokenizedFormulas, err := getTokenizedFormulas(formulas, paramTypes)
	if err != nil {
		return BlackColor, nil, err
	}

	// For the situation where we have formulas without rawSet
	if len(rawSets) == 0 {
		rawSets = []jparser.RawMessageSet{nil}
	}

	formulaResults := make([]FormulaResult, 0, len(rawSets))
	resColor := GreyColor

	for _, rawSet := range rawSets {
		result := FormulaResult{}

		for _, formula := range tokenizedFormulas {
			resToken, err := newParser(formula.Tokens, rawSet).start()
			if err != nil {
				return BlackColor, nil, err
			}

			resValue, err := strconv.ParseBool(resToken.Value)
			if err != nil {
				return BlackColor, nil, err
			}

			if resValue && colorPrecedence[formula.Color] > colorPrecedence[resColor] {
				resColor = formula.Color
			}

			result[formula.Name] = Value{
				Version: formula.Version,
				Color:   formula.Color,
				Result:  resValue,
			}
		}

		if len(result) > 0 {
			formulaResults = append(formulaResults, result)
		}
	}

	return resColor, formulaResults, nil
}

// nolint:gochecknoglobals
var colorPrecedence = map[Color]int{
	BlackColor:  0,
	GreyColor:   1,
	GreenColor:  2,
	YellowColor: 3,
	RedColor:    4,
}

type tokenizedFormula struct {
	Tokens  []core.Token
	Name    string
	Version int64
	Color   Color
}

func getTokenizedFormulas(formulas []Formula, paramTypes map[string]core.ValueType) ([]tokenizedFormula, error) {
	res := make([]tokenizedFormula, 0, len(formulas))

	for _, formula := range formulas {
		if !formula.IsEnable {
			continue
		}

		formulaTokens, err := tokenize(formula.Expression, paramTypes)
		if err != nil {
			return nil, err
		}

		res = append(res, tokenizedFormula{
			Tokens:  formulaTokens,
			Name:    formula.Name,
			Version: formula.Version,
			Color:   formula.Color,
		})
	}

	return res, nil
}
