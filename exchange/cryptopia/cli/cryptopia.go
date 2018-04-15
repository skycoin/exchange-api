package main

import (
	"github.com/spf13/cobra"
	"github.com/skycoin/exchange-api/exchange/cryptopia"
	"os/user"
	"log"
	"path/filepath"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command
	client  *cryptopia.Client
)

func getCommands() map[string]*cobra.Command {
	return map[string]*cobra.Command{
		"get_currencies": {
			Use:   "get_currencies",
			Short: "get_currencies description",
			Run:
		},
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Panicf("failed to execute root cobra command. err: %v", err)
	}
}

func init() {
	currentUser, err := user.Current()
	if err != nil {
		log.Panicf("failed to get the new user. err: %v", err)
	}
	var config = filepath.Join(currentUser.HomeDir, ".exchangectl/config.toml")
	viper.SetConfigFile(config)
	err = viper.ReadInConfig()
	if err != nil {
		log.Panicf("failed to read config from %v. err: %v", config, err)
	}
	key := viper.GetString("cryptopia.key")
	if key == "" {
		log.Panic("cryptopia key is empty")
	}
	secret := viper.GetString("cryptopia.secret")
	if secret == "" {
		log.Panic("cryptopia key is empty")
	}
	client = cryptopia.NewAPIClient(key, secret)
	rootCmd = &cobra.Command{Use: "cryptopia"}

	for _, command := range getCommands() {
		rootCmd.AddCommand(command)
	}
}
