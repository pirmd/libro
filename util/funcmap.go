package util

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	// FilepathFuncMap provides standard functions to build/manipulate path in
	// template.Funcmap's format:
	//  - base            : get the path basename
	//  - ext             : get file extension
	//  - sep             : return a path separator ('\' of '/')
	//  - sanitizePath    : replace special chars that don't usually work well when
	//                      use in path name
	//  - sanitizeFilename: like sanitizePath but also
	//                      removes any path separators.
	//  - nospace         : get rid of spaces
	FilepathFuncMap = template.FuncMap{
		"base":             filepath.Base,
		"ext":              filepath.Ext,
		"sep":              func() string { return string(filepath.Separator) },
		"sanitizePath":     pathSanitizer,
		"sanitizeFilename": filenameSanitizer,
		"nospace":          nospaceSanitizer,
	}

	// StringsFuncMap provides standard functions to manipulate strings in
	// template.Funcmap's format:
	//  - join :  join elements separating them wit the given separator.
	//  - lower:  convert string to lower-case
	//  - title:  convert string to title-case
	//  - upper:  convert string to upper-case
	StringsFuncMap = template.FuncMap{
		"join":  strings.Join,
		"lower": strings.ToLower,
		"title": strings.ToTitle,
		"upper": strings.ToUpper,
	}

	// SerializationFuncMap provides standard functions to serialize/deserialize an
	// interface in template.Funcmap's format:
	//  - toJSON      : converts an interface to JSON representation.
	//  - toPrettyJSON: converts an interface to an easy-to-read JSON representation.
	SerializationFuncMap = template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			output, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(output), nil
		},

		"toPrettyJSON": func(v interface{}) (string, error) {
			output, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return "", err
			}
			return string(output), nil
		},
	}
)

// TmplFuncMap provides standard functions to execute go text templates path in
// template.Funcmap's format:
//  - tmpl  : execute a sub-template by name. Sub-templates are chosen from t
//            namespace.
//            Unlike {{template xXx}}, tmpl allows to pipe its result to
//            another command.
func TmplFuncMap(t *template.Template) template.FuncMap {
	return template.FuncMap{
		"tmpl": execTemplateByName(t),
	}
}

// execTemplateByName applies the parsed template by its registered name.
// execTemplateByName basically has the same behavior than {{ template xXx . }}
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
