package main

import (
	"bufio"
	"container/list"
	"os"
	"regexp"
	"strings"
)

type RefuckInterpreter struct {
	program *list.List
	done    bool
	rds     RDS
	rules   []Rule
}

func (ri *RefuckInterpreter) init(srcPath string) {
	ri.program = ri.parse(srcPath)
	ri.done = false
	ri.rds.init()
	builtInInputRule := BuiltInInputRule{}
	builtInInputRule.init()
	ri.rules = append(ri.rules, &builtInInputRule)
	builtInOutputRule := BuiltInOutputRule{}
	builtInOutputRule.init()
	ri.rules = append(ri.rules, &builtInOutputRule)
	builtInOpRule := BuiltInOpRule{}
	builtInOpRule.init()
	ri.rules = append(ri.rules, &builtInOpRule)
	builtInSliceRule := BuiltInSliceRule{}
	builtInSliceRule.init()
	ri.rules = append(ri.rules, &builtInSliceRule)
}

func (ri *RefuckInterpreter) parse(srcPath string) *list.List {
	program := list.New()
	f, err := os.Open(srcPath)
	checkError(err)
	tokenRegexp := regexp.MustCompile(`\S+`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		hashIndex := strings.IndexRune(line, '#')
		if hashIndex != -1 {
			line = line[:hashIndex]
		}
		tokenStrings := tokenRegexp.FindAllString(line, -1)
		for _, tokenString := range tokenStrings {
			program.PushBack(Token(tokenString))
		}
	}
	return program
}

func (ri *RefuckInterpreter) step() {
	ri.rules = ri.rules[:4]
	ri.done = true

	for programElement := ri.program.Front(); programElement != nil; programElement = programElement.Next() {
		token := programElement.Value.(Token)

		rulesReset := make([]bool, len(ri.rules))
		offsetElement := programElement
		for offsetElement != nil && notAll(rulesReset) {
			for ruleIndex, rule := range ri.rules {
				if rulesReset[ruleIndex] {
					continue
				}
				offsetToken := offsetElement.Value.(Token)
				mr, replaceList := rule.match(&offsetToken)
				if mr == MATCH_RESET {
					rulesReset[ruleIndex] = true
				} else if mr == MATCH_DONE {
					updateSubList(ri.program, programElement, offsetElement, replaceList)
					ri.done = false
					return
				}
			}
			offsetElement = offsetElement.Next()
		}

		mr, newRule := ri.rds.match(&token)
		if mr == MATCH_DONE {
			ri.rules = append(ri.rules, newRule)
		}
	}
}
