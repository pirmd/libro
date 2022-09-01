package htmlutil

import (
	"bytes"
	"io"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// GetRawTextFromHTML extracts text from an HTML document without retaining any
// particular formatting information.
//
// Limitation: GetRawTextFromHTML is only a minimal and naive HTML to text
// extractor, it does not consider any fancy HTML formatting directive nor
// complicated rules related to spaces collapsing and only concentrate into
// getting rid of HTML directive to access meaningful information for scraping
// or searching.
func GetRawTextFromHTML(r io.Reader) (io.Reader, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	txt := getRawTextFromNode(root)
	txt = bytes.TrimSpace(txt)

	return bytes.NewBuffer(txt), nil
}

// getRawTextFromNode returns text from a node and all it's descendant.
func getRawTextFromNode(root *html.Node) []byte {
	txt := []byte{}

	if root.Type == html.TextNode {
		// As a simplification we systematically get rid of '\n' as even if it
		// will hurt formatting will not alter 'meaning' of the extracted text.
		txt = bytes.ReplaceAll([]byte(root.Data), []byte{'\n'}, []byte{' '})
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		switch node.Type {
		case html.ElementNode:
			switch node.DataAtom {
			// Block Elements, extended to atom.Tr, atom.Th and atom.Td.
			case atom.Address, atom.Article, atom.Aside, atom.Blockquote, atom.Canvas,
				atom.Dd, atom.Div, atom.Dl, atom.Dt, atom.Fieldset, atom.Figcaption,
				atom.Figure, atom.Footer, atom.Form, atom.H1, atom.H2, atom.H3, atom.H4,
				atom.H5, atom.H6, atom.Header, atom.Hr, atom.Li, atom.Main, atom.Nav,
				atom.Noscript, atom.Ol, atom.P, atom.Pre, atom.Section, atom.Table,
				atom.Tfoot, atom.Ul, atom.Video, atom.Tr, atom.Th, atom.Td:
				// For Blocks, we retain a really simple formatting approach:
				// trim redundant spaces at the boundaries between child and
				// parent, then we add new line before and after.
				childTxt := getRawTextFromNode(node)
				childTxt = trimLeadingSpace(childTxt)
				txt = trimTrailingSpace(txt)
				txt = append(txt, '\n')
				txt = append(txt, childTxt...)
				txt = append(txt, '\n')
				// Possible 'redundant' trailing newlines should be eliminated
				// afterwards, but even if we have too much of them it will
				// alter formatting but not 'meaning' of text.

				//TODO: tries to extract `alt` attributes content?

			case atom.Script:
				//ignore

				// Mainly Inline elements, maybe more but hopefully not a problem for our use-case
			default:
				// For inline elements, we simply either keep a '\n' from a
				// previous block element, or we make sure that only one space
				// is added (if asked).
				childTxt := getRawTextFromNode(node)
				if endWithEOL(txt) {
					txt = append(trimTrailingSpace(txt), '\n')
					txt = append(txt, trimLeadingSpace(childTxt)...)
				} else if startWithSpace(childTxt) || endWithSpace(txt) {
					txt = append(trimTrailingSpace(txt), ' ')
					txt = append(txt, trimLeadingSpace(childTxt)...)
				} else {
					txt = append(txt, childTxt...)
				}
			}

		case html.CommentNode:
			//ignore

		case html.TextNode:
			// a TextNode is processed just like an inline element
			childTxt := getRawTextFromNode(node)
			if endWithEOL(txt) {
				txt = append(trimTrailingSpace(txt), '\n')
				txt = append(txt, trimLeadingSpace(childTxt)...)
			} else if startWithSpace(childTxt) || endWithSpace(txt) {
				txt = append(trimTrailingSpace(txt), ' ')
				txt = append(txt, trimLeadingSpace(childTxt)...)
			} else {
				txt = append(txt, childTxt...)
			}

		default:
			// ignore
			// TODO: check if something useful can be found in remaining NodeType?
		}
	}

	return txt
}

func trimTrailingSpace(s []byte) []byte {
	return bytes.TrimFunc(s, unicode.IsSpace)
}

func trimLeadingSpace(s []byte) []byte {
	return bytes.TrimFunc(s, unicode.IsSpace)
}

func startWithSpace(s []byte) bool {
	if len(s) == 0 {
		return false
	}

	return unicode.IsSpace(rune(s[0]))
}

func endWithSpace(s []byte) bool {
	if len(s) == 0 {
		return false
	}

	return unicode.IsSpace(rune(s[len(s)-1]))
}

func endWithEOL(s []byte) bool {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			return true
		}

		if !unicode.IsSpace(rune(s[i])) {
			break
		}
	}

	return false
}
