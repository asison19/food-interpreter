package lex

import (
	"bufio"
	"fmt"
	"io"
)

func ScanTokens(reader *bufio.Reader) {
	for {
		c, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(string(c))
	}
}

//func scanToken() {
//
//}
