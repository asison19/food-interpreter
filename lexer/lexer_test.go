package lexer

import (
	"reflect"
	"testing"
)

// TODO these should all fail when including tokens that should be included, but are not.
// TODO speed up testing by not having to make a syscall to open a file for every little thing.
func TestLexerYear(t *testing.T) {
	expected := []Token{
		{YEAR, "y2024"},
		{YEAR, "y1234"},
		{YEAR, "y5678"},
		{YEAR, "y9001"},
		{YEAR, "y1"},
		{YEAR, "y2"},
		{YEAR, "y3"},
		{YEAR, "y4"},
		{YEAR, "y5"},
		{YEAR, "y6"},
		{YEAR, "y7"},
		{YEAR, "y8"},
		{YEAR, "y9"},
		{YEAR, "y0"},
	}
	actual := LexFile("test/data/year.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerSemicolon(t *testing.T) {
	expected := []Token{}
	for range 6 {
		expected = append(expected, Token{SEMICOLON, ";"})
	}
	actual := LexFile("test/data/semicolon.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerMonthAndDay(t *testing.T) {
	expected := []Token{
		{MONTHANDDAY, "12/12"},
		{MONTHANDDAY, "1/2"},
		{MONTHANDDAY, "3/4"},
		{MONTHANDDAY, "5/6"},
		{MONTHANDDAY, "7/8"},
		{MONTHANDDAY, "9/10"},
	}
	actual := LexFile("test/data/monthandday.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerTime(t *testing.T) {
	expected := []Token{
		{TIME, "0000"},
		{TIME, "1234"},
		{TIME, "2359"},
	}
	actual := LexFile("test/data/time.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerFood(t *testing.T) {
	expected := []Token{
		{FOOD, "apple"},
		{FOOD, "banana"},
		{FOOD, "cherry"},
		{FOOD, "durian"},
		{FOOD, "elderberry"},
	}
	actual := LexFile("test/data/food.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerRepeater(t *testing.T) {
	expected := []Token{}
	for range 6 {
		expected = append(expected, Token{REPEATER, ".."})
	}

	actual := LexFile("test/data/repeater.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerComma(t *testing.T) {
	expected := []Token{}
	for range 6 {
		expected = append(expected, Token{COMMA, ","})
	}
	actual := LexFile("test/data/comma.txt")
	assertEqual(t, expected, actual.tokens)
}

func TestLexerSleep(t *testing.T) {
	expected := []Token{
		{SLEEP, "sleep"},
		{SLEEP, "Sleep"},
		{SLEEP, "SLEEP"},
		{SLEEP, "SlEeP"},
		{SLEEP, "sLeEp"},
	}
	actual := LexFile("test/data/sleep.txt")
	assertEqual(t, expected, actual.tokens)
}

// TODO error handling test(s)
func TestLexerError(t *testing.T) {

}

func assertEqual(t *testing.T, a []Token, b []Token) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nExpected lexer result: \n\t%+v\nActual lexer result:\n\t%+v", a, b)
	}
}
