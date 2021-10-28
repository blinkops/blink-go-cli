package formatter

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/go-openapi/loads"
	"github.com/spf13/cobra"
)

type CobraAnnotation struct {
	Use     string
	Aliases []string
	Short   string
	Example string
}

var cobraAnnotationError = fmt.Errorf("x-cobra-command: invalid structure")

var cobraCommands = make(map[string]map[string]CobraAnnotation)

// Format formats the cobra command according to the
// annotations specified in the openapi specification
func Format(root *cobra.Command, spec *loads.Document) {

	// build up the map
	for _, v := range spec.Spec().Paths.Paths {
		if v.Get != nil {
			addOperationToMap(v.Get)
		}
		if v.Put != nil {
			addOperationToMap(v.Put)
		}
		if v.Post != nil {
			addOperationToMap(v.Post)
		}
		if v.Delete != nil {
			addOperationToMap(v.Delete)
		}
		if v.Patch != nil {
			addOperationToMap(v.Patch)
		}
	}

	operations := root.Commands()

	for o := range operations {

		operation := operations[o]

		tag := getOriginalTagName(operation.Use)

		// set the cobra tag name
		cobraTag := GetCobraTag(tag, spec)
		if cobraTag.Use == "" {
			cobraTag.Use = operation.Use
		}

		// set the parent
		operation.Use = cobraTag.Use
		operation.Example = cobraTag.Example
		operation.Short = cobraTag.Short
		operation.Aliases = cobraTag.Aliases

		commands := operation.Commands()
		for c := range commands {
			command := commands[c]
			cobraCommand, ok := cobraCommands[tag][command.Use]
			if !ok {
				continue
			}
			command.Use = cobraCommand.Use
			command.Aliases = cobraCommand.Aliases
			command.Short = cobraCommand.Short
			command.Example = cobraCommand.Example
		}
	}

}

func getOriginalTagName(op string) string {
	spaced := strings.ReplaceAll(op, "_", " ")
	return strings.Title(spaced)
}

func addOperationToMap(operation *spec.Operation) {

	// like always take the first tag
	if len(operation.Tags) == 0 {
		return
	}

	tag := operation.Tags[0]
	if _, found := cobraCommands[tag]; !found {
		cobraCommands[tag] = make(map[string]CobraAnnotation)
	}

	cobra := unmarshalCobraAnnotation(operation.Extensions)
	// if there is no annotation provided use the operation id as the use.
	if cobra.Use == "" {
		cobra.Use = operation.ID
	}

	cobraCommands[tag][operation.ID] = cobra

}

func unmarshalCobraAnnotation(extensions spec.Extensions) CobraAnnotation {

	var cobra CobraAnnotation

	for name, ext := range extensions {

		if name != "x-cobra-command" {
			continue
		}

		cobraParam, ok := ext.(map[string]interface{})
		if !ok {
			// we have problems
			panic("x-cobra-command: invalid structure for operation")
		}

		// populate the extensions
		for k, v := range cobraParam {
			switch k {
			case "use":
				cobra.Use, ok = v.(string)
				if !ok {
					panic(cobraAnnotationError)
				}
			case "aliases":
				aliases, ok := v.([]interface{})
				if !ok {
					panic(cobraAnnotationError)
				}
				for i := range aliases {
					alias, ok := aliases[i].(string)
					if !ok {
						panic(cobraAnnotationError)
					}
					cobra.Aliases = append(cobra.Aliases, alias)
				}
			case "short":
				cobra.Short, ok = v.(string)
				if !ok {
					panic(cobraAnnotationError)
				}
			case "example":
				cobra.Example, ok = v.(string)
				if !ok {
					panic(cobraAnnotationError)
				}
			}
		}
	}

	return cobra
}

func GetCobraTag(name string, spec *loads.Document) CobraAnnotation {
	for _, val := range spec.Spec().Tags {
		if name == val.Name {
			return unmarshalCobraAnnotation(val.Extensions)
		}
	}
	return CobraAnnotation{}
}
