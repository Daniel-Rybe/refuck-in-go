package main

import (
	"container/list"
	"fmt"
)

type BuiltInOutputRule struct {
	stage int
	token *Token
}

func (rule *BuiltInOutputRule) init() {
	rule.reset()
}

func (rule *BuiltInOutputRule) reset() {
	rule.stage = 0
}

func (rule *BuiltInOutputRule) match(token *Token) (MatchResult, *list.List) {
	if rule.stage == 0 {
		rule.token = token
		rule.stage = 1
		return MATCH_OK, nil
	}
	if token.endsWith("!") {
		prompt := ""
		if len(*token) > 1 {
			prompt = string((*token)[:len(*token)-1]) + " "
		}
		fmt.Print(prompt)
		fmt.Println(*rule.token)
		replaceList := list.New()
		replaceList.PushBack(*rule.token)
		rule.reset()
		return MATCH_DONE, replaceList
	}
	rule.reset()
	return MATCH_RESET, nil
}
