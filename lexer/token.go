package lexer

import (
	"fmt"
)

const (
	YEAR = iota
	SEMICOLON
	//EOL
	MONTHANDDAY
	TIME
	FOOD
	REPEATER
	COMMA
	SLEEP
)

type Token struct {
	tokenType int
	//lexeme    string
	lexeme string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token Type: %d, Lexeme: %s", t.tokenType, t.lexeme)
}
