package main

import (
	"html/template"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var exprTemplate = template.Must(template.New("exprDef").Parse(`
// generated code - DO NOT EDIT
package expr

import (
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)

type {{ .Name }} struct {
	expr.Expr
	{{ range .Attributes }}
	{{ .Name }} {{ .Type }}
	{{ end }}
}

func New{{ .Name }}(
	{{ range .Attributes }}
	{{ .Name }} {{ .Type }},
	{{ end }}
) *{{ .Name }} {
	return &{{ .Name }}{
		{{ range .Attributes }}
		{{ .Name }}: {{ .Name}},
		{{ end }}
	}
}
`))

type ExprDefs struct {
	Exprs []ExprDef `yml:"exprs"`
}

type ExprDef struct {
	Name       string     `yaml:"name"`
	Imports    []string   `yaml:"imports"`
	Attributes []ExprAttr `yaml:"attributes"`
}

type ExprAttr struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func main() {
	// Load YAML file into memory
	yamlData, err := os.ReadFile("expr_defs.yml")
	if err != nil {
		log.Fatalf("error reading yaml file: %v", err)
	}

	// Unmarshal YAML data into Go struct
	var data ExprDefs
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		log.Fatalf("error unmarshaling yaml: %v", err)
	}

	log.Println(data)

	os.MkdirAll("../expr/generated", os.ModePerm)

	for _, expr := range data.Exprs {
		out, err := os.Create("../expr/generated/" + expr.Name + ".gen.go")
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		defer out.Close()

		err = exprTemplate.Execute(out, expr)
		if err != nil {
			log.Fatalf("error while executing template: %v", err)
		}
	}
}
