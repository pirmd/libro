package css

import (
	"github.com/gorilla/css/scanner"
)

func str2token(txt string) (tok []*scanner.Token) {
	s := scanner.New(txt)
	for {
		switch token := s.Next(); token.Type {
		case scanner.TokenError:
			panic("fail to parse " + txt)
		case scanner.TokenEOF:
			return
		default:
			tok = append(tok, token)
		}
	}
}

func str2val(txt string) Value {
	tok := str2token(txt)

	val := make(Value, 0)
	for _, t := range tok {
		val.Append(t)
	}

	return val.Scrub()
}

func areSameToken(tok1, tok2 *scanner.Token) bool {
	if tok1 == nil || tok2 == nil {
		return (tok1 == nil) && (tok2 == nil)
	}

	return tok1.Type == tok2.Type && tok1.Value == tok2.Value
}

func areSameSelectors(sel1, sel2 []Value) bool {
	if len(sel1) != len(sel2) {
		return false
	}

	for i := range sel1 {
		//TODO: can do better than that?
		if len(sel1[i]) != len(sel2[i]) || sel1[i].String() != sel2[i].String() {
			return false
		}
	}

	return true
}

func areSameDeclaration(decl1, decl2 *Declaration) bool {
	return (decl1.Property == decl2.Property) &&
		areSameValue(decl1.Value, decl2.Value) &&
		(decl1.IsImportant == decl2.IsImportant)
}

func areSameDeclarations(decl1, decl2 []*Declaration) bool {
	if len(decl1) != len(decl2) {
		return false
	}

	for i := range decl1 {
		if !areSameDeclaration(decl1[i], decl2[i]) {
			return false
		}
	}

	return true
}

func areSameValue(v1, v2 Value) bool {
	if (v1 == nil) || (v2 == nil) {
		return (v1 == nil) && (v2 == nil)
	}

	if len(v1) != len(v2) {
		return false
	}
	// TODO: can do better than that?
	return (v1.String() == v2.String())
}

func areSameRule(r1, r2 *Rule) bool {
	if (r1 == nil) || (r2 == nil) {
		return (r1 == nil) && (r2 == nil)
	}

	return areSameSelectors(r1.Selectors, r2.Selectors) &&
		areSameDeclarations(r1.Declarations, r2.Declarations) &&
		areSameToken(r1.AtKeyword, r2.AtKeyword) &&
		areSameRuleset(r1.EmbeddedRuleset, r2.EmbeddedRuleset)
}

func areSameRuleset(r1, r2 Ruleset) bool {
	if (r1 == nil) || (r2 == nil) {
		return (r1 == nil) && (r2 == nil)
	}

	if len(r1) != len(r2) {
		return false
	}

	for i := range r1 {
		if !areSameRule(r1[i], r2[i]) {
			return false
		}
	}

	return true
}
