package main

import (
	"github.com/egibs/deepwalk/cmd"
	"github.com/spf13/cobra"
)

var (
	defaultValue string
	object       string
	parsedObject map[string]interface{}
	returnValue  string
	searchKey    string
	searchKeys   []string
)

func main() {
	// Commands
	cmdSearch := cmd.CmdSearch(&object, &searchKey, &defaultValue, &returnValue, &parsedObject)
	cmdWalk := cmd.CmdWalk(&object, &searchKeys, &defaultValue, &returnValue, &parsedObject)

	// DeepSearch Flags
	cmdSearch.Flags().StringVar(&object, "object", "", "Object to search")
	cmdSearch.Flags().StringVar(&searchKey, "search-key", "", "Key to search for")
	cmdSearch.Flags().StringVar(&defaultValue, "default-value", "NO_VALUE", "Default value to return if search fails")
	cmdSearch.Flags().StringVar(&returnValue, "return-value", "all", "Value to return if search succeeds")

	// DeepWalk Flags
	cmdWalk.Flags().StringVar(&object, "object", "", "Object to search")
	cmdWalk.Flags().StringSliceVar(&searchKeys, "search-keys", []string{}, "Keys to search for")
	cmdWalk.Flags().StringVar(&defaultValue, "default-value", "NO_VALUE", "Default value to return if search fails")
	cmdWalk.Flags().StringVar(&returnValue, "return-value", "all", "Value to return if search succeeds")

	// Root Command
	rootCmd := &cobra.Command{Use: "deepwalk"}
	rootCmd.AddCommand(cmdWalk, cmdSearch)
	rootCmd.Execute()
}
