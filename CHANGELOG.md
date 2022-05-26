# Changelog
## [0.3.0]
- use SimilarBooks information to improve end-user edition of book's attributes.
- allow to cancel edit operation by emptying the edited book's attributes list.
- extend guessers to support so-called cleaners that can replace existing
  attributes (with their cleaned version).
 
## [0.2.1] - 2022-05-19
- change Books' metadata completion logic by avoiding using googlebooks search
  results if ISBN identifiers are not matching. Libro prefers in this case to
  rely on the end-user choice.
- add a reporting function that draws end-user attention on possible issues
  encountered during Book's processing.
- add support to correct reversed Authors name.
- improve templates support.
- improve documentation.
- simplify code's repository.

## [0.2.0] - 2022-04-20
- add basic Googlebooks API support.
- improve ISBN and Date understanding and handling.

## [0.1.0] - 2022-03-29
- create as a stripped down version of github.com/pirmd/gostore version 0.7.0.


[modeline]: # ( vim: set fenc=utf-8 spell spl=en: )
