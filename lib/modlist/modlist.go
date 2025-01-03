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
package modlist

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const NoModsError = "no valid mods found in modorder file (invalid format?)"

type ModList struct {
	Name     string
	Priority int
}

type ModListing []ModList

func Read(path string) (*ModListing, error) {
	extention := filepath.Ext(path)

	switch extention {
	case ".csv":
		return readCSV(path)
	case ".json":
		return nil, errors.New("reading JSON input is not yet implemented")
	case ".txt":
		return readTEXT(path)
	case ".xml":
		return nil, errors.New("reading XML input is not yet implemented")
	}

	return nil, fmt.Errorf("unsupported file type: %s", extention)
}

func readTEXT(path string) (*ModListing, error) {
	var modListing ModListing
	var mods []string
	var totalMods int

	// Read TXT file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		if err := fileScanner.Err(); err != nil {
			return nil, err
		}

		line := fileScanner.Text()

		// Skip comments and unchecked mods
		if line[0] != '+' {
			continue
		}

		// Skip separators, sometimes they are marked as 'checked'
		if strings.HasSuffix(line, "_separator") {
			continue
		}

		// Skip the 7dtd MO2 plugin
		if strings.Contains(line, "7D2D MO2 Plugin") {
			continue
		}

		mods = append(mods, line[1:])
	}

	totalMods = len(mods)
	if totalMods == 0 {
		return nil, errors.New(NoModsError)
	}

	for index, mod := range mods {
		var priority = totalMods - index

		modList := ModList{
			Priority: priority,
			Name:     mod,
		}
		modListing = append(modListing, modList)
	}

	return &modListing, nil
}

func readCSV(path string) (*ModListing, error) {
	var modListing ModListing

	// Read CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	modlistLines := csv.NewReader(file)

	for {
		data, err := modlistLines.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(data) != 2 {
			return nil, err
		}

		// Skip header
		if data[1] == "#Mod_Name" {
			continue
		}

		priority, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, err
		}

		modList := ModList{
			Priority: priority,
			Name:     data[1],
		}

		modListing = append(modListing, modList)
	}

	if len(modListing) == 0 {
		return nil, errors.New(NoModsError)
	}

	return &modListing, nil
}
