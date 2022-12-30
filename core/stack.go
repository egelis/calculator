package core

type TokenStack []Token

func (s *TokenStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *TokenStack) Push(value Token) {
	*s = append(*s, value)
}

func (s *TokenStack) Pop() (Token, bool) {
	if s.IsEmpty() {
		return Token{}, false
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value, true
}

func (s *TokenStack) Peek() (Token, bool) {
	if s.IsEmpty() {
		return Token{}, false
	}
	return (*s)[len(*s)-1], true
}
