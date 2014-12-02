package qeconv

import (
	"strings"
)

type synConv struct {
}

func (m *synConv) All(f Formula, co *cnv_out) {
	prefix(f, m, "All(", ")", co)
}

func (m *synConv) Ex(f Formula, co *cnv_out) {
	prefix(f, m, "Ex(", ")", co)
}

func (m *synConv) And(f Formula, co *cnv_out) {
	prefix(f, m, "And(", ")", co)
}

func (m *synConv) Or(f Formula, co *cnv_out) {
	prefix(f, m, "Or(", ")", co)
}

func (m *synConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "Not(", ")", co)
}

func (m *synConv) Impl(f Formula, co *cnv_out) {
	prefix(f, m, "Impl(", ")", co)
}

func (m *synConv) Equiv(f Formula, co *cnv_out) {
	prefix(f, m, "Equiv(", ")", co)
}

func (m *synConv) Abs(f Formula, co *cnv_out) {
	co.append("abs(")
	convm(f.args[0], m, co)
	co.append(")")
}

func (m *synConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *synConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *synConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " == ", co)
}

func (m *synConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " <> ", co)
}

func (m *synConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *synConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *synConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *synConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, "*", co)
}
func (m *synConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *synConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *synConv) Ftrue() string {
	return "true"
}

func (m *synConv) Ffalse() string {
	return "false"
}
func (m *synConv) Comment(str string) string {
	return "#" + str
}

func ToSyn(str string) string {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	return conv(fml, new(synConv), l.comment)
}
