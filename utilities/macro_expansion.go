package main

import (

	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"regexp"
)

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

	for i:=0;i<len(lines);i++ {

		res := funcDec.FindStringSubmatch(line[i])
		if len(res) > 0 {		// detection of the func macro
			var macro		:= []string		// the whole function storing
			var funcLines	:= []string		// func style lines
			var varLines	:= []string		// var lines

			macro = append(macro,line)		// storing the func macro detection
			for (j:=i;i<len(lines);j++) {	// macro reading and storing


			}

		}
		fmt.Fprintln(out, line)

		res := re1.FindStringSubmatch(line)
		if len(res) ==  2 {
			fmt.Fprintln(out, res[1] + "if r.Error != nil { goto ErrorTrack }")
			continue
		}

		res = re2.FindStringSubmatch(line)
		if len(res) == 3 {
			fmt.Fprintln(out, res[1] +"if "+ res[2] +".Error != nil { goto ErrorTrack_"+ res[2] +" }")
			continue
		}
		res = re3.FindStringSubmatch(line)
		if len(res) ==  2 {
			fmt.Fprintln(out, res[1] + "if err != nil { goto ErrorTrack_err }")
			continue
		}

		res = re4.FindStringSubmatch(line)
		if len(res) == 3 {
			fmt.Fprintln(out, res[1] +"if "+ res[2] +" != nil { goto ErrorTrack_"+ res[2] +" }")
			continue
		}
	}

	err = out.Close()
	if (err != nil) { panic(err) }

}
