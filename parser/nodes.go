package parser

import (
	"food-interpreter/lexer"
)

type Node interface {
	GetToken() lexer.Token
	GetSubNodes() []Node
}

type Year struct {
	year      lexer.Token
	semicolon Semicolon
}

func (y Year) GetToken() lexer.Token {
	return y.year
}

func (y Year) GetSubNodes() []Node {
	return []Node{y.semicolon}
}

type Semicolon struct {
	semicolon lexer.Token
	time      Time
}

func (s Semicolon) GetToken() lexer.Token {
	return s.semicolon
}

func (s Semicolon) GetSubNodes() []Node {
	return []Node{s.time}
}

type MonthAndDay struct {
	monthAndDay lexer.Token
	time        Time
}

func (m MonthAndDay) GetToken() lexer.Token {
	return m.monthAndDay
}

func (m MonthAndDay) GetSubNodes() []Node {
	return []Node{m.time}
}

type Time struct {
	time  lexer.Token
	right Node // (<food> | <repeater> | <sleep>)
}

func (t Time) GetToken() lexer.Token {
	return t.time
}

func (t Time) GetSubNodes() []Node {
	return []Node{t.right}
}

type Food struct {
	food  lexer.Token
	right Node // (<comma> | <semicolon>)
}

func (f Food) GetToken() lexer.Token {
	return f.food
}

func (f Food) GetSubNodes() []Node {
	return []Node{f.right}
}

type Comma struct {
	comma lexer.Token
	right Node
}

func (c Comma) GetToken() lexer.Token {
	return c.comma
}

func (c Comma) GetSubNodes() []Node {
	return []Node{c.right}
}

type Repeater struct {
	repeater lexer.Token
	right    Node // (<comma> | <semicolon>)
}

func (r Repeater) GetToken() lexer.Token {
	return r.repeater
}

func (r Repeater) GetSubNodes() []Node {
	return []Node{r.right}
}

type Sleep struct {
	sleep lexer.Token
	right Node // (<comma> | <semicolon>)
}

func (s Sleep) GetToken() lexer.Token {
	return s.sleep
}

func (s Sleep) GetSubNodes() []Node {
	return []Node{s.right}
}
