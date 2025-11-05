package nutritionist

import (
	"fmt"
	"food-interpreter/generator"
	"food-interpreter/nutritionist/fdcnal"
	"time"
)

var (
	timezone time.Location = *time.Local // TODO make this configurable
)

// m - Map of date times and their entries of food, sleep, or repeater.
func AddFoodData(m map[time.Time][]generator.Entry) {

	// Map for the foods from the diary and it's most matching food from the FDCNAL DB.
	foods := make(map[string]fdcnal.FdcnalFoodHashed)
	for _, v := range m {
		for _, e := range v {
			if f, ok := e.(generator.FoodEntry); ok {
				foods[f.Name] = fdcnal.FdcnalFoodHashed{}
			}
		}
	}

	fdcnal.SetFoodData(foods)

	arr := []fdcnal.FdcnalFoodHashed{}
	for _, v := range foods {
		arr = append(arr, v)
	}

	total := fdcnal.GetTotalNutrientInfo(arr, 2048)
	fmt.Println(total)
}
