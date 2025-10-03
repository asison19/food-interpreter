package interpreter

import (
	"food-interpreter/generator"
	"food-interpreter/lexer"
	"food-interpreter/parser"
	"os"
	"strings"
	"time"
)

// TODO return the parser for now
func Interpret(diary string) parser.Parser {

	l := lexer.LexString(diary)
	p, nodes := parser.ParseTokens(l.Tokens)

	entries := generator.Generate(nodes)
	writeToFile(entries)

	return p // TODO return diary entries
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO work on pandas/visualizer next to determine what the output should be

// local file output for testing
func writeToFile(entries []generator.DiaryEntry) {
	now := time.Now()
	f, err := os.Create("diaryc-" + now.Format("20060102-150405"))
	check(err)

	f.WriteString("datetime,food\n")
	for _, e := range entries {
		f.WriteString(e.Date.Format(time.RFC3339) + ",\"['" + strings.Join(e.List, "','") + "']\"\n")
	}

	defer f.Close()
}
