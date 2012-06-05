ntriples - a basic ntriples parser in Go

USAGE
=====

	package main

	import (
		"fmt"
		"os"
		"github.com/iand/ntriples"
	)	

	func main() {
		ntfile, err := os.Open("mytriples.nt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(1)
		}
		defer ntfile.Close()


		count := 0
		r := ntriples.NewReader(ntfile)
		
		for triple, err := r.Read(); err == nil;  triple, err = r.Read() {
			count++
			if count % 5000 == 0{
				fmt.Printf("%s\n", triple)
			}
			
		}


	}
