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

type diaryEntry struct {
	date time.Time
	list []string
}

// Nodes - slice of root nodes (YEAR or MONTHANDDAY)
func Generate(nodes []parser.Node) {
	currentDate := currentDate{0, 0, 0}
	diaryEntries := []diaryEntry{}

	for _, node := range nodes {

		// TODO switch?
		if node, ok := node.(parser.Year); ok {
			currentDate.year = handleYear(node)
			fmt.Printf("The current working year is %d\n", currentDate.year)
		}

		if node, ok := node.(parser.MonthAndDay); ok {
			currentDate.month, currentDate.day = handleMonthAndDay(node)
			fmt.Printf("The current working month and day is %d/%d\n", currentDate.month, currentDate.day)

			// Get the times and foods
			subNodes := node.GetSubNodes()
			for _, subNode := range subNodes {
				if timeNode, ok := subNode.(parser.Time); ok {
					hourmin, list := handleTime(timeNode)
					hour := int(hourmin / 100)
					min := int(hourmin % 100)
					entry := diaryEntry{time.Date(currentDate.year, currentDate.month, currentDate.day, hour, min, 0, 0, &timezone), list}
					diaryEntries = append(diaryEntries, entry)
					fmt.Printf("Time is %d\n", hourmin)
					fmt.Printf("The diary entry is %v\n", entry)
				}
			}
		}
		fmt.Println(node)
	}
	fmt.Printf("The diary entries are %v\n", diaryEntries)
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

// Each time should have some food, repeater, or sleep tied to it, or even multiple
// Returns the time and list of items done at that time
func handleTime(node parser.Time) (int, []string) {

	// Get the time
	timeStr := node.GetToken().Lexeme
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Fatal("Invalid time: " + timeStr)
	}

	list := handleSubNodes([]string{}, node.GetSubNodes())
	fmt.Printf("The FRS of time %d, is %s ", time, list)

	return time, list
}

func handleSubNodes(list []string, nodes []parser.Node) []string {
	for _, node := range nodes {
		if node == nil {
			return list
		}

		_, ok1 := node.(parser.Food)
		_, ok2 := node.(parser.Sleep)
		_, ok3 := node.(parser.Repeater)

		if ok1 || ok2 || ok3 {
			list = append(list, node.GetToken().Lexeme)
		}

		return handleSubNodes(list, node.GetSubNodes())
	}
	return list
}
