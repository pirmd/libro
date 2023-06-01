// libro - manage a collection of ebooks.
//
// # DESCRIPTION
//
// `libro` is a command line tool that provides facilities to manage one or more
// collections of ebooks, keeping track of their metadata and/or additional
// information the user wants to record.
// You can think of `libro` as something close to [beets](http://beets.io/) but
// for books.
//
// `libro` sub-commands are developed so that they can be combined (i.e. piped)
// together or with other command-line tools to developed your own books
// management workflows. For example, importing books can be run like:
//
//	libro info "my_book.epub" | vipe --suffix json | libro insert --root=$HOME/books
//
// Or
//
//	libro info -use-guesser -use-googlebooks "my_book.epub" \
//	  | libro edit -auto -default Language=fr -editor=vimdiff \
//	  | libro check -conformity \
//	  | libro insert -root="$HOME/books -rename='{{template "shortname_byauthor.gotmpl" .}}
//
// To communicate through pipes, `libro` outputs book's attributes in JSON
// format that is understood by each `libro`sub-command.
// `libro` detects when output is for a terminal rather than a pipe and will
// switch-over to a more human-friendly format. This behavior can be altered
// using `-format` flag.
//
// # TEMPLATES
//
// Some `libro` features accept user-defined templates.
//
// Templating language is based on golang's built-in template language
// (https://pkg.go.dev/text/template) extended with the following helpers:
//   - path management:
//   - base            : get the path basename
//   - ext             : get file extension
//   - sep             : return a path separator ('\' of '/')
//   - sanitizePath    : replace special chars that don't usually work well
//     when use in path name
//   - sanitizeFilename: like sanitizePath but also removes any path separators.
//   - nospace         : get rid of spaces
//   - strings management:
//   - join :  join elements separating them wit the given separator.
//   - lower:  convert string to lower-case
//   - title:  convert string to title-case
//   - upper:  convert string to upper-case
//   - serialization:
//   - toJSON      : converts an interface to JSON representation.
//   - toPrettyJSON: converts an interface to an easy-to-read JSON representation.
//   - templating:
//   - tmpl  : execute a sub-template by name. Sub-templates are chosen from t namespace.
//     Unlike {{template xXx}}, tmpl allows to pipe its result to another command.
//
// Some pre-defined templates are available:
//
//   - when renaming books (invoke by `-rename={{template "xXx" .}}`):
//
//   - fullname.gotmpl          : Author[0] - [Series SeriesIndex] - Title [LANG].epub
//
//   - fullname_byauthor.gotmpl : Author[0]/Author[0] - [Series SeriesIndex] - Title [LANG].epub
//
//   - shortname.gotmpl         : Title.epub
//
//   - shortname_byauthor.gotmpl: Author[0]/Title.epub
//     In each cases, Author[0] is the first known Authors. Empty Authors is
//     replaced by 'Unknown". Unknown or empty SeriesxXx, or LANG are ignored.
//
//   - when formatting books information (invoke by `-format={{template "xXx" .}}`):
//
//   - book.txt.gotmpl: format book's attributes as plain-text in a key: value
//     fashion. Empty optional attributes are not displayed.
//
// Examples:
//
//	libro info -format='{{toPrettyJSON .}}' my_book.epub
//
// or
//
//	libro info "my_book.epub" | libro insert -rename='{{template "shortname_byauthor.gotmpl" .}}
//
// or
//
//	libro info "my_book.epub" | libro insert -rename='{{ tmpl "shortname.gotmpl" . | nospace }}
//
// User-defined templates can be loaded using specific flags like:
//
//	libro info -format-tmpl=$HOME/books/my_template.gotmpl -format='{{template "my_template.gotmpl" .}}' "my_book.epub"
//
// or
//
//	libro info "my_book.epub" | libro insert -rename-tmpl=$HOME/books/my_template.gotmpl -rename='{{template "my_template.gotmpl" .}}
//
// # BEHAVIOR
//
// `libro` adopts some opinionated behavior when processing book's information.
// Main ones are:
//   - when exact ISBN match is found online, online information is preferred
//     over EPUB's metadata or over guessed ones;
//   - when a match is found online but ISBN are not similar, online information
//     is only used to complete information obtained from EPUB's metadata or
//     guessed;
//   - guessed information are usually not preferred over EPUB's metadata or
//     online information
//
// When editing book's information (using `libro edit`), user is only asked to review information if:
//   - key attributes are not filled,
//   - conflicts or dubious automatic operation have been performed.
//
// This behavior can be altered using `-auto` or `-dont-edit` flags.
//
// # GUESSERS
//
// `libro` can run guessers to complete (and/or confirm) Book's metadata. Current guessers are:
//   - guess Title, Series, Authors or Language from Book's filename,
//   - guess Series information from Book's Title or SubTitle,
//   - guess ISBN by extracting it from the EPUB's content.
//     Use of guessers is governed by the `-use-guesser` flag of `libro info` sub-command.
//
// # CHECKER
//
// `libro` can run different check to verify quality, completeness or conformity of
// information collected about an EPUB or of the EPUB's itself. Findings requiring
// end-user attention are inserted into a specific book's attributes ('Issues')
// for later processing.
//
// `libro` relies on [EPUBcheck](https://www.w3.org/publishing/epubcheck/) tool
// for conformity verification.
package main
