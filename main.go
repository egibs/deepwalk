package main

import (
	"github.com/egibs/deepwalk/v2/cmd"
	"github.com/spf13/cobra"
)

var (
	defaultValue string
	object       string
	parsedObject interface{}
	returnValue  string
	searchKey    string
)

func main() {
	// Commands
	cmdTraverse := cmd.CmdTraverse(object, searchKey, defaultValue, returnValue, parsedObject)

	// Traverse Flags
	cmdTraverse.Flags().StringVar(&object, "object", "", "Object to search")
	cmdTraverse.Flags().StringVar(&searchKey, "search-key", "", "Key to search for")
	cmdTraverse.Flags().StringVar(&defaultValue, "default-value", "NO_VALUE", "Default value to return if search fails")
	cmdTraverse.Flags().StringVar(&returnValue, "return-value", "all", "Value to return if search succeeds")

	// Root Command
	rootCmd := &cobra.Command{Use: "deepwalk"}
	rootCmd.AddCommand(cmdTraverse)
	rootCmd.Execute()
}
