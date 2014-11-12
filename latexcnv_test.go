package qeconv

import (
	"strings"
	"testing"
	"text/scanner"
)

func removeLaTeXComment(s string) string {
	var ret []rune
	var l scanner.Scanner
	l.Init(strings.NewReader(s))

	for l.Peek() != scanner.EOF {
		s := l.Peek()
		if s == '%' {
			for l.Peek() != '\n' {
				l.Next()
			}
			s = l.Next()
		} else {
			l.Next()
		}
		ret = append(ret, s)
	}

	return string(ret)
}

func TestToLaTeX(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"x>0:", "0<x"},
		{"x+y>0:", "0<x+y"},
		{"-x+y<=0:", "-x+y \\leq 0"},
		{"-x+(y*2)<=0:", "-x+y 2 \\leq 0"},
		{"-(x+y/3)=0:", "-(x+y/3)=0"},
		{"-(x+y)*2<=0:", "-(x+y) 2 \\leq 0"},
		{"And(x<=0, y=0):", "x \\leq 0 \\land y = 0"},
		{"Or(x<=0, y<>0):", "x \\leq 0 \\lor y \\neq 0"},
		{"Ex(x, x^2=-1):", "\\exists x (x^{2}=-1)"},
		{"All([x], a*x^2+b*x+c>0):", "\\forall x (0<a x^{2}+b x+c)"},
	}

	for _, p := range data {
		t.Log("inp=%s\n", p.input)
		actual0 := ToLaTeX(p.input)
		t.Log("ac0=%s\n", actual0)
		actual := removeMathComment(actual0)
		t.Log("rem=%s\n", actual)
		if !cmpIgnoreSpace(actual, p.expect) {
			t.Errorf("err actual=%s\nexpect=%s\ninput=%s\n", actual0, p.expect, p.input)
		}
	}
}