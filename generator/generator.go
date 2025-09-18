package generator

// TODO look into gota and gonum/plot?
// or create a csv (or just send the straight up diary/nodes? All the interpreter does it ensure proper grammar?) and send that?

import (
	"fmt"
	"food-interpreter/parser"
)

// Nodes - slice of root nodes (YEAR or MONTHANDDAY)
func Generate(nodes []parser.Node) {
	fmt.Printf("Nodes: %+v\n", nodes)
	for _, node := range nodes {
		if _, ok := node.(parser.Year); ok {
			fmt.Println("YEAR node")
		}
		fmt.Println(node)
	}
}

//func handleYear(node []parser.Node) {
//}
