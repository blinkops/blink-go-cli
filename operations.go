package main

// OpenAPIEnabledOperations
// Provide the exact text as defined in the openapi spec
var OpenAPIEnabledOperations = map[string]map[string]struct{
	use string
}{
	"Connections": {
		"ConnectionFindByFilter": {
			use: "search",
		},
	},
	"Workspaces": {
		"CreateWorkspace": {
			use: "create",
		},
	},
}
