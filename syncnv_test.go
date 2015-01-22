package qeconv

import (
	"strings"
	"testing"
	"text/scanner"
)

func removeLineComment(s string, p rune) string {
	var ret []rune
	var l scanner.Scanner
	l.Init(strings.NewReader(s))

	for l.Peek() != scanner.EOF {
		s := l.Next()
		if s == p {
			s = l.Next()
			for s != '\n' && s != scanner.EOF {
				s = l.Next()
			}
			continue
		}
		ret = append(ret, s)
	}

	return string(ret)
}

func TestToSyn(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"true:", "true:"},
		{"false:", "false:"},
		{"x>0:", "0<x:"},
		{"x+y>0:", "0<x+y:"},
		{"x-2*y<>0:", "x-2*y <> 0:"},
		{"-x+y<=0:", "-x+y<=0:"},
		{"-x+(y*2)<=0:", "-x+y*2<=0:"},
		{"-(x+y/3)=0:", "-(x+y/3) = 0:"},
		{"-(x+y)*2<=0:", "-(x+y)*2<=0:"},
		{"Not(y=0):", "Not(y  =  0):"},
		{"And(x<=0, y=0):", "And(x <= 0, y  =  0):"},
		{"And():", "true:"},
		{"Or():", "false:"},
		{"Or(x<=0, y<>0):", "Or(x<=0, y<>0):"},
		{"And(Or(x*y>0,z>0),Or(x*y<=-z,+z>=0)):", "And(Or(0<x*y, 0<z), Or(x*y<=-z,  0<=+z)):"},
		{"Or(And(x*y>0,z>0),And(x*y<=-z,+z>=0)):", "Or(And(0<x*y, 0<z), And(x*y<=-z,  0<=+z)):"},
		{"Ex(x, x^2=-1):", "Ex([x], x^2 = -1):"},
		{"Ex({x}, x^2=-1):", "Ex([x], x^2 = -1):"},
		{"Ex([x], x^2=-1):", "Ex([x], x^2 = -1):"},
		{"Ex(x, And(x^2=-1, x>0)):", "Ex([x], And(x^2 = -1, 0<x)):"},
		{"All([x], a*x^2+b*x+c>0):", "All([x], 0<a*x^2+b*x+c):"},
		{"All([x], Ex([y], x+y+a=0)):", "All([x],Ex([y],x+y+a = 0)):"},
		{"Equiv(x<0, y=0):", "Equiv(x<0, y = 0):"},
		{"Impl(x<0, y=0):", "Impl(x<0, y = 0):"},
		{"Repl(y=0, x<0):", "Impl(x<0, y = 0):"},
		{"# comment line\n(1+a)*x+(3+b)*y=0:", "(1+a)*x+(3+b)*y  =  0:"},
		{"x+abs(y+z)=0:", "x+abs(y+z) = 0:"},
	}

	for i, p := range data {
		actual0, _ := Convert(p.input, "syn", false, 0)
		actual := removeLineComment(actual0, '#')
		if !cmpIgnoreSpace(actual, p.expect) {
			t.Errorf("err %d\nactual=%s\nexpect=%s\ninput=%s\n", i, actual0, p.expect, p.input)
		}
	}
}
