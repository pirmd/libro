package main

import (
	"encoding/json"
	"text/template"
)

var (
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
