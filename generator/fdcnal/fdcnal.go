package fdcnal

import (
	"encoding/json"
	"flag"
	"food-interpreter/generator/levenshtein"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
)

var (
	fdcnal_api_key = flag.String("fdcnal", os.Getenv("FDCNAL_API_KEY"), "The USDA Food Data Central API key.")
)

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

// m - set of the foods to search for
// Returns a list of foods from the FDCNAL API that closely match the foods from the given set, m
func GetFoodData(set map[string]struct{}) []FdcnalFood {

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

	// Find the foods that most matches the entries.
	var foods []FdcnalFood
	for k := range set {
		foods = append(foods, findFood(k, fdcnal.Foods))
	}
	return foods
}

// Find the food that most matches
//
// f - Food to match
// foods - The FDCNAL foods to compare with
func findFood(f string, foods []FdcnalFood) FdcnalFood {
	d := math.MaxInt
	r := FdcnalFood{}
	for _, e := range foods {
		ef := strings.SplitN(e.Description, ",", 2)[0] // The food is at the first part, the rest are descriptors. TODO check more foods.
		ld := levenshtein.LevenshteinDistance(strings.ToLower(f), strings.ToLower(ef))
		if ld < d {
			d = ld
			r = e
		}
		// TODO search brand information too
	}
	return r
}

// append the keys in the map to a string for use with a url
func appendString(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return "%22" + strings.Join(keys, "%22%20%22") + "%22"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
