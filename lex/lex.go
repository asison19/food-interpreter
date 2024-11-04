package lex

import (
	"bufio"
	"fmt"
	"io"
)

// TODO current location of the file?
type Lexer struct {
	position    int
	currentChar rune
}

func ScanTokens(reader *bufio.Reader) {
	//count := 0
	for {
		c, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Print(string(c))
		//count += lex(c)
		switch c {
		case '0': // TODO use isNumber() instead?
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
			fmt.Println("number.")
		default:
			fmt.Println("Error")
		}
	}
}

// TODO lookahead
func lookahead() {

}

func isNumber() {

}

func isString() {

}
