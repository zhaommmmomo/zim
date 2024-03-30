package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhaommmmomo/zim/gateway"
)

var gatewayCmd = &cobra.Command{
	Use: "gateway",
	Run: gatewayHandle,
}

func gatewayHandle(cmd *cobra.Command, args []string) {
	gateway.Start(configPath)
}

func init() {
	zimCmd.AddCommand(gatewayCmd)
}
