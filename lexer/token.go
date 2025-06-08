package lexer

import (
	"fmt"
	"strings"
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

func GetTokensAsString(tokens []Token) string {

	s := []string{}
	for _, v := range tokens {
		s = append(s, v.String())
	}
	return strings.Join(s, ", ")
}
