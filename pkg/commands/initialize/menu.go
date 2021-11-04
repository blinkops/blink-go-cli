package initialize

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exeName string = filepath.Base(os.Args[0])

func setupMenuOptions(command *cobra.Command, _ []string) (err error) {

	var prompt promptui.Prompt

	prompt = promptui.Prompt{
		Label:   "Hostname",
		Default: "https://app.blinkops.com",
	}

	fullHostname, err := prompt.Run()
	if err != nil {
		return err
	}

	u, err := url.Parse(fullHostname)
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "Blink API Key",
	}

	apiKey, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "Workspace ID",
	}

	workspaceID, err := prompt.Run()
	if err != nil {
		return err
	}

	createConfigFile()
	viper.Set("hostname", u.Host)
	viper.Set("scheme", u.Scheme)
	viper.Set("blink-api-key", apiKey)
	viper.Set("workspace-id", workspaceID)
	viper.WriteConfig()

	fmt.Printf("\nWrote conflig file to %s\n\n", viper.ConfigFileUsed())
	fmt.Println("Try it out - list your playbooks by running the following:")
	fmt.Println("\tblink-cli playbooks list")

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
