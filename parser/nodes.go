package parser

import (
	"food-interpreter/lexer"
)

type Node interface {
	GetToken() lexer.Token
}

type Year struct {
	year      lexer.Token
	semicolon Semicolon
}

func (y Year) GetToken() lexer.Token {
	return y.year
}

type Semicolon struct {
	semicolon lexer.Token
	time      Time
}

func (s Semicolon) GetToken() lexer.Token {
	return s.semicolon
}

type MonthAndDay struct {
	monthAndDay lexer.Token
	time        Time
}

func (m MonthAndDay) GetToken() lexer.Token {
	return m.monthAndDay
}

type Time struct {
	time  lexer.Token
	right Node // (<food> | <repeater> | <sleep>)
}

func (t Time) GetToken() lexer.Token {
	return t.time
}

type Food struct {
	food  lexer.Token
	right Node // (<comma> | <semicolon>)
}

func (f Food) GetToken() lexer.Token {
	return f.food
}

type Comma struct {
	comma lexer.Token
	food  Food
}

func (c Comma) GetToken() lexer.Token {
	return c.comma
}

type Repeater struct {
	repeater lexer.Token
	right    Node // (<comma> | <semicolon>)
}

func (r Repeater) GetToken() lexer.Token {
	return r.repeater
}

type Sleep struct {
	sleep lexer.Token
	right Node // (<comma> | <semicolon>)
}

func (s Sleep) GetToken() lexer.Token {
	return s.sleep
}
