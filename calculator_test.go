// nolint:gochecknoglobals,dupl,revive
package calculator

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/egelis/calculator/core"
	"github.com/egelis/jparser"
)

type args struct {
	formulas    []Formula
	knownParams []jparser.RawMessageSet
	paramTypes  map[string]core.ValueType
}

func TestCalculateSuccess(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name          string
		args          args
		expectedColor Color
		expectedRes   []FormulaResult
	}{
		{
			name: "expression without parameters",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "1=2 AND 1=1 OR 1=1",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: nil,
				paramTypes:  nil,
			},
			expectedColor: GreenColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
				},
			},
		},

		{
			name: "expression without parameters with influence of brackets",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "1=2 AND (1=1 OR 1=1)",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: nil,
				paramTypes:  nil,
			},
			expectedColor: GreyColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
				},
			},
		},

		{
			name: "check expression with parameters",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "(true) AND (s2001 > (s6004 *0.1)) AND (s2001 > stated_capital) AND (s2001 > 1000000)",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
			expectedColor: GreenColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
				},
			},
		},

		{
			name: "'exists' function",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "exists(founder_url)=false AND s2001 > (s2001-s6004) * 0.5",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
			expectedColor: GreenColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
				},
			},
		},

		{
			name: "bool parameter",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "bool_param OR bool_param=true",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
			expectedColor: GreyColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
				},
			},
		},

		{
			name: "unknown parameter",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "exists(founder_url)=false AND s2001 > (s2001-s6004)* unknown_param + 0.5",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
			expectedColor: GreyColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
				},
			},
		},

		{
			name: "disabled and false formulas (grey color)",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "bool_param = true",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
					{
						Name:       "formula_2",
						Expression: "s2001 < 0",
						Color:      YellowColor,
						Version:    1,
						IsEnable:   false,
					},
					{
						Name:       "formula_3",
						Expression: "s2001 >= 0 AND exists(s6004)",
						Color:      RedColor,
						Version:    2,
						IsEnable:   false,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
			expectedColor: GreyColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
				},
			},
		},

		{
			name: "multiple parameters and multiple formulas (red color)",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "bool_param = false",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
					{
						Name:       "formula_2",
						Expression: "s2001 < 0",
						Color:      YellowColor,
						Version:    1,
						IsEnable:   true,
					},
					{
						Name:       "formula_3",
						Expression: "s2001 >= 0 AND exists(s6004)",
						Color:      RedColor,
						Version:    2,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithMultipleElements,
				paramTypes:  types,
			},
			expectedColor: RedColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  false,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  false,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  false,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  true,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  true,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  false,
					},
				},
			},
		},
		{
			name: "multiple parameters and multiple formulas (yellow color)",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "bool_param = false",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
					{
						Name:       "formula_2",
						Expression: "s2001 < 0",
						Color:      YellowColor,
						Version:    1,
						IsEnable:   true,
					},
					{
						Name:       "formula_3",
						Expression: "unknown_param > 0 AND exists(s6004)",
						Color:      RedColor,
						Version:    2,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithMultipleElements,
				paramTypes:  types,
			},
			expectedColor: YellowColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  false,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  false,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  false,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  false,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   YellowColor,
						Result:  true,
					},
					"formula_3": {
						Version: 2,
						Color:   RedColor,
						Result:  false,
					},
				},
			},
		},

		{
			name: "multiple parameters and multiple formulas (green color)",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "s2001 = 0",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
					{
						Name:       "formula_2",
						Expression: "1 = 0",
						Color:      RedColor,
						Version:    1,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithMultipleElements,
				paramTypes:  types,
			},
			expectedColor: GreenColor,
			expectedRes: []FormulaResult{
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
					"formula_2": {
						Version: 1,
						Color:   RedColor,
						Result:  false,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  true,
					},
					"formula_2": {
						Version: 1,
						Color:   RedColor,
						Result:  false,
					},
				},
				{
					"formula_1": {
						Version: 0,
						Color:   GreenColor,
						Result:  false,
					},
					"formula_2": {
						Version: 1,
						Color:   RedColor,
						Result:  false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resColor, formulaRes, err := Calculate(test.args.formulas, test.args.knownParams, test.args.paramTypes)
			if err != nil {
				t.Errorf("Calculate() res error = \"%v\", expected nil", err)
			}

			if !reflect.DeepEqual(formulaRes, test.expectedRes) {
				got, _ := json.MarshalIndent(formulaRes, "", "  ")
				expected, _ := json.MarshalIndent(test.expectedRes, "", "  ")
				t.Errorf("Calculate() got formulaRes = %s\n expected = %s", got, expected)
			}

			if resColor != test.expectedColor {
				t.Errorf("Calculate() got resColor = %s, expected = %s", resColor, test.expectedColor)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args args
	}{
		{
			name: "func exists with num",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "exists(5)",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: nil,
				paramTypes:  nil,
			},
		},

		{
			name: "only arithmetic expression",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "5 + 6 / s6004",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
		},

		{
			name: "wrong number of brackets",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "exists(s2001 = false",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
		},

		{
			name: "different type of operands",
			args: args{
				formulas: []Formula{
					{
						Name:       "formula_1",
						Expression: "bool_param - 1 > 0",
						Color:      GreenColor,
						Version:    0,
						IsEnable:   true,
					},
				},
				knownParams: paramsWithOneElement,
				paramTypes:  types,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resColor, formulaRes, err := Calculate(test.args.formulas, test.args.knownParams, test.args.paramTypes)
			if err == nil {
				t.Errorf("Calculate() got error = nil, expected error")
			}

			if len(formulaRes) != 0 {
				got, _ := json.MarshalIndent(formulaRes, "", "  ")
				t.Errorf("Calculate() got formulaRes = %s\n expected = []", got)
			}

			if resColor != BlackColor {
				t.Errorf("Calculate() got resColor = '%s', expected = 'grey'", resColor)
			}
		})
	}
}

var (
	paramsWithOneElement, _ = jparser.ParseParams(
		json.RawMessage(`
		{
			"s2001": 2000000,
			"s6004": 10,
			"stated_capital": 50,
			"bool_param": false
		}
		`),
		[]jparser.MetaData{
			{Path: "s2001", ParamID: "s2001"},
			{Path: "s6004", ParamID: "s6004"},
			{Path: "stated_capital", ParamID: "stated_capital"},
			{Path: "bool_param", ParamID: "bool_param"},
		},
	)

	paramsWithMultipleElements, _ = jparser.ParseParams(
		json.RawMessage(`
		{
			"stated_capital": 50,
			"bool_param": false,
			"x": [
				{
					"s2001": 10
				},
				{
					"s2001": 0,
					"s6004": 10
				},
				{
					"s2001": -10,
					"s6004": 10
				}
			]
		}
		`),
		[]jparser.MetaData{
			{Path: "x.[].s2001", ParamID: "s2001"},
			{Path: "x.[].s6004", ParamID: "s6004"},
			{Path: "stated_capital", ParamID: "stated_capital"},
			{Path: "bool_param", ParamID: "bool_param"},
		},
	)

	types = map[string]core.ValueType{
		"s2001":          core.NUMBER_TYPE,
		"s6004":          core.NUMBER_TYPE,
		"stated_capital": core.NUMBER_TYPE,
		"bool_param":     core.BOOL_TYPE,
	}
)
