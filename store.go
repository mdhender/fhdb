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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Store struct {
	Systems map[int]*System
	Planets map[string]*Planet
	Sorted  struct {
		Systems []*System
		Planets []*Planet
	}
}

type Item struct {
	Code     string
	Quantity int
	Location string
}

type Planet struct {
	Id                       string
	System                   *System
	Orbit                    int
	Name                     string
	HomeWorld                bool
	AvailablePopulationUnits int
	EconomicEfficiency       int
	LSN                      int
	MiningDifficulty         float64
	ProductionPenalty        int
	sortKey                  string
}

type Ship struct {
	Id string
}

type System struct {
	Id     int
	Coords struct {
		X int
		Y int
		Z int
	}
	Inventory []*Item
	Planets   []*Planet
	Scanned   int
	Shipyards int
	Ships     []*Ship
	Visited   bool
}

func (ds *Store) Load(root string) error {
	log.Printf("loading %q\n", root)
	if ds.Systems == nil {
		ds.Systems = make(map[int]*System)
	}
	if ds.Planets == nil {
		ds.Planets = make(map[string]*Planet)
	}

	type item struct {
		Code     string `json:"code"`
		Quantity int    `json:"qty"`
		Location string `json:"location,omitempty"`
	}
	type ship struct {
		Id        string  `json:"id"`
		Landed    bool    `json:"landed"`
		Orbiting  bool    `json:"orbiting"`
		DeepSpace bool    `json:"deep_space"`
		Hiding    bool    `json:"hiding"`
		Inventory []*item `json:"inventory,omitempty"`
	}
	type planet struct {
		Id                       string  `json:"id"`
		Orbit                    int     `json:"orbit"`
		Name                     string  `json:"name,omitempty"`
		HomeWorld                bool    `json:"home_world,omitempty"`
		AvailablePopulationUnits int     `json:"available_population_units,omitempty"`
		EconomicEfficiency       int     `json:"economic_efficiency,omitempty"`
		Inventory                []*item `json:"inventory,omitempty"`
		LSN                      int     `json:"lsn,omitempty"`
		MiningDifficulty         float64 `json:"mining_difficulty,omitempty"`
		ProductionPenalty        int     `json:"production_penalty,omitempty"`
		Ships                    []*ship `json:"ships,omitempty"`
	}
	type system struct {
		Id     int `json:"id"`
		Coords struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		}
		Inventory []*item   `json:"inventory,omitempty"`
		Planets   []*planet `json:"planets,omitempty"`
		Scanned   int       `json:"scanned,omitempty"`
		Shipyards int       `json:"shipyards,omitempty"`
		Visited   bool      `json:"visited,omitempty"`
		Ships     []*ship   `json:"ships,omitempty"`
	}

	b, err := ioutil.ReadFile(filepath.Join(root, "store.json"))
	if err != nil {
		return err
	}
	var data struct {
		Systems []*system `json:"systems"`
		Ships   []*ship   `json:"ships"`
	}
	if err = json.Unmarshal(b, &data); err != nil {
		return err
	}

	for _, v := range data.Systems {
		log.Printf("loading system %d\n", v.Id)
		sys := &System{
			Id:        v.Id,
			Inventory: nil,
			Planets:   nil,
			Scanned:   v.Scanned,
			Shipyards: v.Shipyards,
			Visited:   v.Visited,
		}
		sys.Coords.X, sys.Coords.Y, sys.Coords.Z = v.Coords.X, v.Coords.Y, v.Coords.Z
		if _, ok := ds.Systems[sys.Id]; ok {
			return fmt.Errorf("system %d: duplicate system", sys.Id)
		}
		ds.Systems[sys.Id] = sys
		ds.Sorted.Systems = append(ds.Sorted.Systems, sys)

		for _, vp := range v.Planets {
			log.Printf("loading system %d planet %d\n", v.Id, vp.Orbit)
			pla := &Planet{
				Id:                       vp.Id,
				System:                   sys,
				Orbit:                    vp.Orbit,
				Name:                     vp.Name,
				HomeWorld:                vp.HomeWorld,
				AvailablePopulationUnits: vp.AvailablePopulationUnits,
				EconomicEfficiency:       vp.EconomicEfficiency,
				LSN:                      vp.LSN,
				MiningDifficulty:         vp.MiningDifficulty,
				ProductionPenalty:        vp.ProductionPenalty,
				sortKey:                  fmt.Sprintf("%4d %4d %4d %2d", sys.Coords.X, sys.Coords.Y, sys.Coords.Z, vp.Orbit),
			}
			if vp.Id == "" {
				pla.Id = fmt.Sprintf("%d %d %d %d", sys.Coords.X, sys.Coords.Y, sys.Coords.Z, vp.Orbit)
			}
			if pla.Id == pla.Name {
				return fmt.Errorf("system: %d: planet: %s: id == name", sys.Id, pla.Id)
			} else if _, ok := ds.Planets[pla.Id]; ok {
				return fmt.Errorf("system: %d: planet: %s: duplicate location", sys.Id, pla.Id)
			} else if _, ok := ds.Planets[pla.Name]; ok {
				return fmt.Errorf("system: %d: planet: %s: duplicate name %q", sys.Id, pla.Id, pla.Name)
			}
			for _, np := range sys.Planets {
				if np.Orbit == pla.Orbit {
					return fmt.Errorf("system %d: planet %s: duplicate orbit %d", sys.Id, pla.Id, pla.Orbit)
				}
			}
			ds.Sorted.Planets = append(ds.Sorted.Planets, pla)
			sys.Planets = append(sys.Planets, pla)
			ds.Planets[pla.Id] = pla
			if pla.Name != "" {
				ds.Planets[pla.Name] = pla
			}
		}
		for i := 0; i < len(sys.Planets); i++ {
			for j := i + 1; j < len(sys.Planets); j++ {
				if sys.Planets[j].Orbit < sys.Planets[i].Orbit {
					sys.Planets[i], sys.Planets[j] = sys.Planets[j], sys.Planets[i]
				}
			}
		}
	}

	for i := 0; i < len(ds.Sorted.Systems); i++ {
		for j := i + 1; j < len(ds.Sorted.Systems); j++ {
			if !ds.Sorted.Systems[i].Less(ds.Sorted.Systems[j]) {
				ds.Sorted.Systems[i], ds.Sorted.Systems[j] = ds.Sorted.Systems[j], ds.Sorted.Systems[i]
			}
		}
	}

	for i := 0; i < len(ds.Sorted.Planets); i++ {
		for j := i + 1; j < len(ds.Sorted.Planets); j++ {
			if !ds.Sorted.Planets[i].Less(ds.Sorted.Planets[j]) {
				ds.Sorted.Planets[i], ds.Sorted.Planets[j] = ds.Sorted.Planets[j], ds.Sorted.Planets[i]
			}
		}
	}

	log.Printf("loaded %6d systems\n", len(ds.Systems))
	log.Printf("loaded %6d planets\n", len(ds.Planets))
	return nil
}

// Less is a helper for sorting
func (s *System) Less(ss *System) bool {
	return s.Id < ss.Id
}

// TaggedAsVisited checks to see if the player manually tagged the system as visited.
func (s *System) TaggedAsVisited() bool {
	return s.Visited && s.Scanned == 0
}

func (p *Planet) Less(pp *Planet) bool {
	return p.sortKey < pp.sortKey
}
