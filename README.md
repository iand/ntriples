ntriples - a basic ntriples parser in Go

INSTALLATION
============

Simply run

	go get github.com/iand/ntriples

USAGE
=====

Example of parsing an ntriples file and printing out every 5000th triple

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

LICENSE
=======
This code and associated documentation is in the public domain.

To the extent possible under law, Ian Davis has waived all copyright
and related or neighboring rights to this file. This work is published from the United Kingdom. 

CREDITS
=======
The design and logic is hugely inspired by Go's standard csv parsing library