/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Source Code file.
  This work is published from the United Kingdom. 
*/
package ntriples

import (
  "bufio"
  "bytes"
  "errors"
  "fmt"
  "io"
  "unicode"
)



// A ParseError is returned for parsing errors.
// The first line is 1.  The first column is 0.
type ParseError struct {
	Line   int   // Line where the error occurred
	Column int   // Column (rune index) where the error occurred
	Err    error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Err)
}

// These are the errors that can be returned in ParseError.Error
var (
	ErrUnexpectedCharacter         = errors.New("unexpected character")
	ErrUnexpectedEOF         = errors.New("unexpected end of file")
	ErrTermCount    = errors.New("wrong number of terms in line")
	ErrUnterminatedIri    = errors.New("unterminated IRI, expecting '>'")
	ErrUnterminatedLiteral    = errors.New("unterminated Literal, expecting '\"'")
)

type Reader struct {
	line             int
	column           int
	r                *bufio.Reader
	buf             bytes.Buffer
}


type Triple struct {
	s RdfTerm
	p RdfTerm
	o RdfTerm
}

const (
	RdfUnknown = iota
	RdfIri
	RdfBlank
	RdfLiteral
)


type RdfTerm struct { 
	value string
	language string
	datatype string
	termtype int
}

func (t Triple) String() string {
	var s, p, o string

	switch t.s.termtype {
	case RdfIri:
		s = fmt.Sprintf("<%s>", t.s.value)

	case RdfBlank:
		s = fmt.Sprintf("_:%s", t.s.value)
	}

	p = fmt.Sprintf("<%s>", t.p.value)

	switch t.o.termtype {
	case RdfIri:
		o = fmt.Sprintf("<%s>", t.o.value)

	case RdfBlank:
		o = fmt.Sprintf("_:%s", t.o.value)
	case RdfLiteral:
		if t.o.language != "" {
			o = fmt.Sprintf("\"%s\"@%s", t.o.value, t.o.language)
		} else if t.o.datatype != "" {
			o = fmt.Sprintf("\"%s\"^^<%s>", t.o.value, t.o.datatype)
		} else {
			o = fmt.Sprintf("\"%s\"", t.o.value)
		}
	}

	return fmt.Sprintf("%s %s %s .", s, p, o)
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:     bufio.NewReader(r),
	}
}


// error creates a new ParseError based on err.
func (r *Reader) error(err error) error {
	return &ParseError{
		Line:   r.line,
		Column: r.column,
		Err:    err,
	}
}

func (r *Reader) Read() (t Triple, e error) {
	r.line++
	r.column=-1

	_, _, err := r.r.ReadRune()
   	if err != nil {
   		return Triple{}, err
	}

	// if r.Comment != 0 && r1 == r.Comment {
	// 	return nil, r.skip('\n')
	// }

	r.r.UnreadRune()

	termCount := 0
	for {
		haveTerm, term, err := r.parseTerm()
		if haveTerm {
			termCount++
			switch termCount {
			case 1:
				t.s = term
			case 2:
				t.p = term
			case 3:
				t.o = term
		
				err = r.readToEOL()
				if err != nil {
					println("FOO: ", err.Error())

					return Triple{}, err
				}

				return t, nil
			default:
				// TODO: error, too many terms
				return Triple{}, r.error(ErrTermCount)
			}


			
		}
		if err != nil {
			return Triple{}, err
		} 
	}
	panic("unreachable")	




	return Triple{}, nil

}

// readRune reads one rune from r, folding \r\n to \n and keeping track
// of how far into the line we have read.  r.column will point to the start
// of this rune, not the end of this rune.
func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()

	// Handle \r\n here.  We make the simplifying assumption that
	// anytime \r is followed by \n that it can be folded to \n.
	// We will not detect files which contain both \r\n and bare \n.
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	r.column++
	return r1, err
}


// unreadRune puts the last rune read from r back.
func (r *Reader) unreadRune() {
	r.r.UnreadRune()
	r.column--
}

func (r *Reader) parseTerm() (haveField bool, term RdfTerm, err error) {
	r.buf.Reset()
	
	r1, err := r.readRune()
   	if err != nil {
		return false, term, err
	}

	// Skip whitespace
	for r1 != '\n' && unicode.IsSpace(r1) {
		r1, err = r.readRune()
		if err != nil {
			return false, term, err
		}
	}

	switch r1 {
	case '<':
		// Read an IRI
		for {
			r1, err = r.readRune()
   			if err != nil {
   				if err == io.EOF {
   					return false, term, r.error(ErrUnexpectedEOF)
   				}
   				return false, term, err
   			}
			switch r1 {
				case '>':
					return true, RdfTerm{value:r.buf.String(), termtype:RdfIri}, nil
			}
			r.buf.WriteRune(r1)
		}
	case '_':
		// Read a blank node
		r1, err = r.readRune()
		if err != nil {
			if err == io.EOF {
				return false, term, r.error(ErrUnexpectedEOF)
			}
			return false, term, err
		}
		if r1 != ':' {
			return false, term, r.error(ErrUnexpectedCharacter)
		}
		for {
			r1, err = r.readRune()
   			if err != nil {
   				if err == io.EOF {
   					return false, term, r.error(ErrUnexpectedEOF)
   				}
   				return false, term, err
   			}
			if r1 == '.' || unicode.IsSpace(r1) {
				return true, RdfTerm{value:r.buf.String(), termtype:RdfBlank}, nil
			}
			r.buf.WriteRune(r1)
		}

	case '"':
		// Read a literal
		for {
			r1, err = r.readRune()
   			if err != nil {
   				if err == io.EOF {
   					return false, term, r.error(ErrUnexpectedEOF)
   				}
   				return false, term, err
   			}
			switch r1 {
				case '"':
					r1, err = r.readRune()
	   				if err == io.EOF {
	   					return false, term, r.error(ErrUnexpectedEOF)
	   				}
	   				if r1 == '.' || unicode.IsSpace(r1) {
	   					r.unreadRune()
						return true, RdfTerm{value:r.buf.String(), termtype:RdfLiteral}, nil
	   				}
	   				if r1 == '@' {
	   					tmpterm := RdfTerm{value:r.buf.String(), termtype:RdfLiteral}
	   					r.buf.Reset()

	   					for {
							r1, err = r.readRune()
				   			if err != nil {
				   				if err == io.EOF {
				   					return false, term, r.error(ErrUnexpectedEOF)
				   				}
				   				return false, term, err
				   			}
	   						if r1 == '.' || unicode.IsSpace(r1) {
	   							if r.buf.Len() == 0 {
									return false, term, r.error(ErrUnexpectedCharacter)
	   							}
	   							tmpterm.language = r.buf.String()
				   				return true, tmpterm, nil
	   						}
	   						if r1 == '-' || (r1 >= 'a' && r1 <= 'z') || (r1 >= '0' && r1 <= '9') {
								r.buf.WriteRune(r1)
							} else {
								return false, term, r.error(ErrUnexpectedCharacter)
							}
	   					}


	   				}
					return false, term, r.error(ErrUnexpectedCharacter)



				case '\\':
				r1, err = r.readRune()
	   			if err != nil {
	   				if err == io.EOF {
	   					return false, term, r.error(ErrUnexpectedEOF)
	   				}
	   				return false, term, err
	   			}
				switch r1 {
				case '\\', '"':
				case 't':
					r1 = '\t'
				case 'r':
					r1 = '\r'
				case 'n':
					r1 = '\n'
				case 'u', 'U':

					codepoint := rune(0)

					for i:=3; i >= 0; i-- {
						r1, err = r.readRune()
			   				
			   			if err != nil {
			   				if err == io.EOF {
			   					return false, term, r.error(ErrUnexpectedEOF)
			   				}
			   				return false, term, err
			   			}

			   			if r1 >= '0' && r1 <= '9' {
			   				codepoint += (1 << uint32(4*i)) * (r1 - '0')
			   			} else if r1 >= 'a' && r1 <= 'f' {
			   				codepoint += (1 << uint32(4*i)) * (r1 - 'a' + 10)
			   			} else if r1 >= 'A' && r1 <= 'F' {
			   				codepoint += (1 << uint32(4*i)) * (r1 - 'A' + 10)
			   			} else {
							return false, term, r.error(ErrUnexpectedCharacter)
			   			}

					}
					r1 = codepoint

				default:
	   					return false, term, r.error(ErrUnexpectedCharacter)
				}
			}
			r.buf.WriteRune(r1)
		}
	default:
		// TODO: raise error, unexpected character
		return false, term, r.error(ErrUnexpectedCharacter)

	}

	panic("unreachable")	

}


func (r *Reader) readToEOL() (err error) {
	r1, err := r.readRune()
   	if err != nil {
		if err == io.EOF {
			return r.error(ErrUnexpectedEOF)
		}
		return err
	}

	// Skip whitespace
	for unicode.IsSpace(r1) {
		r1, err = r.readRune()
		if err != nil {
			if err == io.EOF {
				return r.error(ErrUnexpectedEOF)
			}
			return err
		}
	}

	if r1 != '.' {
		return r.error(ErrUnexpectedCharacter)
	}

	r1, err = r.readRune()
   	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	// Skip whitespace
	for unicode.IsSpace(r1) {
		if r1 == '\n' {
			return nil
		}
		r1, err = r.readRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	if r1 != '\n' {
		println("FOOOO", r1)
		return r.error(ErrUnexpectedCharacter)
	}

	return nil

}