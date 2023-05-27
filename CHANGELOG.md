# Changelog
## [0.4.2] - not released
- fix Language cleaning function that wasn't working.
- update Report by downgrading PublishedDate's issue to a warning.
- add a new PublishedDate format ("2006-01-02T15:04:05")
- add a new Series format ("[<Series>-<SeriesIndex>] <SeriesTitle>")
- add a cleaner step for inappropriate Language ("un")
- add new cleaners for Authors (clean UPPER-CASES names, detect name/surname
  inversion) 
- add extended guessers to extract Book's information from filenames.
- add new Book's Title cleaners.

## [0.4.1] - 2023-05-19
- fix lacking book.Path information after edition.
- add PublishedYear attribute and use it in fullname.gotmpl.
- add new guesser regexp for (incomplete) Series information in book's Title.
- add verification that given PublishedDate makes sense before adopting it.
- add naive Language normalization logic.
- add an intermediate issue report level (warning) to help further limit manual
  intervention.
- tweak guessers' regexp to capture more cases.

## [0.4.0] - 2022-09-02
- add string edit distances to improve detection of inconsistencies tracking.
- add cleaning of Book's Description to remove HTML tags for easier reading in
  text mode.
- add a new libro 'check' command to review Book's information completeness.
  libro 'edit' command with flag `-auto` does not check quality anymore.
- add an EPUB conformity check relying on w3.org epubcheck tool.
- add initial support for checking EPUB content for possible HTML/CSS security
  issues.

## [0.3.0] - 2022-06-25
- add a guesser that looks for ISBN numbers within the EPUB's content.
- use SimilarBooks information to improve end-user edition of book's attributes.
- allow to cancel edit operation by emptying the edited book's attributes list.
- extend guessers to support cleaners that can replace existing attributes with
  their 'cleaned' version.
 
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
