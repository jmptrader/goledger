package main

import "time"

import "github.com/prataprc/goparsec"

type Transprefix byte
type Transcode string

type Transaction struct {
	// start
	date     time.Time
	edate    time.Time
	prefix   byte
	code     string
	desc     string
	postings []*Posting
	note     string

	db *Datastore // read-only copy
}

func NewTransaction(db *Datastore) *Transaction {
	trans := &Transaction{db: db}
	return trans
}

func (trans *Transaction) Y() parsec.Parser {
	// DATE
	ydate := Ydate(
		trans.db.Year(), trans.db.Month(), trans.db.Dateformat(),
	)
	// [=EDATE]
	yequal := parsec.Token("=", "TRANSEQUAL")
	yedate := parsec.Maybe(
		maybenode,
		parsec.And(
			func(nodes []parsec.ParsecNode) parsec.ParsecNode {
				return nodes[1] // EDATE
			},
			yequal,
			ydate,
		),
	)
	// [*|!]
	yprefix := parsec.Maybe(
		func(nodes []parsec.ParsecNode) parsec.ParsecNode {
			s := string(nodes[0].(*parsec.Terminal).Value)
			return Transprefix(s[0])
		},
		parsec.Token(`\*|!`, "TRANSPREFIX"),
	)
	// [(CODE)]
	ycode := parsec.Maybe(
		func(nodes []parsec.ParsecNode) parsec.ParsecNode {
			code := string(nodes[0].(*parsec.Terminal).Value)
			ln := len(code)
			return Transcode(code[1 : ln-1])
		},
		parsec.Token(`\(.*\)`, "TRANSCODE"),
	)
	// DESC
	ydesc := parsec.Token(".+", "TRANSDESC")

	y := parsec.And(
		func(nodes []parsec.ParsecNode) parsec.ParsecNode {
			n := 0
			trans.date = nodes[n].(time.Time)
			n++
			if edate, ok := nodes[n].(time.Time); ok {
				trans.edate = edate
				n++
			}
			if prefix, ok := nodes[n].(Transprefix); ok {
				trans.prefix = byte(prefix)
				n++
			}
			if code, ok := nodes[n].(Transcode); ok {
				trans.code = string(code)
				n++
			}
			trans.desc = nodes[n].(string)
			return trans
		},
		ydate, yedate, yprefix, ycode, ydesc,
	)
	return y
}

func (trans *Transaction) Parse(scanner parsec.Scanner) parsec.Scanner {
	var bs []byte
	var node parsec.ParsecNode

	for {
		if bs, scanner = scanner.SkipWS(); len(bs) == 0 {
			return scanner
		}
		node, scanner = NewPosting(trans.db).Y()(scanner)
		trans.postings = append(trans.postings, node.(*Posting))
	}
	return scanner
}

func maybenode(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return nodes[0]
}