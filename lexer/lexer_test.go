package lexer

import (
	"testing"
)

// TODO these should all fail when including tokens that should be included, but are not
func TestLexerYear(t *testing.T) {
	LexFile("../test/data/year.txt")
}

func TestLexerSemicolon(t *testing.T) {
	LexFile("../test/data/semicolon.txt")
}

func TestLexerMonthAndDay(t *testing.T) {
	LexFile("../test/data/monthandday.txt")
}

func TestLexerTime(t *testing.T) {
	LexFile("../test/data/time.txt")
}

func TestLexerFood(t *testing.T) {
	LexFile("../test/data/food.txt")
}

func TestLexerRepeater(t *testing.T) {
	LexFile("../test/data/repeater.txt")
}

func TestLexerComma(t *testing.T) {
	LexFile("../test/data/comma.txt")
}

func TestLexerSleep(t *testing.T) {
	LexFile("../test/data/sleep.txt")
}

// TODO error handling test(s)
func TestLexerError(t *testing.T) {

}
