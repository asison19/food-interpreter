package generator

// TODO look into gota and gonum/plot?
// or create a csv (or just send the straight up diary/nodes? All the interpreter does it ensure proper grammar?) and send that?

import (
	"encoding/json"
	"flag"
	"fmt"
	"food-interpreter/parser"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	timezone       time.Location = *time.Local // TODO make this configurable
	fdcnal_api_key               = flag.String("fdcnal", os.Getenv("FDCNAL_API_KEY"), "The USDA Food Data Central API key.")
)

// Current working date
type currentDate struct {
	year  int
	month time.Month
	day   int
}

//type DiaryEntry struct {
//	Date time.Time
//	List []string
//}

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

// FDCNAL API JSON structs
type Fdcnal struct {
	Foods []FdcnalFood `json:"foods"`
}

type FdcnalFood struct {
	FdcId         int                   `json:"fdcId"`
	Description   string                `json:"description"`
	ServingSize   float64               `json:"servingSize"`
	FoodNutrients []FdcnalFoodNutrients `json:"foodNutrients"`
}

type FdcnalFoodNutrients struct {
	NutrientId   int    `json:"nutrientId"`
	NutrientName string `json:"nutrientName"`
}

// Nodes - slice of root nodes (YEAR or MONTHANDDAY)
func Generate(nodes []parser.Node) map[time.Time][]entry {
	currentDate := currentDate{0, 0, 0}

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

// m - hashmap of the times and entries
func addFoodData(m map[time.Time][]entry) {

	// Get the food entries
	set := make(map[string]struct{})
	for _, v := range m {
		for _, e := range v {
			if f, ok := e.(foodEntry); ok {
				set[f.name] = struct{}{}
			}
		}
	}

	// Get the nutritional data
	flag.Parse()
	query := appendString(set)

	// TODO best dataType?
	// TODO pagesize?
	//dataType := "Branded,Foundation,SR%20Legacy"
	dataType := "Foundation"
	url := "https://api.nal.usda.gov/fdc/v1/foods/search?query=" + query + "&dataType=" + dataType + "&pageSize=250&pageNumber=1&sortBy=dataType.keyword&sortOrder=asc&api_key=" + *fdcnal_api_key

	req, e := http.NewRequest("GET", url, nil)
	check(e)
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, e := client.Do(req)
	check(e)
	defer resp.Body.Close()

	body, e := io.ReadAll(resp.Body)
	check(e)

	var fdcnal Fdcnal
	json.Unmarshal(body, &fdcnal)
	//fmt.Println(fdcnal.Foods)

	for _, f := range fdcnal.Foods {
		fmt.Println(f.Description)
	}

	// Find the food that most matches
	var foods []FdcnalFood
	for k, _ := range set {
		foods = append(foods, findFood(k, fdcnal))
	}
	fmt.Println(foods)
}

// Find the food that most matches
//
// f - Food to match
// fdcnal - The FDCNAL API call result
func findFood(f string, fdcnal Fdcnal) FdcnalFood {
	// TODO fuzzy matching
	d := math.MaxInt
	r := FdcnalFood{}
	for _, e := range fdcnal.Foods {
		ld := levenshteinDistance(f, e.Description)
		if ld < d {
			d = ld
			r = e
		}
		// TODO search brand information too
	}
	return r
}

// TODO add some tests
func levenshteinDistance(a string, b string) int {
	n := len(a)
	m := len(b)

	p := make([]int, m+1)
	c := make([]int, m+1)

	for j := range m {
		p[j] = j
	}
	for i := range n {
		c[0] = i + 1

		for j := range m {
			sc := 0
			if a[i] != b[j] {
				sc = 1
			}
			c[j+1] = min(
				p[j+1]+1, // deletion
				c[j]+1,   // insertion
				p[j]+sc,  // substitution
			)
		}
		copy(p, c)
	}
	return p[m]
}

// append the keys in the map to a string for use with a url
func appendString(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return "%22" + strings.Join(keys, "%22%20%22") + "%22"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
