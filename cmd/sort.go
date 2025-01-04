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
	"path/filepath"
	"regexp"

	"github.com/donovanmods/7dtd-modtools/lib/modinfo"
	"github.com/donovanmods/smoot/lib/modlist"
	"github.com/fatih/color"
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
	var dryrun = viper.GetBool("dryrun")
	var directory = cmd.Flag("dir").Value.String()
	var modorder = cmd.Flag("modorder").Value.String()
	var blue = color.New(color.FgBlue).SprintFunc()
	var red = color.New(color.FgRed).SprintFunc()
	var green = color.New(color.FgGreen).SprintFunc()
	var yellow = color.New(color.FgYellow).SprintFunc()

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Fatal("FATAL Directory ", directory, "does not exist")
	}

	if _, err := os.Stat(modorder); os.IsNotExist(err) {
		log.Fatal(red("FATAL Modorder file ", modorder, "does not exist"))
	}

	if verbosity > 0 {
		fmt.Printf("Reading modinfo files in %s ...\n", blue(directory))
	}

	modInfos, err := modinfo.ParseDir(modinfo.ParseOpts{Directory: directory, Verbosity: verbosity})
	if err != nil {
		log.Fatal(red("FATAL Error while reading modinfo files from", directory, ":", err))
	}

	if len(*modInfos) == 0 {
		log.Fatal(red("FATAL No modinfo.xml files found in", directory))
	}

	if verbosity > 0 {
		fmt.Printf("Parsing modorder file %s ...\n", blue(modorder))
	}

	modListing, err := modlist.Read(modorder)
	if err != nil {
		log.Fatal(red("FATAL Error while reading modorder file", modorder, ":", err))
	}

	if verbosity > 0 {
		fmt.Println("Sorting mods ...")
	}
	for _, mod := range *modListing {
		modInfo, found := findByMO2Name(modInfos, mod.Name)

		if verbosity > 2 && found {
			fmt.Printf("Found mod %s with priority %s\n", blue(mod.Name), blue("%-0.4d", mod.Priority))
		}

		if !found {
			log.Printf("%s %s%s\n", yellow("WARNING Did NOT find Mod"), blue(mod.Name), yellow("... skipped"))
			continue
		}

		newFilename := filepath.Join(modInfo.Dir(), generateNewFilename(modInfo, mod.Priority))

		if verbosity > 1 {
			fmt.Printf("Renaming %s from %s to %s\n", green(mod.Name), blue(modInfo.Filename()), blue(filepath.Base(newFilename)))
		}

		if !dryrun {
			err := os.Rename(modInfo.Path(), newFilename)
			if err != nil {
				log.Printf("%s: %v\n", red("ERROR Error moving %s to %s", modInfo.Path(), newFilename), err)
			}
		}
	}

	if verbosity > 0 {
		fmt.Println(green("Completed Successfully"))
	}
}

func generateNewFilename(modInfo *modinfo.ModInfo, priority int) string {
	var re = regexp.MustCompile(`^\d{4}-`)
	var newFilename = re.ReplaceAllString(modInfo.Filename(), "")

	return fmt.Sprintf("%-0.4d-%s", priority, newFilename)
}

func findByMO2Name(modInfos *modinfo.ModInfos, name string) (*modinfo.ModInfo, bool) {
	for _, modInfo := range *modInfos {
		MO2Name := filepath.Base(modInfo.Dir())

		if MO2Name == name {
			return modInfo, true
		}
	}

	return nil, false
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().StringP("dir", "d", "", "The directory to be sorted (required)")
	sortCmd.Flags().StringP("modorder", "m", "", "The modorder file to read for load order")

	_ = sortCmd.MarkFlagRequired("dir")
	_ = sortCmd.MarkFlagRequired("modorder")
}
