package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func genService() {
	t, _ := template.New("api").Funcs(template.FuncMap{"lower": strings.ToLower}).Parse(templ)

	for _, s := range []string{"ECS", "RDS", "OTS", "SLB", "RAM", "OSS"} {
		dirpath := filepath.Join("api", s)
		err := os.MkdirAll(strings.ToLower(dirpath), 0755)
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		f, err := os.Create(filepath.Join(dirpath, "api.go"))
		defer f.Close()
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		b := bytes.NewBuffer([]byte{})
		if err := t.Execute(b, map[string]string{"service": s}); err != nil {
			log.Fatalf("%v\n", err)
		}

		bb, err := format.Source(b.Bytes())
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		f.Write(bb)
	}
}

const templ = `package {{ lower .service }}

import (
	api "github.com/tmtk75/aliyun-sdk-go/api"
)

type {{ .service }} struct {
	config *api.Config
}

func New(c *api.Config) *{{ .service }} {
	return &{{ .service }}{config: c}
}

func (s *{{ .service }}) Request(req interface{}) {
	a := api.Fill(s.config, req)
	api.Request(s.config, s, a)
}

func (s *{{ .service }}) Version() string {
	return API_VERSION
}

func (s *{{ .service }}) Endpoint(name string) string {
	return ENDPOINT + methods[name].Uri
}

func (s *{{ .service }}) Method(opename string) string {
	return methods[opename].Method
}
`
