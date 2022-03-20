package libro

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	// FilepathFuncMap provides standard functions to build/manipulate path in
	// template.Funcmap's format:
	//  - ext         : get file extension
	//  - sanitizePath: replace special chars that don't usually work well when
	//                  use in path name
	//  - nospace     : get rid of spaces
	FilepathFuncMap = template.FuncMap{
		"base":             filepath.Base,
		"ext":              filepath.Ext,
		"sanitizePath":     pathSanitizer,
		"sanitizeFilename": filenameSanitizer,
		"nospace":          nospaceSanitizer,
	}
)

// TmplFuncMap provides standard functions to execute go text templates path in
// template.Funcmap's format:
//  - tmpl  : execute a sub-template by name. Sub-templates are chosen from t
//            namespace
func TmplFuncMap(t *template.Template) template.FuncMap {
	return template.FuncMap{
		"tmpl": execTemplateByName(t),
	}
}

// execTemplateByName applies the parsed template by its registered name.
// execTemplateByName basically has th esame behavior than {{ template xXx . }}
// directive but allow to pipe outcome to another function.
func execTemplateByName(t *template.Template) func(string, interface{}) (string, error) {
	return func(name string, v interface{}) (string, error) {
		buf := &bytes.Buffer{}
		err := t.ExecuteTemplate(buf, name, v)
		return buf.String(), err
	}
}

// pathSanitizer rewrites string to remove non-standard path characters
func pathSanitizer(path string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsMark(r) ||
			r == '.' ||
			r == '_' ||
			r == '-' ||
			r == '%' ||
			r == '#' ||
			r == ' ' ||
			r == os.PathSeparator {
			return r
		}

		if unicode.IsSpace(r) {
			return ' '
		}

		if unicode.In(r, unicode.Hyphen) {
			return '-'
		}

		return -1
	}, path)
}

// filenameSanitizer rewrites string to remove non-standard filename characters
func filenameSanitizer(filename string) string {
	return strings.Map(func(r rune) rune {
		if r == os.PathSeparator {
			return '_'
		}

		return r
	}, pathSanitizer(filename))
}

// nospaceSanitizer rewrites string to remove unreasonable path characters
func nospaceSanitizer(path string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsMark(r) ||
			r == '.' ||
			r == '_' ||
			r == '-' ||
			r == '%' ||
			r == '#' ||
			r == os.PathSeparator {
			return r
		}

		if unicode.IsSpace(r) {
			return '_'
		}

		if unicode.In(r, unicode.Hyphen) ||
			r == '\'' {
			return '-'
		}

		return -1
	}, path)
}
