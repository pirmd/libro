// libro - manage a collection of ebooks.
//
// Description
//
// libro is a command line utility that provides facilities for managing ebooks
// collections: finding information about a book, storing with books according
// to a given naming scheme, searching the book collection...
//
// libro can output book's information using a user-supplied template.
// Templating language is based on golang's built-in template language
// (https://pkg.go.dev/text/template) extended with the following helpers: -
// toJSON      : converts an interface to JSON representation.  - toPrettyJSON:
// converts an interface to an easy-to-read JSON representation.
//
// By default, libro either pretty-prints an ebook's metadata set in a
// human readable key/value form or in its json format depending if libro
// output is a terminal or not (e.g. when piping result to another command).
package main
