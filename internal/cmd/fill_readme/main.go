package main

import (
	"context"
	"os"
	"os/exec"
	"text/template"
)

func godoc(tgt string) string {
	out, err := exec.CommandContext(context.Background(), "go", "doc", tgt).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

type Param struct {
	GoDoc GoDoc
}
type GoDoc struct {
	Collection    string
	Hiter         string
	HiterIterable string
}

func main() {
	p := Param{
		GoDoc: GoDoc{
			Collection:    godoc("./collection"),
			Hiter:         godoc("./hiter"),
			HiterIterable: godoc("./hiter/iterable"),
		},
	}

	w, err := os.Create("./README.md")
	if err != nil {
		panic(err)
	}

	err = template.Must(template.ParseFiles("./template.readme.md")).Execute(w, p)
	if err != nil {
		panic(err)
	}
}
