package css

import (
	"fmt"
	"strings"

	"github.com/gorilla/css/scanner"
)

// TODO: merge Value.Append() and Value.Scrub()?

const (
	// tab is used to indent nested rules, declarations and cie
	tab = "    "
)

// Ruleset represents a series of rules.
// NOTA: Ruleset is slightly different from
// https://www.w3.org/People/howcome/foo.dir/grammar.html definition as it
// extend the notion to also capture at-rules nested rules.
type Ruleset []*Rule

func (rules Ruleset) String() string {
	var s []string
	for _, r := range rules {
		s = append(s, r.String())
	}

	return strings.Join(s, "\n")
}

// Rule represents a CSS rule
type Rule struct {
	//TODO: find a better name for Selector that does fit qualified rule case
	//but not really at-rule cases

	// Selectors contains the list of tokens of pre-Block declaration that is to say:
	//  - list of selectors for qualified rules,
	//  - list of at-rule specific properties (anything before a ';' or '{'),
	// Selectors are split during parsing operation using ',' as separator,
	Selectors []Value

	// Declarations contains the list of style declarations for the given Rule.
	Declarations []*Declaration

	// AtKeyword is the at-keyword introducing teh Rule. If the Rule is a
	// qualified rule, AtKeyword is nil.
	AtKeyword *scanner.Token

	// EmbeddedRuleset contains the nested rules for at-rules that expects it,
	// otherwise is nil.
	EmbeddedRuleset Ruleset
}

func (r Rule) String() string {
	switch {
	case r.AtKeyword == nil:
		return r.sprintQualifiedRule()

	case expectNestedRulesetAfter(r.AtKeyword):
		return r.sprintAtRuleWithNestedRules()

	case expectRulesetAfter(r.AtKeyword):
		return r.sprintAtRuleWithRuleset()

	case expectSelectorAfter(r.AtKeyword):
		return r.sprintAtRuleWithSelector()

	default:
		return fmt.Sprintf("!Err(unknown at-rule %#v)", r)
	}
}

func (r Rule) sprintQualifiedRule() string {
	var s strings.Builder

	for i, sel := range r.Selectors {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(sel.String())
	}

	if len(r.Declarations) != 0 {
		s.WriteString(" {\n")
		for i, decl := range r.Declarations {
			if i > 0 {
				s.WriteString("\n")
			}
			s.WriteString(tab)
			s.WriteString(decl.String())
			s.WriteString(";")
		}
		s.WriteString("\n}")
	}

	return s.String()
}

func (r Rule) sprintAtRuleWithNestedRules() string {
	var s strings.Builder

	s.WriteString(r.AtKeyword.Value)
	if len(r.Selectors) > 0 {
		s.WriteString(" ")
	}

	for i, sel := range r.Selectors {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(sel.String())
	}

	if len(r.EmbeddedRuleset) != 0 {
		s.WriteString(" {\n")
		s.WriteString(tab)
		ns := r.EmbeddedRuleset.String()
		ns = strings.ReplaceAll(ns, "\n", "\n"+tab)
		ns = strings.TrimSuffix(ns, tab) // previous Replace might add unneeded ending spaces
		s.WriteString(ns)
		s.WriteString("\n}")
	}

	return s.String()
}

func (r Rule) sprintAtRuleWithRuleset() string {
	var s strings.Builder

	s.WriteString(r.AtKeyword.Value)
	if len(r.Selectors) > 0 {
		s.WriteString(" ")
	}

	s.WriteString(r.sprintQualifiedRule())

	return s.String()
}

func (r Rule) sprintAtRuleWithSelector() string {
	var s strings.Builder

	s.WriteString(r.AtKeyword.Value)
	if len(r.Selectors) > 0 {
		s.WriteString(" ")
	}

	for i, sel := range r.Selectors {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(sel.String())
	}

	return s.String()
}

// expectNestedRuleset is true it the given at-keywords requires a
// set of nested-rules.
func expectNestedRulesetAfter(atkeyword *scanner.Token) bool {
	// acc. to https://developer.mozilla.org/en-US/docs/Web/CSS/At-rule
	return atkeyword.Value == "@media" || atkeyword.Value == "@supports" ||
		atkeyword.Value == "@document" || atkeyword.Value == "@keyframes" ||
		atkeyword.Value == "@font-feature-values"
}

// expectRuleset is true it the given at-keywords requires a set of rules.
func expectRulesetAfter(atkeyword *scanner.Token) bool {
	// acc. to https://developer.mozilla.org/en-US/docs/Web/CSS/At-rule
	return atkeyword.Value == "@page" || atkeyword.Value == "@font-face" || atkeyword.Value == "@viewport" ||
		atkeyword.Value == "@counter-style" || atkeyword.Value == "@property" || atkeyword.Value == "@color-profile" ||
		// only in a @font-feature-values acc. to https://developer.mozilla.org/en-US/docs/Web/CSS/At-rule
		atkeyword.Value == "@swash" || atkeyword.Value == "@annotation" || atkeyword.Value == "@ornaments" ||
		atkeyword.Value == "@stylistic" || atkeyword.Value == "@styleset" || atkeyword.Value == "@character-variant"
}

// expectSelectors is true it the given at-keywords requires a simple property
// specifications.
func expectSelectorAfter(atkeyword *scanner.Token) bool {
	// acc. to https://developer.mozilla.org/en-US/docs/Web/CSS/At-rule
	return atkeyword.Value == "@charset" ||
		atkeyword.Value == "@import" ||
		atkeyword.Value == "@namespace"
}

// Declaration represents a style definition.
type Declaration struct {
	// Property is the style property's name.
	Property string

	// Value is the style property's value as a list of token.
	Value Value

	// IsImportant is true if the property is tagged as '!important'.
	IsImportant bool

	// mightBeImportant is an internal flag triggered when a "!" is encountered
	// to detect !important flag
	mightBeImportant *scanner.Token
}

func (d *Declaration) String() string {
	if d.IsImportant {
		return fmt.Sprintf("%s: %s !important", d.Property, d.Value)
	}
	return fmt.Sprintf("%s: %s", d.Property, d.Value)
}

// AppendToValue adds a new token to a property. It makes sure that
// '!important' tag are captured.
func (d *Declaration) AppendToValue(tok *scanner.Token) error {
	if tok == nil {
		return nil
	}

	if d.IsImportant {
		return fmt.Errorf("token '%s' added to declaration's Value after !important tag", tok)
	}

	if d.mightBeImportant != nil {
		if tok.Type == scanner.TokenIdent && strings.ToLower(tok.Value) == "important" {
			d.IsImportant, d.mightBeImportant = true, nil
			return nil
		}

		d.Value.Append(d.mightBeImportant)
	}

	if tok.Type == scanner.TokenChar && tok.Value == "!" {
		d.mightBeImportant = tok
	} else {
		d.mightBeImportant = nil
		d.Value.Append(tok)
	}

	return nil
}

// Value represents a tokenized element (without comments).
type Value []*scanner.Token

// Append adds a new token to a Value.
func (val *Value) Append(tok *scanner.Token) {
	if tok == nil {
		return
	}
	*val = append(*val, tok)
}

func (val Value) String() string {
	var s strings.Builder

	for _, v := range val {
		s.WriteString(v.Value)
	}

	return s.String()
}

// Scrub is called once all token have been scanned for a given Value so
// that the Value's content is cleaned from:
//  - Trailing or Tailing token of type scanner.TokenS;
//  - Token of type scanner.TokenComment;
//  - "inline" comments are detected and Value is rescanned without them to
//  address of tricky injections like "expr/*XSS*/ession(alert('XSS'))".
func (val Value) Scrub() Value {
	var clean Value
	var needRescan bool // true if inline comment is detected to rescan the Value

	for _, tok := range val {
		if tok.Type == scanner.TokenComment {
			// TODO: Can maybe limit rescan to cases where previous token is an
			// IDENT, but I cannot be sure at this point, better safe than
			// sorry. Need more investigations.
			needRescan = true
			continue
		}

		// Trailing TokenS
		if tok.Type == scanner.TokenS {
			if len(clean) == 0 {
				continue
			}
		}

		clean = append(clean, tok)
	}

	if len(clean) == 0 {
		return clean
	}

	// Trailing spaces
	for clean[len(clean)-1].Type == scanner.TokenS {
		clean = clean[:len(clean)-1]
	}

	if needRescan {
		clean = clean.rescan()
	}

	return clean
}

func (val Value) rescan() Value {
	var newVal Value

	s := scanner.New(val.String())
	for {
		switch token := s.Next(); token.Type {
		case scanner.TokenError:
			panic("fail to parse again value " + val.String())
		case scanner.TokenEOF:
			return newVal
		default:
			newVal.Append(token)
		}
	}
}
