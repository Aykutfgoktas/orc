package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Aykutfgoktas/orc/cfile"
	"github.com/Aykutfgoktas/orc/client"
	"github.com/Aykutfgoktas/orc/config"
	"github.com/Aykutfgoktas/orc/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var set bool
var list bool
var add string
var remove bool

var s *spinner.Spinner
var conf config.Config
var ghc client.IGithubClient
var confService config.Service

var spinnerChoice = 9
var spinnerDuration = 100 * time.Millisecond

var pageSize = 100

var configFile = ".orc.conf.json"

var version = "0.8.0"
var use = "orc"
var description = "List repositories in a GitHub organization and clone the selected repository"
var example = "orc -l"

func init() {
	RootCmd.PersistentFlags().StringVarP(&add, "add", "a", "", "add organization")
	RootCmd.PersistentFlags().BoolVarP(&list, "list", "l", false, "list organizations")
	RootCmd.PersistentFlags().BoolVarP(&set, "set", "s", false, "set default organization")
	RootCmd.PersistentFlags().BoolVarP(&remove, "remove", "r", false, "remove organization")

	s = spinner.New(spinner.CharSets[spinnerChoice], spinnerDuration)
	home, _ := os.UserHomeDir()

	cfile := cfile.New(home + "/" + configFile)

	confService = config.New(cfile)

	isOk := confService.CheckConfigFile()

	if !isOk {
		conf = readInput()

		path, err := confService.Create(conf.APIKey, conf.DefaultOrganization)

		if err != nil {
			fmt.Printf("Error on creating the config: %v \n", err)
		} else {
			fmt.Printf("Config file created to here: %s \n", path)
		}
	} else {
		c, _ := confService.Read()
		conf = *c
	}

	ghc = client.NewGithubClient(conf.APIKey)
}

var RootCmd = &cobra.Command{
	Use:   use,
	Short: description,
	RunE: func(cmd *cobra.Command, args []string) error {

		if list {
			listOrganization()
			return nil
		}

		if add != "" {
			addOrganization(add)
			return nil
		}

		if set {
			setDefaultOrganization()
			return nil
		}

		if remove {
			deleteOrganization()
			return nil
		}

		listRepo(conf.DefaultOrganization)
		return nil
	},
	Example: example,
	Version: version,
}

func readInput() config.Config {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the organization name: ")

	org, err := reader.ReadString('\n')

	org = strings.TrimSuffix(org, "\n")

	if err != nil {
		log.Fatal("Error reading input:", err)
	}

	fmt.Print("Enter the Github API Key: ")

	key, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		log.Fatal("Error reading input:", err)
	}

	return config.Config{
		APIKey:              string(key),
		DefaultOrganization: org,
		Organizations:       []string{org},
	}
}

func setDefaultOrganization() {
	fmt.Printf("Current default organization is: %s \n", conf.DefaultOrganization)

	prompt := &survey.Select{
		Message: "Select default organization:",
		Options: conf.Organizations,
	}

	var selectedOrg string

	if err := survey.AskOne(prompt, &selectedOrg, survey.WithPageSize(pageSize)); err != nil {
		log.Fatal("Error selecting repository:", "error", err)
	}

	if err := confService.UpdateDefaultOrganization(selectedOrg); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Organization has been selected as default: %s \n", selectedOrg)
	}
}

func addOrganization(org string) {
	isOk, _ := confService.AddOrganization(org)
	if isOk {
		fmt.Printf("Oranization %s is already in the list, operation will be ignored", org)
	} else {
		fmt.Printf("Oranization %s successfully added", org)
	}
}

func listOrganization() {
	var org string

	prompt := &survey.Select{
		Message: "Select a organization:",
		Options: conf.Organizations,
	}

	if err := survey.AskOne(prompt, &org, survey.WithPageSize(pageSize)); err != nil {
		log.Fatal("Error selecting organization:", "error", err)
	}

	listRepo(org)
}

func listRepo(org string) {
	utils.ClearTerminal()

	s.Prefix = "Getting the list of repositories from " + org + " "

	s.Start()
	repos, err := ghc.Repositories(org)
	s.Stop()

	if err != nil {
		fmt.Printf("Error while getting the repositories from %s: %s \n", org, err.Error())
		return
	}

	var selectedRepo string

	prompt := &survey.Select{
		Message: "Select a repository to clone:",
		Options: repos.RepositoryNames(),
	}

	if err := survey.AskOne(prompt, &selectedRepo, survey.WithPageSize(pageSize)); err != nil {
		log.Fatal("Error selecting repository:", "error", err)
	}

	repo := repos.FindRepoByName(selectedRepo)
	s.Prefix = "Cloning the repository " + repo.Name + " "
	s.Start()
	err = repo.Clone()
	s.Stop()

	if err != nil {
		fmt.Printf("Error while cloning the repo %s, error: %v \n", repo.Name, err)
	} else {
		fmt.Printf("Repository successfully cloned %s \n", repo.Name)
	}
}

func deleteOrganization() {
	var org string

	prompt := &survey.Select{
		Message: "Select a organization to delete:",
		Options: conf.Organizations,
	}

	if err := survey.AskOne(prompt, &org, survey.WithPageSize(pageSize)); err != nil {
		log.Fatal("Error selecting organization:", "error", err)
	}

	err := confService.DeleteOrganization(org)

	if err != nil {
		fmt.Printf("Error while deleting the organization %s, error: %v \n", org, err)
	} else {
		fmt.Printf("Organization successfully deleted %s \n", org)
	}
}
