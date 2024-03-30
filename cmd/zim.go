package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	configPath string
	zimCmd     = &cobra.Command{
		Use: "zim",
		Run: zimHandle,
	}
)

func zimHandle(cmd *cobra.Command, args []string) {
	fmt.Println("hello zim!")
}

func Execute() {
	if err := zimCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	zimCmd.PersistentFlags().StringVar(&configPath, "config", "./zim.yaml", "config file path")
}
