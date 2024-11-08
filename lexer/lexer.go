package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
	lexer := Lexer{0, 0, "", "", nil}

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

		if hasText {
			fmt.Println(c)
		} else {
			fmt.Println("No hasText")
			break
		}

		//switch {
		//case isNumber(byte(c)):
		//	str := c
		//	//tokens = append(tokens, addToken(c))
		//	ahead, err := lexer.lookahead(2)
		//	if err == bufio.ErrBufferFull {
		//		fmt.Println("End of file reading number", bufio.ErrBufferFull)
		//	}

		//	// MonthAndDay
		//	// TODO allow for month and day separator other than '/'?
		//	if strings.Contains(string(ahead), "/") {
		//		// TODO read 2 runes
		//		a, err := lexer.reader.ReadBytes(byte('/'))
		//		str += string(a)
		//		if err != nil {
		//			fmt.Println(err)
		//			return
		//		}
		//	}

		//case isLetter(byte(c)):
		//	fmt.Println("letter.")
		//default:
		//	fmt.Println(string(c))
		//}

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
}

func (l *Lexer) advance(s *bufio.Scanner) (string, bool) {
	l.position++
	if hasText := s.Scan(); hasText {
		return s.Text(), true
	}
	return "", false
}

// TODO lookahead
func (l *Lexer) lookahead(amount int) (string, bool) {
	if l.position+amount > len(l.line) {
		return "", false
	}
	return l.line[l.position : l.position+amount], true
}

func (l *Lexer) addToken(token Token) {
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) scanningError(token Token) {
	fmt.Println("Error at line ", l.position)
}

func isNumber(c byte) bool {
	return '0' <= c && c <= '9'
}

func isNumbers(s []byte) bool {
	for _, c := range s {
		if '0' > c && c > '9' {
			return false
		}
	}
	return true
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isAlphaNumeric(c byte) bool {
	return isLetter(c) && isAlphaNumeric(c)
}
