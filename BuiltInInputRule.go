package main

import (
	"container/list"
	"fmt"
)

type BuiltInInputRule struct {
}

func (rule *BuiltInInputRule) init() {
}

func (rule *BuiltInInputRule) match(token *Token) (MatchResult, *list.List) {
	if token.endsWith("?") {
		prompt := ""
		if len(*token) > 1 {
			prompt = string((*token)[:len(*token)-1]) + " "
		}
		fmt.Print(prompt)
		var inputString string
		fmt.Scan(&inputString)
		replaceList := list.New()
		replaceList.PushBack(Token(inputString))
		return MATCH_DONE, replaceList
	}
	return MATCH_RESET, nil
}
