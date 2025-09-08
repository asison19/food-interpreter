package parser

import (
	"food-interpreter/lexer"
)

type Node interface{}

type Year struct {
	year      lexer.Token
	semicolon Semicolon
}

type Semicolon struct {
	semicolon lexer.Token
	time      Time
}

type MonthAndDay struct {
	monthAndDay lexer.Token
	time        Time
}

type Time struct {
	time  lexer.Token
	right Node // (<food> | <repeater> | <sleep>)
}

type Food struct {
	food  lexer.Token
	right Node // (<comma> | <semicolon>)
}

type Comma struct {
	comma lexer.Token
	food  Food
}

type Repeater struct {
	repeater lexer.Token
	right    Node // (<comma> | <semicolon>)
}

type Sleep struct {
	sleep lexer.Token
	right Node // (<comma> | <semicolon>)
}
