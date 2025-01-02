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
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileType int

const (
	CSV FileType = iota
	JSON
	TEXT
	XML
)

type ModList struct {
	Name     string
	Priority string
	Type     FileType
}

type ModListing []ModList

func Read(path string) (*ModListing, error) {
	extention := filepath.Ext(path)

	switch extention {
	case ".csv":
		return readCSV(path)
	case ".json":
		return nil, errors.New("Reading JSON input is not yet implemented")
	case ".txt":
		return nil, errors.New("ot yet implemented")
	case ".xml":
		return nil, errors.New("Reading XML input is not yet implemented")
	}

	return nil, fmt.Errorf("unsupported file type: %s", extention)
}

func readCSV(path string) (*ModListing, error) {
	var modListing ModListing

	// Read CSV file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	modlistLines := csv.NewReader(file)

	for {
		data, err := modlistLines.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if len(data) != 2 {
			log.Fatal(`Modlist file must have two columns, priority and name, in the following format: ["0000","Name"]`)
			os.Exit(1)
		}

		// Skip header
		if data[1] == "#Mod_Name" {
			continue
		}

		modList := ModList{
			Priority: data[0],
			Name:     data[1],
			Type:     CSV,
		}

		modListing = append(modListing, modList)
	}

	return &modListing, nil
}
