package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/tmtk75/aliyun-sdk-go/api"
	"github.com/tmtk75/cli"
)

type OpeFilter func(k string, o api.Ope) bool

func makeOpeFilter(fl string) OpeFilter {
	re := regexp.MustCompile(fmt.Sprintf("(?i)^%v.*", fl))
	ff := func(k string, o api.Ope) bool { return re.MatchString(k) }
	return ff
}

type options struct {
	dir    string
	filter string
	debug  bool
}

func newOptions(c *cli.Context) options {
	return options{
		dir:    c.String("dir"),
		filter: c.String("filter"),
		debug:  c.Bool("debug"),
	}
}

func (o options) filename(dirname, basename string) string {
	return fmt.Sprintf("%v/%v/%v.go", o.dir, strings.ToLower(dirname), strings.ToLower(basename))
}

func genOpe(fn string, opt options) {
	ff := makeOpeFilter(opt.filter)
	body, j := Genfile(fn, ff)
	if opt.debug {
		fmt.Printf("%v", string(body))
		return
	}

	writeFile(body, opt.filename(j.ServiceAbbreviation, j.ServiceAbbreviation))
}

func writeFile(body []byte, path string) {
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	f.Write(body)

}

func Genfile(path string, ff OpeFilter) (body []byte, a api.APIDef) {
	j := load(path)

	keys := make([]string, 0, len(j.Operations))
	for k := range j.Operations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	im := make(map[string]bool)
	ob := bytes.NewBuffer([]byte{})
	for _, k := range keys {
		v := j.Operations[k]
		if !ff(k, v) {
			continue
		}
		t := genStruct(v, im)
		ob.Write(t)
	}

	t, err := template.New("head").Funcs(template.FuncMap{"lower": strings.ToLower, "filter": ff}).Parse(tmplhead)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	bf := bytes.NewBuffer([]byte{})
	j.Imports = im
	t.Execute(bf, j)
	//fmt.Println(string(ob.Bytes()))

	bf.Write(ob.Bytes())
	fs, err := format.Source(bf.Bytes())
	if err != nil {
		log.Fatalf("Source: %v\n", err)
	}
	return fs, j
}

const tmplhead = `//
package {{ lower .ServiceAbbreviation }}
import (
	{{range $k, $v := .Imports}}
	"{{ $k }}"
	{{end}}
	"github.com/tmtk75/aliyun-sdk-go/api"
)
const API_VERSION= "{{ .APIVersion }}"
const ENDPOINT = "https://{{ .EndpointPrefix }}.aliyuncs.com"

var methods = make(map[string]api.Http)

func init() {
	{{- range $k, $v := .Operations}}
	{{- if filter $k $v}}
	methods["{{ .Name }}"] = api.Http{Method:"{{ $v.Http.Method }}", Uri:"{{ $v.Http.Uri }}"}
	{{- end}}
	{{- end}}
}
`

func load(path string) api.APIDef {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	var j api.APIDef
	json.Unmarshal(b, &j)
	return j
}

func isTimeField(key string) bool {
	re := regexp.MustCompile("^(StartTime|EndTime)$")
	return re.MatchString(key)
}

func conform(m api.Member, opename string, key string) (api.Member, []string) {
	empty := []string{}
	switch m.Type {
	case "integer":
		return api.Member{Type: "int", Required: m.Required}, empty
	case "boolean":
		return api.Member{Type: "bool", Required: m.Required}, empty
	case "structure":
		return api.Member{Type: "interface{}", Required: m.Required}, empty
	case "":
		if m.Location == "" {
			return api.Member{Type: "interface{}", Required: m.Required, Name: m.Name}, empty
		}
	}
	if m.Location == "uri" || m.Location == "header" {
		return api.Member{Type: "string", Required: m.Required, Name: m.Name}, empty
	}

	//
	if isTimeField(key) {
		f := "2006-01-02T15:04:03Z"
		if opename == "DescribeDBInstancePerformance" {
			f = "2006-01-02T15:04Z"
		}
		return api.Member{Type: "time.Time", Required: m.Required, Format: f}, []string{"time"}
	}

	return m, empty
}

func mktags(key string, m api.Member) string {
	a := []string{
		fmt.Sprintf(`required:"%v"`, m.Required),
	}
	if m.Format != "" {
		a = append(a, fmt.Sprintf(`format:"%v"`, m.Format))
	}
	return fmt.Sprintf("`%v`", strings.Join(a, " "))
}

func genStruct(o api.Ope, imports map[string]bool) []byte {
	valid := regexp.MustCompile("^[a-zA-Z]+$")
	for k, _ := range o.Input.Members {
		if unicode.IsLower(rune(k[0])) || !valid.MatchString(k) {
			delete(o.Input.Members, k)
			continue
		}
		m := o.Input.Members[k]
		km, im := conform(m, o.Name, k)
		o.Input.Members[k] = km
		for _, k := range im {
			imports[k] = true
		}
	}

	tmpl, err := template.New("operation").Funcs(template.FuncMap{"mktags": mktags}).Parse(opetmpl)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	bf := bytes.NewBufferString("")
	err = tmpl.Execute(bf, o)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return bf.Bytes()
}

const opetmpl = `
// {{ .Http.Method }} {{ .Http.Uri }}
type {{ .Name }} struct {
	{{- range $k, $v := .Input.Members}}
	{{- if ne $k "Action"}}
	{{ $k }} {{ $v.Type }} {{ mktags $k $v }} 
	{{- end}}{{end}}
}
`
