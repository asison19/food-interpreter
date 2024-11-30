package lexer

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TODO current location of the file?
type Lexer struct {
	position     int
	nextPosition int
	linePosition int
	lexeme       string
	line         string
	tokens       []Token
}

func ScanTokens(scanner *bufio.Scanner) {
	lexer := Lexer{0, 0, 0, "", "", []Token{}}

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println("Line:", string(line))
		scanLine(&lexer, string(line))
	}

}

func scanLine(lexer *Lexer, line string) {
	lexer.line = line

	// TODO error handling
	// TODO do I need this for loop?
L:
	for {
		c := lexer.advance()

		if c == "" {
			break
		}

		switch {
		case isNumber(c): // MONTHANDDAY or TIME
			lexer.addToken(lexer.number())
			continue L
		}

		switch c {
		// Ignoreables
		case " ":
		case "\n":
			lexer.clearLexeme()
			continue L
		case "y": // YEAR
			lexer.addToken(lexer.year())
			continue L
		case ";": // SEMICOLON
			lexer.addToken(Token{
				tokenType: SEMICOLON,
				lexeme:    lexer.lexeme,
			})
		case ",": // COMMA
			lexer.addToken(Token{
				tokenType: COMMA,
				lexeme:    lexer.lexeme,
			})
			continue L
		// TODO
		case "(":
			continue L
		case ".":
			lexer.addToken(lexer.reapeater())
			continue L
		}

		// Food, variables, sleep, etc.
		if isAlpha(c) {
			lexer.addToken(lexer.identifier())
		}
	}
	lexer.clearLine()
	fmt.Println("End of line, lexer.tokens", lexer.tokens)
}

func (l *Lexer) advance() string {
	defer func() { l.position += 1 }()
	if l.position >= len(l.line) {
		return ""
	}
	c := string(l.line[l.position])
	l.lexeme += c
	return c
}

// Used to ignore certain tokens
// TODO potentially just use a lookahead?
func (l *Lexer) clearLexeme() {
	l.lexeme = ""
}

func (l *Lexer) clearLine() {
	l.line = "" // TODO probably don't need, but leave just in case?
	l.position = 0
}

func (l *Lexer) lookahead(amount int) string {
	if l.position+amount > len(l.line) {
		return "\n" // TODO return empty string?
	}
	return l.line[l.position : l.position+amount]
}

func (l *Lexer) addToken(token Token) {
	l.tokens = append(l.tokens, token)
	l.lexeme = ""
}

func (l *Lexer) scanningError(token Token) {
	fmt.Println("Error at line ", l.position)
}

func (l *Lexer) reapeater() Token {
	ahead := l.lookahead(1)
	if ahead == "." {
		l.advance()
		return Token{
			tokenType: REPEATER,
			lexeme:    l.lexeme,
		}
	}
	return Token{}
}

func (l *Lexer) year() Token {
	ahead := l.lookahead(1)

	for isNumber(ahead) {
		l.advance()
		ahead = l.lookahead(1)
	}
	return Token{
		tokenType: YEAR,
		lexeme:    l.lexeme,
	}
}

func (l *Lexer) identifier() Token {
	ahead := l.lookahead(1)

	for isAlphaNumeric(ahead) {
		l.advance()
		ahead = l.lookahead(1)
	}

	tokenType, ok := reservedWords[l.lexeme]

	// FOOD
	if !ok {
		return Token{
			tokenType: FOOD,
			lexeme:    l.lexeme,
		}
	}

	return Token{
		tokenType: tokenType,
		lexeme:    l.lexeme,
	}
}

// TODO certain dates should not be possible
// E.g. 0/0, 1/32, etc.
func (l *Lexer) number() Token {
	ahead := l.lookahead(1)

	for isNumber(ahead) {
		l.advance()
		ahead = l.lookahead(1)
	}

	// Month
	if strings.Contains(ahead, "/") {
		l.advance()
		day := l.lookahead(1)
		if day == "\n" {
			fmt.Println("End of file reading number", bufio.ErrBufferFull)
		}

		if !isNumber(day) {
			fmt.Println("Error, no day given in MonthAndDay token")
			return Token{}
		}
		for isNumber(day) {
			l.advance()
			day = l.lookahead(1)
		}
		return Token{
			tokenType: MONTHANDDAY,
			//lexeme:
			lexeme: l.lexeme,
		}
	}
	return Token{
		tokenType: TIME,
		//lexeme:
		lexeme: l.lexeme,
	}
}

func isNumber(c string) bool {
	if _, err := strconv.Atoi(c); err == nil {
		return true
	}
	return false
}

func isAlpha(c string) bool {
	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(c)
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isNumber(c)
}
