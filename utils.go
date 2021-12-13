package main

import (
	"container/list"
	"strconv"
	"unicode"
)

type Token string

const NIL_TOKEN Token = ""

func (t *Token) isUpper() bool {
	hasLetters := false
	for _, r := range *t {
		if unicode.IsLower(r) {
			return false
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}
	return hasLetters
}

func (t *Token) endsWith(s string) bool {
	diff := len(*t) - len(s)
	if diff < 0 {
		return false
	}
	for i := len(s) - 1; i >= 0; i-- {
		if (*t)[i+diff] != s[i] {
			return false
		}
	}
	return true
}

type MatchResult int

const (
	MATCH_OK    MatchResult = iota
	MATCH_DONE  MatchResult = iota
	MATCH_RESET MatchResult = iota
)

type Rule interface {
	match(*Token) (MatchResult, *list.List)
}

func isInt(val float64) bool {
	return val == float64(int(val))
}

type Operation int

const (
	PLUS  Operation = iota
	MINUS Operation = iota
)

var opsMap = map[Token]Operation{
	"+": PLUS,
	"-": MINUS,
}

func applyOpFloat(f1, f2 float64, op Operation) float64 {
	switch op {
	case PLUS:
		return f1 + f2
	case MINUS:
		return f1 - f2
	}
	return 0
}

func applyOpInt(i1, i2 int64, op Operation) int64 {
	switch op {
	case PLUS:
		return i1 + i2
	case MINUS:
		return i1 - i2
	}
	return 0
}

func getResultToken(ops []Token, op Operation) Token {
	i1, erri1 := strconv.ParseInt(string(ops[0]), 10, 64)
	i2, erri2 := strconv.ParseInt(string(ops[1]), 10, 64)
	f1, errf1 := strconv.ParseFloat(string(ops[0]), 64)
	f2, errf2 := strconv.ParseFloat(string(ops[1]), 64)

	if erri1 != nil && errf1 != nil || erri2 != nil && errf2 != nil {
		return NIL_TOKEN
	}

	if erri1 == nil && erri2 == nil {
		return Token(strconv.FormatInt(applyOpInt(i1, i2, op), 10))
	}

	result := applyOpFloat(f1, f2, op)

	if erri1 == nil {
		result = applyOpFloat(float64(i1), f2, op)
	} else if erri2 == nil {
		result = applyOpFloat(f1, float64(i2), op)
	}

	if isInt(result) {
		return Token(strconv.FormatInt(int64(result), 10))
	}

	resultToken := Token(strconv.FormatFloat(result, 'f', 10, 64))
	trailingZeroStart := len(resultToken)
	for resultToken[trailingZeroStart-1] == '0' {
		trailingZeroStart--
	}
	return resultToken[:trailingZeroStart]
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func notAll(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return true
		}
	}
	return false
}

func updateSubList(l *list.List, start, end *list.Element, s *list.List) {
	if start != end {
		for start.Next() != end {
			l.Remove(start.Next())
		}
		l.Remove(end)
	}
	start.Value = s.Front().Value
	for e := s.Back(); e != s.Front(); e = e.Prev() {
		l.InsertAfter(e.Value, start)
	}
}

/*
func printList(l *list.List) {
	print("[")
	for e := l.Front(); e != nil; e = e.Next() {
		print(e.Value.(Token))
		if e.Next() != nil {
			print(" ")
		}
	}
	print("]")
}
*/
