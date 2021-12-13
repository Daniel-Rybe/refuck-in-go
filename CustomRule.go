package main

import (
	"container/list"
)

type CustomRule struct {
	matchList   *list.List
	replaceList *list.List
	varValues   map[Token]Token
	e           *list.Element
	lastToken   Token
}

func (rule *CustomRule) init(matchList *list.List, replaceList *list.List) {
	rule.matchList = matchList
	rule.replaceList = replaceList
	rule.varValues = map[Token]Token{}
	rule.reset()
}

func (rule *CustomRule) reset() {
	for t := range rule.varValues {
		delete(rule.varValues, t)
	}
	rule.e = rule.matchList.Front()
	rule.lastToken = rule.e.Value.(Token)
}

func (rule *CustomRule) match(token *Token) (MatchResult, *list.List) {
	if rule.lastToken.isUpper() {
		value, ok := rule.varValues[rule.lastToken]
		if ok && value != *token {
			rule.reset()
			return MATCH_RESET, nil
		}
		rule.varValues[rule.lastToken] = *token
	} else if rule.lastToken != *token {
		rule.reset()
		return MATCH_RESET, nil
	}

	if rule.e.Next() != nil {
		rule.e = rule.e.Next()
		rule.lastToken = rule.e.Value.(Token)
		return MATCH_OK, nil
	}
	replaceList := list.New()
	for e := rule.replaceList.Front(); e != nil; e = e.Next() {
		currentToken := e.Value.(Token)
		value, ok := rule.varValues[currentToken]
		if ok {
			replaceList.PushBack(value)
		} else {
			replaceList.PushBack(currentToken)
		}
	}
	rule.reset()
	return MATCH_DONE, replaceList
}
