package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhaommmmomo/zim/ipconf"
)

var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: ipConfHandle,
}

func ipConfHandle(cmd *cobra.Command, args []string) {
	ipconf.Start(configPath)
}

func init() {
	zimCmd.AddCommand(ipConfCmd)
}
