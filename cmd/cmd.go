package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/egibs/deepwalk/internal/util"
	"github.com/egibs/deepwalk/pkg/traverse"
	"github.com/spf13/cobra"
)

func CmdTraverse(
	object string,
	searchKey string,
	defaultValue string,
	returnValue string,
	parsedObject interface{},
) *cobra.Command {
	return &cobra.Command{
		Use:   "traverse --object <filename or JSON string> --search-key <search key> --default-value <default value> --return-value <return value>",
		Short: "Traverse an object for the specified key and return the value associated with that key",
		Long: `traverse utilizes the Traverse function which requires a ReturnControl value to determine what to return.
		First will return the first value found, Last will return the last value found, and All will return all values found.`,
		Run: func(cmd *cobra.Command, args []string) {
			var returnValueType traverse.ReturnControl
			switch returnValue {
			case "first":
				returnValueType = traverse.First
			case "last":
				returnValueType = traverse.Last
			case "all":
				returnValueType = traverse.All
			}

			err := util.HandleObjectInput(object, parsedObject)
			if err != nil {
				log.Fatalf("Error handling object input: %s", err)
			}
			result, err := traverse.Traverse(parsedObject, searchKey, defaultValue, returnValueType)
			if err != nil {
				log.Fatalf("Error traversing object: %s", err)
			}
			formattedResult, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}

			fmt.Println(string(formattedResult))
		},
	}
}
