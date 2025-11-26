package cmd

import (
	"errors"
	"fmt"

	"github.com/coffeemakingtoaster/excel-diff/pkg/diff"
	"github.com/coffeemakingtoaster/excel-diff/pkg/display"
	fileutils "github.com/coffeemakingtoaster/excel-diff/pkg/file_utils"
	"github.com/spf13/cobra"
)

var columns = []string{}
var columnIds = []int{}

func getDiffCmd() *cobra.Command {
	var diffCmd = &cobra.Command{
		Use:   "diff",
		Short: "Get diff of two excel files",
		Long:  `Bla bla`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Need exactly TWO filepaths")
			}
			return nil
			for _, p := range args {
				if !fileutils.FileExists(p) {
					return fmt.Errorf("Path %s does not resolve to file", p)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			differ, err := diff.LoadExcelDiff(args[0], args[1], columns, columnIds)
			if err != nil {
				return err
			}
			added := differ.GetAddedLines()
			removed := differ.GetRemovedLines()

			fmt.Printf("-- THIS: %s\n-- THAT: %s\n", args[0], args[1])
			fmt.Printf("-- added (%d) --\n", len(added))
			display.PrintColoredBlock(added, fmt.Sprintf("%s+", display.GREEN))
			fmt.Printf("-- removed (%d) --\n", len(removed))
			display.PrintColoredBlock(removed, fmt.Sprintf("%s-", display.RED))
			return nil
		},
	}
	diffCmd.Flags().StringArrayVarP(&columns, "columns", "c", []string{}, "Specify all column names that should be a unique tuple")
	diffCmd.Flags().IntSliceVarP(&columnIds, "column-ids", "i", []int{}, "Specify all column ids that should be a unique tuple")

	return diffCmd
}
