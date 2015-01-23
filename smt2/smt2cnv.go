package qeconv

import (
	"errors"
	. "github.com/hiwane/qeconv/def"
	"strconv"
)

type Smt2Conv struct {
	err error
}

func (m *Smt2Conv) quantifier(f Formula, co *CnvOut, qstr string) {
	q := f.Args()[0]
	co.Append("(" + qstr + " (")
	for i := 0; i < len(q.Args()); i++ {
		co.Append(" (")
		Conv2(q.Arg(i), m, co)
		co.Append(" Real)")
	}
	co.Append(" ) ")
	Conv2(f.Args()[1], m, co)
	co.Append(" )")
}

func (m *Smt2Conv) All(f Formula, co *CnvOut) {
	m.quantifier(f, co, "forall")
}

func (m *Smt2Conv) Ex(f Formula, co *CnvOut) {
	m.quantifier(f, co, "exists")
}

func (m *Smt2Conv) And(f Formula, co *CnvOut) {
	Prefixm(f, m, "(and ", " ", ")", co)
}

func (m *Smt2Conv) Or(f Formula, co *CnvOut) {
	Prefixm(f, m, "(or ", " ", ")", co)
}

func (m *Smt2Conv) Not(f Formula, co *CnvOut) {
	Prefixm(f, m, "(not ", " ", ")", co)
}

func (m *Smt2Conv) Impl(f Formula, co *CnvOut) {
	Prefixm(f, m, "(=> ", " ", ")", co)
}

func (m *Smt2Conv) Equiv(f Formula, co *CnvOut) {
	co.Append("(and (=> ")
	Conv2(f.Args()[0], m, co)
	co.Append(" ")
	Conv2(f.Args()[1], m, co)
	co.Append(") (=> ")
	Conv2(f.Args()[1], m, co)
	co.Append(" ")
	Conv2(f.Args()[0], m, co)
	co.Append("))")
}

func (m *Smt2Conv) Abs(f Formula, co *CnvOut) {
	m.err = errors.New("unsupport the abs function")
}

func (m *Smt2Conv) Leop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(<= ", " ", ")", co)
}

func (m *Smt2Conv) Ltop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(< ", " ", ")", co)
}

func (m *Smt2Conv) Eqop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(= ", " ", ")", co)
}

func (m *Smt2Conv) Neop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(not (= ", " ", "))", co)
}

func (m *Smt2Conv) List(f Formula, co *CnvOut) {
	Prefixm(f, m, "(", " ", ")", co)
}

func (m *Smt2Conv) Plus(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(+ ", " ", ")", co)
}

func (m *Smt2Conv) Minus(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(- ", " ", ")", co)
}

func (m *Smt2Conv) Mult(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(* ", " ", ")", co)
}

func (m *Smt2Conv) Div(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(/ ", " ", ")", co)
}

func (m *Smt2Conv) Pow(fml Formula, co *CnvOut) {
	exp := fml.Arg(1)
	if exp.Cmd() != NUMBER {
		m.err = errors.New("unsupport rational exponential")
	}
	co.Append("(*")
	n, _ := strconv.Atoi(exp.String())
	for i := 0; i < n; i++ {
		co.Append(" ")
		Conv2(fml.Args()[0], m, co)
	}
	co.Append(")")
}

func (m *Smt2Conv) Uniop(fml Formula, ope string, co *CnvOut) {
	co.Append("(" + ope + " 0 ")
	Conv2(fml.Args()[0], m, co)
	co.Append(")")
}

func (m *Smt2Conv) Ftrue() string {
	return "(= 0 0)"
}

func (m *Smt2Conv) Ffalse() string {
	return "(= 0 1)"
}

func (m *Smt2Conv) Comment(str string) string {
	return ";" + str
}

func smt2header(fml Formula) string {
	var str string
	if fml.IsQff() {
		str = "(set logic NRA)\n"
	} else {
		str = "(set logic QF_NRA)\n"
	}

	vs := fml.FreeVars()
	for i := 0; i < vs.Len(); i++ {
		v := vs.Get(i)
		str += "(declare-fun " + v + " () Real)\n"
	}

	return str
}

func smt2footer(fml Formula) string {
	return "(check-sat)\n"
}

func (m *Smt2Conv) Convert(fml Formula, co *CnvOut) (string, error) {
	if fml.IsList() {
		return "", errors.New("unsupported input")
	}
	qc := new(Smt2Conv)
	qc.err = nil

	header := smt2header(fml)
	Conv2(fml, qc, co)
	header += "(assert " + co.String() + ")\n"
	header += smt2footer(fml)
	return header, qc.err

}

func (m *Smt2Conv) Sep() string {
	return "\n"
}