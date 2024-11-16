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
	literal      string
	line         string
	tokens       []Token
}

func ScanTokens(reader *bufio.Reader) {
	lexer := Lexer{0, 0, "", "", []Token{}} // TODO add stuff to lexer

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

	fmt.Println("Tokens", lexer.tokens)
}

func scanLine(lexer Lexer, line string) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	lexer.line = line

	for {
		c, hasText := lexer.advance(scanner)
		fmt.Println("c: ", c)

		if hasText {
			//fmt.Println(c)
		} else {
			//fmt.Println("No hasText")
			break
		}

		switch {
		case isNumber(c): // MonthAndDay or Time
			lexer.tokens = append(lexer.tokens, lexer.number(c, scanner))

		case isLetter(c):
			fmt.Println("letter.")
		default:
			//fmt.Println(string(c))
		}

		//// Single rune checking
		//switch c {
		//case 'y': // Year or potential start of food, variable, or value
		//	ahead, err := lexer.lookahead(1)
		//	if err == bufio.ErrBufferFull {
		//		fmt.Println("End of file reading number", bufio.ErrBufferFull)
		//	}

		//case '/':
		//	if lexer.lookahead(a) == '/' {

		//	} else {

		//	}
		//}
	}
	fmt.Println(lexer.tokens)
}

func (l *Lexer) advance(s *bufio.Scanner) (string, bool) {
	l.position++
	if hasText := s.Scan(); hasText {
		l.literal += s.Text()
		return s.Text(), true
	}
	return "", false
}

// TODO lookahead
func (l *Lexer) lookahead(amount int) string {
	if l.position+amount >= len(l.line) {
		return "\n"
	}
	return l.line[l.position : l.position+amount]
}

func (l *Lexer) addToken(token Token) {
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) scanningError(token Token) {
	fmt.Println("Error at line ", l.position)
}

// TODO
func (l *Lexer) number(c string, s *bufio.Scanner) Token {
	ahead := l.lookahead(1)
	lit := ""

	for isNumber(ahead) {
		adv, _ := l.advance(s)
		lit += adv
		fmt.Println(lit)
		ahead = l.lookahead(1)
	}

	// Month
	if strings.Contains(ahead, "/") {
		l.advance(s)
		day := l.lookahead(1)
		if day == "\n" {
			fmt.Println("End of file reading number", bufio.ErrBufferFull) // TODO what is this again, ErrBufferFull?
		}

		if !isNumber(day) {
			fmt.Println("Error, no day given in MonthAndDay token")
			return Token{}
		}
		for isNumber(day) {
			l.advance(s)
		}
		return Token{
			tokenType: MONTHANDDAY,
			//lexeme:
			literal: lit,
		}
	}
	var tok = Token{
		tokenType: TIME,
		//lexeme:
		literal: lit,
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

//func isNumbers(s []byte) bool {
//	for _, c := range s {
//		if '0' > c && c > '9' {
//			return false
//		}
//	}
//	return true
//}

//func isInt(s string) bool {
//	for _, c := range s {
//		if !unicode.IsDigit(c) {
//			return false
//		}
//	}
//	return true
//}

func isLetter(c string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	return isAlpha(c)
}

func isAlphaNumeric(c string) bool {
	return isLetter(c) && isNumber(c)
}
