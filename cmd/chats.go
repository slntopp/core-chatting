package cmd

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:     "chats",
	Aliases: []string{"chat", "c", "cht"},
	Short:   "Manage Core-Chatting Chats",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List Chats",
	RunE: func(cmd *cobra.Command, args []string) error {
		chats, err := chat_client.List(context.Background(), connect.NewRequest(&cc.Empty{}))
		if err != nil {
			return err
		}

		for _, chat := range chats.Msg.GetChats() {
			fmt.Printf("%s | %s\n", chat.Uuid, chat.Role.String())
			for _, user := range chat.Users {
				fmt.Printf(" -> %s\n", user)
			}
			for _, user := range chat.Admins {
				fmt.Printf(" -> %s\n", user)
			}
		}

		return nil
	},
}

func init() {
	chatCmd.AddCommand(listCmd)

	rootCmd.AddCommand(chatCmd)
}
