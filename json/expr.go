//line expr.y:13
package main

import __yyfmt__ "fmt"

//line expr.y:14
import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/big"
	"os"
	"unicode/utf8"
)

//line expr.y:28
type exprSymType struct {
	yys       int
	stringVal *string
	boolVal   bool
	numVal    *big.Rat
	f         func(int)
	f_p       func()
}

const STRING = 57346
const BOOL = 57347
const NUM = 57348
const NULL = 57349

var exprToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'['",
	"']'",
	"'{'",
	"'}'",
	"','",
	"':'",
	"STRING",
	"BOOL",
	"NUM",
	"NULL",
}
var exprStatenames = [...]string{}

const exprEofCode = 1
const exprErrCode = 2
const exprInitialStackSize = 16

//line expr.y:117

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
		case '[', ']', ',', '{', '}', ':':
			return int(c)
		//Num Value
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.getNumValue(c, yylval)

		//TODO: TODOS ESTOS 3
		//Bool value
		case 't', 'f':
			yylval.boolVal = 't' == c
			return BOOL //x.getBoolValue(c, yylval)
		//Null value
		case 'n':
			return NULL //,x.getNullValue(c, yylval)
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
L:
	for {
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

//line yacctab:1
var exprExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const exprPrivate = 57344

const exprLast = 35

var exprAct = [...]int{

	6, 9, 11, 5, 10, 4, 8, 20, 25, 15,
	17, 16, 18, 5, 7, 4, 22, 8, 19, 15,
	17, 16, 18, 23, 24, 5, 26, 4, 21, 13,
	2, 12, 3, 1, 14,
}
var exprPact = [...]int{

	21, -1000, -1000, -1000, 7, -1, 11, -1000, -2, 23,
	-1000, 8, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	9, -1000, 9, 0, -1000, -4, -1000,
}
var exprPgo = [...]int{

	0, 31, 1, 2, 34, 33, 29, 0,
}
var exprR1 = [...]int{

	0, 5, 5, 6, 6, 7, 7, 1, 1, 2,
	2, 3, 3, 3, 4, 4, 4, 4,
}
var exprR2 = [...]int{

	0, 1, 1, 3, 2, 3, 5, 3, 2, 1,
	3, 1, 1, 1, 1, 1, 1, 1,
}
var exprChk = [...]int{

	-1000, -5, -6, -1, 6, 4, -7, 7, 10, -2,
	5, -3, -1, -6, -4, 10, 12, 11, 13, 7,
	9, 5, 8, -3, -2, 8, -7,
}
var exprDef = [...]int{

	0, -2, 1, 2, 0, 0, 0, 4, 0, 0,
	8, 9, 11, 12, 13, 14, 15, 16, 17, 3,
	0, 7, 0, 5, 10, 0, 6,
}
var exprTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 8, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 9, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 4, 3, 5, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 6, 3, 7,
}
var exprTok2 = [...]int{

	2, 3, 10, 11, 12, 13,
}
var exprTok3 = [...]int{
	0,
}

var exprErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	exprDebug        = 0
	exprErrorVerbose = false
)

type exprLexer interface {
	Lex(lval *exprSymType) int
	Error(s string)
}

type exprParser interface {
	Parse(exprLexer) int
	Lookahead() int
}

type exprParserImpl struct {
	lval  exprSymType
	stack [exprInitialStackSize]exprSymType
	char  int
}

func (p *exprParserImpl) Lookahead() int {
	return p.char
}

func exprNewParser() exprParser {
	return &exprParserImpl{}
}

const exprFlag = -1000

func exprTokname(c int) string {
	if c >= 1 && c-1 < len(exprToknames) {
		if exprToknames[c-1] != "" {
			return exprToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func exprStatname(s int) string {
	if s >= 0 && s < len(exprStatenames) {
		if exprStatenames[s] != "" {
			return exprStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func exprErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !exprErrorVerbose {
		return "syntax error"
	}

	for _, e := range exprErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + exprTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := exprPact[state]
	for tok := TOKSTART; tok-1 < len(exprToknames); tok++ {
		if n := base + tok; n >= 0 && n < exprLast && exprChk[exprAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if exprDef[state] == -2 {
		i := 0
		for exprExca[i] != -1 || exprExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; exprExca[i] >= 0; i += 2 {
			tok := exprExca[i]
			if tok < TOKSTART || exprExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if exprExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += exprTokname(tok)
	}
	return res
}

func exprlex1(lex exprLexer, lval *exprSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = exprTok1[0]
		goto out
	}
	if char < len(exprTok1) {
		token = exprTok1[char]
		goto out
	}
	if char >= exprPrivate {
		if char < exprPrivate+len(exprTok2) {
			token = exprTok2[char-exprPrivate]
			goto out
		}
	}
	for i := 0; i < len(exprTok3); i += 2 {
		token = exprTok3[i+0]
		if token == char {
			token = exprTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = exprTok2[1] /* unknown char */
	}
	if exprDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", exprTokname(token), uint(char))
	}
	return char, token
}

func exprParse(exprlex exprLexer) int {
	return exprNewParser().Parse(exprlex)
}

func (exprrcvr *exprParserImpl) Parse(exprlex exprLexer) int {
	var exprn int
	var exprVAL exprSymType
	var exprDollar []exprSymType
	_ = exprDollar // silence set and not used
	exprS := exprrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	exprstate := 0
	exprrcvr.char = -1
	exprtoken := -1 // exprrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		exprstate = -1
		exprrcvr.char = -1
		exprtoken = -1
	}()
	exprp := -1
	goto exprstack

ret0:
	return 0

ret1:
	return 1

exprstack:
	/* put a state and value onto the stack */
	if exprDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", exprTokname(exprtoken), exprStatname(exprstate))
	}

	exprp++
	if exprp >= len(exprS) {
		nyys := make([]exprSymType, len(exprS)*2)
		copy(nyys, exprS)
		exprS = nyys
	}
	exprS[exprp] = exprVAL
	exprS[exprp].yys = exprstate

exprnewstate:
	exprn = exprPact[exprstate]
	if exprn <= exprFlag {
		goto exprdefault /* simple state */
	}
	if exprrcvr.char < 0 {
		exprrcvr.char, exprtoken = exprlex1(exprlex, &exprrcvr.lval)
	}
	exprn += exprtoken
	if exprn < 0 || exprn >= exprLast {
		goto exprdefault
	}
	exprn = exprAct[exprn]
	if exprChk[exprn] == exprtoken { /* valid shift */
		exprrcvr.char = -1
		exprtoken = -1
		exprVAL = exprrcvr.lval
		exprstate = exprn
		if Errflag > 0 {
			Errflag--
		}
		goto exprstack
	}

exprdefault:
	/* default state action */
	exprn = exprDef[exprstate]
	if exprn == -2 {
		if exprrcvr.char < 0 {
			exprrcvr.char, exprtoken = exprlex1(exprlex, &exprrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if exprExca[xi+0] == -1 && exprExca[xi+1] == exprstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			exprn = exprExca[xi+0]
			if exprn < 0 || exprn == exprtoken {
				break
			}
		}
		exprn = exprExca[xi+1]
		if exprn < 0 {
			goto ret0
		}
	}
	if exprn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			exprlex.Error(exprErrorMessage(exprstate, exprtoken))
			Nerrs++
			if exprDebug >= 1 {
				__yyfmt__.Printf("%s", exprStatname(exprstate))
				__yyfmt__.Printf(" saw %s\n", exprTokname(exprtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for exprp >= 0 {
				exprn = exprPact[exprS[exprp].yys] + exprErrCode
				if exprn >= 0 && exprn < exprLast {
					exprstate = exprAct[exprn] /* simulate a shift of "error" */
					if exprChk[exprstate] == exprErrCode {
						goto exprstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if exprDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", exprS[exprp].yys)
				}
				exprp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if exprDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", exprTokname(exprtoken))
			}
			if exprtoken == exprEofCode {
				goto ret1
			}
			exprrcvr.char = -1
			exprtoken = -1
			goto exprnewstate /* try again in the same state */
		}
	}

	/* reduction by production exprn */
	if exprDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", exprn, exprStatname(exprstate))
	}

	exprnt := exprn
	exprpt := exprp
	_ = exprpt // guard against "declared and not used"

	exprp -= exprR2[exprn]
	// exprp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if exprp+1 >= len(exprS) {
		nyys := make([]exprSymType, len(exprS)*2)
		copy(nyys, exprS)
		exprS = nyys
	}
	exprVAL = exprS[exprp+1]

	/* consult goto table to find next state */
	exprn = exprR1[exprn]
	exprg := exprPgo[exprn]
	exprj := exprg + exprS[exprp].yys + 1

	if exprj >= exprLast {
		exprstate = exprAct[exprg]
	} else {
		exprstate = exprAct[exprj]
		if exprChk[exprstate] != -exprn {
			exprstate = exprAct[exprg]
		}
	}
	// dummy call; replaced with literal code
	switch exprnt {

	case 2:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:48
		{
			exprDollar[1].f(0)
		}
	case 7:
		exprDollar = exprS[exprpt-3 : exprpt+1]
		//line expr.y:61
		{
			exprVAL.f = func(level int) {
				exprDollar[2].f(level)
			}
		}
	case 8:
		exprDollar = exprS[exprpt-2 : exprpt+1]
		//line expr.y:66
		{
		}
	case 9:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:70
		{
			exprVAL.f = func(level int) {
				exprDollar[1].f(level)
			}
		}
	case 10:
		exprDollar = exprS[exprpt-3 : exprpt+1]
		//line expr.y:76
		{
			exprVAL.f = func(level int) {
				exprDollar[1].f(level)
				exprDollar[3].f(level)
			}
		}
	case 11:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:86
		{
			exprVAL.f = func(level int) {
				for i = 0; i < level; i++ {
					fmt.Print(" ")
				}
				fmt.Print(" - ")
				exprDollar[1].f(level + 1)
			}
		}
	case 12:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:95
		{
		}
	case 13:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:97
		{
			exprVAL.f = func(level int) {
				for i = 0; i < level; i++ {
					fmt.Print(" ")
				}
				fmt.Print(" - ")
				exprDollar[1].f()
			}
		}
	case 14:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:109
		{
			exprVAL.f = func() {
				fmt.Println(exprDollar[1].stringVal, f_p.stringVal)
			}
		}
	case 15:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:114
		{
		}
	case 16:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:115
		{
		}
	case 17:
		exprDollar = exprS[exprpt-1 : exprpt+1]
		//line expr.y:116
		{
		}
	}
	goto exprstack /* stack new state and value */
}
