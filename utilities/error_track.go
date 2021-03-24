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
	// I remove the _.go of the input name and I replace in the output file for .go
	out, err := os.Create(argsWithoutProg[0][:len(argsWithoutProg[0])-4] + "__.go")
    if err != nil { panic(err) }



	file	:= string(fileToParse)

	lines	:= strings.Split(file,"\n")
	// obj.Error code generation
	// detection of the "r" assignament
	re1 := regexp.MustCompile(`(\s\s*)r\s*=`)				// one group to capture : the spaces
	// detection of variables ended "_r" assignament
	re2 := regexp.MustCompile(`(\s\s*)([A-Za-z]+_r)\s*=`)	// two groups to capture: the spaces and var name
	// generic error code generation
	// detection of the "err" assignament
	re3 := regexp.MustCompile(`(\s\s*)err\s*=`)				// one group to capture : the spaces
	// detection of variables ended "_err" assignament
	re4 := regexp.MustCompile(`(\s\s*)([A-Za-z]+_err)\s*=`)	// two groups to capture: the spaces and var name

	for _,line := range lines {
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
