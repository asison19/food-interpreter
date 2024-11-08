package lexer

const (
	YEAR = iota
	SEMICOLON
	EOL
	MONTHANDDAY
	TIME
	FOOD
	REPEATER
	COMMA
	SLEEP
)

type Token struct {
	tokenType int
	lexeme    string
	literal   string
}
