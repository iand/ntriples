/*
  This is free and unencumbered software released into the public domain. For more
  information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package ntriples

import (
	"bytes"
	"strings"
	"testing"
)

var testCases = map[string]Triple{
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> .": {
		S: RdfTerm{Value: "http://example.org/resource1", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "http://example.org/resource2", TermType: RdfIri},
	},

	"_:anon <http://example.org/property> <http://example.org/resource2> .": {
		S: RdfTerm{Value: "anon", TermType: RdfBlank},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "http://example.org/resource2", TermType: RdfIri},
	},

	"<http://example.org/resource1> <http://example.org/property> _:anon .": {
		S: RdfTerm{Value: "http://example.org/resource1", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "anon", TermType: RdfBlank},
	},

	" 	 <http://example.org/resource3> 	 <http://example.org/property>	 <http://example.org/resource2> 	.": {
		S: RdfTerm{Value: "http://example.org/resource3", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "http://example.org/resource2", TermType: RdfIri},
	},

	"<http://example.org/resource7> <http://example.org/property> \"simple literal\" .": {
		S: RdfTerm{Value: "http://example.org/resource7", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "simple literal", TermType: RdfLiteral},
	},

	`<http://example.org/resource8> <http://example.org/property> "backslash:\\" .`: {
		S: RdfTerm{Value: "http://example.org/resource8", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "backslash:\\", TermType: RdfLiteral},
	},

	`<http://example.org/resource9> <http://example.org/property> "dquote:\"" .`: {
		S: RdfTerm{Value: "http://example.org/resource9", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "dquote:\"", TermType: RdfLiteral},
	},

	`<http://example.org/resource10> <http://example.org/property> "newline:\n" .`: {
		S: RdfTerm{Value: "http://example.org/resource10", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "newline:\n", TermType: RdfLiteral},
	},

	`<http://example.org/resource11> <http://example.org/property> "return\r" .`: {
		S: RdfTerm{Value: "http://example.org/resource11", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "return\r", TermType: RdfLiteral},
	},

	`<http://example.org/resource12> <http://example.org/property> "tab:\t" .`: {
		S: RdfTerm{Value: "http://example.org/resource12", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "tab:\t", TermType: RdfLiteral}},

	`<http://example.org/resource16> <http://example.org/property> "\u00E9" .`: {
		S: RdfTerm{Value: "http://example.org/resource16", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "\u00E9", TermType: RdfLiteral},
	},

	`<http://example.org/resource30> <http://example.org/property> "chat"@fr .`: {
		S: RdfTerm{Value: "http://example.org/resource30", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "chat", Language: "fr", TermType: RdfLiteral},
	},

	`<http://example.org/resource31> <http://example.org/property> "chat"@en .`: {
		S: RdfTerm{Value: "http://example.org/resource31", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "chat", Language: "en", TermType: RdfLiteral},
	},

	"# this is a comment \n<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> .": {
		S: RdfTerm{Value: "http://example.org/resource1", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "http://example.org/resource2", TermType: RdfIri},
	},

	"# this is a comment \n   # another comment \n<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> .": {
		S: RdfTerm{Value: "http://example.org/resource1", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "http://example.org/resource2", TermType: RdfIri},
	},

	"<http://example.org/resource7> <http://example.org/property> \"typed literal\"^^<http://example.org/DataType1> .": {
		S: RdfTerm{Value: "http://example.org/resource7", TermType: RdfIri},
		P: RdfTerm{Value: "http://example.org/property", TermType: RdfIri},
		O: RdfTerm{Value: "typed literal", DataType: "http://example.org/DataType1", TermType: RdfLiteral},
	},
}

var negativeCases = map[string]error{
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> ":   ErrUnterminatedTriple,
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> ,":  ErrUnexpectedCharacter,
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> ..": ErrUnexpectedCharacter,
	"http://example.org/resource1> <http://example.org/property> <http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1 <http://example.org/property> <http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1><http://example.org/property> <http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1> <http://example.org/property><http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1> http://example.org/property> <http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1> <http://example.org/property <http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1> <http://example.org/property> http://example.org/resource2>.":    ErrUnexpectedCharacter,
	"<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2.":    ErrUnexpectedEOF,
	"<http://example.org/resource1> \n<http://example.org/property> <http://example.org/resource2>.": ErrUnexpectedCharacter,
	"_:foo\n <http://example.org/property> <http://example.org/resource2>.":                          ErrUnexpectedCharacter,
	"_:0abc <http://example.org/property> <http://example.org/resource2>.":                           ErrUnexpectedCharacter,
	"_abc <http://example.org/property> <http://example.org/resource2>.":                             ErrUnexpectedCharacter,
	"_:a-bc <http://example.org/property> <http://example.org/resource2>.":                           ErrUnexpectedCharacter,
	"_:abc<http://example.org/property> <http://example.org/resource2>.":                             ErrUnexpectedCharacter,
	"_:abc <http://example.org/property> \"foo\"@ .":                                                 ErrUnexpectedCharacter,
	"_:abc <http://example.org/property> \"foo\"^ .":                                                 ErrUnexpectedCharacter,
	"_:abc <http://example.org/property> \"foo\"^^< .":                                               ErrUnexpectedCharacter,
	"_:abc <http://example.org/property> \"foo\"^^<> .":                                              ErrUnexpectedCharacter,
	"_:abc <> _:abc .":                                                                               ErrUnexpectedCharacter,
	"_:abc < > _:abc .":                                                                              ErrUnexpectedCharacter,
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

func TestReadErrors(t *testing.T) {

	for ntriple, expected := range negativeCases {
		r := NewReader(strings.NewReader(ntriple))
		_, err := r.Read()

		if err == nil {
			t.Errorf("Expected %s for %s but no error reported", expected, ntriple)
		} else if err.(*ParseError).Err != expected {
			t.Errorf("Expected %s for %s but got error %s", expected, ntriple, err.(*ParseError).Err)
		}
	}
}
