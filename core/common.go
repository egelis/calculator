package core

type TokenType string

const (
	ARITH_OP    TokenType = "arithmeticOp"
	COMP_OP     TokenType = "comparisonOp"
	LOG_OP      TokenType = "logicOp"
	LBR         TokenType = "leftBranch"
	RBR         TokenType = "rightBranch"
	NUMBER      TokenType = "number"
	BOOL        TokenType = "boolWord"
	IDENT       TokenType = "identificator"
	EXISTS_FUNC TokenType = "existsFunc"
)

type ValueType string

const (
	NUMBER_TYPE  ValueType = "number"
	BOOL_TYPE    ValueType = "bool"
	UNKNOWN_TYPE ValueType = "unknown"
	// DATE etc.
)

type Token struct {
	Type      TokenType
	Value     string
	ValueType ValueType
}
