package lexer

import (
	"fmt"
)

type TokenType int

const (
	YEAR TokenType = iota
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
	Type TokenType
	//lexeme    string
	Lexeme string
}

func (t TokenType) String() string {
	return [...]string{"YEAR", "SEMICOLON", "MONTHANDDAY", "TIME", "FOOD", "REPEATER", "COMMA", "SLEEP"}[t]
}

func (t *Token) String() string {
	return fmt.Sprintf("Token Type: %d, Lexeme: %s", t.Type, t.Lexeme)
}
