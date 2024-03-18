package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var zimCmd = &cobra.Command{
	Use: "zim",
	Run: zimHandle,
}

func zimHandle(cmd *cobra.Command, args []string) {
	fmt.Println("hello zim!")
}

func Execute() {
	if err := zimCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
