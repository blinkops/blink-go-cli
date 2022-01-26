package menu

import (
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/blinkops/blink-go-cli/gen/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exeName = filepath.Base(os.Args[0])

func Setup(_ *cobra.Command, _ []string) (err error) {
	type config struct {
		hostname  string
		apiKey    string
		workspaceId string
		workspaceName string
	}

	configValues := config{
		hostname:  consts.DefaultBlinkHostname,
		apiKey:    "",
		workspaceId: "",
	}

	// check if file already exists and get previously configured values
	configPath := viper.ConfigFileUsed()
	if _, err := os.Stat(configPath); err == nil {
		scheme := viper.GetString(consts.SchemeEntry)
		configValues.hostname = scheme + "://" + viper.GetString(consts.HostnameEntry)
		configValues.apiKey = viper.GetString(consts.ApiKeyEntry)
		configValues.workspaceId = viper.GetString(consts.WorkspaceIdEntry)
	}

	var prompt promptui.Prompt

	prompt = promptui.Prompt{
		Label:   "Hostname",
		Default: configValues.hostname,
	}
	fullHostname, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: fmt.Sprintf(
			"Blink API Key (Obtain key by accessing %s/api/v1/apikey in your webbrowser), leave blank to use previously configured value", fullHostname,
		),
	}
	apiKey, err := prompt.Run()
	if err != nil {
		return err
	}

	// use previously configured value
	if apiKey == "" {
		apiKey = configValues.apiKey
	}

	u, err := url.Parse(fullHostname)
	if err != nil {
		return err
	}

	r := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	r.DefaultAuthentication = httptransport.Compose(
		httptransport.APIKeyAuth(consts.ApiKeyHeader, "header", apiKey),
	)

	userDetails, err := client.New(r, strfmt.Default).UserInfo.GetUserDetails(nil, nil)
	if err != nil {
		// handle this error and prompt the user for manual?
		return err
	}

	var workspaces []string

	for _, val := range userDetails.Payload.Workspaces {
		// convert previous workspace ID to name and don't append it
		if val.ID == configValues.workspaceId {
			configValues.workspaceName = val.Name
		} else {
			workspaces = append(workspaces, val.Name)
		}
	}

	// if previous workspace was found, prepend it to the list of workspaces
	if configValues.workspaceName != "" {
		workspaces = append([]string{configValues.workspaceName}, workspaces...)
	}

	promptSelect := promptui.Select{
		Label: "Workspace ID",
		Items: workspaces,
	}

	_, workspaceName, err := promptSelect.Run()
	if err != nil {
		return err
	}

	// look up the workspace id
	var workspaceID string
	for key, val := range userDetails.Payload.Workspaces {
		if val.Name == workspaceName {
			workspaceID = key
		}
	}

	createConfigFile()
	viper.Set(consts.HostnameEntry, u.Host)
	viper.Set(consts.SchemeEntry, u.Scheme)
	viper.Set(consts.ApiKeyEntry, apiKey)
	viper.Set(consts.WorkspaceIdEntry, workspaceID)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

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
