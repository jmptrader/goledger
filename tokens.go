package main

import "github.com/prataprc/goparsec"

var ytok_equal = parsec.Token("=", "EQUAL")

//---- Transaction tokens
type Transprefix byte
type Transcode string
type Transnote string

var ytok_accname = parsec.Token("[a-zA-Z][a-zA-Z: ~.,;?/-]*", "FULLACCNM")
var ytok_vaccname = parsec.Token(`\([a-zA-Z][a-zA-Z: ~.,;?/-]*\)`, "VFULLACCNM")
var ytok_baccname = parsec.Token(`\[[a-zA-Z][a-zA-Z: ~.,;?/-]*\]`, "BFULLACCNM")

var ytok_prefix = parsec.Maybe(
	func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		s := string(nodes[0].(*parsec.Terminal).Value)
		return Transprefix(s[0])
	},
	parsec.Token(`\*|!`, "TRANSPREFIX"),
)
var ytok_code = parsec.Maybe(
	func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		code := string(nodes[0].(*parsec.Terminal).Value)
		ln := len(code)
		return Transcode(code[1 : ln-1])
	},
	parsec.Token(`\(.*\)`, "TRANSCODE"),
)
var ytok_desc = parsec.Token(".+", "TRANSDESC")
var ytok_persnote = parsec.Token(";[^;]+", "TRANSPNOTE")

//---- Posting tokens

var ytok_postamount = parsec.Token("[^;]+", "AMOUNT")
var ytok_postnote = parsec.Token(";[^;]+", "TRANSNOTE")

//---- Directives
var ytok_account = parsec.Token("account", "DRTV_ACCOUNT")
var ytok_note = parsec.Token("note", "DRTV_ACCOUNT_NOTE")
var ytok_alias = parsec.Token("alias", "DRTV_ACCOUNT_ALIAS")
var ytok_payee = parsec.Token("payee", "DRTV_ACCOUNT_PAYEE")
var ytok_check = parsec.Token("check", "DRTV_ACCOUNT_CHECK")
var ytok_assert = parsec.Token("assert", "DRTV_ACCOUNT_ASSERT")
var ytok_eval = parsec.Token("eval", "DRTV_ACCOUNT_EVAL")
var ytok_default = parsec.Token("default", "DRTV_ACCOUNT_DEFAULT")
var ytok_value = parsec.Token(".*", "DRTV_VALUE")

var ytok_apply = parsec.Token("apply", "DRTV_APPLY")

var ytok_aliasname = parsec.Token("[^=]+", "DRTV_ALIASNAME")

//
func maybenode(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return nodes[0]
}

func vector2scalar(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return nodes[0]
}
