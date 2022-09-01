package htmlutil

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/gorilla/css/scanner"
	"github.com/pirmd/libro/book/htmlutil/css"
)

var (
	// reSpaceSeparatedNames matches only names separated by spaces.
	reSpaceSeparatedNames = regexp.MustCompile(`^([\s\p{L}\p{N}_-]+)$`)

	// reSlashSeparatedNames matches only names separated by slash.
	reSlashSeparatedNames = regexp.MustCompile(`^([\p{L}\p{N}-/]+)$`)

	// reConditionalOrSSIComment matches suspicious directives hidden in
	// comments like any SSI or conditional comment.
	reConditionalOrSSIComment = regexp.MustCompile(`\A[\[#]`)

	// reAnonymousHost matches suspicious HOST identified by an IP
	// address only.
	reAnonymousHost = regexp.MustCompile(`^[\\.0-9]+$`)
)

// Scanner represents an HTML/CSS scanner that looks for possible security
// risks.
// This scanner is only oriented to check existing untrusted HTML/CSS in
// EPUB, it does not properly managed all injection cases, notably obfuscated
// strings relying for example on strange characters encodings.
type Scanner struct {
	// AllowedTags is the white-list of allowed tags. For each tag, allowed
	// attributes can be expressed as a pattern:
	//  - **         : all attributes (except data-xxx and onxxx events that
	//                 should be explicitly allowed) are allowed.
	//  - *          : all attributes (except data-xxx and onxxx events that
	//                 should be explicitly allowed) are allowed if their value
	//                 is a space-separated list  of names (letters, numbers, _
	//                 or -).
	//  - a=**       : attribute 'a' whatever its value is.
	//  - a or a=*   : attribute 'a' whose value is a space-separated list of
	//                 names (letters, numbers, _ or -).
	//  - a=_MIME    : attribute 'a' whose value is a mimetype specification
	//                 (names separated by '/' like 'text/css')
	//  - a=key      : attribute 'a' whose value is 'key'.
	//  - a=__URL    : attribute 'a' whose value is an 'allowed' URL.
	//                 An allowed URL is parsable with a scheme matching
	//                 SafeSchemes.
	//                 'Anonymous' hosts (no recorded domain name) are not
	//                 accepted.
	//                 Absolute URL for style-sheets are not accepted.
	//                 Absolute URL with target=_blank but without rel="noopener"
	//                 URL's query are not accepted except if _? suffix
	//                 is added.
	//  - a=__REL_URL: like __URL but only relative URL are accepted.
	//
	// Several patterns can be listed for a given attribute's name, knowing
	// that patterns are checked against in their declaration order (first
	// matching will pass/first non-matching check will fail).
	// LIMITATION: Be extra-careful when using catch-all patterns. For instance
	// {http-equiv=refresh, '*'} as a result will allow any http-equiv to be
	// accepted, so catch-all patterns are actually quite tedious to use.
	// TODO: As off now, it is a "good enough" approach but probably needs further
	// polishing/rework to make something acceptable out of this.
	AllowedTags map[atom.Atom][]string

	// AllowedURLSchemes is the white-list of allowed schemes in URL.
	// "*" allows any schemes.
	AllowedURLSchemes []string

	// CSS property governs allowed styles' declaration, should they be allowed
	// in AllowedTags.

	// AllowAbsoluteURLinCSS, when set to true, accepts using external URL.
	// (by default only relative URL or local URL are considered). Queries are
	// not accepted.
	AllowAbsoluteURLinCSS bool

	// AllowedCSSProperties is the white-list of accepted CSS properties.
	// "*" allows any property, "!xxx" failed immediately for property xxx even
	// if property xxx is allowed afterwards.
	AllowedCSSProperties []string

	// AllowedCSSFunctions is the white-list of accepted CSS functions.
	// "*" allows any functions, "!xxx" failed immediately for function xxx even
	// if function xxx is allowed afterwards.
	AllowedCSSFunctions []string

	// AllowedCSSAtKeywords is the white-list of accepted at-keywords.
	// "*" allows any keywords, "!xxx" failed immediately for keyword xxx even
	// if keyword xxx is allowed afterwards.
	AllowedCSSAtKeywords []string
}

// NewMinimalScanner creates a new scanner that allows only minimal HTML
// features, no CSS nor JS.
// globalAttr are added to the list of allowed attributes of all atom.
func NewMinimalScanner(globalAttr ...string) *Scanner {
	// globalAttr are added at the end of the allowed default attributes list,
	// essentially to ensure that 'catch-all' attributes patterns like '*' or
	// '**' are processed at the end.
	return &Scanner{
		AllowedTags: map[atom.Atom][]string{
			atom.A:          append([]string{"href=__URL"}, globalAttr...),
			atom.B:          globalAttr,
			atom.Big:        globalAttr,
			atom.Blockquote: append([]string{"cite=__URL"}, globalAttr...),
			atom.Body:       globalAttr,
			atom.Br:         globalAttr,
			atom.Caption:    globalAttr,
			atom.Cite:       globalAttr,
			atom.Col:        append([]string{"span"}, globalAttr...),
			atom.Colgroup:   append([]string{"span"}, globalAttr...),
			atom.Dd:         globalAttr,
			atom.Del:        append([]string{"cite=__URL"}, globalAttr...),
			atom.Dfn:        globalAttr,
			atom.Div:        globalAttr,
			atom.Em:         globalAttr,
			atom.Figure:     globalAttr,
			atom.H1:         globalAttr,
			atom.H2:         globalAttr,
			atom.H3:         globalAttr,
			atom.H4:         globalAttr,
			atom.H5:         globalAttr,
			atom.H6:         globalAttr,
			atom.Head:       {},
			atom.Hr:         globalAttr,
			atom.Html:       append([]string{"xml:lang", "xmlns=__URL"}, globalAttr...),
			atom.I:          globalAttr,
			atom.Img:        append([]string{"height", "width", "src=__URL", "alt=*"}, globalAttr...),
			atom.Ins:        append([]string{"cite=__URL"}, globalAttr...),
			atom.Li:         globalAttr,
			// TODO: for Link using type=_MIME might mask type precision in globalAttr (like type=text/css)
			atom.Link:    append([]string{"href=__URL", "type=_MIME", "rel"}, globalAttr...),
			atom.Meta:    append([]string{"http-equiv=content-type", "charset", "name", "content=**"}, globalAttr...),
			atom.Ol:      globalAttr,
			atom.P:       globalAttr,
			atom.Q:       append([]string{"cite=__URL"}, globalAttr...),
			atom.S:       globalAttr,
			atom.Section: globalAttr,
			atom.Small:   globalAttr,
			atom.Span:    globalAttr,
			atom.Strong:  globalAttr,
			atom.Sub:     globalAttr,
			atom.Sup:     globalAttr,
			atom.Table:   append([]string{"summary=**"}, globalAttr...),
			atom.Tbody:   globalAttr,
			atom.Td:      append([]string{"colspan", "rowspan"}, globalAttr...),
			atom.Tfoot:   globalAttr,
			atom.Th:      append([]string{"abbr", "colspan", "rowspan"}, globalAttr...),
			atom.Thead:   globalAttr,
			atom.Title:   globalAttr,
			atom.Tr:      globalAttr,
			atom.U:       globalAttr,
			atom.Ul:      globalAttr,
		},

		AllowedURLSchemes: []string{"http", "https"},
	}
}

// NewScannerWithStyle creates a new scanner that extends NewStrictScanner to
// allow use of common CSS properties (but not CSS functions or external CSS
// resources).
// If globalAttr are provided, these attributes will be allowed for each atom
// in addition to "style" and "class" attributes that are accepted by default.
func NewScannerWithStyle(globalAttr ...string) *Scanner {
	s := NewMinimalScanner(append(globalAttr, "class", "style")...)
	s.AllowedTags[atom.Style] = []string{"type=text/css"}
	s.AllowedCSSProperties = []string{
		"background", "background-color", "border", "border-bottom",
		"border-collapse", "border-color", "border-radius", "border-style",
		"border-width", "clear", "color", "cursor", "direction", "display",
		"flex", "float", "font", "font-family", "font-size", "font-style",
		"font-variant", "font-weight", "grid", "height", "left",
		"letter-spacing", "line-height", "list-style", "margin",
		"margin-bottom", "margin-left", "margin-right", "margin-top",
		"max-height", "max-width", "min-height", "min-width", "overflow",
		"overflow-x", "overflow-y", "padding", "padding-bottom",
		"padding-left", "padding-right", "padding-top", "page-beak-after",
		"page-break-before", "position", "right", "src", "text-align",
		"table-layout", "text-decoration", "text-indent", "top",
		"vertical-align", "visibility", "white-space", "width", "word-spacing",
		"z-index", "zoom",
	}

	return s
}

// NewPermissiveScanner creates a new scanner that allows any attributes or CSS
// properties but fetching external resources/URL.
func NewPermissiveScanner() *Scanner {
	s := NewMinimalScanner("title=**", "*")

	// Allow style
	s.AllowedTags[atom.Style] = []string{"type=text/css"}
	s.AllowedCSSProperties = []string{"*"}
	s.AllowedCSSAtKeywords = []string{"!@import", "*"}

	// pattern definition for attributes is quite tricky, we need to
	// specifically redefine Meta expectation, otherwise the use of catch-all
	// pattern will allow any http-equiv to be accepted
	s.AllowedTags[atom.Meta] = []string{
		"http-equiv=content-type",
		"http-equiv=Content-Type",
		"charset",
		"name",
		"content=**",
	}

	return s
}

// Scan checks that io.Reader contains only allowed tags or attributes.
// Scan returns a list of messages describing encountered issues.
func (s *Scanner) Scan(r io.Reader) ([]string, error) {
	var issues []string
	reportIssue := func(format string, a ...interface{}) {
		issues = append(issues, fmt.Sprintf(format, a...))
	}

	var inStyleNode bool

	tokenizer := html.NewTokenizer(r)
	for {
		if tokenizer.Next() == html.ErrorToken {
			if err := tokenizer.Err(); err != nil {
				if err == io.EOF {
					// Tokenizer seems to simply ignore bad formatted HTML.
					// Some XSS examples seems to consider that it can be a way
					// to abuse parsers, so here we try to detect such
					// situations.
					if notparsed := tokenizer.Raw(); len(notparsed) > 0 {
						reportIssue("Unparsed HTML found: %s", string(notparsed))
					}
					return issues, nil
				}
				reportIssue("Parsing error: %v", err)
				return issues, err
			}
		}

		switch token := tokenizer.Token(); token.Type {
		case html.DoctypeToken:

		case html.CommentToken:
			if reConditionalOrSSIComment.MatchString(token.Data) {
				reportIssue("Suspicious directive hidden in a comment: %s", token.Data)
			}

		case html.StartTagToken, html.SelfClosingTagToken:
			if token.DataAtom == 0 {
				reportIssue("Tag '%s' is unknown", token.Data)
				continue
			}

			inStyleNode = (token.DataAtom == atom.Style)

			if _, isAllowed := s.AllowedTags[token.DataAtom]; !isAllowed {
				reportIssue("Tag '%s' is not allowed", token.Data)
				continue
			}

			for _, attr := range token.Attr {
				for _, issue := range s.inspectAttr(token, attr) {
					reportIssue("%s=%s: %s", attr.Key, attr.Val, issue)
				}
			}

		case html.EndTagToken:
			if token.DataAtom == 0 {
				reportIssue("Tag '%s' is unknown", token.Data)
				continue
			}

			if token.DataAtom == atom.Style {
				inStyleNode = false
			}

			// EndTagToken has no attributes (ignored by tokenizer) but it
			// seems that some XSS can abuse that so we check that it is really
			// the case (example from bluemonday's test cases).
			if len(token.String()) != len(string(tokenizer.Raw())) {
				reportIssue("Closing tag seems to contain unexpected data: %s", string(tokenizer.Raw()))
			}

		case html.TextToken:
			if inStyleNode {
				cssIssues, err := s.inspectCSS(token.Data)
				if err != nil {
					reportIssue("fail to inspect CSS declaration '%s': %v", token.Data, err)
				} else if len(cssIssues) > 0 {
					issues = append(issues, cssIssues...)
				}

			}

		default:
			reportIssue("Unknown token: %v", string(tokenizer.Raw()))
		}
	}
}

// ScanCSS checks that io.Reader contains only allowed CSS style declarations.
// ScanCSS returns a list of messages describing encountered issues.
func (s *Scanner) ScanCSS(r io.Reader) ([]string, error) {
	cssTxt := new(strings.Builder)
	_, err := io.Copy(cssTxt, r)
	if err != nil {
		return nil, err
	}

	return s.inspectCSS(cssTxt.String())
}

func (s *Scanner) inspectAttr(inTag html.Token, attr html.Attribute) (issues []string) {
	// It seems that tokenizer does not detect properly empty attribute values, so
	// this workaround might be better than nothing
	if p := strings.SplitN(attr.Val, "=", 2); len(p) == 2 {
		if a := strings.ToLower(p[0]); atom.Lookup([]byte(a)) != 0 {
			issues = append(issues, "value seems (possibly purposely) empty and parser might misuse next attr")
			return
		}
	}

	for _, pattern := range s.AllowedTags[inTag.DataAtom] {
		switch p := strings.SplitN(pattern, "=", 2); len(p) {
		case 1:
			switch {
			case (pattern == "**" || pattern == "*" || pattern == "style") && attr.Key == "style":
				cssIssues, err := s.inspectInlineCSS(attr.Val)
				if err != nil {
					issues = append(issues, fmt.Sprintf("inline CSS parsing error: %v", err))
					return
				}

				issues = append(issues, cssIssues...)
				return

			case attr.Key == pattern:
				if !reSpaceSeparatedNames.MatchString(attr.Val) {
					issues = append(issues, "only space-separated names are allowed")
				}
				return

			case strings.HasPrefix(attr.Key, "on") || strings.HasPrefix(attr.Key, "data"):
				// '*' or '**' cannot trigger onxxx or dataxxx

			case pattern == "*":
				if !reSpaceSeparatedNames.MatchString(attr.Val) {
					issues = append(issues, "only space-separated names are allowed")
				}
				return

			case pattern == "**":
				return
			}

		case 2:
			if attr.Key != p[0] {
				continue
			}

			switch p[1] {
			case "**": // Attribute value is accepted whatever its value is
				return

			case "*": // Attribute value is accepted for space-separated names
				if !reSpaceSeparatedNames.MatchString(attr.Val) {
					issues = []string{"only space-separated names are allowed"}
				}
				return

			case "_MIME": // Attribute value is accepted for slash-separated names (like text/css)
				if !reSlashSeparatedNames.MatchString(attr.Val) {
					issues = []string{"only slash-separated names are allowed"}
				}
				return

			case "__URL", "__URL_?":
				if issue := s.inspectURL(inTag, attr, true, p[1] == "__URL_?"); issue != "" {
					issues = []string{issue}
				}
				return

			case "__REL_URL", "__REL_URL_?":
				if issue := s.inspectURL(inTag, attr, false, p[1] == "__REL_URL_?"); issue != "" {
					issues = []string{issue}
				}
				return

			case attr.Val: // we have an exact match
				return
			}
		}
	}

	issues = []string{"attribute is not an allowed pattern"}
	return
}

func (s *Scanner) inspectURL(inTag html.Token, attr html.Attribute, allowAbsoluteURL bool, allowURLQuery bool) (issue string) {
	u, err := url.Parse(attr.Val)
	if err != nil {
		return fmt.Sprintf("%s=%s: non-parsable url are not allowed in '%s' (%v)", attr.Key, attr.Val, inTag.DataAtom, err)
	}

	if u.Scheme != "" && !isInList(u.Scheme, s.AllowedURLSchemes) {
		return fmt.Sprintf("%s=%s: url scheme '%s' is not allowed in '%s'", attr.Key, attr.Val, u.Scheme, inTag.DataAtom)
	}

	if !allowURLQuery && len(u.RawQuery) > 0 {
		return fmt.Sprintf("%s=%s: url query '%s' is not allowed in '%s'", attr.Key, attr.Val, u.RawQuery, inTag.DataAtom)
	}

	if u.IsAbs() {
		if !allowAbsoluteURL {
			return fmt.Sprintf("%s=%s: absolute URL are not allowed in '%s'", attr.Key, attr.Val, inTag.DataAtom)
		}

		if reAnonymousHost.MatchString(u.Host) {
			return fmt.Sprintf("%s=%s: anonymous hosts are not allowed in '%s'", attr.Key, attr.Val, inTag.DataAtom)
		}

		var rel, target string
		for _, a := range inTag.Attr {
			if a.Key == "rel" {
				rel = a.Val
			}
			if a.Key == "target" {
				target = a.Val
			}
		}

		if strings.Contains(rel, "stylesheet") {
			return fmt.Sprintf("%s=%s with rel=stylesheet: external stylesheets are not allowed in '%s'", attr.Key, attr.Val, inTag.DataAtom)
		}

		if strings.Contains(target, "_blank") && !strings.Contains(rel, "noopener") {
			return fmt.Sprintf("%s=%s with target=_blank: external link are not allowed without rel=noopener in '%s'", attr.Key, attr.Val, inTag.DataAtom)
		}
	}
	return
}

func (s *Scanner) inspectCSS(cssTxt string) (issues []string, err error) {
	ruleset, err := css.Parse(cssTxt)
	if err != nil {
		return nil, err
	}

	for _, rule := range ruleset {
		issues = append(issues, s.inspectCSSRule(rule)...)
	}

	return issues, nil
}

func (s *Scanner) inspectInlineCSS(cssTxt string) (issues []string, err error) {
	ruleset, err := css.ParseInline(cssTxt)
	if err != nil {
		return nil, err
	}

	for _, rule := range ruleset {
		issues = append(issues, s.inspectCSSRule(rule)...)
	}

	return issues, nil
}

func (s *Scanner) inspectCSSRule(rule *css.Rule) (issues []string) {
	if rule.AtKeyword != nil {
		if !isInList(rule.AtKeyword.Value, s.AllowedCSSAtKeywords) {
			return []string{fmt.Sprintf("%s is not an allowed at-keyword", rule.AtKeyword.Value)}
		}
	}

	for _, val := range rule.Selectors {
		if issue := s.inspectCSSValue(val); issue != "" {
			return append(issues, issue)
		}
	}

	for _, decl := range rule.Declarations {
		if !isInList(decl.Property, s.AllowedCSSProperties) {
			return append(issues, fmt.Sprintf("%s is not an allowed CSS property", decl.Property))
		}
		if issue := s.inspectCSSValue(decl.Value); issue != "" {
			return append(issues, issue)
		}
	}

	for _, r := range rule.EmbeddedRuleset {
		issues = append(issues, s.inspectCSSRule(r)...)
	}

	return issues
}

func (s *Scanner) inspectCSSValue(val css.Value) string {
	for _, tok := range val {
		switch tok.Type {
		case scanner.TokenURI:
			if issue := s.inspectCSSURL(tok.Value, s.AllowAbsoluteURLinCSS, false); issue != "" {
				return issue
			}

		case scanner.TokenFunction:
			if !isInList(strings.TrimSuffix(tok.Value, "("), s.AllowedCSSFunctions) {
				return fmt.Sprintf("%s function's call is not allowed", tok.Value)
			}
		}
	}

	return ""
}

func (s *Scanner) inspectCSSURL(val string, allowAbsoluteURL bool, allowURLQuery bool) (issue string) {
	u, err := url.Parse(val)
	if err != nil {
		return fmt.Sprintf("non-parsable URL: %v", err)
	}

	if u.Scheme != "" && !isInList(u.Scheme, s.AllowedURLSchemes) {
		return fmt.Sprintf("URL scheme '%s' is not allowed in style declaration", u.Scheme)
	}

	if !allowURLQuery && len(u.RawQuery) > 0 {
		return fmt.Sprintf("URL with query '%s' is not allowed in style declaration", u.RawQuery)
	}

	if u.IsAbs() {
		if !allowAbsoluteURL {
			return "absolute URL are not allowed in style declaration"
		}

		if reAnonymousHost.MatchString(u.Host) {
			return "anonymous hosts are not allowed in style declaration"
		}
	}

	return
}

func isInList(s string, allowed []string) bool {
	for _, a := range allowed {
		switch {
		case a == "*":
			return true
		case s == a:
			return true
		case len(a) > 0 && (a[0] == '!' && s == a[1:]):
			return false
		}
	}

	return false
}
