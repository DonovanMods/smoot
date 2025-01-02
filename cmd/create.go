/*
Copyright © 2025 Donovan C. Young <dyoung522@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/donovanmods/7dtd-modtools/lib/modinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a load order file based upon modinfo files read",
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	var verbosity = viper.GetInt("verbosity")
	var directory = cmd.Flag("dir").Value.String()
	var outputFile = cmd.Flag("output").Value.String()
	var parseOpts = modinfo.ParseOpts{Directory: directory, Verbosity: verbosity}

	if verbosity > 2 {
		log.Printf("create called on directory:%q\n", directory)
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Fatal("Directory", directory, "does not exist")
		os.Exit(1)
	}

	modInfos, err := modinfo.ParseDir(parseOpts)
	if err != nil {
		log.Fatal("Error while reading modinfo files from", directory, ":", err)
	}

	if len(*modInfos) == 0 {
		log.Fatal("No modinfo.xml files found in", directory)
		os.Exit(1)
	}

	var f *os.File

	if outputFile != "" {
		var err error

		f, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error opening output file", outputFile, ":", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	for _, modInfo := range *modInfos {
		if outputFile != "" {
			_, _ = f.WriteString(modInfo.Name.Value + "\n")
		} else {
			fmt.Println(modInfo.Name.Value)
		}
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("dir", "d", "", "The directory to scan for modfiles (required)")
	createCmd.Flags().StringP("output", "o", "", "The file to write the modlist to (default: stdout)")

	_ = sortCmd.MarkFlagRequired("dir")
}
