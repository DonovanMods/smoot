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
	"log"
	"os"

	"github.com/donovanmods/7dtd-modtools/lib/modinfo"
	"github.com/donovanmods/smoot/lib/modlist"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sorts the given mod directory based on the input list",
	Run:   runSort,
}

func runSort(cmd *cobra.Command, args []string) {
	var verbosity = viper.GetInt("verbosity")
	var modorder = viper.GetString("modorder")
	var directory = cmd.Flag("dir").Value.String()

	if verbosity > 2 {
		log.Printf("sort called with modorder:%q and directory:%q\n", modorder, directory)
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Fatal("Directory", directory, "does not exist")
		os.Exit(1)
	}

	if _, err := os.Stat(modorder); os.IsNotExist(err) {
		log.Fatal("Modlist file", modorder, "does not exist")
		os.Exit(1)
	}

	modInfos := *modinfo.ParseDir(modinfo.ParseOpts{Directory: directory, Verbosity: verbosity})
	if len(modInfos) == 0 {
		log.Fatal("No modinfo.xml files found in", directory)
		os.Exit(1)
	}

	modListing, err := modlist.Read(modorder)
	if err != nil {
		log.Fatal("Error reading modlist file", modorder, ":", err)
		os.Exit(1)
	}

	log.Print("Mod Listing:\n", *modListing)
	for _, mod := range *modListing {
		modInfo, found := modInfos.Get(mod.Name)

		if verbosity >= 2 && found {
			log.Printf("Found mod %q\n", mod.Name)
		}

		if !found {
			log.Printf("Did NOT find Mod %q\n", mod.Name)
			continue
		}

		newFilename := fmt.Sprintf("%s-%s", mod.Priority, modInfo.Filename())

		if verbosity > 0 {
			log.Printf("Moving %q to %q\n", mod.Name, newFilename)
		}
	}
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().StringP("dir", "d", "", "The directory to be sorted (required)")
	sortCmd.PersistentFlags().StringP("modorder", "m", "", "The modorder file to read for load order")

	_ = sortCmd.MarkFlagRequired("dir")
	_ = sortCmd.MarkFlagFilename("modorder")
	_ = sortCmd.MarkFlagRequired("modorder")
}
