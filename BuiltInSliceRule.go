package main

import (
	"container/list"
	"strconv"
)

type BuiltInSliceRule struct {
	stage int
	token *Token
	slice int64
}

func (rule *BuiltInSliceRule) init() {
	rule.reset()
}

func (rule *BuiltInSliceRule) reset() {
	rule.stage = 0
}

func (rule *BuiltInSliceRule) match(token *Token) (MatchResult, *list.List) {
	if rule.stage == 0 {
		rule.token = token
		rule.stage = 1
		return MATCH_OK, nil
	}
	if rule.stage == 1 {
		tokenInt, err := strconv.ParseInt(string(*token), 10, 64)
		if err != nil {
			rule.reset()
			return MATCH_RESET, nil
		}
		rule.slice = tokenInt
		rule.stage = 2
		return MATCH_OK, nil
	}
	if *token != "/" {
		rule.reset()
		return MATCH_RESET, nil
	}
	tokenLen := int64(len(*rule.token))
	var replaceToken Token
	if rule.slice > tokenLen || rule.slice < -tokenLen {
		replaceToken = *rule.token
	} else if rule.slice >= 0 {
		replaceToken = (*rule.token)[:rule.slice]
	} else {
		replaceToken = (*rule.token)[tokenLen+rule.slice:]
	}
	replaceList := list.New()
	replaceList.PushBack(replaceToken)
	rule.reset()
	return MATCH_DONE, replaceList
}
