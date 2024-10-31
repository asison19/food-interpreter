package token

const (
	YEAR = iota
	SEMICOLON
	EOL
	MONTHANDDAY
	TIME
	FOOD
	REPEATER
	COMMA
)

type token struct {
	tokenType int
	lexeme    string
}
