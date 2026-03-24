package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	argFirst     int
	argLast      int
	argSelect    int
	argHideIndex bool
	argCsv       bool
	header       table.Row
	rows         []table.Row
)

func verifyArgs(args *[]string) error {
	if len(*args) == 0 {
		return errors.New("path to the file is required")
	}

	if argFirst != 0 && argLast != 0 {
		return errors.New("first and last should not be used at the same time")
	}

	if argSelect != 0 && (argFirst != 0 || argLast != 0) {
		return errors.New("select cannot be combined with first or last flags")
	}

	return nil
}

func readRows(reader *csv.Reader) error {
	var recordIndex int
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
			if !argHideIndex {
				header = append(header, "#")
			}
			for _, col := range record {
				header = append(header, col)
			}
		} else {
			if argFirst != 0 && recordIndex > argFirst {
				break
			}

			r := table.Row{}

			if !argHideIndex {
				r = append(r, recordIndex)
			}
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
			var termWidth int

			fd := int(os.Stdout.Fd())
			isTerminal := term.IsTerminal(fd)
			if isTerminal {
				width, _, err := term.GetSize(fd)
				if err != nil {
					PrintFatal("failed to get terminal size")
				}
				termWidth = width
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
			file.Close()

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			var finalRows []table.Row
			var rowsModified bool

			if argSelect != 0 {
				if len(rows) >= argSelect {
					rows = []table.Row{rows[argSelect-1]}
					if !argHideIndex {
						rows[0][0] = 1
					}
				} else {
					rows = []table.Row{}
				}
			}

			if argLast != 0 {
				finalRows = append(finalRows, rows[max(0, len(rows)-argLast):]...)
				if !argHideIndex {
					for i := range finalRows {
						finalRows[i][0] = i + 1
					}
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
			if isTerminal {
				t.Style().Size.WidthMax = termWidth
			}

			if argCsv {
				t.RenderCSV()
			} else {
				t.Render()
			}
		},
	}

	rootCmd.Flags().IntVar(&argFirst, "first", 0, "select some amount of rows from top")
	rootCmd.Flags().IntVar(&argLast, "last", 0, "select some amount of rows from bottom")
	rootCmd.Flags().IntVar(&argSelect, "select", 0, "select specific row")
	rootCmd.Flags().BoolVar(&argCsv, "csv", false, "display data as CSV")
	rootCmd.Flags().BoolVar(&argHideIndex, "hide-index", false, "do not show index in final table")

	return rootCmd
}
