package consts

const (
	CompletionAutoGen = "completion [Bash|Zsh|Fish|Powershell]"
	HelpAutoGen       = "help"
	InitExtra         = "init"
)

func AllowedOperations() []string {
	return []string{CompletionAutoGen, HelpAutoGen, InitExtra}
}
