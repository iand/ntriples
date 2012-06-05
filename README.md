ntriples - a basic ntriples parser in Go

Notes
-----

This parser is incomplete. Known limitations:

 * Does not parse literals with a datatype (datatype is not parsed)
 * No checking of IRI syntax
 * No checking of valid bnode labels
 * It's likely that malformed triples will be parsed without error
 * No comment parsing
