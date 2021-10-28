//go:generate swagger generate cli -f ../specification/swagger.yaml

package main

import (
	"fmt"
	"github.com/go-openapi/loads"
	"os"
	"path/filepath"
	"text/template"
)

func main(){

	var d struct {
		Spec string
	}

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

	d.Spec = string(doc.Raw())

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
	if err != nil{
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
