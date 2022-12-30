package core

import (
	"fmt"
	"strconv"
)

type operatorFunc func(x, y Token) (res *Token, err error)

var operatorFuncs = map[string]operatorFunc{
	"OR":  orOperator,
	"AND": andOperator,
	"=":   equalOperator,
	"!=":  notEqualOperator,
	">":   moreOperator,
	"<":   lessOperator,
	">=":  moreEqualOperator,
	"<=":  lessEqualOperator,
	"+":   addOperator,
	"-":   subOperator,
	"*":   mulOperator,
	"/":   divOperator,
}

func addOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      NUMBER,
			Value:     fmt.Sprintf("%f", op1+op2),
			ValueType: NUMBER_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func subOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      NUMBER,
			Value:     fmt.Sprintf("%f", op1-op2),
			ValueType: NUMBER_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func mulOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      NUMBER,
			Value:     fmt.Sprintf("%f", op1*op2),
			ValueType: NUMBER_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func divOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		if op2 == 0 {
			return nil, &CalculationError{Reason: errDivisionByZero}
		}

		return &Token{
			Type:      NUMBER,
			Value:     fmt.Sprintf("%f", op1/op2),
			ValueType: NUMBER_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func lessEqualOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 <= op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func moreEqualOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 >= op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func moreOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 > op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func lessOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 < op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func orOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == BOOL_TYPE {
		op1, op2, err := parseBoolOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 || op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func andOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == BOOL_TYPE {
		op1, op2, err := parseBoolOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 && op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func equalOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == BOOL_TYPE {
		op1, op2, err := parseBoolOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 == op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 == op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func notEqualOperator(x, y Token) (res *Token, err error) {
	if x.ValueType != y.ValueType {
		return nil, &CalculationError{
			Reason: errDifferentType,
			Value:  fmt.Sprintf("'%s' '%s'", x.Value, y.Value),
		}
	}

	if x.ValueType == BOOL_TYPE {
		op1, op2, err := parseBoolOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 != op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	if x.ValueType == NUMBER_TYPE {
		op1, op2, err := parseFloatOperands(x.Value, y.Value)
		if err != nil {
			return nil, err
		}

		return &Token{
			Type:      BOOL,
			Value:     fmt.Sprintf("%t", op1 != op2),
			ValueType: BOOL_TYPE,
		}, nil
	}

	return nil, &CalculationError{
		Reason: errInvalidOperatorForType,
		Value:  fmt.Sprintf("'%s'", x.ValueType),
	}
}

func parseBoolOperands(value1, value2 string) (op1, op2 bool, err error) {
	op1, err = strconv.ParseBool(value1)
	if err != nil {
		return false, false, &CalculationError{
			Reason: errTypeCast,
			Value:  fmt.Sprintf("'%s' failed cast to '%s'", value1, BOOL_TYPE),
		}
	}

	op2, err = strconv.ParseBool(value2)
	if err != nil {
		return false, false, &CalculationError{
			Reason: errTypeCast,
			Value:  fmt.Sprintf("'%s' failed cast to '%s'", value2, BOOL_TYPE),
		}
	}

	return op1, op2, nil
}

func parseFloatOperands(value1, value2 string) (op1, op2 float64, err error) {
	op1, err = strconv.ParseFloat(value1, 64)
	if err != nil {
		return 0, 0, &CalculationError{
			Reason: errTypeCast,
			Value:  fmt.Sprintf("'%s' failed cast to '%s'", value1, NUMBER_TYPE),
		}
	}

	op2, err = strconv.ParseFloat(value2, 64)
	if err != nil {
		return 0, 0, &CalculationError{
			Reason: errTypeCast,
			Value:  fmt.Sprintf("'%s' failed cast to '%s'", value2, NUMBER_TYPE),
		}
	}

	return op1, op2, nil
}
