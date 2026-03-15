package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	first       int
	last        int
	file        string
	columns     string
	selectRow   int
	recordIndex int
	header      table.Row
	rows        []table.Row
)

func verifyArgs(args *[]string) error {
	if len(*args) == 0 {
		return errors.New("path to the file is required")
	}

	if first != 0 && last != 0 {
		return errors.New("first and last should not be used at the same time")
	}

	if selectRow != 0 && (first != 0 || last != 0) {
		return errors.New("select cannot be combined with first or last flags")
	}

	return nil
}

func readRows(reader *csv.Reader) error {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("cannot read row: %s", err)
		}

		// Consider first row as header
		if recordIndex == 0 {
			header = append(header, "#")
			for _, col := range record {
				header = append(header, col)
			}
		} else {
			if first != 0 && recordIndex > first {
				break
			}

			r := table.Row{}

			r = append(r, recordIndex)
			for _, col := range record {
				r = append(r, col)
			}
			rows = append(rows, r)
		}

		recordIndex++
	}

	return nil
}

func BuildCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "csvv [file]",
		Short: "CSV Viewer",
		Long:  "A CLI tool to inspect CSV data.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := verifyArgs(&args); err != nil {
				PrintFatal(err.Error())
			}

			file, err := os.Open(args[0])
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					PrintFatal("file not found")
				}

				PrintFatal(fmt.Sprintf("cannot open file: %v", err))
			}

			reader := csv.NewReader(file)
			if err := readRows(reader); err != nil {
				PrintFatal(err.Error())
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			var finalRows []table.Row
			var rowsModified bool

			if selectRow != 0 {
				if len(rows) > selectRow {
					rows = []table.Row{rows[selectRow-1]}
					rows[0][0] = 1
				} else {
					rows = []table.Row{}
				}
			}

			if last != 0 {
				finalRows = append(finalRows, rows[max(0, len(rows)-last):]...)
				for i := range finalRows {
					finalRows[i][0] = i + 1
				}
				rowsModified = true
			}

			if rowsModified {
				t.AppendRows(finalRows)
			} else {
				t.AppendRows(rows)
			}

			t.AppendHeader(header)

			t.SetStyle(table.StyleRounded)
			t.Render()
		},
	}

	rootCmd.Flags().IntVar(&first, "first", 0, "select some amount of rows from top")
	rootCmd.Flags().IntVar(&last, "last", 0, "select some amount of rows from bottom")
	rootCmd.Flags().IntVar(&selectRow, "select", 0, "select specific row")
	rootCmd.Flags().StringVar(&columns, "columns", "", "columns to select for display")

	return rootCmd
}
