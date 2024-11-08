package lexer

type ReservedWords struct {
	words map[string]int
}

var reservedWords ReservedWords

func GetReservedWords() *map[string]int {
	return &reservedWords.words
}

func (rw *ReservedWords) NewReservedWords() {
	rw.words = map[string]int{
		"sleep": SLEEP,
	}
}
