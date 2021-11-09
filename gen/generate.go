//go:generate swagger generate cli -f ../specification/swagger.yaml

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-openapi/loads"
)

type Specification struct {
	Spec string
}

func main() {
	fmt.Println("Running generate spec")

	currenDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	specPath := filepath.Join(currenDir, "specification", "swagger.yaml")

	doc, err := loads.Spec(specPath)
	if err != nil {
		panic(err)
	}

	d := Specification{
		Spec: string(doc.Raw()),
	}

	specGenDirPath := filepath.Join(currenDir, "gen", "spec")
	err = os.MkdirAll(specGenDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(specGenDirPath, "docs.go"))
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("spec").Parse(specTemplate))
	err = t.Execute(f, d)
	if err != nil {
		panic(err)
	}

	cobraGenDir := filepath.Join(currenDir, "gen", "cli")
	err = os.MkdirAll(cobraGenDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(cobraGenDir, "cobra_init_config.go"), []byte(cobraInitConfig), 0644)
	if err != nil {
		panic(err)
	}
}

var specTemplate = `
package spec
import "github.com/go-openapi/loads"
var SwaggerSpec = ` + "`{{ .Spec }}`" + `

func GetSwaggerSpec()(*loads.Document, error){
	return loads.Analyzed([]byte(SwaggerSpec), "")
}

`

// add package accessible  function, so we can call the exact cobra config
var cobraInitConfig = `
package cli
func InitViperConfigs(){
	initViperConfigs()
}
`
