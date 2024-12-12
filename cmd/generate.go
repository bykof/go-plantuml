package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/bykof/go-plantuml/astParser"
	"github.com/bykof/go-plantuml/domain"
	"github.com/bykof/go-plantuml/formatter"
)

var (
	outPath     string
	directories []string
	files       []string
	exclusion   string
	leftToRight bool
	recursive   bool
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a plantuml diagram from given paths",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var packages domain.Packages
			for _, file := range files {
				packages = append(packages, astParser.ParseFile(file))
			}

			var options []astParser.ParserOptionFunc
			if recursive {
				options = append(options, astParser.WithRecursive())
			}

			if exclusion != "" {
				options = append(options, astParser.WithFileExclusion(exclusion))
			}

			for _, directory := range directories {
				packages = append(packages, astParser.ParseDirectory(directory, options...)...)
			}

			var formatterOptions = formatter.FormatterOptions{LeftToRight: leftToRight}
			formattedPlantUML := formatter.FormatPlantUML(packages, formatterOptions)
			err := os.WriteFile(outPath, []byte(formattedPlantUML), 0644)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	generateCmd.Flags().StringSliceVarP(
		&directories,
		"directories",
		"d",
		[]string{"."},
		"the go source directories",
	)
	generateCmd.Flags().StringSliceVarP(
		&files,
		"files",
		"f",
		[]string{},
		"the go source files",
	)
	generateCmd.Flags().StringVarP(
		&outPath,
		"out",
		"o",
		"graph.puml",
		"the graphfile",
	)
	generateCmd.Flags().BoolVarP(
		&recursive,
		"recursive",
		"r",
		false,
		"traverse the given directories recursively",
	)
	generateCmd.Flags().BoolVarP(
		&leftToRight,
		"left-to-right",
		"l",
		false,
		"display the UML diagram left to right (default: top to bottom)",
	)
	generateCmd.Flags().StringVarP(
		&exclusion,
		"exclude",
		"x",
		"",
		"exclude file matching given regex expression, not used if using -f flag",
	)
	rootCmd.AddCommand(generateCmd)
}
