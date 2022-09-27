package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/guionardo/todo-cli/internal"
	"github.com/guionardo/todo-cli/pkg/git"
	"github.com/guionardo/todo-cli/pkg/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	SetupCommand = &cli.Command{
		Name:     "setup",
		Usage:    "Setup todo app",
		Aliases:  []string{"s"},
		Category: "Setup",
		// Action:  ActionSetup,
		Subcommands: []*cli.Command{
			{
				Name:   "show",
				Usage:  "Show current setup",
				Action: ActionShowSetup,
			},
			{
				Name:   "new",
				Usage:  "Create a new setup",
				Action: ActionNewSetup,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Name of the todo list",
						Value:   fmt.Sprintf("%s's TODO", getUser()),
					},
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "Github token. Create a new github token at https://github.com/settings/tokens/new with gist permission",
						Required: true,
					},
				},
			},
			{
				Name:   "shell",
				Usage:  "Setup shell integration",
				Action: ActionSetupShell,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "remove",
						Aliases: []string{"r"},
						Usage:   "Remove shell integration",
					},
				},
			},
		},
	}
)

func getUser() string {
	gitUser, err := git.GetCurrentGitUser()
	if err == nil && len(gitUser.Name) > 0 {
		return gitUser.Name
	}
	user := os.Getenv("USER")
	if len(user) == 0 {
		user = "UNDEFINED"
	}
	return user
}

func getConfigFile(c *cli.Context) string {
	configFile := c.String("config")
	context := internal.GetRunningContext(c)
	if len(configFile) == 0 {
		configFile = internal.DefaultCollectionFilePath
		context.DebugLog("Using default config file: %s", configFile)
	} else {
		context.DebugLog("Using config file: %s", configFile)
	}
	return configFile
}

func ActionShowSetup(c *cli.Context) error {
	configFile := getConfigFile(c)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("Config file %s does not exist", configFile)
	}
	existentConfig, err := internal.ParseCollectionFile(configFile)
	if err != nil {
		return fmt.Errorf("Error reading config file %s: %v", configFile, err)
	}
	fmt.Printf("%v", existentConfig)

	return nil
}

func ActionSetupShell(c *cli.Context) error{
	context := internal.GetRunningContext(c)
	if c.Bool("remove") {
		return internal.RemoveShellIntegration(context)
	}
	return internal.SetupShellIntegration(context)

}

func tryGetConfigFromGist(auth string, debugMode bool) (config *internal.ToDoCollection) {
	// Try to read configuration from gist
	api := github.NewGitHubGistAPI(string(auth), debugMode)
	err := api.GetToDoConfigFileGist()

	if err != nil {
		if err.Error() == "Invalid token" {
			return nil
		}
		config = &internal.ToDoCollection{}
	} else {

		if err = api.GetConfigFileContent(); err != nil {
			return nil
		}
		config, _ = internal.ParseCollectionData(api.ConfigFileContent)
	}

	return
}

func ActionNewSetup(c *cli.Context) error {
	auth := c.String("token")
	if len(auth) == 0 {
		return fmt.Errorf("Github token is required")
	}
	context := internal.GetRunningContext(c)
	localConfigFile := getConfigFile(c)
	gistConfig := tryGetConfigFromGist(auth, context.DebugMode)
	localConfig, _ := internal.ParseCollectionFile(localConfigFile)

	if gistConfig == nil {
		return fmt.Errorf("Invalid token")
	}
	if len(gistConfig.Config.Authorization) > 0 {
		if askYesNo(true, "Found existing configuration in gist %s.\nDo you want to use it? [Y/n]:", gistConfig) {
			if localConfig != nil {

				err := localConfig.GISTSync(context.DebugMode)
				if err != nil {
					return fmt.Errorf("Error merging configuration from gist: %v", err)
				}
				fmt.Printf("Found existing configuration in %s. It was merged with data from gist\n", localConfigFile)
				return nil
			}
			if err := gistConfig.Save(localConfigFile); err != nil {
				return fmt.Errorf("Error saving configuration to %s: %v", localConfigFile, err)
			}
			fmt.Printf("Configuration saved to %s\n", localConfigFile)
			return nil
		} else {
			if askYesNo(false, "After syncing, the existing configuration in gist will be overwritten.\nDo you want to continue? [y/N]:") {
				return nil
			}
		}
	}

	if localConfig != nil {
		fmt.Printf("Found existing configuration in %s\n", localConfigFile)
		if localConfig.Config.Authorization != auth {
			fmt.Printf("Updating authorization token\n")

			localConfig.Config.Authorization = auth
		}
	} else {
		localConfig = &internal.ToDoCollection{
			Config: internal.Config{
				Authorization:  auth,
				ToDoListName:   c.String("name"),
				ConfigFileName: localConfigFile,
			},
		}
	}
	if err := localConfig.Save(localConfigFile); err != nil {
		return fmt.Errorf("Error saving configuration to %s: %v", localConfigFile, err)
	}
	if err := localConfig.GISTSync(context.DebugMode); err != nil {
		return fmt.Errorf("Error syncing configuration to gist: %v", err)
	}
	fmt.Printf("Configuration saved to %s\n and synced to gist", localConfigFile)
	return nil
}

func ActionSetup(c *cli.Context) error {
	fmt.Println("Starting todo setup")
	context := internal.GetRunningContext(c)

	configFile := getConfigFile(c)

	defaultCollectionName := fmt.Sprintf("%s's TODO", getUser())

	existentConfig, err := internal.ParseCollectionFile(configFile)
	previousAuth := ""
	previousGistId := ""

	if err == nil {
		previousAuth = existentConfig.Config.Authorization
		previousGistId = existentConfig.Config.GistId
		defaultCollectionName = existentConfig.Config.ToDoListName
	}

	collectionName := inputText("Collection name", defaultCollectionName)

	fmt.Println("Create a new github token at https://github.com/settings/tokens/new with gist permission")
	fmt.Print("Github authentication token: ")
	if len(previousAuth) > 0 {
		fmt.Print("[ENTER to keep previous] ")
	}
	auth, err := term.ReadPassword(int(syscall.Stdin))
	if len(auth) == 0 {
		if len(previousAuth) > 0 {
			auth = []byte(previousAuth)
		} else {
			err = fmt.Errorf("Authentication token is required")
		}
	}
	if err != nil {
		return err
	}

	config := internal.Config{
		Authorization:  string(auth),
		ToDoListName:   collectionName,
		ConfigFileName: configFile,
		GistId:         previousGistId,
	}
	collection := internal.ToDoCollection{
		Config: config,
	}

	if err == nil {
		if len(existentConfig.Items) > 0 {
			collection.Items = existentConfig.Items
		}
	}

	err = collection.Save(configFile)
	if err != nil {
		return err
	}

	err = collection.GISTSync(context.DebugMode)
	return err

}
