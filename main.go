package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	deepsearch "github.com/egibs/deepwalk/pkg/deepsearch"
	deepwalk "github.com/egibs/deepwalk/pkg/deepwalk"
	"github.com/spf13/cobra"
)

func main() {
	var fileObject string
	var stringObject string
	parsedObject := make(map[string]interface{})
	var searchKeys []string
	var searchKey string
	var defaultValue string
	var returnValue string

	cmdSearch := &cobra.Command{
		Use:   "search --file-object <file name> | --string-object <JSON string> --search-key <search key> --default-value <default value> --return-value <return value>",
		Short: "Naively search an object for the specified key",
		Long: `search utilizes the DeepSearch function which does not need to know the structure of the data.
		It will search the entire object for the specified key and return the value associated with that key.`,
		Run: func(cmd *cobra.Command, args []string) {
			if stringObject != "" {
				err := json.Unmarshal([]byte(stringObject), &parsedObject)
				if err != nil {
					fmt.Println(err)
				}
			}
			if fileObject != "" {
				contents, err := os.ReadFile(fileObject)
				if err != nil {
					fmt.Println(err)
				}
				err = json.Unmarshal([]byte(contents), &parsedObject)
				if err != nil {
					fmt.Println(err)
				}
			}
			result, err := deepsearch.DeepSearch(parsedObject, searchKey, defaultValue, returnValue)
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}

			formattedResult, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}

			fmt.Printf("Search result: %s\n", formattedResult)
		},
	}

	cmdWalk := &cobra.Command{
		Use:   "walk --file-object <file name> | --string-object <JSON string> --search-keys <search key 1, search key 2> --default-value <default value> --return-value <return value>",
		Short: "Walk an object for the specified keys and return the value associated with the last key",
		Long: `walk utilizes the DeepWalk function which requires a traversal path with the last element of the slice being
		the desired key to retrieve a value for.`,
		Run: func(cmd *cobra.Command, args []string) {
			if stringObject != "" {
				err := json.Unmarshal([]byte(stringObject), &parsedObject)
				if err != nil {
					fmt.Println(err)
				}
			}
			if fileObject != "" {
				contents, err := os.ReadFile(fileObject)
				if err != nil {
					fmt.Println(err)
				}
				err = json.Unmarshal([]byte(contents), &parsedObject)
				if err != nil {
					fmt.Println(err)
				}
			}
			result, err := deepwalk.DeepWalk(parsedObject, searchKeys, defaultValue, returnValue)
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}
			formattedResult, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}

			fmt.Printf("Search result: %s\n", formattedResult)
		},
	}

	// DeepSearch Flags
	cmdSearch.Flags().StringVar(&fileObject, "file-object", "", "File name containing object to search")
	cmdSearch.Flags().StringVar(&stringObject, "string-object", "", "String object to search")
	cmdSearch.Flags().StringVar(&searchKey, "search-key", "", "Key to search for")
	cmdSearch.Flags().StringVar(&defaultValue, "default-value", "NO_VALUE", "Default value to return if search fails")
	cmdSearch.Flags().StringVar(&returnValue, "return-value", "all", "Value to return if search succeeds")

	// DeepWalk Flags
	cmdWalk.Flags().StringVar(&fileObject, "file-object", "", "File name containing object to search")
	cmdWalk.Flags().StringVar(&stringObject, "string-object", "", "String object to search")
	cmdWalk.Flags().StringSliceVar(&searchKeys, "search-keys", []string{}, "Keys to search for")
	cmdWalk.Flags().StringVar(&defaultValue, "default-value", "NO_VALUE", "Default value to return if search fails")
	cmdWalk.Flags().StringVar(&returnValue, "return-value", "all", "Value to return if search succeeds")

	// Root Command
	rootCmd := &cobra.Command{Use: "deepwalk"}
	rootCmd.AddCommand(cmdWalk, cmdSearch)
	rootCmd.Execute()
}
