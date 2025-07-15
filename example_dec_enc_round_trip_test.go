package goiteratorhelper_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/encodingiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
)

func Example_dec_enc_round_trip() {
	src := []byte(`
	{"foo":"foo"}
	{"bar":"bar"}
	{"baz":"baz"}
	`)

	rawDec := json.NewDecoder(bytes.NewReader(src))
	dec := errbox.New(encodingiter.Decode[map[string]string](rawDec))

	enc := json.NewEncoder(os.Stdout)

	err := encodingiter.Encode(
		enc,
		hiter.Map(
			func(m map[string]string) map[string]string {
				return maps.Collect(
					hiter.Map2(
						func(k, v string) (string, string) { return k + k, v + v },
						maps.All(m),
					),
				)
			},
			dec.IntoIter(),
		),
	)

	fmt.Printf("dec error = %v\n", dec.Err())
	fmt.Printf("enc error = %v\n", err)
	// Output:
	// {"foofoo":"foofoo"}
	// {"barbar":"barbar"}
	// {"bazbaz":"bazbaz"}
	// dec error = <nil>
	// enc error = <nil>
}
