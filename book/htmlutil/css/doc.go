// Package css - parse CSS
//
// DESCRIPTION
//
// `css` is a library offering basic CSS parsing.
// CSS support is minimal enough to support github.com/pirmd/libro needs,
// essentially identifying possible CSS security issues based on a CSS
// white-list approach.
//
// `css` is based on github.com/gorilla/css CSS scanner.
//
// Even though, I've tried hard enough to stick to CSS standard and detecting
// possible CSS syntax issues, it is not a target at this time.
// You'd probably better go with github.com/aymerick/douceur/. Several test
// cases are borrowed for its implementation so `css` might be close in term of
// parsing "correctness". Besides implementation choices and type format, main
// difference lies in the fact that `css` kept github.com/gorilla/css tokens
// when describing CSS Rules for easier rules inspection.
package css
