package goOpenstackAuth

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/foize/go.sgr"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func StringDiff(text1, text2 string) string {
	// find diffs
	dmp := diffmatchpatch.New()
	test := dmp.DiffMain(text1, text2, false)
	diffs := dmp.DiffCleanupSemantic(test)

	// output with colors
	var buffer bytes.Buffer
	for _, v := range diffs {
		// scape text
		v.Text = strings.Replace(v.Text, "[", "[[", -1)
		v.Text = strings.Replace(v.Text, "]", "]]", -1)

		if v.Type == 0 {
			buffer.WriteString(v.Text)
		} else if v.Type == -1 {
			buffer.WriteString("[bg-red bold]")
			buffer.WriteString(v.Text)
			buffer.WriteString("[reset]")
		} else if v.Type == 1 {
			buffer.WriteString("[bg-blue bold]")
			buffer.WriteString(v.Text)
			buffer.WriteString("[reset]")
		}
	}
	// parse to set colors
	colorDiff := sgr.MustParseln(buffer.String())
	return colorDiff
}

func TestServer(code int, body string, headers map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(code) // keep the code after setting headers. If not they will disapear...
		fmt.Fprintln(w, body)
	}))
	return server
}
