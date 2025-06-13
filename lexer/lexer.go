package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO foods should be able to have a space in them.

type Lexer struct {
	position     int
	nextPosition int
	linePosition int
	lexeme       string
	line         string
	Tokens       []Token
}

func LexFile(filePath string) Lexer {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return Lexer{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return ScanTokens(scanner)
}

func LexString(s string) Lexer {
	scanner := bufio.NewScanner(strings.NewReader(s))
	return ScanTokens(scanner)
}

func ScanTokens(scanner *bufio.Scanner) Lexer {
	lexer := Lexer{0, 0, 0, "", "", []Token{}}

	for scanner.Scan() {
		line := scanner.Text()
		lexer.linePosition += 1

		scanLine(&lexer, string(line))
	}

	return lexer
}

func scanLine(lexer *Lexer, line string) {
	lexer.line = line

	// TODO error handling
L:
	for {
		c := lexer.advance()

		if c == "" {
			break
		}

		switch {
		case isNumber(c): // MONTHANDDAY or TIME
			lexer.addToken(lexer.number())
			continue L
		}

		switch c {
		case " ", "\n", "\r": // ignorables
			lexer.clearLexeme()
			continue L
		case "y": // YEAR
			lexer.addToken(lexer.year())
			continue L
		case ";": // SEMICOLON
			lexer.addToken(Token{
				Type:   SEMICOLON,
				Lexeme: lexer.lexeme,
			})
			continue L
		case ",": // COMMA
			lexer.addToken(Token{
				Type:   COMMA,
				Lexeme: lexer.lexeme,
			})
			continue L
		// TODO
		case "(":
			continue L
		case ".":
			lexer.addToken(lexer.reapeater())
			continue L

			// TODO
			// case "/":
		}

		// Food, variables, sleep, etc.
		if isAlpha(c) {
			lexer.addToken(lexer.identifier())
			continue L
		}

		lexer.reportError(fmt.Sprintf("Invalid token. lexeme: %s", lexer.lexeme))
		break
	}
	lexer.clearLine()
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
	l.line = ""
	l.position = 0
}

func (l *Lexer) lookahead(amount int) string {
	if l.position+amount > len(l.line) {
		return "\n" // TODO return empty string?
	}
	return l.line[l.position : l.position+amount]
}

func (l *Lexer) addToken(token Token) {
	l.Tokens = append(l.Tokens, token)
	l.lexeme = ""
}

func (l *Lexer) reapeater() Token {
	ahead := l.lookahead(1)
	if ahead == "." {
		l.advance()
		return Token{
			Type:   REPEATER,
			Lexeme: l.lexeme,
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
		Type:   YEAR,
		Lexeme: l.lexeme,
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
			Type:   FOOD,
			Lexeme: l.lexeme,
		}
	}

	return Token{
		Type:   tokenType,
		Lexeme: l.lexeme,
	}
}

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
		if month > 12 || month < 1 {
			l.reportError(fmt.Sprintf("Invalid token, %d. Month must be between 1 and 12. lexeme: %s", month, l.lexeme))
			return Token{}
		}

		// day
		l.advance()
		day := l.lookahead(1)
		if day == "\n" {
			log.Println("End of file reading number: ", bufio.ErrBufferFull)
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
			Type: MONTHANDDAY,
			//lexeme:
			Lexeme: l.lexeme,
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
		Type: TIME,
		//lexeme:
		Lexeme: l.lexeme,
	}
}

// TODO this is going to skip some valid tokens that don't need to seperated by spaces such as commas
func (l *Lexer) reportError(err string) {
	log.Printf("Lexical error: %s in line: %d, index: %d", err, l.linePosition, l.position)

	// Advance to the next intentional token
	for {
		c := l.advance()
		if c == "" {
			break
		}
		switch c {
		case " ", "\n", "\r": // TODO make a list for these
			l.clearLexeme()
			return
		}
	}
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
