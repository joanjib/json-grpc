package main

import (

	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"regexp"
)
// macro line
type mLine struct {
	line	string
	subs	bool		// must be inserted here a line
}
func generateCode(out *os.File,macro []mLine,funcLines []string,varLines []string) {

	for i,_:= range funcLines { // for each function to expand:
		first := true
		for _,line := range macro {

			if first {	// begining of a function expansion
				if line.subs {
					fmt.Fprintln(out, funcLines[i])
					first = false
				} else {
					fmt.Fprintln(out, line.line)
				}
			} else {	// var expansions
				if line.subs {
					fmt.Fprintln(out, varLines[i])
				} else {
					fmt.Fprintln(out, line.line)
				}
			}
		}
	}
}

func main () {

	if len(os.Args) != 2 {
		fmt.Println("Utitility program to add error control on the lines with 'r:=' var or '_r:=' assignaments")
		panic("This program only accept one parameter")
	}
	// input file
	argsWithoutProg := os.Args[1:]

	fileToParse,err := ioutil.ReadFile(argsWithoutProg[0])
    if err != nil {
        panic(err)
    }

	// output file:
	// I remove the __.go of the input name and I replace in the output file for .go
	out, err := os.Create(argsWithoutProg[0][:len(argsWithoutProg[0])-5] + ".go")
    if err != nil { panic(err) }

	file	:= string(fileToParse)

	lines	:= strings.Split(file,"\n")

	// detection of the "func macro" assignament
	funcDet := regexp.MustCompile(`^//<<(func.*)>>`)				// Captures the function detection line
	// detection of the var macros //<<var
	varDet  := regexp.MustCompile(`^//<<(var.*)>>`)	// capturing var variables into the functions
	// detection of the end of the macro
	macroEnd:= regexp.MustCompile(`^//<<end>>`)	// capturing var variables into the functions

	macro		:= []mLine{}		// the whole function storing and each line to be substituted masked with a boolean
	funcLines	:= []string{}		// func style lines
	varLines	:= []string{}		// var lines
	beginMacro	:= false

	for i:=0;i<len(lines);i++ {

		res := funcDet.FindStringSubmatch(lines[i])
		if len(res) == 2  {		// detection of the func macro
			beginMacro = true
			// while there are macro head lines we must store them:
			macro		= append(macro,mLine{lines[i],false})		// storing the func macro detection
			funcLines	= append(funcLines,res[1])
			fmt.Fprintln(out, lines[i])
			i++
			for j:=i;j<len(lines);j++ {	// macro reading and storing
				res = funcDet.FindStringSubmatch(lines[j])
				if len(res)!= 2 {		// no more matches -> break the macro 
					macro		= append(macro,mLine{lines[j],true})		// storing the "hole" for the function expansion
					fmt.Fprintln(out, lines[j])
					i=j+1
					break
				} else {
					macro		= append(macro,mLine{lines[j],false})		// storing the func macro detection
					funcLines	= append(funcLines,res[1])
					fmt.Fprintln(out, lines[j])
				}
			}
			// end storing macro head lines

		}// we exit this block with the first "actual" func line (not a macro one) stored in the macro slice

		res = varDet.FindStringSubmatch(lines[i])
		if len(res) == 2 {
			// while there are macro variables lines we must store them:
			macro		= append(macro,mLine{lines[i],false})
			varLines	= append(varLines,res[1])
			fmt.Fprintln(out, lines[i])
			i++
			for j:=i;j<len(lines);j++ {	// macro reading and storing
				res = varDet.FindStringSubmatch(lines[j])
				if len(res) != 2 {		// no more matches -> break the macro 
					macro		= append(macro,mLine{lines[j],true})
					fmt.Fprintln(out, lines[j])
					i=j+1
					break
				} else {
					macro		= append(macro,mLine{lines[j],false})
					varLines	= append(varLines,res[1])
					fmt.Fprintln(out, lines[j])
				}
			}
			// end storing macro var lines

		}// we exit this block with the first "actual" var line (not a macro one) stored in the macro slice

		res = macroEnd.FindStringSubmatch(lines[i])
		if len(res) == 1 {
			beginMacro	= false
			fmt.Fprintln(out, lines[i])
			generateCode(out,macro,funcLines,varLines)
		}
		if beginMacro	{			// inside the macro processing => no output lines
			macro		= append(macro,mLine{lines[i],false})		// storing actual code
			fmt.Fprintln(out, lines[i])
		}else			{			// no inside a macro processing => writting to the general output.
			fmt.Fprintln(out, lines[i])
		}
	}

	err = out.Close()
	if (err != nil) { panic(err) }

}
