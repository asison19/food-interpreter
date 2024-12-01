package errorhandler

import (
	"log"
)

func ReportErrorLexer(err string, line int, index int) {
	log.Printf("Lexical error: %s in line: %d, index: %d", err, line, index)
}
