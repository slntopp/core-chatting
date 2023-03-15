/*
Copyright © 2023 Mikita Iwanowski info@slnt-opp.xyz
*/
package cmd

import (
	"net/http"
	"os"

	"github.com/slntopp/core-chatting/cc/ccconnect"
	"github.com/spf13/cobra"
)

var chat_client = ccconnect.NewChatsAPIClient(
	http.DefaultClient, "http://localhost:8000",
)
var msg_client = ccconnect.NewMessagesAPIClient(
	http.DefaultClient, "http://localhost:8000",
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "core-chatting",
	Aliases: []string{"cc"},
	Short:   "Core Chatting CLI Client",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.core-chatting.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
