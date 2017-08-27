ntriples - a basic ntriples parser in Go

[![Build Status](https://travis-ci.org/iand/ntriples.svg?branch=master)](https://travis-ci.org/iand/ntriples)

INSTALLATION
============

Simply run

	go get github.com/iand/ntriples

Documentation is at [http://go.pkgdoc.org/github.com/iand/ntriples](http://go.pkgdoc.org/github.com/iand/ntriples)

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
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.

CREDITS
=======
The design and logic is hugely inspired by Go's standard csv parsing library
