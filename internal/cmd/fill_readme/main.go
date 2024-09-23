package main

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

var (
	hasIterArg = regexp.MustCompile(`func\s.*(\[.*\])?\(.*iter\.Seq2?\[.*\].*\)`)
	hasIterRet = regexp.MustCompile(`\)\s\(?iter\.Seq2?.*\]\)?$`)
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
	Hiter         Hiter
	HiterIterable string
	HiterErrbox   string
}

type Hiter struct {
	Source    []string
	Adapter   []string
	Collector []string
}

func main() {

	var hiterData Hiter

	hiterDoc := strings.Split(godoc("./hiter"), "\n")

	for _, line := range hiterDoc {
		line = strings.TrimSpace(line)
		if slices.ContainsFunc(
			[]string{
				"func IndexAccessible",
			},
			func(s string) bool {
				return strings.HasPrefix(line, s)
			},
		) {
			hiterData.Source = append(hiterData.Source, line)
			continue
		}
		if slices.ContainsFunc(
			[]string{
				"func Omit[",
				"func Omit2[",
			},
			func(s string) bool { return strings.HasPrefix(line, s) },
		) {
			hiterData.Adapter = append(hiterData.Adapter, line)
			continue
		}
		if hasIterArg.MatchString(line) {
			if hasIterRet.MatchString(line) {
				hiterData.Adapter = append(hiterData.Adapter, line)
			} else {
				hiterData.Collector = append(hiterData.Collector, line)
			}
		} else if strings.HasPrefix(line, "func") && hasIterRet.MatchString(line) {
			hiterData.Source = append(hiterData.Source, line)
		}
	}

	p := Param{
		GoDoc: GoDoc{
			Hiter:         hiterData,
			HiterIterable: godoc("./hiter/iterable"),
			HiterErrbox:   godoc("./hiter/errbox"),
		},
	}

	w, err := os.Create("./README.md")
	if err != nil {
		panic(err)
	}

	err = template.Must(
		template.New("").
			Funcs(template.FuncMap{"join": strings.Join}).
			ParseFiles("./template.readme.md"),
	).ExecuteTemplate(w, "template.readme.md", p)
	if err != nil {
		panic(err)
	}
}
