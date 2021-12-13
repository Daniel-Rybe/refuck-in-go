package main

import "container/list"

type BuiltInOpRule struct {
	operands [2]Token
	stage    int
}

func (rule *BuiltInOpRule) init() {
	rule.reset()
}

func (rule *BuiltInOpRule) reset() {
	rule.stage = 0
}

func (rule *BuiltInOpRule) match(token *Token) (MatchResult, *list.List) {
	if rule.stage < 2 {
		rule.operands[rule.stage] = *token
		rule.stage++
		return MATCH_OK, nil
	}
	op, ok := opsMap[*token]
	if !ok {
		rule.reset()
		return MATCH_RESET, nil
	}
	resultToken := getResultToken(rule.operands[:], op)
	rule.reset()
	if resultToken == NIL_TOKEN {
		return MATCH_RESET, nil
	}
	replaceList := list.New()
	replaceList.PushBack(resultToken)
	return MATCH_DONE, replaceList
}
