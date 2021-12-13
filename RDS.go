package main

import "container/list"

type RDS struct {
	stage       int
	nestedDepth int
	matchList   *list.List
	replaceList *list.List
}

func (rds *RDS) init() {
	rds.matchList = list.New()
	rds.replaceList = list.New()
	rds.reset()
}

func (rds *RDS) reset() {
	rds.stage = 0
	rds.nestedDepth = 0
	rds.matchList.Init()
	rds.replaceList.Init()
}

func (rds *RDS) match(token *Token) (MatchResult, *CustomRule) {
	if *token == "{" {
		rds.nestedDepth++
	} else if *token == "}" {
		rds.nestedDepth--
	}

	if rds.stage == 0 {
		if rds.nestedDepth == 0 {
			return MATCH_RESET, nil
		}
		rds.stage = 1
		return MATCH_OK, nil
	}
	if rds.stage == 1 {
		if *token == ":" && rds.nestedDepth == 1 {
			rds.stage = 2
			return MATCH_OK, nil
		}
		rds.matchList.PushBack(*token)
		return MATCH_OK, nil
	}
	if *token == "}" && rds.nestedDepth == 0 {
		newRule := CustomRule{}
		newMatchList := list.New()
		for e := rds.matchList.Front(); e != nil; e = e.Next() {
			newMatchList.PushBack(e.Value.(Token))
		}
		newReplaceList := list.New()
		for e := rds.replaceList.Front(); e != nil; e = e.Next() {
			newReplaceList.PushBack(e.Value.(Token))
		}
		newRule.init(newMatchList, newReplaceList)
		rds.reset()
		return MATCH_DONE, &newRule
	}
	rds.replaceList.PushBack(*token)
	return MATCH_OK, nil
}
