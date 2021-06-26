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

package memory

import "log"

func (ds *Store) Write(root string) error {
	log.Printf("not saving %q\n", root)
	//filename := filepath.Join(root, "wstore.json")
	//
	//// convert in-memory structures to json-file structures
	//var data struct {
	//	Turn    int                  `json:"turn"`
	//	Species map[string]*wSpecies `json:"species"`
	//	Systems []*wSystem           `json:"systems"`
	//}
	//data.Turn = ds.Turn
	//data.Species = make(map[string]*wSpecies)
	//for _, s := range ds.Sorted.Species {
	//	sp := &wSpecies{
	//		Id:            s.Id,
	//		EconomicUnits: s.EconomicUnits,
	//		TechLevels:    make(map[string]*wTechLevel),
	//	}
	//	for k, v := range s.TechLevels {
	//		val := &wTechLevel{Value: v.Value}
	//		switch k {
	//		case "BI":
	//			sp.TechLevels["biology"] = val
	//		case "GV":
	//			sp.TechLevels["gravitics"] = val
	//		case "LS":
	//			sp.TechLevels["life_support"] = val
	//		case "MA":
	//			sp.TechLevels["manufacturing"] = val
	//		case "MI":
	//			sp.TechLevels["mining"] = val
	//		case "ML":
	//			sp.TechLevels["military"] = val
	//		}
	//	}
	//	data.Species[s.Name] = sp
	//}
	//for _, s := range ds.Sorted.Systems {
	//	//if s.Empty && (len(s.Planets) == 0 && len(s.Ships) == 0) {
	//	//	continue
	//	//}
	//	sys := &wSystem{
	//		Id:      fmt.Sprintf("%d %d %d", s.X, s.Y, s.Z),
	//		Empty:   s.Empty,
	//		Scanned: s.Scanned,
	//		Ships:   nil,
	//		Visited: s.Visited,
	//	}
	//	for _, p := range s.Planets {
	//		pla := &wPlanet{
	//			Orbit:                    p.Orbit,
	//			Name:                     p.Name,
	//			HomeWorld:                p.HomeWorld,
	//			AvailablePopulationUnits: p.AvailablePopulationUnits,
	//			EconomicEfficiency:       p.EconomicEfficiency,
	//			Inventory:                nil,
	//			LSN:                      p.LSN,
	//			MiningDifficulty:         p.MiningDifficulty,
	//			ProductionPenalty:        p.ProductionPenalty,
	//			Ships:                    nil,
	//			Shipyards:                p.Shipyards,
	//		}
	//		sys.Planets = append(sys.Planets, pla)
	//	}
	//	data.Systems = append(data.Systems, sys)
	//}
	//b, err := json.MarshalIndent(&data, "", "  ")
	//if err != nil {
	//	return err
	//} else if err = ioutil.WriteFile(filename, b, 0644); err != nil {
	//	return err
	//}
	//
	//log.Printf("saved %6d systems\n", len(data.Systems))
	return nil
}

type wSpecies struct {
	Id            int                    `json:"id"`
	EconomicUnits int                    `json:"economic_units"`
	TechLevels    map[string]*wTechLevel `json:"tech_levels"`
}
type wTechLevel struct {
	Value int `json:"value"`
}
type wSystem struct {
	Id        string     `json:"id"`
	Empty     bool       `json:"empty,omitempty"`
	Inventory []*wItem   `json:"inventory,omitempty"`
	Planets   []*wPlanet `json:"planets,omitempty"`
	Scanned   int        `json:"scanned,omitempty"`
	Ships     []*wShip   `json:"ships,omitempty"`
	Visited   bool       `json:"visited,omitempty"`
}
type wPlanet struct {
	Orbit                    int      `json:"orbit"`
	Name                     string   `json:"name,omitempty"`
	HomeWorld                bool     `json:"home_world,omitempty"`
	AvailablePopulationUnits int      `json:"available_population_units,omitempty"`
	EconomicEfficiency       int      `json:"economic_efficiency,omitempty"`
	Inventory                []*wItem `json:"inventory,omitempty"`
	LSN                      int      `json:"lsn,omitempty"`
	MiningDifficulty         float64  `json:"mining_difficulty,omitempty"`
	ProductionPenalty        int      `json:"production_penalty,omitempty"`
	Ships                    []*wShip `json:"ships,omitempty"`
	Shipyards                int      `json:"shipyards,omitempty"`
}
type wShip struct {
	Id        string   `json:"id"`
	Landed    bool     `json:"landed"`
	Orbiting  bool     `json:"orbiting"`
	DeepSpace bool     `json:"deep_space"`
	Hiding    bool     `json:"hiding"`
	Inventory []*wItem `json:"inventory,omitempty"`
}
type wItem struct {
	Code     string `json:"code"`
	Quantity int    `json:"qty"`
	Location string `json:"location,omitempty"`
}
