package nutritionist

import (
	"food-interpreter/generator"
	"food-interpreter/nutritionist/fdcnal"
	"log"
	"slices"
	"time"
)

var (
	timezone    time.Location = *time.Local // TODO make this configurable
	calorie_ids               = []int{2048, 2047, 1008}
)

type FoodEntry struct {
	Food    fdcnal.FdcnalFoodHashed
	Details string // TODO strings of details about the food such as amount, brand, etc.
}

// m - Map of date times and their entries of food, sleep, or repeater.
func AddFoodData(m map[time.Time][]generator.Entry) map[time.Time]FoodEntry {

	// Map for the foods from the diary and it's most matching food from the FDCNAL DB.
	fdcnalFoods := make(map[string]fdcnal.FdcnalFoodHashed)
	for _, v := range m {
		for _, e := range v {
			if f, ok := e.(generator.FoodEntry); ok {
				fdcnalFoods[f.Name] = fdcnal.FdcnalFoodHashed{}
			}
		}
	}

	fdcnal.SetFoodData(fdcnalFoods)

	foods := make(map[time.Time]FoodEntry)
	for time, entries := range m {
		for _, entry := range entries {
			foods[time] = FoodEntry{Food: fdcnalFoods[entry.GetName()], Details: entry.GetDetails()}
		}
	}

	return foods

	//total := fdcnal.GetTotalNutrientInfo(arr, 2048)
	//fmt.Println(total)
}

// foods - Set of foods from the FDCNAL API whose nutrients have been hashed.
// id    - ID of the type of nutrient to gather
func GetNutrition(m map[time.Time]FoodEntry, id int) float64 {
	if slices.Contains(calorie_ids, id) {
		return GetCalories(m)
	}

	nutritionAmt := 0.0
	for _, entry := range m {
		nutritionAmt += entry.Food.FoodNutrients[id].Value
	}
	return nutritionAmt
}

// start - The beginning of the period to get the nutritional data from.
// end   - The ending of the period to get the nutritional data from.
// foods - Set of foods from the FDCNAL API whose nutrients have been hashed.
// id    - ID of the type of nutrient to gather
func GetDateNutrition(start, end time.Time, m map[time.Time]FoodEntry, id int) float64 {
	nutritionAmt := 0.0
	for time, entry := range m {
		if inTimeSpan(start, end, time) {
			nutritionAmt += entry.Food.FoodNutrients[id].Value
		}
	}
	return nutritionAmt
}

func GetCalories(m map[time.Time]FoodEntry) float64 {
	amt := 0.0
	for _, f := range m {
		v, ok := findIdValue(f, calorie_ids)
		if ok {
			amt += v
		} else {
			log.Printf("Calories not found for food %v.\n", f)
		}
	}
	return amt
}

// Check that the given time is between the given start and end
func inTimeSpan(start, end, time time.Time) bool {
	return time.After(start) && time.Before(end)
}

// Find the value of the passed in FdcnalFoodHashed given a set of IDs, returning the first found.
func findIdValue(f FoodEntry, ids []int) (float64, bool) {
	for _, id := range ids {
		n, ok := f.Food.FoodNutrients[id]
		if ok {
			return n.Value, true
		}
	}
	return 0, false
}
