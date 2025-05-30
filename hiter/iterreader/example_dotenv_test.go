package iterreader_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterreader"
)

// as per https://www.gnu.org/software/bash/manual/html_node/Double-Quotes.html
const escapeChars = "\"$`\\\n"

var escapeCharsReplacer = strings.NewReplacer(
	`"`, `\"`,
	`$`, `\$`,
	"`", "\\`",
	`\`, `\\`,
	"\n", "\\\n",
)

func quoteVariable(s string) string {
	if !strings.ContainsAny(s, escapeChars) {
		return s
	}
	return `"` + escapeCharsReplacer.Replace(s) + `"`
}

func Example_writing_dot_env() {
	values := map[string]string{
		"FOO": "BAR",
		"BAZ": "`\"${YAY}\"`\nOh?",
	}
	r := iterreader.Reader(
		func(s string) ([]byte, error) { return []byte(s), nil },
		hiter.Unify(
			func(k, v string) string {
				return k + "=" + quoteVariable(v) + "\n"
			},
			hiter.MapsSorted(values),
		),
	)

	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)

	_, err := io.Copy(gw, r)
	if err != nil {
		panic(err)
	}
	if err := gw.Close(); err != nil {
		panic(err)
	}

	gr, err := gzip.NewReader(buf)
	if err != nil {
		panic(err)
	}

	bin, err := io.ReadAll(gr)
	if err != nil {
		panic(err)
	}
	if err := gr.Close(); err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", bin)

	// Output:
	// BAZ="\`\"\${YAY}\"\`\
	// Oh?"
	// FOO=BAR
}
