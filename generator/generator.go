package generator

// TODO look into gota and gonum/plot?
// or create a csv (or just send the straight up diary/nodes? All the interpreter does it ensure proper grammar?) and send that?

import (
	"fmt"
	"food-interpreter/parser"
	"log"
	"strconv"
	"strings"
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
	currentDate := currentDate{0, 0, 0}

	for _, node := range nodes {

		// TODO switch?
		if node, ok := node.(parser.Year); ok {
			currentDate.year = handleYear(node)
			fmt.Printf("The current working year is %d\n", currentDate.year)
		}

		// TODO, the rest of the tokens
		if node, ok := node.(parser.MonthAndDay); ok {
			currentDate.month, currentDate.day = handleMonthAndDay(node)
			fmt.Printf("The current working month and day is %d/%d\n", currentDate.month, currentDate.day)
		}

		// In the case of semicolon, do nothing

		fmt.Println(node)
	}
}

// Given a year node, return the int value of the year.
func handleYear(node parser.Year) int {
	token := node.GetToken()
	number, err := strconv.Atoi(trimFirstRune(token.Lexeme))
	if err != nil {
		log.Fatal("Invalid Year token: " + token.Lexeme)
	}

	return number
}

// Send the result of removing the first element in the passed in string
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	fmt.Println(i)
	return s[i:]
}

func handleMonthAndDay(node parser.MonthAndDay) (time.Month, int) {
	mad := node.GetToken().Lexeme
	madSli := strings.Split(mad, "/")

	if len(madSli) != 2 {
		log.Fatal("Invalid Month and Day token: " + mad)
	}

	m, err := strconv.Atoi(madSli[0])
	if err != nil {
		log.Fatal("Invalid Month token: " + mad)
	}

	d, err := strconv.Atoi(madSli[1])
	if err != nil {
		log.Fatal("Invalid Day token: " + mad)
	}

	return time.Month(m), d

}
