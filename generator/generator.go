package generator

// TODO look into gota and gonum/plot?
// or create a csv (or just send the straight up diary/nodes? All the interpreter does it ensure proper grammar?) and send that?

import (
	"fmt"
	"food-interpreter/parser"
	"log"
	"strconv"
	"time"
	"unicode/utf8"
)

// TODO make this configurable
var timezone time.Location = *time.Local

// Current working date
type currentDate struct {
	year  int
	month time.Month
	day   int
}

// Nodes - slice of root nodes (YEAR or MONTHANDDAY)
func Generate(nodes []parser.Node) {
	fmt.Printf("Nodes: %+v\n", nodes)
	for _, node := range nodes {

		// TODO switch?
		if yearNode, ok := node.(parser.Year); ok {
			fmt.Println("YEAR node")
			year := handleYear(yearNode)
			fmt.Printf("The working year is %d\n", year)
		}

		// TODO, the rest of the tokens

		// In the case of semicolon, do nothing

		fmt.Println(node)
	}
}

// Given a year node, return the int value of the year.
func handleYear(yearNode parser.Year) int {
	token := yearNode.GetToken()
	number, err := strconv.Atoi(trimFirstRune(token.Lexeme))
	if err != nil {
		log.Fatal("Year token is not an int: " + token.Lexeme)
	}

	return number
}

// Send the result of removing the first element in the passed in string
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	fmt.Println(i)
	return s[i:]
}
