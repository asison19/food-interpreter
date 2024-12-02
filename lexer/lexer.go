package lexer

import (
	"bufio"
	"fmt"
	"food-interpreter/errorhandler"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO current location of the file?
type Lexer struct {
	position     int
	nextPosition int
	linePosition int
	lexeme       string
	line         string
	tokens       []Token
}

// TODO error handling
func LexFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ScanTokens(scanner)
}

func ScanTokens(scanner *bufio.Scanner) {
	lexer := Lexer{0, 0, 0, "", "", []Token{}}

	for scanner.Scan() {
		line := scanner.Text()
		lexer.linePosition += 1

		fmt.Println("Line:", string(line))
		scanLine(&lexer, string(line))
	}

}

func scanLine(lexer *Lexer, line string) {
	lexer.line = line

	// TODO error handling
	// TODO do I need this for loop?

L:
	for {
		c := lexer.advance()

		if c == "" {
			break
		}

		switch {
		case isNumber(c): // MONTHANDDAY or TIME
			if !lexer.addToken(lexer.number()) { // TODO is there a better way of doing this conditional?
				break
			}
			continue L
		}

		switch c {
		// Ignoreables
		case " ":
		case "\n":
		case "\r":
			lexer.clearLexeme()
			continue L
		case "y": // YEAR
			if !lexer.addToken(lexer.year()) {
				break
			}
			continue L
		case ";": // SEMICOLON
			if !lexer.addToken(Token{
				tokenType: SEMICOLON,
				lexeme:    lexer.lexeme,
			}) {
				break
			}
			continue L
		case ",": // COMMA
			if !lexer.addToken(Token{
				tokenType: COMMA,
				lexeme:    lexer.lexeme,
			}) {
				break
			}
			continue L
		// TODO
		case "(":
			continue L
		case ".":
			if !lexer.addToken(lexer.reapeater()) {
				break
			}
			continue L
		}

		// Food, variables, sleep, etc.
		if isAlpha(c) {
			if !lexer.addToken(lexer.identifier()) {
				break
			}
		}
		// TODO end of line reporting incorrectly
		lexer.reportError(fmt.Sprintf("Invalid token. lexeme: %s", lexer.lexeme))
		break
	}
	lexer.clearLine()
	fmt.Println("lexer.tokens", lexer.tokens)
}

func (l *Lexer) advance() string {
	defer func() { l.position += 1 }()
	if l.position >= len(l.line) {
		return ""
	}
	c := string(l.line[l.position])
	l.lexeme += c
	return c
}

func (l *Lexer) clearLexeme() {
	l.lexeme = ""
}

func (l *Lexer) clearLine() {
	l.line = "" // TODO probably don't need, but leave just in case?
	l.position = 0
}

func (l *Lexer) lookahead(amount int) string {
	if l.position+amount > len(l.line) {
		return "\n" // TODO return empty string?
	}
	return l.line[l.position : l.position+amount]
}

// TODO if token is empty, it's invalid and the entirety of the line should be skipped
func (l *Lexer) addToken(token Token) bool {
	if token == (Token{}) {
		return false
	}
	l.tokens = append(l.tokens, token)
	l.lexeme = ""
	return true
}

func (l *Lexer) scanningError(token Token) {
	fmt.Println("Error at line ", l.position)
}

func (l *Lexer) reapeater() Token {
	ahead := l.lookahead(1)
	if ahead == "." {
		l.advance()
		return Token{
			tokenType: REPEATER,
			lexeme:    l.lexeme,
		}
	}
	l.reportError(fmt.Sprintf("Unexpected character %s when lexing token REPEATER, lexeme: ", l.lexeme))
	return Token{}
}

func (l *Lexer) year() Token {
	ahead := l.lookahead(1)

	for isNumber(ahead) {
		l.advance()
		ahead = l.lookahead(1)
	}
	return Token{
		tokenType: YEAR,
		lexeme:    l.lexeme,
	}
}

func (l *Lexer) identifier() Token {
	ahead := l.lookahead(1)

	for isAlphaNumeric(ahead) {
		l.advance()
		ahead = l.lookahead(1)
	}

	// To simplify end user experience, reserved words are kept case insensitive
	tokenType, ok := reservedWords[strings.ToLower(l.lexeme)]

	// FOOD
	if !ok {
		return Token{
			tokenType: FOOD,
			lexeme:    l.lexeme,
		}
	}

	return Token{
		tokenType: tokenType,
		lexeme:    l.lexeme,
	}
}

// TODO certain dates should not be possible
// E.g. 0/0, 1/32, etc.
// Certain times are not possible:
//   0000 - 2359 only
func (l *Lexer) number() Token {
	ahead := l.lookahead(1)

	d := 1
	for isNumber(ahead) {
		d += 1
		l.advance()
		ahead = l.lookahead(1)
	}

	if d > 4 {
		l.reportError(fmt.Sprintf("Invalid token. Invalid digit count of %d found. lexeme: %s", d, l.lexeme))
		return Token{}
	}

	// MONTHANDDAY
	if strings.Contains(ahead, "/") {
		// month
		month, _ := strconv.Atoi(l.lexeme)
		if month > 12 || month < 0 {
			l.reportError(fmt.Sprintf("Invalid token, %d. Month must be between 1 and 12. lexeme: %s", month, l.lexeme))
			return Token{}
		}

		// day
		l.advance()
		day := l.lookahead(1)
		if day == "\n" {
			fmt.Println("End of file reading number", bufio.ErrBufferFull)
		}
		if !isNumber(day) {
			l.reportError(fmt.Sprintf("Invalid token, %s. Day must be a number between 1 and 12. lexeme: %s", day, l.lexeme))
			return Token{}
		}
		for isNumber(day) {
			l.advance()
			day = l.lookahead(1)
		}

		d, _ := strconv.Atoi(day)

		if d > 31 {
			l.reportError(fmt.Sprintf("Invalid day over 32, found %s while lexing day portion of token, MONTHANDDAY.", l.lexeme))
			return Token{}
		}

		return Token{
			tokenType: MONTHANDDAY,
			//lexeme:
			lexeme: l.lexeme,
		}
	}

	// TIME
	time, err := strconv.Atoi(l.lexeme)
	if err != nil {
		l.reportError(fmt.Sprintf("Invalid token, %d. Found %s when lexing token TIME, lexeme: ", time, l.lexeme))
		return Token{}
	}

	if time < 0 || time > 2359 {
		l.reportError(fmt.Sprintf("Invalid token, %d. Time must be between 0 and 2359. lexeme: %s ", time, l.lexeme))
		return Token{}
	}

	if minutes := time % 100; minutes < 0 || minutes > 59 {
		l.reportError(fmt.Sprintf("Invalid token %d. Minutes of time must be between 0 and 59 lexeme: %s ", time, l.lexeme))
		return Token{}
	}

	return Token{
		tokenType: TIME,
		//lexeme:
		lexeme: l.lexeme,
	}
}

// TODO invalid line is sent
func (l *Lexer) reportError(err string) {
	// Any error we get, just skip the rest of the line
	errorhandler.ReportErrorLexer(err, l.linePosition, l.position)
	l.clearLexeme()
}

func isNumber(c string) bool {
	if _, err := strconv.Atoi(c); err == nil {
		return true
	}
	return false
}

func isAlpha(c string) bool {
	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(c)
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isNumber(c)
}
