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
	lexeme       string
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

		if err == io.EOF {
			fmt.Println("End of File.")
			break
		}
		scanLine(&lexer, line)
	}

}

func scanLine(lexer *Lexer, line string) {
	lexer.line = line
	fmt.Println("lexer.line: ", lexer.line)

L:
	for {
		c := lexer.advance()

		if c == "" {
			break
		}

		switch {
		case isNumber(c): // MonthAndDay or Time
			lexer.addToken(lexer.number(c))
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

		//switch c {
		//case "y":

		//default:
		//}
		//fmt.Println("End of for, lexer.tokens", lexer.tokens)
		fmt.Println("Not found for:", c)
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
	if l.position+amount >= len(l.line) {
		return "\n"
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
			lexeme: l.lexeme,
		}

		fmt.Println("returning token for MONTHANDDAY", tok)
		return tok
	}
	var tok = Token{
		tokenType: TIME,
		//lexeme:
		lexeme: l.lexeme,
	}
	fmt.Println("returning token for TIME", tok)
	//return Token{
	//	tokenType: TIME,
	//	//lexeme:
	//	lexeme: lexeme,
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
