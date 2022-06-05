package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

func GetFileContentType(content []byte) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	var buffer []byte
	if len(content) > 512 {
		buffer = content[:512]
	} else {
		buffer = content
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	if strings.HasPrefix(contentType, "text/") {
		contentType = detectTextContentType(content)
	}

	return contentType, nil
}

func detectTextContentType(content []byte) string {
	// Try parse json
	var js json.RawMessage
	if json.Unmarshal(content, &js) == nil {
		return "application/json"
	}

	// Try parse xml
	var xm xml.StartElement
	if xml.Unmarshal(content, &xm) == nil {
		return "application/xml"
	}

	// Try parse yaml
	var ym yaml.MapSlice
	if yaml.Unmarshal(content, &ym) == nil {
		return "application/yaml"
	}
	return "text/plain"
}