# GOSTORE

[![GoDoc](https://godoc.org/github.com/pirmd/libro?status.svg)](https://godoc.org/github.com/pirmd/libro)&nbsp; 
[![Go Report Card](https://goreportcard.com/badge/github.com/pirmd/libro)](https://goreportcard.com/report/github.com/pirmd/libro)&nbsp;

`libro` is a command line tool aiming at providing facilities to manage one
or more collections of ebooks, keeping track of their metadata and/or
additional information the user wants to record.

You can think of `libro` as something close to [beets](http://beets.io/) but
for books.

## USAGE
Usage with some help guidance can be obtained from `libro`'s command line by running:
``` shell
libro -help
```

`libro` sub-commands are developed so that they can be combined (i.e. piped)
together or with other command-line tools to developed your own books
management workflows. For example, importing books can be run like:
``` shell
libro info "my favorite book.epub" | vipe --suffix json | libro --root=$HOME/books add
̀̀``

## INSTALLATION
With golang binary installed on your system, you just need to run:
̀``shell
go install github.com/pirmd/libro
```

## MAIN GOALS
Beside bug hunting and improved user experience, main functions planned to be
developed (in no special order):
    - scrapers to retrieve metadata from known remote sites (like googlebooks);
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
