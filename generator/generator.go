package generator

// TODO look into gota and gonum/plot?
// or create a csv (or just send the straight up diary/nodes? All the interpreter does it ensure proper grammar?) and send that?

import (
	"fmt"
	"food-interpreter/generator/fdcnal"
	"food-interpreter/parser"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	timezone time.Location = *time.Local // TODO make this configurable
)

// Current working date
type currentDate struct {
	year  int
	month time.Month
	day   int
}

type entry interface {
}

// TODO the rest
type foodEntry struct {
	name     string
	calories int // TODO The rest of the pertinent nutritional data and then comments in the future
}

type sleepEntry struct {
	// TODO sleep hygiene?
	details string
}

type repeaterEntry struct {
	// TODO this entry will potentially have a mix of stuff. Should deal with it before getting to here?
	details string
}

// m - Map of date times and their entries of food, sleep, or repeater.
func addFoodData(m map[time.Time][]entry) {

	// Map for the foods from the diary and it's most matching food from the FDCNAL DB.
	foods := make(map[string]fdcnal.FdcnalFood)
	for _, v := range m {
		for _, e := range v {
			if f, ok := e.(foodEntry); ok {
				foods[f.name] = fdcnal.FdcnalFood{}
			}
		}
	}

	fdcnal.SetFoodData(foods)

	arr := []fdcnal.FdcnalFood{}
	for _, v := range foods {
		arr = append(arr, v)
	}

	total := fdcnal.GetTotalNutrientInfo(arr, 1089)
	fmt.Println(total)

}

// Nodes - slice of root nodes (YEAR or MONTHANDDAY)
func Generate(nodes []parser.Node) map[time.Time][]entry {
	currentDate := currentDate{0, 0, 0}

	// A map of date times and their entries of food, sleep, or repeater.
	m := make(map[time.Time][]entry)
	for _, node := range nodes {

		// Get the year
		if node, ok := node.(parser.Year); ok {
			currentDate.year = handleYear(node)
		}

		// Get the month and day along with the times and foods
		if node, ok := node.(parser.MonthAndDay); ok {
			currentDate.month, currentDate.day = handleMonthAndDay(node)

			// Get the times and foods
			subNodes := node.GetSubNodes()
			for _, subNode := range subNodes {
				if timeNode, ok := subNode.(parser.Time); ok {
					hourmin, list := handleTime(timeNode)
					hour := int(hourmin / 100)
					min := int(hourmin % 100)
					time := time.Date(currentDate.year, currentDate.month, currentDate.day, hour, min, 0, 0, &timezone)
					//diaryEntries = append(diaryEntries, entry)
					m[time] = list
				}
			}
		}
	}
	addFoodData(m)
	return m
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
// Returns the time and list of entries done at that time
func handleTime(node parser.Time) (int, []entry) {

	// Get the time
	timeStr := node.GetToken().Lexeme
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Fatal("Invalid time: " + timeStr)
	}

	list := handleSubNodes([]entry{}, node.GetSubNodes())

	return time, list
}

func handleSubNodes(list []entry, nodes []parser.Node) []entry {
	for _, node := range nodes {
		if node == nil {
			return list
		}

		// TODO at this point do the lexeme parsing for user inputted information?
		if n, ok := node.(parser.Food); ok {
			entry := foodEntry{n.GetToken().Lexeme, 0}
			list = append(list, entry)
		} else if n, ok := node.(parser.Sleep); ok {
			entry := sleepEntry{n.GetToken().Lexeme}
			list = append(list, entry)
		} else if n, ok := node.(parser.Repeater); ok {
			entry := repeaterEntry{n.GetToken().Lexeme}
			list = append(list, entry)
		}

		return handleSubNodes(list, node.GetSubNodes())
	}
	return list
}
