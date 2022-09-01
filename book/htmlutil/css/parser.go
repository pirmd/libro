package css

import (
	"fmt"

	"github.com/gorilla/css/scanner"
)

// Specification from https://www.w3.org/People/howcome/foo.dir/cover.html

type parser struct {
	scanner *scanner.Scanner

	// token is the current token
	token *scanner.Token
}

func newParser(css string) *parser {
	p := &parser{
		scanner: scanner.New(css),
	}
	p.nextToken() // move to the first token
	return p
}

func (p *parser) nextToken() {
	p.token = p.scanner.Next()
}

// Parse builds the list of corresponding Ruleset contains in the parser's CSS
// text.
func (p *parser) Parse() (Ruleset, error) {
	ruleset, err := p.parseForRuleset()
	if err != nil {
		return nil, err
	}

	if p.token.Type != scanner.TokenEOF {
		return nil, fmt.Errorf("'%s' if an unexpected style-sheet end", p.token)
	}

	return ruleset, nil
}

// parseForRuleset parses a CSS text looking for a list of qualified or at-rules.
// parseForRuleset ends if 'EOF' is encountered but also if '}' is encountered.
// For the later case, it is up to the caller to decide if it is acceptable.
func (p *parser) parseForRuleset() (Ruleset, error) {
	var ruleset Ruleset

	for {
		switch p.token.Type {
		case scanner.TokenError:
			return nil, fmt.Errorf("parsing error at %s", p.token)

		case scanner.TokenEOF:
			return ruleset, nil

		case scanner.TokenCDO, scanner.TokenCDC, scanner.TokenComment, scanner.TokenS:

		case scanner.TokenAtKeyword:
			rule, err := p.parseForAtRule()
			if err != nil {
				return nil, err
			}
			ruleset = append(ruleset, rule)

		case scanner.TokenChar:
			if p.token.Value == "}" {
				return ruleset, nil
			}
			fallthrough // other TokenChar have no special behaviour, handle as any other token

		default:
			rule, err := p.parseForQualifiedRule()
			if err != nil {
				return nil, err
			}
			ruleset = append(ruleset, rule)
		}

		p.nextToken()
	}
}

// parseForAtRule parses CSS text for an at-rule.
func (p *parser) parseForAtRule() (*Rule, error) {
	if p.token.Type != scanner.TokenAtKeyword {
		panic("parseForAtRule called for a non at-keyword token: " + p.token.String())
	}

	switch {
	case expectNestedRulesetAfter(p.token): // @media, @document, @keyframes...
		return p.parseForAtRuleWithNestedRuleset()

	case expectRulesetAfter(p.token): // @page, @font-faces,...
		return p.parseForAtRuleWithRuleset()

	case expectSelectorAfter(p.token): // @charset, @importe...
		return p.parseForAtRuleWithSelector()

	default:
		return nil, fmt.Errorf("at-keyword '%s' is unknown", p.token)
	}
}

func (p *parser) parseForAtRuleWithNestedRuleset() (*Rule, error) {
	atKeyword := p.token
	p.nextToken()

	sel, err := p.parseForSelector()
	if err != nil {
		return nil, err
	}

	if p.token.Value != "{" {
		return nil, fmt.Errorf("at-rule '%s' needs a block defining a set of rules", atKeyword)
	}
	p.nextToken()

	rules, err := p.parseForRuleset()
	if err != nil {
		return nil, err
	}

	if p.token.Value != "}" {
		return nil, fmt.Errorf("at-rule '%s''s nested rules block does not terminate with a '}'", atKeyword)
	}
	p.nextToken()

	return &Rule{
		AtKeyword:       atKeyword,
		Selectors:       sel,
		EmbeddedRuleset: rules,
	}, nil

}

func (p *parser) parseForAtRuleWithRuleset() (*Rule, error) {
	atKeyword := p.token
	p.nextToken()

	rule, err := p.parseForQualifiedRule()
	if err != nil {
		return nil, err
	}

	return &Rule{
		AtKeyword:    atKeyword,
		Selectors:    rule.Selectors,
		Declarations: rule.Declarations,
	}, nil
}

func (p *parser) parseForAtRuleWithSelector() (*Rule, error) {
	atKeyword := p.token
	p.nextToken()

	sel, err := p.parseForSelector()
	if err != nil {
		return nil, err
	}

	if p.token.Value != ";" {
		return nil, fmt.Errorf("at-rule '%s' selector not ending with ';' but with '%s'", atKeyword.Value, p.token.Value)
	}
	p.nextToken()

	return &Rule{
		AtKeyword: atKeyword,
		Selectors: sel,
	}, nil
}

// parseForQualifiedRule parses CSS text for a qualified rule.
func (p *parser) parseForQualifiedRule() (*Rule, error) {
	rule := new(Rule)

	sel, err := p.parseForSelector()
	if err != nil {
		return nil, err
	}
	rule.Selectors = sel

	decl, err := p.parseForDeclarations()
	if err != nil {
		return nil, err
	}
	rule.Declarations = decl

	return rule, nil
}

// parseForSelector parses a selector (before reaching an '{' as well as an
// at-rule 'selector' (before a ; or a {)
// TODO: it may be wiser to separate them into two functions (on for at-rule,
// one for qualified rule).
// TODO: if separate, maybe some token are not allowed/expected to be
// encountered in a qualified rule selector but can be in an at-rule situation.
func (p *parser) parseForSelector() ([]Value, error) {
	var selectors []Value
	var selector Value

	for {
		switch p.token.Type {
		case scanner.TokenError:
			return nil, fmt.Errorf("parsing error at %s", p.token)

		case scanner.TokenEOF:
			return nil, fmt.Errorf("reach EOF before end of block's selectors declaration")

		case scanner.TokenChar:
			switch p.token.Value {
			case "{", ";":
				if selector = selector.Scrub(); len(selector) == 0 {
					return selectors, nil
				}
				return append(selectors, selector), nil

			case ",":
				if selector = selector.Scrub(); len(selector) == 0 {
					return nil, fmt.Errorf("find an empty selector in a list of selectors")
				}
				selectors = append(selectors, selector)
				selector = nil

			default:
				// other TokenChar have no special behaviour, handle as any other token
				selector.Append(p.token)
			}

		default:
			selector.Append(p.token)
		}

		p.nextToken()
	}
}

// parseForDeclarations parses CSS text for a list of style delcarations.
func (p *parser) parseForDeclarations() ([]*Declaration, error) {
	var declarations []*Declaration
	var declaration *Declaration

	if p.token.Type != scanner.TokenChar || p.token.Value != "{" {
		return nil, fmt.Errorf("block of declarations does not start with a '{' (%s)", p.token)
	}
	p.nextToken()

	for {
		switch p.token.Type {
		case scanner.TokenError:
			return nil, fmt.Errorf("parsing error at %s", p.token)

		case scanner.TokenEOF:
			return nil, fmt.Errorf("reach EOF before end of block declaration")

		case scanner.TokenIdent:
			switch {
			case declaration == nil || declaration.Property == "":
				declaration = &Declaration{Property: p.token.Value}
			case declaration != nil && declaration.Property != "" && declaration.Value != nil:
				if err := declaration.AppendToValue(p.token); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("token '%s' found in an unexpected location (outside of property or value)", p.token)
			}

		case scanner.TokenChar:
			switch p.token.Value {
			case "{":
				return nil, fmt.Errorf("encounter a '{' in a block that does not support embedded blocks")

			case "}": // end of block
				if declaration == nil {
					p.nextToken()
					return declarations, nil
				}

				if declaration.Value == nil {
					return nil, fmt.Errorf("end-of-block without a fully defined declaration")
				}

				if declaration.Value = declaration.Value.Scrub(); len(declaration.Value) == 0 {
					return nil, fmt.Errorf("end-of-block  and declaration without an empty property value")
				}

				declarations = append(declarations, declaration)
				p.nextToken()
				return declarations, nil

			case ";": // end of declaration
				if declaration.Property == "" || declaration.Value == nil {
					return nil, fmt.Errorf("end-of-declaration without a defined declaration")
				}
				if declaration.Value = declaration.Value.Scrub(); len(declaration.Value) == 0 {
					return nil, fmt.Errorf("end-of-declaration  without an empty value")
				}
				declarations = append(declarations, declaration)
				declaration = nil // reset declaration parsing

			case ":": // end of property
				if declaration == nil || declaration.Property == "" {
					return nil, fmt.Errorf("find end-of-property without a defined property")
				}
				declaration.Value = make(Value, 0)

			default:
				// other TokenChar have no special behaviour, handle as any other token
				if declaration == nil || declaration.Property == "" || declaration.Value == nil {
					return nil, fmt.Errorf("'%s' non-authorized token in property declaration", p.token)
				}
				if err := declaration.AppendToValue(p.token); err != nil {
					return nil, err
				}
			}

		case scanner.TokenNumber, scanner.TokenString, scanner.TokenHash, scanner.TokenFunction,
			scanner.TokenUnicodeRange, scanner.TokenURI, scanner.TokenDimension, scanner.TokenPercentage:

			if declaration == nil || declaration.Property == "" || declaration.Value == nil {
				return nil, fmt.Errorf("'%s' non-authorized token in property declaration", p.token)
			}
			if err := declaration.AppendToValue(p.token); err != nil {
				return nil, err
			}

		case scanner.TokenS, scanner.TokenComment:
			if declaration != nil && declaration.Value != nil {
				declaration.AppendToValue(p.token)
			}

		default:
			return nil, fmt.Errorf("non-recognized token '%s' in a block declaration", p.token)
		}

		p.nextToken()
	}
}

// Parse parses CSS text for its list of Rules.
func Parse(s string) (Ruleset, error) {
	return newParser(s).Parse()
}

// ParseInline parses inline CSS text for its list of Rules.
func ParseInline(s string) (Ruleset, error) {
	//TODO: add check for consistency with inline CSS (not at-keyword, no
	//selector, no embedded-rule)
	return newParser("{" + s + "}").Parse()
}
