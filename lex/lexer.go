package lexer

import (
	"bufio"
	"fmt"
	//"food-interpreter/token"
	"io"
	"strings"
)

// TODO current location of the file?
type Lexer struct {
	reader       *bufio.Reader
	position     int
	nextPosition int
	literal      string
	tokens       []Token
}

func ScanTokens(reader *bufio.Reader) {
	//count := 0

	lexer := Lexer{reader, 0, 0, "", nil}

	for {
		c, err := lexer.advance()

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		//
		switch {
		case isNumber(byte(c)):
			str := string(c)
			//tokens = append(tokens, addToken(c))
			ahead, err := lexer.lookahead(2)
			if err == bufio.ErrBufferFull {
				fmt.Println("End of file reading number", bufio.ErrBufferFull)
			}

			// Year
			// A year token must have at least 4 digits.
			// TODO have year always be right after a start of line because time is also 3-4 digits
			// TODO make year just start with a y?
			if isNumbers(ahead) {
				// consume until we reach a space?
				for {
					a, _, err := lexer.reader.ReadRune()
					if err == io.EOF {
						break
					}

					if !isNumber(byte(a)) {
						fmt.Println(lexer.position, " - Not a valid year.")
						break
					}

					switch a {
					//case ';':
					case '\n':
					case '\r':
						lexer.addToken(Token{YEAR, "YEAR", lexer.literal})
					}
					lexer.literal += string(a)
				}
			}

			// MonthAndDay
			// TODO allow for month and day separator other than '/'?
			if strings.Contains(string(ahead), "/") {
				// TODO read 2 runes
				a, err := lexer.reader.ReadBytes(byte('/'))
				str += string(a)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

		case isLetter(byte(c)):
			fmt.Println("letter.")
		default:
			fmt.Println(string(c))
		}
		// Single rune checking
		switch c {
		case 'y': // Year or potential start of food, variable, or value
			ahead, err := lexer.lookahead(1)
			if err == bufio.ErrBufferFull {
				fmt.Println("End of file reading number", bufio.ErrBufferFull)
			}

		case '/':
			//if ()
		}
	}

	fmt.Println("Tokens", lexer.tokens)
}

func (l *Lexer) advance() (rune, error) {
	c, _, err := l.reader.ReadRune()
	l.literal = string(c)
	return c, err
}

// TODO lookahead
func (l *Lexer) lookahead(amount int) ([]byte, error) {
	return l.reader.Peek(amount)
}

func (l *Lexer) addToken(token Token) {
	l.tokens = append(l.tokens, token)
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
