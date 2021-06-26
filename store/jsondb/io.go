/*
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package jsondb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func Read(filename string) (*Store, error) {
	var ds Store
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &ds); err != nil {
		return nil, err
	}

	// validate planet data
	for _, planet := range ds.Planets {
		if planet.Id < 1 {
			return nil, fmt.Errorf("invalid planet id %d", planet.Id)
		}
		var total int
		for gas, percentage := range planet.Gases {
			switch gas {
			case "NH3": // Ammonia
			case "CO2": // Carbon Dioxide
			case "Cl2": // Chlorine
			case "F2": // Fluorine
			case "He": // Helium
			case "H2": // Hydrogen
			case "HCl": // Hydrogen Chloride
			case "H2S": // Hydrogen Sulfide
			case "CH4": // Methane
			case "N2": // Nitrogen
			case "O2": // Oxygen
			case "SO2": // Sulfur Dioxide
			case "H2O": // Water or Steam
			default:
				return nil, fmt.Errorf("unknown gas %q on planet %d", gas, planet.Id)
			}
			if percentage < 1 || percentage > 100 {
				return nil, fmt.Errorf("invalid percentage %d for gas %q on planet %d", percentage, gas, planet.Id)
			}
			total += percentage
		}
		if total > 0 && total != 100 {
			return nil, fmt.Errorf("invalid percentage %d for gases on planet %d", total, planet.Id)
		}
	}

	// validate and normalize species data
	for key, sp := range ds.Species {
		if sp.Id < 1 {
			return nil, fmt.Errorf("invalid species id %d", sp.Id)
		}
		if key != fmt.Sprintf("SP%02d", sp.Id) {
			return nil, fmt.Errorf("invalid key %q for species %d", key, sp.Id)
		}
		sp.Key = key
		sp.Aliens = make(map[int]string)

		// parse Contact, then Ally, then Neutral, then Enemy on the off-chance
		// that someone edits the JSON and assigns the species to multiple categories.
		for _, a := range sp.Contacts {
			if !strings.HasPrefix(a, "SP") {
				continue
			}
			id, err := strconv.Atoi(a[2:])
			if err != nil {
				continue
			} else if id < 1 || MAX_SPECIES < id {
				continue
			}
			sp.Aliens[id] = "neutral"
		}
		for _, a := range sp.Allies {
			if !strings.HasPrefix(a, "SP") {
				continue
			}
			id, err := strconv.Atoi(a[2:])
			if err != nil {
				continue
			} else if id < 1 || MAX_SPECIES < id {
				continue
			}
			sp.Aliens[id] = "ally"
		}
		for _, a := range sp.Enemies {
			if !strings.HasPrefix(a, "SP") {
				continue
			}
			id, err := strconv.Atoi(a[2:])
			if err != nil {
				continue
			} else if id < 1 || MAX_SPECIES < id {
				continue
			}
			sp.Aliens[id] = "enemy"
		}
	}

	return &ds, nil
}

func (ds *Store) Write(filename string) error {
	b, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return err
	}
	if filename == "*stdout*" {
		fmt.Println(string(b))
		return nil
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
