package lexer

import (
	"fmt"
)

const (
	YEAR = iota
	SEMICOLON
	MONTHANDDAY
	TIME
	FOOD
	REPEATER
	COMMA
	SLEEP
	// TODO
	//LEFT_PAREN
	//RIGHT_PAREN
)

type Token struct {
	tokenType int
	//lexeme    string
	lexeme string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token Type: %d, Lexeme: %s", t.tokenType, t.lexeme)
}
