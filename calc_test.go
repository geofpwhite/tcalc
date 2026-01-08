package main

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	s := state{
		variables: make(map[string]int64),
	}
	err := s.Exec("1 + 3 + 2")
	if err != nil {
		t.Fail()
	}
	if s.ans != 6 {
		t.Fail()
	}
	err = s.Exec("1 * (3 + 2)")
	if err != nil || s.ans != 5 {
		t.Fail()
	}

	fmt.Println(s.ans)
	err = s.Exec("2 * (3 + 2)")
	fmt.Println(s.ans)
	if err != nil || s.ans != 10 {
		t.Fail()
	}
	err = s.Exec("1 << 5")

	fmt.Println(s.ans)
	if err != nil || s.ans != 1<<5 {
		t.Fail()
	}
	tokens := s.Tokenize("2**5")
	fmt.Println(tokens)
}

func TestDraw(t *testing.T) {
	fmt.Println(decimalDisplayString(63))
	fmt.Println(hexDisplayString(63))
	fmt.Println(binaryDisplayString(63))
}
