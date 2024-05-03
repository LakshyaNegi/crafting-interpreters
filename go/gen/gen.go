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

type Visitor{{ .Name }} interface {
	{{- range .Data }}
	Visit{{ . }} ({{ . | ToLower }} *{{ . }}) (interface{}, error)
	{{- end }}
}
`))

type visitorDef struct {
	Name string
	Data []string
}

var iExprTemplate = template.Must(template.New("iExpr").Parse(`
// generated code - DO NOT EDIT
package generated

type Expr interface {
	Accept(VisitorExpr) (interface{}, error)
}
`))

var iStmtTemplate = template.Must(template.New("iStmt").Parse(`
// generated code - DO NOT EDIT
package generated

type Stmt interface {
	Accept(VisitorStmt) (interface{}, error)
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
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }}
	{{- end }}
}

func New{{ .Name }}(
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }},
	{{- end }}
) *{{ .Name }} {
	return &{{ .Name }} {
		{{- range .Attributes }}
		{{ .Name }}: {{ .Name}},
		{{- end }}
	}
}

func (x *{{ .Name }}) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.Visit{{ .Name }}(x)
}
`))

var stmtTemplate = template.Must(template.New("stmtDef").Parse(`
// generated code - DO NOT EDIT
package generated

import (
	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
)

type {{ .Name }} struct {
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }}
	{{- end }}
}

func New{{ .Name }}(
	{{- range .Attributes }}
	{{ .Name }} {{ .Type }},
	{{- end }}
) *{{ .Name }} {
	return &{{ .Name }} {
		{{- range .Attributes }}
		{{ .Name }}: {{ .Name}},
		{{- end }}
	}
}

func (x *{{ .Name }}) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.Visit{{ .Name }}(x)
}
`))

type ExprDefs struct {
	Name  string    `yml:"name"`
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

type StmtDefs struct {
	Name  string    `yml:"name"`
	Stmts []StmtDef `yml:"stmts"`
}

type StmtDef struct {
	Name       string     `yaml:"name"`
	Imports    []string   `yaml:"imports"`
	Attributes []StmtAttr `yaml:"attributes"`
}

type StmtAttr struct {
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

	var enames []string
	for _, expr := range data.Exprs {
		out, err := os.Create("../generated/" +
			strings.ToLower(expr.Name) + "_expr" + ".gen.go")
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		defer out.Close()

		err = exprTemplate.Execute(out, expr)
		if err != nil {
			log.Fatalf("error while executing template: %v", err)
		}

		enames = append(enames, expr.Name)
	}

	// expr interface
	out, err := os.Create("../generated/" +
		"expr" + ".gen.go")
	if err != nil {
		log.Fatalf("error creating expr interface file: %v", err)
	}
	defer out.Close()

	err = iExprTemplate.Execute(out, nil)
	if err != nil {
		log.Fatalf("error while executing expr interface template: %v", err)
	}

	// Load YAML file into memory
	yamlData, err = os.ReadFile("stmt_defs.yml")
	if err != nil {
		log.Fatalf("error reading yaml file: %v", err)
	}

	// Unmarshal YAML data into Go struct
	var sdata StmtDefs
	err = yaml.Unmarshal(yamlData, &sdata)
	if err != nil {
		log.Fatalf("error unmarshaling yaml: %v", err)
	}

	var snames []string
	for _, stmt := range sdata.Stmts {
		out, err := os.Create("../generated/" +
			strings.ToLower(stmt.Name) + ".gen.go")
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		defer out.Close()

		err = stmtTemplate.Execute(out, stmt)
		if err != nil {
			log.Fatalf("error while executing template: %v", err)
		}

		snames = append(snames, stmt.Name)
	}

	// stmt interface
	out, err = os.Create("../generated/" +
		"stmt" + ".gen.go")
	if err != nil {
		log.Fatalf("error creating stmt interface file: %v", err)
	}
	defer out.Close()

	err = iStmtTemplate.Execute(out, nil)
	if err != nil {
		log.Fatalf("error while executing stmt interface template: %v", err)
	}

	vds := []*visitorDef{
		{
			Name: "Expr",
			Data: enames,
		},
		{
			Name: "Stmt",
			Data: snames,
		},
	}

	for _, vd := range vds {
		out, err = os.Create("../generated/" +
			"visitor_" + strings.ToLower(vd.Name) + ".gen.go")
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		defer out.Close()

		err = iVisitorTemplate.Execute(out, vd)
		if err != nil {
			log.Fatalf("error while executing template: %v", err)
		}
	}
}
