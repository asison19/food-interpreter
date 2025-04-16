package parser

import (
	"bufio"
	"fmt"
	"food-interpreter/lexer"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return Lexer{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return Parse(scanner * bufio.Scanner)
}

func Parse(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()

	}
}

type Parser struct {
	tokens  []lexer.Token
	current int 
}

func newParser(tokens []lexer.Token) Parser {
	return Parser{
		tokens = tokens,
		current = 0,
	}
}

func year() Year {
	match()
	return Year
}

func semicolon() {

}

type YearNt interface {

}

// TODO tree nodes should be structs?
type Year struct {
	Year      lexer.Token
	Semicolon lexer.Token
}

type Semicolon struct {
	Semicolon lexer.Token
}

//type MonthAndDay struct {
//}
