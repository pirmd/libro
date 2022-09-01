package css

import (
	"testing"
)

// Some tests borrowed from:
// - https://www.w3.org/Style/CSS/Test/CSS3/Selectors/current/html/index.html
// - https://github.com/aymerick/douceur/blob/master/parser/parser_test.go

func TestParseForSelectors(t *testing.T) {
	testCases := []struct {
		in        string
		want      []Value
		shallFail bool
	}{
		{
			in:   `li,p { background-color : lime }`,
			want: []Value{str2val("li"), str2val("p")},
		},

		{
			in:   `>#foo { background-color : lime }`,
			want: []Value{str2val(">#foo")},
		},
		{
			in:   `>* { color : lime }`,
			want: []Value{str2val(">*")},
		},
		{
			in:   `p.t1.t2 { background: green; color: white; }`,
			want: []Value{str2val("p.t1.t2")},
		},
		{
			in:   `>button:active { background: green; color: white; }`,
			want: []Value{str2val(">button:active")},
		},
		{
			in:   `#fail#test { background: red; color: yellow; }`,
			want: []Value{str2val("#fail#test")},
		},
		{
			in:   `div.stub *:not([title^="si on"]) { color : lime }`,
			want: []Value{str2val("div.stub *:not([title^=\"si on\"])")},
		},
		{
			in:   `blockquote + div ~ p { color: green; }`,
			want: []Value{str2val("blockquote + div ~ p")},
		},
		{
			in:   `foo & address, p { background: red; }`,
			want: []Value{str2val("foo & address"), str2val("p")},
		},
		{
			in:   `.\13 { color: red; }`,
			want: []Value{str2val(".\\13 ")},
		},
		{
			in:   `.foo..quux { background: red; color: yellow; }`,
			want: []Value{str2val(".foo..quux")},
		},
		{
			in:   `p:not([class^=""]) { color: lime; }`,
			want: []Value{str2val("p:not([class^=\"\"])")},
		},
		{
			in:   `table[class="body"] .footer__content td { color: lime; }`,
			want: []Value{str2val("table[class=\"body\"] .footer__content td")},
		},
	}

	for _, tc := range testCases {
		p := newParser(tc.in)

		got, err := p.parseForSelector()
		if err != nil {
			if !tc.shallFail {
				t.Errorf("fail to parse selectors for '%s': %v", tc.in, err)
			}
			continue
		}

		if tc.shallFail {
			t.Errorf("parse selectors for '%s' does not fail", tc.in)
			continue
		}

		if !areSameSelectors(tc.want, got) {
			t.Errorf("fail to parse selectors for '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}

func TestParseForDeclaration(t *testing.T) {
	testCases := []struct {
		in        string
		want      []*Declaration
		shallFail bool
	}{
		{
			in: `{ background-color : lime }`,
			want: []*Declaration{
				{Property: "background-color", Value: str2val("lime")},
			},
		},
		{
			in: `{
                         padding: 2px;
                         box-sizing: border-box;
                         background-color: #4CAF50;
                         /*margin:5px;*/
                       }`,
			want: []*Declaration{
				{Property: "padding", Value: str2val("2px")},
				{Property: "box-sizing", Value: str2val("border-box")},
				{Property: "background-color", Value: str2val("#4CAF50")},
			},
		},
		{
			in: `{
                           border: 2px solid black !important;
                           font: 14px/22px normal helvetica, sans-serif;
                       }`,
			want: []*Declaration{
				{Property: "border", Value: str2val("2px solid black"), IsImportant: true},
				{Property: "font", Value: str2val("14px/22px normal helvetica, sans-serif")},
			},
		},
		{
			in: `{
                           color: rgba(244, 145, 14, 0.80); /* bright orange */
                        }`,
			want: []*Declaration{
				{Property: "color", Value: str2val("rgba(244, 145, 14, 0.80)")},
			},
		},
		{
			in: `{
                            display: block;
                            background-image: url(images/closequote1.gif);
                            background-repeat: no-repeat;
                            background-position: bottom right;
                         }`,
			want: []*Declaration{
				{Property: "display", Value: str2val("block")},
				{Property: "background-image", Value: str2val("url(images/closequote1.gif)")},
				{Property: "background-repeat", Value: str2val("no-repeat")},
				{Property: "background-position", Value: str2val("bottom right")},
			},
		},
		{
			in: `{ font-family: "Source Sans Pro", Arial, sans-serif; font-size: 27px; }`,
			want: []*Declaration{
				{Property: "font-family", Value: str2val("\"Source Sans Pro\", Arial, sans-serif")},
				{Property: "font-size", Value: str2val("27px")},
			},
		},
	}

	for _, tc := range testCases {
		p := newParser(tc.in)

		got, err := p.parseForDeclarations()
		if err != nil {
			if tc.shallFail {
				continue
			}
			t.Errorf("fail to parse block's declarations for '%s': %v", tc.in, err)
		}

		if tc.shallFail {
			t.Errorf("parse block's declarations for '%s' does not fail", tc.in)
			continue
		}

		if !areSameDeclarations(tc.want, got) {
			t.Errorf("fail to parse block's declarations for '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}

func TestParseForRuleset(t *testing.T) {
	testCases := []struct {
		in        string
		want      Ruleset
		shallFail bool
	}{
		{
			in: `
                p { background: green; color: white; }
                .t1.fail { background: red; color: yellow; }
                .fail.t1 { background: red; color: yellow; }
                .t2.fail { background: red; color: yellow; }
                .fail.t2 { background: red; color: yellow; }
                /* Note: This is a valid test even per CSS1, since in CSS1 those rules
                         are invalid and should be dropped. */
            `,
			want: Ruleset{
				&Rule{
					Selectors: []Value{str2val("p")},
					Declarations: []*Declaration{
						{Property: "background", Value: str2val("green")},
						{Property: "color", Value: str2val("white")},
					},
				},
				&Rule{
					Selectors: []Value{str2val(".t1.fail")},
					Declarations: []*Declaration{
						{Property: "background", Value: str2val("red")},
						{Property: "color", Value: str2val("yellow")},
					},
				},
				&Rule{
					Selectors: []Value{str2val(".fail.t1")},
					Declarations: []*Declaration{
						{Property: "background", Value: str2val("red")},
						{Property: "color", Value: str2val("yellow")},
					},
				},
				&Rule{
					Selectors: []Value{str2val(".t2.fail")},
					Declarations: []*Declaration{
						{Property: "background", Value: str2val("red")},
						{Property: "color", Value: str2val("yellow")},
					},
				},
				&Rule{
					Selectors: []Value{str2val(".fail.t2")},
					Declarations: []*Declaration{
						{Property: "background", Value: str2val("red")},
						{Property: "color", Value: str2val("yellow")},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		p := newParser(tc.in)

		got, err := p.parseForRuleset()
		if err != nil {
			if tc.shallFail {
				continue
			}
			t.Errorf("fail to parse '%s': %v", tc.in, err)
		}

		if tc.shallFail {
			t.Errorf("parse '%s' does not fail", tc.in)
			continue
		}

		if !areSameRuleset(tc.want, got) {
			t.Errorf("fail to parse '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}

func TestParseForAtRule(t *testing.T) {
	testCases := []struct {
		in        string
		want      *Rule
		shallFail bool
	}{
		{
			in: `@charset "UTF-8";`,
			want: &Rule{
				AtKeyword: str2token("@charset")[0],
				Selectors: []Value{str2val("\"UTF-8\"")},
			},
		},

		{
			in: `@counter-style footnote {
                             system: symbolic;
                             symbols: '*' ⁑ † ‡;
                             suffix: '';
                        }`,
			want: &Rule{
				AtKeyword: str2token("@counter-style")[0],
				Selectors: []Value{str2val("footnote")},
				Declarations: []*Declaration{
					{Property: "system", Value: str2val("symbolic")},
					{Property: "symbols", Value: str2val("'*' ⁑ † ‡")},
					{Property: "suffix", Value: str2val("''")},
				},
			},
		},

		{
			in: `@document url(http://www.w3.org/),
            url-prefix(http://www.w3.org/Style/),
            domain(mozilla.org),
            regexp("https:.*")
            {
                /* CSS rules here apply to:
                + The page "http://www.w3.org/".
                + Any page whose URL begins with "http://www.w3.org/Style/"
                + Any page whose URL's host is "mozilla.org" or ends with ".mozilla.org"
                + Any page whose URL starts with "https:" */
                /* make the above-mentioned pages really ugly */
                body { color: purple; background: yellow; }
            }`,
			want: &Rule{
				AtKeyword: str2token("@document")[0],
				Selectors: []Value{str2val("url(http://www.w3.org/)"), str2val("url-prefix(http://www.w3.org/Style/)"), str2val("domain(mozilla.org)"), str2val("regexp(\"https:.*\")")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						Selectors: []Value{str2val("body")},
						Declarations: []*Declaration{
							{Property: "color", Value: str2val("purple")},
							{Property: "background", Value: str2val("yellow")},
						},
					},
				},
			},
		},

		{
			in: `@font-face {
                  font-family: MyHelvetica;
                  src: local("Helvetica Neue Bold"), local("HelveticaNeue-Bold"), url(MgOpenModernaBold.ttf);
                  font-weight: bold;
              }`,
			want: &Rule{
				AtKeyword: str2token("@font-face")[0],
				Declarations: []*Declaration{
					{Property: "font-family", Value: str2val("MyHelvetica")},
					{Property: "src", Value: str2val("local(\"Helvetica Neue Bold\"), local(\"HelveticaNeue-Bold\"), url(MgOpenModernaBold.ttf)")},
					{Property: "font-weight", Value: str2val("bold")},
				},
			},
		},

		{
			in: `@font-feature-values Font Two { /* How to activate nice-style in Font Two */
                            @styleset {
                                nice-style: 4;
                            }
                        }`,
			want: &Rule{
				AtKeyword: str2token("@font-feature-values")[0],
				Selectors: []Value{str2val("Font Two")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						AtKeyword: str2token("@styleset")[0],
						Declarations: []*Declaration{
							{Property: "nice-style", Value: str2val("4")},
						},
					},
				},
			},
		},

		{
			in: `@import "my-styles.css";`,
			want: &Rule{
				AtKeyword: str2token("@import")[0],
				Selectors: []Value{str2val("\"my-styles.css\"")},
			},
		},
		{
			in: `@import url('landscape.css') screen and (orientation:landscape);`,

			want: &Rule{
				AtKeyword: str2token("@import")[0],
				Selectors: []Value{str2val("url('landscape.css') screen and (orientation:landscape)")},
			},
		},

		{
			in: `@keyframes identifier {
                            0% { top: 0; left: 0; }
                            100% { top: 100px; left: 100%; }
                        }`,
			want: &Rule{
				AtKeyword: str2token("@keyframes")[0],
				Selectors: []Value{str2val("identifier")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						Selectors: []Value{str2val("0%")},
						Declarations: []*Declaration{
							{Property: "top", Value: str2val("0")},
							{Property: "left", Value: str2val("0")},
						},
					},
					&Rule{
						Selectors: []Value{str2val("100%")},
						Declarations: []*Declaration{
							{Property: "top", Value: str2val("100px")},
							{Property: "left", Value: str2val("100%")},
						},
					},
				},
			},
		},

		{
			in: `@media screen, print {
                          body { line-height: 1.2 }
                        }`,
			want: &Rule{
				AtKeyword: str2token("@media")[0],
				Selectors: []Value{str2val("screen"), str2val("print")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						Selectors: []Value{str2val("body")},
						Declarations: []*Declaration{
							{Property: "line-height", Value: str2val("1.2")},
						},
					},
				},
			},
		},

		{
			in: `@namespace svg url(http://www.w3.org/2000/svg);`,
			want: &Rule{
				AtKeyword: str2token("@namespace")[0],
				Selectors: []Value{str2val("svg url(http://www.w3.org/2000/svg)")},
			},
		},

		{
			in: `@page :left {
                              margin-left: 4cm;
                              margin-right: 3cm;
                        }`,
			want: &Rule{
				AtKeyword: str2token("@page")[0],
				Selectors: []Value{str2val(":left")},
				Declarations: []*Declaration{
					{Property: "margin-left", Value: str2val("4cm")},
					{Property: "margin-right", Value: str2val("3cm")},
				},
			},
		},

		{
			in: `@supports (animation-name: test) {
                            /* specific CSS applied when animations are supported unprefixed */
                            @keyframes { /* @supports being a CSS conditional group at-rule, it can includes other relevent at-rules */
                                  0% { top: 0; left: 0; }
                                  100% { top: 100px; left: 100%; }
                             }
                        }`,
			want: &Rule{
				AtKeyword: str2token("@supports")[0],
				Selectors: []Value{str2val("(animation-name: test)")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						AtKeyword: str2token("@keyframes")[0],
						EmbeddedRuleset: Ruleset{
							&Rule{
								Selectors: []Value{str2val("0%")},
								Declarations: []*Declaration{
									{Property: "top", Value: str2val("0")},
									{Property: "left", Value: str2val("0")},
								},
							},
							&Rule{
								Selectors: []Value{str2val("100%")},
								Declarations: []*Declaration{
									{Property: "top", Value: str2val("100px")},
									{Property: "left", Value: str2val("100%")},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		p := newParser(tc.in)

		got, err := p.parseForAtRule()
		if err != nil {
			if tc.shallFail {
				continue
			}
			t.Errorf("fail to parse '%s': %v", tc.in, err)
		}

		if tc.shallFail {
			t.Errorf("parse '%s' does not fail", tc.in)
			continue
		}

		if !areSameRule(tc.want, got) {
			t.Errorf("fail to parse '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}
