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

package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func (ds *Store) Read(root string) error {
	log.Printf("loading %q\n", root)
	ds.Planets = make(map[string]*Planet)
	ds.Species = make(map[string]*Species)
	ds.Systems = make(map[string]*System)

	filename := filepath.Join(root, "store.json")
	var data struct {
		Species map[string]*rSpecies
		Systems []*rSystem
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	} else if err = json.Unmarshal(b, &data); err != nil {
		return err
	}

	// inject system coordinates into planets and ships
	for _, s := range data.Systems {
		fields := strings.Split(s.Id, " ")
		if len(fields) != 3 {
			return fmt.Errorf("system %s: invalid system id", s.Id)
		}
		var err error
		if s.X, err = strconv.Atoi(fields[0]); err != nil {
			return fmt.Errorf("system %s: invalid x: %+v", s.Id, err)
		} else if s.Y, err = strconv.Atoi(fields[1]); err != nil {
			return fmt.Errorf("system %s: invalid y: %+v", s.Id, err)
		} else if s.Z, err = strconv.Atoi(fields[2]); err != nil {
			return fmt.Errorf("system %s: invalid z: %+v", s.Id, err)
		}
		log.Println(fields, s.X, s.Y, s.Z)
		for _, p := range s.Planets {
			p.X, p.Y, p.Z = s.X, s.Y, s.Z
		}
	}

	// reject input if we find duplicate systems, planets, or names
	systems := make(map[string]*rSystem)
	planets := make(map[string]*rPlanet)
	species := make(map[int]*rSpecies)
	names := make(map[string]bool)
	for name, sp := range data.Species {
		if _, ok := species[sp.Id]; ok {
			return fmt.Errorf("species %q: id %d: found species with same id", name, sp.Id)
		}
		for k := range sp.TechLevels {
			switch k {
			case "biology":
			case "gravitics":
			case "life_support":
			case "manufacturing":
			case "military":
			case "mining":
			default:
				return fmt.Errorf("species %q: invalid tech %q", name, k)
			}
		}
	}
	for _, s := range data.Systems {
		s.Id = fmt.Sprintf("%d %d %d", s.X, s.Y, s.Z)
		if _, ok := systems[s.Id]; ok {
			return fmt.Errorf("system %q: found system with same coordinates", s.Id)
		}
		systems[s.Id] = s

		for _, p := range s.Planets {
			p.Id = fmt.Sprintf("%d %d %d %d", p.X, p.Y, p.Z, p.Orbit)
			if _, ok := planets[p.Id]; ok {
				return fmt.Errorf("planet %q: found planet with same orbit", p.Id)
			}
			planets[p.Id] = p
			if p.Name != "" {
				if names[p.Name] {
					return fmt.Errorf("planet %q: name %q: found another thing with same name", p.Id, p.Name)
				}
				planets[p.Name] = p
			}
		}
	}

	// move the data into our in-memory data structures
	for name, sp := range data.Species {
		species := &Species{
			Id:            sp.Id,
			Name:          name,
			EconomicUnits: sp.EconomicUnits,
			TechLevels:    make(map[string]*TechLevel),
		}
		for name, tl := range sp.TechLevels {
			v := &TechLevel{Value: tl.Value}
			switch name {
			case "biology":
				species.TechLevels["BI"] = v
			case "gravitics":
				species.TechLevels["GV"] = v
			case "life_support":
				species.TechLevels["LS"] = v
			case "manufacturing":
				species.TechLevels["MA"] = v
			case "military":
				species.TechLevels["ML"] = v
			case "mining":
				species.TechLevels["MI"] = v
			}
		}
		ds.Species[species.Name] = species
		ds.Sorted.Species = append(ds.Sorted.Species, species)
	}
	for _, v := range data.Systems {
		log.Printf("loading system %s\n", v.Id)
		sys := &System{
			Id:      v.Id,
			Empty:   v.Empty,
			Scanned: v.Scanned,
			Visited: v.Visited,
		}
		sys.X, sys.Y, sys.Z = v.X, v.Y, v.Z
		if _, ok := ds.Systems[sys.Id]; ok {
			return fmt.Errorf("system %s: duplicate system", sys.Id)
		}
		ds.Systems[sys.Id] = sys
		ds.Sorted.Systems = append(ds.Sorted.Systems, sys)

		for _, vp := range v.Planets {
			log.Printf("loading system %q planet %d\n", v.Id, vp.Orbit)
			pla := &Planet{
				Id:                       vp.Id,
				System:                   sys,
				Orbit:                    vp.Orbit,
				Name:                     vp.Name,
				HomeWorld:                vp.HomeWorld,
				AvailablePopulationUnits: vp.AvailablePopulationUnits,
				EconomicEfficiency:       vp.EconomicEfficiency,
				Inventory:                make(map[string]*Item),
				LSN:                      vp.LSN,
				MiningDifficulty:         vp.MiningDifficulty,
				ProductionPenalty:        vp.ProductionPenalty,
				Shipyards:                vp.Shipyards,
				Ships:                    make(map[string]*Ship),
			}
			if pla.Id == pla.Name {
				return fmt.Errorf("system %s: planet: %s: id == name", sys.Id, pla.Id)
			} else if _, ok := ds.Planets[pla.Id]; ok {
				return fmt.Errorf("system %s: planet: %s: duplicate location", sys.Id, pla.Id)
			} else if _, ok := ds.Planets[pla.Name]; ok {
				return fmt.Errorf("system %s: planet: %s: duplicate name %q", sys.Id, pla.Id, pla.Name)
			}
			for _, np := range sys.Planets {
				if np.Orbit == pla.Orbit {
					return fmt.Errorf("system %s: planet %s: duplicate orbit %d", sys.Id, pla.Id, pla.Orbit)
				}
			}
			for code, item := range vp.Inventory {
				pla.Inventory[code] = &Item{
					Code:     code,
					Location: item.Location,
					Quantity: item.Quantity,
				}
			}
			for name, ship := range vp.Ships {
				sh := &Ship{
					Id:              name,
					Age:             ship.Age,
					Capacity:        ship.Capacity,
					FTL:             !ship.SubLight,
					Inventory: make(map[string]*Item),
					Location:        ship.Location,
					MaintenanceCost: ship.MaintenanceCost,
				}
				if ship.Location == "(C)" {
					sh.Landed = true
				} else if ship.Location != "" {
					loc, orb := ship.Location[:len(ship.Location)-1], ship.Location[len(ship.Location)-1:]
					if orb != fmt.Sprintf("%d", pla.Orbit) {
						return fmt.Errorf("system %s: planet %s: ship %q: invalid orbit %q", sys.Id, pla.Id, name, orb)
					}
					switch loc {
					case "D":
						sh.DeepSpace = true
					case "FJ":
						sh.DeepSpace = true
						sh.ForcedJump = true
					case "L":
						sh.Landed = true
					case "O":
						sh.Orbiting = true
					case "WD":
						sh.DeepSpace = true
						sh.WithdrewFromCombat = true
					}
				}
				pla.Ships[name] = sh
			}

			ds.Sorted.Planets = append(ds.Sorted.Planets, pla)
			sys.Planets = append(sys.Planets, pla)
			ds.Planets[pla.Id] = pla
			if pla.Name != "" {
				ds.Planets[pla.Name] = pla
			}
		}
	}

	for i := 0; i < len(ds.Sorted.Planets); i++ {
		for k := i + 1; k < len(ds.Sorted.Planets); k++ {
			if ds.Sorted.Planets[k].Less(ds.Sorted.Planets[i]) {
				ds.Sorted.Planets[i], ds.Sorted.Planets[k] = ds.Sorted.Planets[k], ds.Sorted.Planets[i]
			}
		}
	}
	for i := 0; i < len(ds.Sorted.Species); i++ {
		for k := i + 1; k < len(ds.Sorted.Species); k++ {
			if ds.Sorted.Species[k].Id < ds.Sorted.Species[i].Id {
				ds.Sorted.Species[i], ds.Sorted.Species[k] = ds.Sorted.Species[k], ds.Sorted.Species[i]
			}
		}
	}
	for i := 0; i < len(ds.Sorted.Systems); i++ {
		for k := i + 1; k < len(ds.Sorted.Systems); k++ {
			if ds.Sorted.Systems[k].Less(ds.Sorted.Systems[i]) {
				ds.Sorted.Systems[i], ds.Sorted.Systems[k] = ds.Sorted.Systems[k], ds.Sorted.Systems[i]
			}
		}
	}

	log.Printf("loaded %6d species\n", len(ds.Species))
	log.Printf("loaded %6d systems\n", len(ds.Systems))
	log.Printf("loaded %6d planets\n", len(ds.Planets))
	return nil
}

type rSystem struct {
	Id        string     `json:"id"`
	X         int        `json:"x"`
	Y         int        `json:"y"`
	Z         int        `json:"z"`
	Empty     bool       `json:"empty,omitempty"`
	Inventory []*rItem   `json:"inventory,omitempty"`
	Planets   []*rPlanet `json:"planets,omitempty"`
	Scanned   int        `json:"scanned,omitempty"`
	Ships     []*rShip   `json:"ships,omitempty"`
	Visited   bool       `json:"visited,omitempty"`
}

type rPlanet struct {
	Id                       string            `json:"id"`
	X                        int               `json:"-"`
	Y                        int               `json:"-"`
	Z                        int               `json:"-"`
	Orbit                    int               `json:"orbit"`
	Name                     string            `json:"name,omitempty"`
	HomeWorld                bool              `json:"home_world,omitempty"`
	AvailablePopulationUnits int               `json:"available_population_units,omitempty"`
	EconomicEfficiency       int               `json:"economic_efficiency,omitempty"`
	Inventory                map[string]*rItem `json:"inventory,omitempty"`
	LSN                      int               `json:"lsn,omitempty"`
	MiningDifficulty         float64           `json:"mining_difficulty,omitempty"`
	ProductionPenalty        int               `json:"production_penalty,omitempty"`
	Ships                    map[string]*rShip `json:"ships,omitempty"`
	Shipyards                int               `json:"shipyards,omitempty"`
}

type rShip struct {
	Id              string            `json:"id"`
	Age             int           `json:"age"`
	Location        string            `json:"location,omitempty"`
	Capacity        int               `json:"capacity,omitempty"`
	MaintenanceCost int               `json:"maintenance_cost,omitempty"`
	SubLight        bool              `json:"sub_light,omitempty"`
	Inventory       map[string]*rItem `json:"inventory,omitempty"`
}

type rItem struct {
	Quantity int    `json:"qty"`
	Location string `json:"location,omitempty"`
}

type rSpecies struct {
	Id            int                    `json:"id"`
	EconomicUnits int                    `json:"economic_units,omitempty"`
	TechLevels    map[string]*rTechLevel `json:"tech_levels"`
}

type rTechLevel struct {
	Value int `json:"value"`
}
