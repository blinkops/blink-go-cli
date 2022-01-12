package normalizer

import (
	"fmt"
	"strings"

	"github.com/blinkops/blink-go-cli/pkg/consts"

	"github.com/go-openapi/spec"

	"github.com/go-openapi/loads"
	"github.com/spf13/cobra"
)

type CobraAnnotation struct {
	Use     string
	Aliases []string
	Short   string
	Example string
	Hidden  bool
}

type CommandContext struct {
	Cobra CobraAnnotation
	Op    *spec.Operation
}

var cobraAnnotationError = fmt.Errorf("x-cobra-command: invalid structure")

var cobraCommands = make(map[string]map[string]CommandContext)
var cobraOperations = make(map[string]CobraAnnotation)

type CommandNormalizer struct {
	doc *loads.Document
}

func NewCommandNormalizer(doc *loads.Document) *CommandNormalizer {
	return &CommandNormalizer{
		doc: doc,
	}
}

// Normalize Normalizes and formats the cobra command according to the
// annotations specified in the openapi specification
func (normalizer *CommandNormalizer) Normalize(root *cobra.Command) error {

	for k, v := range normalizer.doc.Spec().Extensions {
		if k != "x-cobra-command-operations" {
			continue
		}
		ops, ok := v.(map[string]interface{})
		if !ok {
			return cobraAnnotationError
		}
		for opName, opValue := range ops {
			cobraObj, err := toCobraAnnotation(opValue)
			cobra.CheckErr(err)
			cobraOperations[opName] = cobraObj
		}
	}

	var operationsToAdd []*spec.Operation

	// build up the commands map
	for _, v := range normalizer.doc.Spec().Paths.Paths {
		if v.Get != nil {
			operationsToAdd = append(operationsToAdd, v.Get)
		}
		if v.Put != nil {
			operationsToAdd = append(operationsToAdd, v.Put)
		}
		if v.Post != nil {
			operationsToAdd = append(operationsToAdd, v.Post)
		}
		if v.Delete != nil {
			operationsToAdd = append(operationsToAdd, v.Delete)
		}
		if v.Patch != nil {
			operationsToAdd = append(operationsToAdd, v.Patch)
		}
	}

	for _, op := range operationsToAdd {
		err := addOperationToMap(op)
		cobra.CheckErr(err)
	}

	operations := root.Commands()

	for o := range operations {

		operation := operations[o]
		tag := getOriginalTagName(operation.Use)

		if shouldHideOperation(tag) {
			operation.Hidden = true
			continue
		}

		if isAllowedOperations(tag) {
			continue
		}

		// set the cobra tag name
		cobraTag := GetCobraTag(tag)
		if cobraTag.Use == "" {
			cobraTag.Use = operation.Use
		}

		// set the parent, the operation
		operation.Use = cobraTag.Use
		operation.Example = cobraTag.Example
		operation.Short = cobraTag.Short
		operation.Aliases = cobraTag.Aliases
		operation.Hidden = cobraTag.Hidden

		commands := operation.Commands()
		for c := range commands {
			command := commands[c]
			cobraCommand, ok := cobraCommands[tag][command.Use]
			if !ok {
				continue
			}
			command.Use = cobraCommand.Cobra.Use
			command.Aliases = cobraCommand.Cobra.Aliases
			command.Short = cobraCommand.Cobra.Short
			command.Example = cobraCommand.Cobra.Example
			command.Hidden = cobraCommand.Cobra.Hidden
		}
	}
	return nil
}

func getOriginalTagName(op string) string {
	spaced := strings.ReplaceAll(op, "_", " ")
	return strings.Title(spaced)
}

func addOperationToMap(operation *spec.Operation) error {

	if len(operation.Tags) == 0 {
		return nil
	}

	// like always take the first tag
	tag := operation.Tags[0]

	if _, found := cobraCommands[tag]; !found {
		cobraCommands[tag] = make(map[string]CommandContext)
	}

	cobraObj, err := unmarshalCobraAnnotation(operation.Extensions)
	if err != nil {
		return err
	}

	// if there is no annotation provided
	// use the corresponding swagger annotation

	if cobraObj.Use == "" {
		cobraObj.Use = operation.ID
	}

	if cobraObj.Short == "" {
		cobraObj.Short = operation.Description
	}

	cobraCommands[tag][operation.ID] = CommandContext{
		Cobra: cobraObj, Op: operation,
	}
	return nil

}

func unmarshalCobraAnnotation(extensions spec.Extensions) (CobraAnnotation, error) {

	var cobraObj CobraAnnotation

	for name, ext := range extensions {
		if name != "x-cobra-command" {
			continue
		}
		return toCobraAnnotation(ext)
	}

	return cobraObj, nil
}

func toCobraAnnotation(ext interface{}) (cobraObj CobraAnnotation, err error) {
	cobraParam, ok := ext.(map[string]interface{})
	if !ok {
		// we have problems
		return cobraObj, cobraAnnotationError
	}

	// populate the extensions
	for k, v := range cobraParam {
		switch k {
		case "use":
			cobraObj.Use, ok = v.(string)
			if !ok {
				return cobraObj, cobraAnnotationError
			}
		case "aliases":
			aliases, ok := v.([]interface{})
			if !ok {
				return cobraObj, cobraAnnotationError
			}
			for i := range aliases {
				alias, ok := aliases[i].(string)
				if !ok {
					return cobraObj, cobraAnnotationError
				}
				cobraObj.Aliases = append(cobraObj.Aliases, alias)
			}
		case "short":
			cobraObj.Short, ok = v.(string)
			if !ok {
				return cobraObj, cobraAnnotationError
			}
		case "example":
			cobraObj.Example, ok = v.(string)
			if !ok {
				return cobraObj, cobraAnnotationError
			}
		case "hidden":
			cobraObj.Hidden = v.(bool)
		}
	}
	return cobraObj, nil
}

func GetCobraTag(name string) CobraAnnotation {
	op, found := cobraOperations[name]
	if !found {
		return CobraAnnotation{}
	}
	return op
}

func shouldHideOperation(name string) bool {
	// we want to hide  only the operations that are allowed as specified
	// or that

	if isAllowedOperations(name) {
		return false
	}
	_, found := cobraOperations[name]
	return !found
}

func isAllowedOperations(name string) bool {

	for _, v := range consts.AllowedOperations() {
		if getOriginalTagName(v) == name {
			return true
		}
	}

	return false
}
