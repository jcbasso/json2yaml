// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an example of a goyacc program.
// To build it:
// goyacc -p "expr" expr.y (produces y.go)
// go build -o expr y.go
// expr
// > <type an expression>

%{

package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"math/big"
	"unicode/utf8"
)

%}

%union {
	stringVal *string
	boolVal bool
	numVal *big.Rat
	f func(int)
	f_p func()
}

%type <f> A A_1 V VI

%token '[' ']' '{' '}' ',' ':'
%token <stringVal, f_p> STRING
%token <boolVal, f_p> BOOL
%token <numVal, f_p> NUM
%token <f_p> NULL

%%
J:
	O
|	A
	{
		$1(0)
	}

O:
	'{' O_1 '}'
|	'{' '}'

O_1:
	STRING ':' V
|	STRING ':' V ',' O_1
A:
	'[' A_1 ']'
	{
		$$ = func(level int) {
			$2(level)
		}
	}
|	'[' ']' {}

A_1:
	V
	{
		$$ = func(level int) {
			$1(level)
		}
	}
|	V ',' A_1
	{
		$$ = func(level int) {
			$1(level)
			$3(level)
		}
	}


V: 
	A 
	{
		$$ = func(level int) {
			for i = 0; i < level; i++ {
				fmt.Print(" ")
			}
			fmt.Print(" - ")
			$1(level + 1)
		}
	}
|	O {}
|	VI 
	{
		$$ = func(level int) {
				for i = 0; i < level; i++ {
					fmt.Print(" ")
				}
				fmt.Print(" - ")
				$1()
		}
	}

VI:
	STRING
	{
		$$ = func() {
			fmt.Println($1.stringVal)
		}
	}
|	NUM {}
|	BOOL {}
|	NULL {}
%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *exprSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '[',']',',','{','}',':':
			return int(c)
		//Num Value
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.getNumValue(c, yylval)

		//TODO: TODOS ESTOS 3
		//Bool value
		case 't','f':
			yylval.boolVal = 't' == c
			return BOOL //x.getBoolValue(c, yylval)
		//Null value
		case 'n':
			return NULL//,x.getNullValue(c, yylval)
		//String value
		case '"':
			 return x.readStringValue(c, yylval)
		case ' ', '\t', '\n', '\r':
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

func (x *exprLex) getNumValue(c rune, yylval *exprSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.numVal = &big.Rat{}
	_, ok := yylval.numVal.SetString(b.String())
	if !ok {
		log.Printf("bad number %q", b.String())
		return eof
	}
	return NUM
}

func (x *exprLex) readStringValue(c rune, yylval *exprSymType) int {
    var b bytes.Buffer

	b.WriteRune(c)
	c = x.next()
	for c != eof && c != '"' {
		b.WriteRune(c)
		c = x.next()
	}

	if c != '"' {
		log.Printf("String malformed")
		return eof
	}

	b.WriteRune(c)
	res := b.String()
	yylval.stringVal = &res
	return STRING
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
	log.Printf("parse error: %s", s)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		exprParse(&exprLex{line: line})
	}
}
