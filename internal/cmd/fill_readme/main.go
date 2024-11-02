package main

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"
	"unicode"
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
	HiterAsync    string
	HiterErrbox   string
	HiterIterable string
	HiterSh       string
}

type Hiter struct {
	Source    []string
	Adapter   []string
	Collector []string
}

func main() {

	var hiterData Hiter

	hiterDoc := strings.Split(parse("./hiter"), "\n")

	for _, line := range hiterDoc {
		line = strings.TrimSpace(line)
		if slices.ContainsFunc(
			[]string{
				"func WrapSeqIterable",
				"func WrapSeqIterable2",
			},
			func(s string) bool {
				return strings.HasPrefix(line, s)
			},
		) {
			continue
		}
		if slices.ContainsFunc(
			[]string{
				"func AtterIndices",
				"func AtterAll",
				"func AtterRange",
				"func MapKeys",
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
			HiterAsync:    godoc("./hiter/async"),
			HiterErrbox:   godoc("./hiter/errbox"),
			HiterIterable: godoc("./hiter/iterable"),
			HiterSh:       godoc("./hiter/sh"),
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

func parse(s string) string {
	fset := token.NewFileSet()
	dirents, err := os.ReadDir(s)
	if err != nil {
		panic(err)
	}
	var files []*ast.File
	for _, dirent := range dirents {
		if !dirent.Type().IsRegular() {
			continue
		}
		if strings.HasSuffix(dirent.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(
			fset,
			filepath.Join(s, dirent.Name()),
			nil,
			parser.AllErrors|parser.ParseComments,
		)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}

	var buf bytes.Buffer
	for _, f := range files {
		for _, d := range f.Decls {
			fnDec, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !unicode.IsUpper(rune(fnDec.Name.Name[0])) {
				continue
			}
			if fnDec.Recv != nil {
				continue
			}
			buf.WriteString("func ")
			err := printer.Fprint(&buf, fset, fnDec.Name)
			if err != nil {
				panic(err)
			}
			fnDec.Type.Func = 0
			zeroPos(fnDec.Type.TypeParams)
			zeroPos(fnDec.Type.Params)
			zeroPos(fnDec.Type.Results)

			var s strings.Builder
			err = printer.Fprint(&s, fset, fnDec.Type)
			if err != nil {
				panic(err)
			}
			buf.WriteString(foldInterface(s.String())[len("func"):])
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func zeroPos(fl *ast.FieldList) {
	if fl == nil {
		return
	}
	fl.Opening = 0
	for _, f := range fl.List {
		for _, n := range f.Names {
			n.NamePos = 0
		}
		i, ok := f.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		i.Interface = 0
		if i.Methods == nil {
			continue
		}
		for _, i := range i.Methods.List {
			for _, n := range i.Names {
				n.NamePos = 0
			}
		}
	}
	fl.Closing = 0
}

func foldInterface(s string) string {
	scanning := false
	var starting, ending int
	for i := 0; i < len(s); i++ {
		if strings.HasPrefix(s[i:], "interface") {
			scanning = true
			starting = i + len("interface {")
		}
		if scanning && strings.HasPrefix(s[i:], "}") {
			ending = i
		}
	}
	if ending == 0 {
		return s
	}
	ss := strings.TrimSpace(s[starting:ending])
	var sss string
	for i, line := range strings.Split(ss, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if i > 0 {
			sss += "; "
		}
		sss += line
	}

	return s[:starting] + " " + sss + " " + s[ending:]
}
