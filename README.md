# GOSTORE

[![GoDoc](https://godoc.org/github.com/pirmd/libro?status.svg)](https://godoc.org/github.com/pirmd/libro)&nbsp; 
[![Go Report Card](https://goreportcard.com/badge/github.com/pirmd/libro)](https://goreportcard.com/report/github.com/pirmd/libro)&nbsp;

`libro` is a command line tool that provides facilities to manage one or more
collections of ebooks, keeping track of their metadata and/or additional
information the user wants to record.

You can think of `libro` as something close to [beets](http://beets.io/) but
for books.

## INSTALLATION
With golang binary installed on your system, you just need to run:
̀``shell
go install github.com/pirmd/libro
```

## USAGE
Usage with some help guidance can be obtained from `libro`'s command line by running:
``` shell
libro -help
```

`libro` sub-commands are developed so that they can be combined (i.e. piped)
together or with other command-line tools to developed your own books
management workflows. For example, importing books can be run like:
``` shell
libro info "my_book.epub" | vipe --suffix json | libro --root=$HOME/books add
̀̀``

To communicate through pipes, `libro` outputs book's attributes in JSON format
that is understood by each `libro`sub-command.  `libro` detects when output is
for a terminal rather than a pipe and will switch-over to a more human-friendly
format. This behavior can be altered using `-format` flag.

## TEMPLATES
Some `libro` features accept user-defined templates.

Templating language is based on golang's built-in template language
(https://pkg.go.dev/text/template) extended with the following helpers:
- path management:
   * base            : get the path basename
   * ext             : get file extension
   * sep             : return a path separator ('\' of '/')
   * sanitizePath    : replace special chars that don't usually work well when
                       use in path name
   * sanitizeFilename: like sanitizePath but also removes any path separators.
   * nospace         : get rid of spaces
- strings management:
   * join :  join elements separating them wit the given separator.
   * lower:  convert string to lower-case
   * title:  convert string to title-case
   * upper:  convert string to upper-case
- serialization:
   * toJSON      : converts an interface to JSON representation.
   * toPrettyJSON: converts an interface to an easy-to-read JSON representation.
- templating:
  * tmpl  : execute a sub-template by name. Sub-templates are chosen from t namespace.
            Unlike {{template xXx}}, tmpl allows to pipe its result to another command.

Some pre-defined templates are available:
- when renaming books (invoke by `-rename={{template "xXx" .}}`): 
    * default         : Author[0] - [Series SeriesIndex] - Title [LANG].epub
    * default_byauthor: Author[0]/Author[0] - [Series SeriesIndex] - Title [LANG].epub
    * short           : Title.epub
    * short_byauthor  : Author[0]/Title.epub
    * donothing       : keep original book's path
  In each cases, Author[0] is the first known Authors. Empty Authors is
  replaced by 'Unknown" Unknown or empty SeriesxXx, or LANG are ignored.
- when formatting books information (invoke by `-format={{template "xXx" .}}`):
    * plaintext: format book's attributes as plaintext in a key: value fashion.
                 Empty optional attributes are not displayed.

Examples:
``` shell
libro info -format='{{template "plaintext" .}}' my_book.epub
```
or
``` shell
libro add -rename='{{template "short_byauthor" .}} my_book.epub
```
or
``` shell
libro add -rename='{{ tmpl "short" . | nospace }} my_book.epub
```

User-defined templates can be loaded using specific flags like:
``` shell
libro info -format-tmpl=$HOME/books/my_template.gotmpl -format='{{template "my_template.gotmpl" .}}' my_book.epub
```
or
``` shell
libro add -rename-tmpl=$HOME/books/my_template.gotmpl -rename='{{template "my_template.gotmpl" .}} my_book.epub
```

## BEHAVIOR
`libro` adopts some opinionated behavior when processing book's information. Main ones are:
- when exact ISBN match is found online, online information is preferred over
  epub's metadata or over guessed ones;
- when a match is found online but ISBN are not similar, online information is
  only used to complete information obtained from epub's metadata or guessed;
- guessed information are usually not preferred over epub's metadata or online
  information

When editing book's information (using `libro edit`), user is only asked to review information if:
- key attributes are not filled,
- conflicts or dubious automatic operation have been performed.
This behavior can be altered using `-auto` or `-dont-edit` flags.

## BOOK ATTRIBUTES
`libro` uses the following attributes for a Book:
- Path         : Path is the location of the book's file in the file-system.
- Title        : Title is the book's title.
- Authors      : Authors is the list names of the authors and/or editors for this book.
- ISBN         : ISBN is the unique industry standard identifier for this book.
                 `libro` tends to prefer ISBN_13 format when available or when
                 it can be derived from an ISBN_10.
                 ISBN10 and ISBN13 methods can be invoked to convert from one
                 format to the other.
- SubTitle     : SubTitle is the book's sub-title.
- Publisher    : Publisher is the publisher of this book.
- PublishedDate: PublishedDate is the date of publication of this book.
                 `libro` tries to normalize dates using '2006-01-02' format.
                 When 'precision' of date is not enough to capture known month or days, date is
                 cut to '2006-01' or simply to '2006'.
- Description  : Description is the synopsis of the book. The text of the
                 description is formatted in HTML and includes simple
                 formatting elements.
- Series       : Series is the series to which this book belongs to.
- SeriesIndex  : SeriesIndex is the position in the series to which the book
                 belongs to.
- SeriesTitle  : SeriesTitle is the book's title in the series (without Series
                 nor SubTitle information).
- Language     : Language is the book's language. It is the two-letter
                 ISO 639-1 code such as 'fr', 'en'.
- PageCount    : PageCount is total number of pages of this book.
- Subject      : Subject is the list of subject categories, such as "Fiction",
                 "Suspense".


## MAIN GOALS
Beside bug hunting and improved user experience, main functions planned to be
developed (in no special order):
- more scrapers to retrieve metadata from known remote sites (like
  googlebooks);
- offering more record's metadata processing allowing further cleaning and
  quality of collection content; 
- allowing syncing file's embedded metadata with cleaned and completed
  metadata stored in the collection;
- tweak output template to issue static html description of the collection;
- improve batch operation (add several media at a time);
- add book indexing support for getting fancy search features ;
- new media family to be supported (like mp3).

## CONTRIBUTION
If you feel like to contribute, just follow github guidelines on
[forking](https://help.github.com/articles/fork-a-repo/) then [send a pull
request](https://help.github.com/articles/creating-a-pull-request/)


[modeline]: # ( vim: set fenc=utf-8 spell spl=en: )
