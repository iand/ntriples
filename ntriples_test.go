/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Source Code file.
  This work is published from the United Kingdom. 
*/
package ntriples

import (
  	"bytes"
	"testing"
	"strings"
)

var testCases = map[string]Triple{
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> ." : Triple{ s:RdfTerm{value:"http://example.org/resource1", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"http://example.org/resource2", termtype:RdfIri}, },

	"_:anon <http://example.org/property> <http://example.org/resource2> ." : Triple{ s:RdfTerm{value:"anon", termtype:RdfBlank},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"http://example.org/resource2", termtype:RdfIri}, },

	"<http://example.org/resource1> <http://example.org/property> _:anon ." : Triple{ s:RdfTerm{value:"http://example.org/resource1", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"anon", termtype:RdfBlank}, },

	" 	 <http://example.org/resource3> 	 <http://example.org/property>	 <http://example.org/resource2> 	." : Triple{ s:RdfTerm{value:"http://example.org/resource3", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"http://example.org/resource2", termtype:RdfIri}, },

	"<http://example.org/resource7> <http://example.org/property> \"simple literal\" ." : Triple{ s:RdfTerm{value:"http://example.org/resource7", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"simple literal", termtype:RdfLiteral}, },



	`<http://example.org/resource8> <http://example.org/property> "backslash:\\" .` : Triple{ s:RdfTerm{value:"http://example.org/resource8", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"backslash:\\", termtype:RdfLiteral}, },

	`<http://example.org/resource9> <http://example.org/property> "dquote:\"" .` : Triple{ s:RdfTerm{value:"http://example.org/resource9", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"dquote:\"", termtype:RdfLiteral}, },

	`<http://example.org/resource10> <http://example.org/property> "newline:\n" .` : Triple{ s:RdfTerm{value:"http://example.org/resource10", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"newline:\n", termtype:RdfLiteral}, },

	`<http://example.org/resource11> <http://example.org/property> "return\r" .` : Triple{ s:RdfTerm{value:"http://example.org/resource11", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"return\r", termtype:RdfLiteral}, },

	`<http://example.org/resource12> <http://example.org/property> "tab:\t" .` : Triple{ s:RdfTerm{value:"http://example.org/resource12", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"tab:\t", termtype:RdfLiteral}, },

	`<http://example.org/resource16> <http://example.org/property> "\u00E9" .` : Triple{ s:RdfTerm{value:"http://example.org/resource16", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"\u00E9", termtype:RdfLiteral}, },

	`<http://example.org/resource30> <http://example.org/property> "chat"@fr .` : Triple{ s:RdfTerm{value:"http://example.org/resource30", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"chat", language:"fr", termtype:RdfLiteral}, },

	`<http://example.org/resource31> <http://example.org/property> "chat"@en .` : Triple{ s:RdfTerm{value:"http://example.org/resource31", termtype:RdfIri},
																											  p:RdfTerm{value:"http://example.org/property", termtype:RdfIri},
																											  o:RdfTerm{value:"chat", language:"en", termtype:RdfLiteral}, },




	}




func TestRead(t *testing.T) {


	for ntriple, expected := range testCases {
		r := NewReader(strings.NewReader(ntriple))
		triple, err := r.Read()
		if err != nil {
			t.Errorf("Expected %s but got error %s", expected, err)
		}

		if triple != expected {
			t.Errorf("Expected %s but got %s", expected, triple)
		}
	}
}


func TestReadMultiple(t *testing.T) {
	var ntriples bytes.Buffer
	var triples []Triple

	for ntriple, triple := range testCases {
		ntriples.WriteString(ntriple)
		ntriples.WriteRune('\n')
		triples = append(triples, triple)
	}

	count := 0
	r := NewReader(strings.NewReader(ntriples.String()))
	triple, err := r.Read()
	for err == nil {
		if triple != triples[count] {
			t.Errorf("Expected %s but got %s", triples[count], triple)
			break
		}

		count++
		triple, err = r.Read()
	}

	if count != len(triples) {
		t.Errorf("Expected %d but only parsed %d triples", len(triples), count)

	}

}

