package lexer

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	//"unicode"
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
	lexer := Lexer{0, 0, 0, "", "", []Token{}} // TODO add stuff to lexer

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println("Line:", string(line))
		scanLine(&lexer, string(line))
	}

}

func scanLine(lexer *Lexer, line string) {
	lexer.line = line

	// TODO do I need this for loop?
L:
	for {
		c := lexer.advance()

		if c == "" {
			break
		}

		switch {
		case isNumber(c): // MonthAndDay or Time
			lexer.addToken(lexer.number())
			continue L
		default:
			//fmt.Println(string(c))
		}

		switch c {
		case " ":
		case "\n":
			lexer.clearLexeme()
			continue L
		}

		switch c {
		case "y":
			lexer.addToken(lexer.year())
			continue L
		case ",":
			lexer.addToken(Token{
				tokenType: COMMA,
				lexeme:    lexer.lexeme,
			})
			continue L
		case "s":
		case "S":
			//sleep
		case "(":
		//case ";":
		//	continue L
		default:
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
			fmt.Println("End of file reading number", bufio.ErrBufferFull) // TODO what is this again, ErrBufferFull?
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

func isLetter(c string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(c)
}

func isAlphaNumeric(c string) bool {
	return isLetter(c) && isNumber(c)
}
