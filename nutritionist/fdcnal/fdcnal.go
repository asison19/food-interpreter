package fdcnal

import (
	"encoding/json"
	"flag"
	"fmt"
	"food-interpreter/levenshtein"
	"io"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
)

var (
	fdcnal_api_key = flag.String("fdcnal", os.Getenv("FDCNAL_API_KEY"), "The USDA Food Data Central API key.")
	calorie_ids    = []int{2048, 2047, 1008}
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
	NutrientId   int     `json:"nutrientId"`
	NutrientName string  `json:"nutrientName"`
	Value        float64 `json:"value"`
}

type FdcnalFoodHashed struct {
	FdcId         int
	Description   string
	ServingSize   float64
	FoodNutrients map[int]Nutrient
}

type Nutrient struct {
	NutrientName string
	Value        float64
}

// m - map of the foods to search for and add FdcnalFood to
// Modifies m, adding matching FdcnalFood as a value to the key.
// Returns a list of foods from the FDCNAL API that closely match the foods from the given set, m
func SetFoodData(m map[string]FdcnalFoodHashed) []FdcnalFoodHashed {

	// Get the nutritional data
	flag.Parse()
	query := appendString(m)

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
	var foods []FdcnalFoodHashed
	for k := range m {
		matchingFood := findFood(k, fdcnal.Foods)
		f := FdcnalFoodHashed{
			matchingFood.FdcId,
			matchingFood.Description,
			matchingFood.ServingSize,
			foodNutrientsToHashMap(matchingFood.FoodNutrients),
		}

		foods = append(foods, f)
		m[k] = f
	}
	return foods
}

// Get the total nutrient info by id of a given set of foods.
// Some nutrients, like calories (it's a nutrient in FDCNAL), have different versions, such as
//
//	1008, 2047, 2048,
//
// foods - FDCNAL foods whose nutrients to total.
// id    - FDCNAL ID of the nutrient.
func GetTotalNutrientInfo(foods []FdcnalFoodHashed, id int) float64 {
	if slices.Contains(calorie_ids, id) {
		return GetTotalCalories(foods)
	}

	total := 0.0
	for _, f := range foods {
		total += f.FoodNutrients[id].Value
	}
	return total
}

// Get the total amount of calories in the given slice of foods.
func GetTotalCalories(foods []FdcnalFoodHashed) float64 { // TODO this is bugged when calling from interpreter with any calorie ID
	total := 0.0
	for _, f := range foods {
		v, ok := findIdValue(f, calorie_ids)
		if ok {
			total += v
		} else {
			fmt.Printf("Calories not found for food %v.\n", f)
		}
	}
	return total
}

// Find the value of the passed in FdcnalFoodHashed given a set of IDs, returning the first found.
func findIdValue(f FdcnalFoodHashed, ids []int) (float64, bool) {
	for _, id := range ids {
		n, ok := f.FoodNutrients[id]
		if ok {
			return n.Value, true
		}
	}
	return 0, false
}

func foodNutrientsToHashMap(fn []FdcnalFoodNutrients) map[int]Nutrient {
	m := make(map[int]Nutrient)
	for _, n := range fn {
		m[n.NutrientId] = Nutrient{n.NutrientName, n.Value}
	}
	return m
}

// Find the food that most matches
//
// f     - Food to match
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
func appendString(m map[string]FdcnalFoodHashed) string {
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
