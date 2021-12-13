package main

import (
	"flag"
	"fmt"
)

func main() {
	s := flag.Bool("s", false, "step by step mode")
	l := flag.Bool("l", false, "log intermediate state")
	i := flag.String("i", "", "specify input file")
	flag.Parse()

	ri := RefuckInterpreter{}
	ri.init(*i)

	for !ri.done {
		if *s || *l {
			for e := ri.program.Front(); e != nil; e = e.Next() {
				fmt.Print(e.Value.(Token), " ")
			}
		}
		if *s {
			fmt.Scanln()
		} else if *l {
			fmt.Print("\n")
		}

		ri.step()
	}
}
