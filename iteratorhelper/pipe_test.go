package iteratorhelper

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
	"testing"
)

func TestPipe(t *testing.T) {
	buf := new(bytes.Buffer)
 	_, _ = io.CopyN(buf, rand.Reader, 37 * 1024 + 20439)

	content := hex.EncodeToString(buf.Bytes())

	r := Pipe(func(w CloserWithError) {
		io.Copy(w, strings.NewReader(content))
	})

	bin, err := io.ReadAll(r)
	
	t.Logf("%s", bin)
	if err != nil {
		t.Fatalf("err is not nil. err =%v", err)
	}
	if string(bin) != content {
		t.Fatalf("content is not an expected one, expected foobarbaz, but is %s", bin)
	}
}