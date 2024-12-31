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
	TokenType int
	//lexeme    string
	Lexeme string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token Type: %d, Lexeme: %s", t.TokenType, t.Lexeme)
}
