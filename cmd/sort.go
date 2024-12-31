/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>

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
	"os"

	"github.com/donovanmods/7dtd-modtools/lib/modinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sorts the mod directory based on input file",
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	var verbosity = viper.GetInt("verbosity")
	var modlist = viper.GetString("modlist")
	var directory = cmd.Flag("dir").Value.String()

	if verbosity > 0 {
		fmt.Printf("sort called with modlist:%q and directory:%q\n", modlist, directory)
	}

	modInfos := *modinfo.ParseDir(directory)
	if len(modInfos) == 0 {
		fmt.Println("No modinfo.xml files found in", directory)
		os.Exit(1)
	}

	if modlist == "STDIN" {
		fmt.Println("STDIN not yet supported, please provide a modlist file")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().StringP("dir", "d", "", "The directory to be sorted (required)")
	_ = sortCmd.MarkFlagRequired("dir")
}
