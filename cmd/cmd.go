package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/egibs/deepwalk/internal/util"
	"github.com/egibs/deepwalk/pkg/deepsearch"
	"github.com/egibs/deepwalk/pkg/deepwalk"
	"github.com/spf13/cobra"
)

func CmdSearch(
	object *string,
	searchKey *string,
	defaultValue *string,
	returnValue *string,
	parsedObject *map[string]interface{},
) *cobra.Command {
	return &cobra.Command{
		Use:   "search --object <filename or JSON string> --search-key <search key> --default-value <default value> --return-value <return value>",
		Short: "Naively search an object for the specified key",
		Long: `search utilizes the DeepSearch function which does not need to know the structure of the data.
		It will search the entire object for the specified key and return the value associated with that key.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := util.HandleObjectInput(*object, parsedObject)
			if err != nil {
				fmt.Println(err)
			}
			result, err := deepsearch.DeepSearch(*parsedObject, *searchKey, *defaultValue, *returnValue)
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
}

func CmdWalk(
	object *string,
	searchKeys *[]string,
	defaultValue *string,
	returnValue *string,
	parsedObject *map[string]interface{},
) *cobra.Command {
	return &cobra.Command{
		Use:   "walk --object <filename or JSON string> --search-keys <search key 1, search key 2> --default-value <default value> --return-value <return value>",
		Short: "Walk an object for the specified keys and return the value associated with the last key",
		Long: `walk utilizes the DeepWalk function which requires a traversal path with the last element of the slice being
		the desired key to retrieve a value for.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := util.HandleObjectInput(*object, parsedObject)
			if err != nil {
				fmt.Println(err)
			}
			result, err := deepwalk.DeepWalk(*parsedObject, *searchKeys, *defaultValue, *returnValue)
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
}
