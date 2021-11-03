package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exeName string = filepath.Base(os.Args[0])

func SetupMenuOptions(command *cobra.Command, _ []string) error {

	var results []string
	var prompts = []promptui.Prompt{
		{
			Label: "Blink API Key",
		},
		{
			Label: "Workspace ID",
		},
	}

	for _, prompt := range prompts {
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	createConfigFile()
	viper.Set("blink-api-key", results[0])
	viper.Set("workspace-id", results[1])
	viper.WriteConfig()

	fmt.Printf("Wrote conflig file to %s\n", viper.ConfigFileUsed())

	return nil
}

func createConfigFile() {
	// look for default config
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	configPath := path.Join(home, ".config", exeName)
	filePath := path.Join(configPath, "config.json")
	_, err = os.Stat(configPath)
	if !os.IsExist(err) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			cobra.CheckErr(err)
		}
	}

	if _, err := os.Create(filePath); err != nil { // perm 0666
		cobra.CheckErr(err)
	}
}
