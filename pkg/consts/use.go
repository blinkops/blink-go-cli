package consts

const (
	CompletionAutoGen = "completion [Bash|Zsh|Fish|Powershell]"
	HelpAutoGen       = "help"
)

func AllowedOperations() []string {
	return allowedList
}

func AddAllowedOperation(name string) {
	allowedList = append(allowedList, name)
}

var allowedList = []string{CompletionAutoGen, HelpAutoGen}
