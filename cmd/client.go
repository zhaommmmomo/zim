package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhaommmmomo/zim/client"
)

var clientCmd = &cobra.Command{
	Use: "client",
	Run: clientHandle,
}

func clientHandle(cmd *cobra.Command, args []string) {
	client.InitCui()
}

func init() {
	zimCmd.AddCommand(clientCmd)
}
