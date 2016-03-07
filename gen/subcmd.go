package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/tmtk75/aliyun-sdk-go/api"
)

func genSubcmd(fn string, opt options) {
	ff := makeOpeFilter(opt.filter)
	def := load(fn)
	t, err := template.New("subcmd").Funcs(template.FuncMap{"lower": strings.ToLower}).Parse(tmplcmd)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	cmds := make([]string, 0)
	for k, ope := range def.Operations {
		if !ff(k, ope) {
			continue
		}
		delete(ope.Input.Members, "Action")
		v := map[string]interface{}{
			"service": strings.ToLower(def.ServiceAbbreviation),
			"api":     def,
			"ope":     ope,
			"name":    camelToKebab(k),
			"args":    mkargs(ope.Input.Members),
		}
		b := bytes.NewBuffer([]byte{})
		t.Execute(b, v)
		cmds = append(cmds, b.String())
	}
	sort.Strings(cmds)

	t, err = template.New("file").Funcs(template.FuncMap{"lower": strings.ToLower}).Parse(tmplfile)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fb := bytes.NewBuffer([]byte{})
	t.Execute(fb, map[string]interface{}{
		"api":      def,
		"commands": strings.Join(cmds, ","),
	})
	//fmt.Println(fb)

	r, err := format.Source(fb.Bytes())
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if opt.debug {
		fmt.Println(string(r))
		return
	}

	writeFile(r, opt.filename(def.ServiceAbbreviation, "cmd"))
}

const tmplfile = `//
package {{ lower .api.ServiceAbbreviation }}

import (
	"github.com/tmtk75/cli"
	"github.com/tmtk75/aliyun-sdk-go/api"
	"github.com/tmtk75/aliyun-sdk-go/cli/helper"
	"github.com/tmtk75/aliyun-sdk-go/api/{{ lower .api.ServiceAbbreviation }}"
)

func init() {
	helper.Atoi("0") // to avoid the error "imported and not used"
}

var defaultCommands = []cli.Command{
	{{- .commands }},
}
`

const tmplcmd = `
{
	Name: "{{ .name }}",
	{{- if .args}}
	Args: "{{ .args }}",
	{{- end}}
	Action: func(c *cli.Context) {
		conf := &api.Config{RegionId: c.GlobalString("region")}
		{{- range $k, $v := .args}}
		{{- if $v.Varname }}
		{{ $v.Varname }}, _ := c.ArgFor("{{ $v.Argname }}")
		{{- end}}
		{{- end}}
		{{ .service }}.New(conf).Request(&{{ .service }}.{{ .ope.Name }}{
			{{- range $k, $v := .args}}
			{{ $k }}: {{ $v.RightValue }},
			{{- end}}
		})
	},
}`

type ft struct {
	Varname    string
	RightValue string
	Argname    string
}

type Args map[string]ft

func (a Args) String() string {
	s := make([]string, 0)
	for k, _ := range a {
		if !isTimeField(k) {
			s = append(s, "<"+camelToKebab(k)+">")
		}
	}
	sort.Strings(s)
	return strings.Join(s, " ")
}

func mkargs(members map[string]api.Member) Args {
	if len(members) == 0 {
		return nil
	}
	m := make(map[string]ft)
	for k, v := range members {
		if k == "RegionId" {
			continue
		} else if v.Required && k == "EndTime" {
			m[k] = ft{
				RightValue: fmt.Sprintf(`helper.EndTime(c.GlobalString("%v"))`, camelToKebab(k)),
			}
		} else if v.Required && k == "StartTime" {
			m[k] = ft{
				RightValue: `helper.StartTime(c)`,
			}
		} else if v.Required && !isTimeField(k) {
			m[k] = ft{
				RightValue: rightValue(k, v),
				Varname:    camelToSnake(k),
				Argname:    camelToKebab(k),
			}
		}
	}
	return m
}

func rightValue(key string, m api.Member) string {
	v := camelToSnake(key)
	switch m.Type {
	case "integer":
		return fmt.Sprintf(`helper.Atoi(%v)`, v)
	}
	return v
}

func camelToKebab(a string) string {
	return uncamelize(a, "-")
}

func camelToSnake(a string) string {
	return uncamelize(a, "_")
}

func uncamelize(a, deli string) string {
	re := regexp.MustCompile("([a-z]+)?([A-Z]+)?")
	r := re.FindAllStringSubmatch(a, -1)
	s := make([]string, 0)
	var last byte
	for i, e := range r {
		// ignore e[0]
		//fmt.Println(reflect.TypeOf(e).Name())
		p := fmt.Sprintf("%c%v", last, e[1])
		if last == 0 { // take care of initial chunk
			p = e[1]
		}
		s = append(s, p)

		q := e[2]
		if len(q) > 0 {
			last = q[len(q)-1]
			//fmt.Printf("%c\n", last)
			if len(q) > 1 {
				if i == len(r)-1 {
					s = append(s, q[0:len(q)])
				} else {
					s = append(s, q[0:len(q)-1])
				}
			}
		}
	}
	//fmt.Println(r)
	//fmt.Println(s)
	//fmt.Println("--")
	g := make([]string, 0)
	for _, e := range s {
		if len(e) > 0 {
			g = append(g, e)
		}
	}
	return strings.ToLower(strings.Join(g, deli))
}
