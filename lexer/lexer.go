package lexer

import (
	"bufio"
	"fmt"
	"io"
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
	literal      string
	line         string
	tokens       []Token
}

func ScanTokens(reader *bufio.Reader) {
	lexer := Lexer{0, 0, 0, "", "", []Token{}} // TODO add stuff to lexer

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		if err == io.EOF { // TODO broken
			fmt.Println("End of File.")
			break
		}
		scanLine(lexer, line)
	}

}

func scanLine(lexer Lexer, line string) {
	lexer.line = line

	//for i, ch := range line {
	for {
		c := lexer.advance()

		if c != "" {
		} else {
			break
		}

		switch {
		case isNumber(c): // MonthAndDay or Time
			lexer.tokens = append(lexer.tokens, lexer.number(c))
		default:
			//fmt.Println(string(c))
		}
	}
	fmt.Println(lexer.tokens)
}

// Returns position based off of a 1 indexed "line"
func (l *Lexer) advance() string {
	defer func() { l.position += 1 }()
	if l.position >= len(l.line) {
		return ""
	}
	c := string(l.line[l.position])
	l.literal += c
	return c
}

func (l *Lexer) lookahead(amount int) string {
	if l.position+amount >= len(l.line) {
		return "\n"
	}
	return l.line[l.position : l.position+amount]
}

func (l *Lexer) addToken(token Token) {
	l.tokens = append(l.tokens, token)
	l.literal = ""
}

func (l *Lexer) scanningError(token Token) {
	fmt.Println("Error at line ", l.position)
}

func (l *Lexer) number(c string) Token {
	ahead := l.lookahead(1)

	for isNumber(ahead) {
		l.advance()
		// fmt.Println(lit)
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
		var tok = Token{
			tokenType: MONTHANDDAY,
			//lexeme:
			literal: l.literal,
		}

		fmt.Println("returning token ", tok)
		return tok

	}
	var tok = Token{
		tokenType: TIME,
		//lexeme:
		literal: l.literal,
	}
	fmt.Println("returning token ", tok)
	//return Token{
	//	tokenType: TIME,
	//	//lexeme:
	//	literal: literal,
	//}
	return tok
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
