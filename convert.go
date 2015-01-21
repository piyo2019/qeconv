package qeconv

import (
	"errors"
	"strings"
)

type cnv_out struct {
	str     string
	lno     int
	comment []Comment
}

type CnvInfMathOp interface {
	/* math op */
	Plus(f Formula, co *cnv_out)
	Minus(f Formula, co *cnv_out)
	Mult(f Formula, co *cnv_out)
	Div(f Formula, co *cnv_out)
	Pow(f Formula, co *cnv_out)
	Comment(str string) string
	Abs(f Formula, co *cnv_out)
	uniop(f Formula, ope string, co *cnv_out)

	/* atom */
	Ftrue() string
	Ffalse() string
}

type CnvInf interface {
	CnvInfMathOp

	/* quantifier */
	All(f Formula, co *cnv_out)
	Ex(f Formula, co *cnv_out)

	/* logical operator
	 * Mathematica: Implies/Equivalent
	 * redlog     : impl/repl/equiv
	 * qepcad     : =>/<=/<=>
	 */
	And(f Formula, co *cnv_out)
	Or(f Formula, co *cnv_out)
	Not(f Formula, co *cnv_out)
	Impl(f Formula, co *cnv_out)
	Equiv(f Formula, co *cnv_out)

	/* comparator */
	Leop(f Formula, co *cnv_out)
	Ltop(f Formula, co *cnv_out)
	Eqop(f Formula, co *cnv_out)
	Neop(f Formula, co *cnv_out)

	List(f Formula, co *cnv_out)
}

type Formula struct {
	cmd      int
	str      string
	args     []Formula
	priority int
	lineno   int
}

func (self *Formula) Cmd() int {
	return self.cmd
}
func (self *Formula) IsList() bool {
	return self.cmd == LIST
}
func (self *Formula) IsQuantifier() bool {
	return self.cmd == ALL || self.cmd == EX
}
func (self *Formula) IsAtom() bool {
	return self.cmd == LEOP || self.cmd == LTOP ||
		self.cmd == EQOP || self.cmd == NEOP
}
func (self *Formula) IsBool() bool {
	return self.cmd == F_TRUE || self.cmd == F_FALSE
}
func (self *Formula) String() string {
	return conv(*self, new(synConv), make([]Comment, 0))
}
func (self *Formula) Args() []Formula {
	return self.args
}
func (self *Formula) Arg(id int) Formula {
	return self.args[id]
}

func (c *cnv_out) append(s string) {
	//	fmt.Printf("append [%s]\n", s)
	c.str += s
}

func SynToFml(str string) Formula {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	return tofml(stack)
}

func cmpfml(f1, f2 Formula) bool {
	if f1.cmd != f2.cmd {
		return false
	}

	if len(f1.args) != len(f2.args) {
		return false
	}

	for i := 0; i < len(f1.args); i++ {
		if !cmpfml(f1.args[i], f2.args[i]) {
			return false
		}
	}

	return true
}

func rmdup(fml Formula) Formula {
	if fml.cmd != LIST {
		return fml
	}

	for i := 1; i < len(fml.args); i++ {
		for j := 0; j < i; j++ {
			if cmpfml(fml.args[i], fml.args[j]) {
				args := fml.args
				fml.args = make([]Formula, len(args) - 1)
				for k := 0; k < i; k++ {
					fml.args[k] = args[k]
				}
				for k := i+1; k < len(args); k++ {
					fml.args[k] = args[k+1]
				}
				break
			}
		}
	}

	return fml
}

func tofml(s *Stack) Formula {
	n, _ := s.pop()
	fml := Formula{}
	fml.cmd = n.cmd
	fml.str = n.str
	fml.lineno = n.lineno
	fml.priority = n.priority
	fml.args = make([]Formula, n.val)
	if (fml.cmd == OR || fml.cmd == AND) && n.val == 1 {
		return tofml(s)
	}
	if n.rev {
		for i := 0; i < n.val; i++ {
			fml.args[i] = tofml(s)
		}
	} else {
		for i := 0; i < n.val; i++ {
			fml.args[n.val-i-1] = tofml(s)
		}
	}
	return fml
}

func mop(fml Formula, cinf CnvInfMathOp, op string, co *cnv_out) {
	for i := 0; i < len(fml.args); i++ {
		if i != 0 {
			co.append(op)
		}
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			co.append("(")
			convm(fml.args[i], cinf, co)
			co.append(")")
		} else {
			convm(fml.args[i], cinf, co)
		}
	}
}
func uniop(fml Formula, cinf CnvInfMathOp, op string, co *cnv_out) {
	co.append(op)
	if fml.priority > 0 && fml.priority < fml.args[0].priority {
		co.append("(")
		convm(fml.args[0], cinf, co)
		co.append(")")
	} else {
		convm(fml.args[0], cinf, co)
	}
}

func conv(fml Formula, cinf CnvInf, comment []Comment) string {
	var co *cnv_out
	co = &cnv_out{str: "", lno: 1, comment: comment}
	conv2(fml, cinf, co)
	return co.str
}

func skipcomment(fml Formula, cinf CnvInfMathOp, co *cnv_out) {
	for co.lno < fml.lineno {
		if len(co.comment) > 0 && co.comment[0].lineno == co.lno {
			co.append(cinf.Comment(co.comment[0].str))
			co.comment = co.comment[1:len(co.comment)]
		}
		co.append("\n")
		co.lno++
	}
}

func convm(fml Formula, cinf CnvInfMathOp, co *cnv_out) {

	skipcomment(fml, cinf, co)

	switch fml.cmd {
	case PLUS:
		cinf.Plus(fml, co)
	case MINUS:
		cinf.Minus(fml, co)
	case MULT:
		cinf.Mult(fml, co)
	case DIV:
		cinf.Div(fml, co)
	case POW:
		cinf.Pow(fml, co)
	case ABS:
		cinf.Abs(fml, co)
	case NAME, NUMBER:
		co.append(fml.str)
	case F_TRUE:
		co.append(cinf.Ftrue())
	case F_FALSE:
		co.append(cinf.Ffalse())
	case UNARYMINUS:
		cinf.uniop(fml, "-", co)
	case UNARYPLUS:
		cinf.uniop(fml, "+", co)
	default:
		errors.New("unknown type")
	}
}

func conv2(fml Formula, cinf CnvInf, co *cnv_out) {
	//	fmt.Printf("fml.cmd=%d,lineno=%d/%d str=%s\n", fml.cmd, fml.lineno, co.lno, fml.str)
	skipcomment(fml, cinf, co)

	switch fml.cmd {
	case ALL:
		cinf.All(fml, co)
	case EX:
		cinf.Ex(fml, co)
	case AND:
		cinf.And(fml, co)
	case OR:
		cinf.Or(fml, co)
	case NOT:
		cinf.Not(fml, co)
	case IMPL:
		cinf.Impl(fml, co)
	case EQUIV:
		cinf.Equiv(fml, co)
	case LEOP:
		cinf.Leop(fml, co)
	case LTOP:
		cinf.Ltop(fml, co)
	case EQOP:
		cinf.Eqop(fml, co)
	case NEOP:
		cinf.Neop(fml, co)
	case LIST:
		cinf.List(fml, co)
	default:
		convm(fml, cinf, co)
	}
}

func prefixm(fml Formula, cinf CnvInf, left, mid, right string, co *cnv_out) {
	co.append(left)
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		co.append(sep)
		conv2(fml.args[i], cinf, co)
		sep = mid
	}
	co.append(right)
}

func prefix(fml Formula, cinf CnvInf, left, right string, co *cnv_out) {
	prefixm(fml, cinf, left, ",", right, co)
}

func infixm(fml Formula, cinf CnvInf, op string, co *cnv_out, str, end string) {
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		co.append(sep)
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			co.append(str)
			conv2(fml.args[i], cinf, co)
			co.append(end)
		} else {
			conv2(fml.args[i], cinf, co)
		}
		sep = op
	}
}

func infix(fml Formula, cinf CnvInf, op string, co *cnv_out) {
	infixm(fml, cinf, op, co, "(", ")")
}
