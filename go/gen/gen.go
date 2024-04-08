package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var templateFuncMap = template.FuncMap{
	"ToLower": strings.ToLower,
}

var iVisitorTemplate = template.Must(template.New("iVisitor").Funcs(templateFuncMap).Parse(`
// generated code - DO NOT EDIT
package generated

type Visitor[T any] interface {
	{{- range .Exprs }}
	visit{{ .Name }} ({{ .Name | ToLower }} Expr) T
	{{- end }}
}
`))

var iExprTemplate = template.Must(template.New("iExpr").Parse(`
// generated code - DO NOT EDIT
package generated

type Expr interface {
	accept(Visitor[any]) any
}
`))

var exprTemplate = template.Must(template.New("exprDef").Parse(`
// generated code - DO NOT EDIT
package generated

import (
	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
)

type {{ .Name }} struct {
	Expr
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }}
	{{- end }}
}

func New{{ .Name }}(
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }},
	{{- end }}
) *{{ .Name }} {
	return &{{ .Name }}{
		{{- range .Attributes }}
		{{ .Name }}: {{ .Name}},
		{{- end }}
	}
}

func (x *{{ .Name }}) accept(visitor Visitor[any]) any {
	return visitor.visit{{ .Name }}(x)
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

	// generate expr defs
	os.MkdirAll("../generated", os.ModePerm)

	for _, expr := range data.Exprs {
		out, err := os.Create("../generated/" +
			strings.ToLower(expr.Name) + ".gen.go")
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		defer out.Close()

		err = exprTemplate.Execute(out, expr)
		if err != nil {
			log.Fatalf("error while executing template: %v", err)
		}
	}

	out, err := os.Create("../generated/" +
		"expr" + ".gen.go")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer out.Close()

	err = iExprTemplate.Execute(out, nil)
	if err != nil {
		log.Fatalf("error while executing template: %v", err)
	}

	out, err = os.Create("../generated/" +
		"visitor" + ".gen.go")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer out.Close()

	err = iVisitorTemplate.Execute(out, data)
	if err != nil {
		log.Fatalf("error while executing template: %v", err)
	}

}
