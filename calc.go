package main

import (
	"errors"
	"math"
	"slices"
	"strconv"
)

type calcNode struct {
	left       *calcNode
	right      *calcNode
	value      *string
	assignment *assignment
}

type state struct {
	variables map[string]int64
	ans       int64
}

func newState() *state {
	return &state{
		variables: make(map[string]int64),
	}
}

type (
	operator           rune
	doubleRuneOperator string
)

const (
	SUM    operator = '+'
	SUB    operator = '-'
	PROD   operator = '*'
	DIV    operator = '/'
	XOR    operator = '^'
	OR     operator = '|'
	NOT    operator = '~'
	AND    operator = '&'
	MOD    operator = '%'
	ASSIGN operator = '='
	LPAREN operator = '('
	RPAREN operator = ')'
	// two rune operators
	LEFTSHIFT  doubleRuneOperator = "<<"
	RIGHTSHIFT doubleRuneOperator = ">>"
	EXP        doubleRuneOperator = "**"
)

var length1operatorsInfix = []operator{
	SUM, SUB, PROD, DIV, XOR, AND, MOD, ASSIGN, OR,
}

var length1operatorsPrefix = []operator{
	NOT,
}

var length2operators = []doubleRuneOperator{
	LEFTSHIFT, RIGHTSHIFT,
}

var negative = []string{"("}

func (s *state) Tokenize(input string) []string {
	tokens := make([]string, 0, len(input))
	cur := ""
	for _, char := range input {
		if char == '(' || char == ')' || slices.Contains(length1operatorsInfix, operator(char)) || slices.Contains(length1operatorsPrefix, operator(char)) {
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
			tokens = append(tokens, string(char))
			continue
		}
		switch char {
		case ' ', '\r', '\n':
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
		case '>', '<', '*':
			if cur == string(char) {
				tokens = append(tokens, cur+string(char))
				cur = ""
				continue
			}
			if len(cur) > 0 {
				tokens = append(tokens, cur)
			}
			cur = string(char)
		default:
			cur = cur + string(char)
			if slices.Contains(length2operators, doubleRuneOperator(cur)) {
				tokens = append(tokens, cur)
				cur = ""
			}
		}
	}
	tokens = tokens[:len(tokens):len(tokens)]
	if cur != "" {
		tokens = append(tokens, cur)
	}
	for i, token := range tokens {
		if token == "-" {
			before := tokens[:i]
			after := tokens[i+2:]
			s := append([]string{"(", "0"}, []string{token, tokens[i+1]}...)
			s = append(s, ")")
			tokens = append(before, s...)
			tokens = append(tokens, after...)
		}
	}
	return tokens
}

func (s *state) Exec(input string) error {
	tokens := s.Tokenize(input)
	node, err := s.Parse(tokens)
	if err != nil {
		return err
	}
	value, err := s.Eval(node)
	if err != nil {
		return err
	}
	s.ans = value
	return nil
}

func (s *state) Eval(curNode calcNode) (int64, error) {
	if curNode.assignment != nil {
		num, err := s.Eval(curNode.assignment.right)
		if err != nil {
			return -1, err
		}
		s.variables[curNode.assignment.name] = num
		return num, nil
	}
	if curNode.value == nil {
		return -1, errors.New("bad value")
	}

	if slices.Contains(length1operatorsInfix, operator((*curNode.value)[0])) {
		l, err := s.Eval(*curNode.left)
		if err != nil {
			return 0, err
		}
		r, err := s.Eval(*curNode.right)
		if err != nil {
			return 0, err
		}
		switch *curNode.value {
		case "+":
			return l + r, nil
		case "-":
			return l - r, nil
		case "*":
			return l * r, nil
		case "/":
			return l / r, nil
		case "&":
			return l & r, nil
		case "^":
			return l ^ r, nil
		case "|":
			return l | r, nil
		default:
			return -1, errors.New("invalid operator")
		}
	}
	if slices.Contains(length1operatorsPrefix, operator((*curNode.value)[0])) {
		num, err := s.Eval(*curNode.right)
		if err != nil {
			return -1, err
		}
		switch *curNode.value {
		case "~":
			return ^num, nil
		default:
			return -1, errors.New("bad prefix operator")
		}
	}
	if slices.Contains(length2operators, doubleRuneOperator(*curNode.value)) {
		l, err := s.Eval(*curNode.left)
		if err != nil {
			return 0, err
		}
		r, err := s.Eval(*curNode.right)
		if err != nil {
			return 0, err
		}
		switch *curNode.value {
		case "<<":
			return l << r, nil
		case ">>":
			return l >> r, nil
		case "**":
			f := math.Pow(float64(l), float64(r))
			return int64(f), nil
		default:
			return -1, errors.New("bad double rune operator")
		}
	}
	num, err := strconv.ParseInt(*curNode.value, 10, 64)
	if err != nil {
		if *curNode.value == "_ans_" {
			return s.ans, nil
		}
		return s.variables[*curNode.value], nil
	}
	return num, nil
}

func (s *state) Parse(tokens []string) (calcNode, error) {
	node := s.parse(tokens, 0, nil)
	if node == nil {
		return calcNode{}, errors.New("nil error")
	}
	return *node, nil
}

func (s *state) parse(tokens []string, index int, cur *calcNode) *calcNode {
	if index >= len(tokens) {
		return cur
	}
	token := tokens[index]
	new := calcNode{value: &token}
	if slices.Contains(length1operatorsInfix, operator(token[0])) && token[0] != '=' {
		new.left = cur
		new.right = s.parse(tokens, index+1, nil)
		return &new
	}
	switch token {
	case "=":
		if index == 0 {
			return nil
		}
		name := *cur.value
		new = *s.parse(tokens[index+1:], 0, nil)
		return &calcNode{assignment: &assignment{name: name, right: new}}

	case "<<", ">>":
		new.left = cur
		new.right = s.parse(tokens, index+1, nil)
		return &new
	case "(":
		rParenIndex := slices.Index(tokens[index:], ")")
		if rParenIndex == -1 {
			return nil
		}
		inner := tokens[index+1:][:slices.Index(tokens[index:], ")")]
		node := s.parse(inner, 0, nil)
		return s.parse(tokens[index+1+len(inner):], 0, node)
	case ")":
		return s.parse(tokens, index+1, cur)
	default:
		return s.parse(tokens, index+1, &new)
	}
}

type assignment struct {
	name  string
	right calcNode
}
