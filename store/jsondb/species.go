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
	"fmt"
)

func (ds *Store) xFilterSpecies(fn func(*Species) bool) []*Species {
	var result []*Species
	for _, s := range ds.Species {
		if fn(s) {
			result = append(result, s)
		}
	}
	if result == nil {
		return []*Species{}
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Name > result[j].Name {

			}
		}
	}
	return result
}

func (sp *Species) xLess(sp2 *Species) bool {
	return sp.Key < sp2.Key
}

func xSpeciesById(roles map[string]bool, ids ...int) func(*Species) bool {
	return func(sp *Species) bool {
		for _, id := range ids {
			if id == sp.Id {
				return roles[sp.Key]
			}
		}
		return false
	}
}

func xSpeciesByName(roles map[string]bool, names ...string) func(*Species) bool {
	return func(sp *Species) bool {
		for _, name := range names {
			if name == sp.Name {
				return roles[fmt.Sprintf("SP%02d", sp.Id)]
			}
		}
		return false
	}
}

type Species struct {
	Id         int    `json:"id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	Government struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"government"`
	Homeworld struct {
		Key    string `json:"key"`
		Coords Coords `json:"coords"`
		Orbit  int    `json:"orbit"`
	} `json:"homeworld"`
	Gases struct {
		Required map[string]*GasMinMax `json:"required"`
		Poison   map[string]bool       `json:"poison"`
	} `json:"gases"`
	AutoOrders bool `json:"auto_orders"`
	Tech       struct {
		Biology       Technology `json:"BI"`
		Gravitics     Technology `json:"GV"`
		LifeSupport   Technology `json:"LS"`
		Manufacturing Technology `json:"MA"`
		Mining        Technology `json:"MI"`
		Military      Technology `json:"ML"`
	} `json:"tech"`
	BankedEconUnits  int                     `json:"econ_units"`
	HpOriginalBase   int                     `json:"hp_original_base"`
	FleetCost        int                     `json:"fleet_cost"`
	FleetPercentCost int                     `json:"fleet_percent_cost"`
	Contacts         []string                `json:"contacts"`
	Allies           []string                `json:"allies"`
	Enemies          []string                `json:"enemies"`
	NamedPlanets     map[string]*NamedPlanet `json:"namplas"`
	Ships            map[string]*Ship        `json:"ships"`
	// Aliens is a map of SPxx to AlienRelationship
	Aliens map[int]string `json:"aliens"`
}
