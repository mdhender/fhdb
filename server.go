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
	"fmt"
	"github.com/mdhender/fhdb/way"
	"net/http"
	"strconv"
)

type Server struct {
	http.Server
	DtFmt     string // format string for timestamps in responses
	Router    *way.Router
	Templates struct {
		Root string
	}
	debug bool
	ds    *Store
}

func (s *Server) handleGetPlanet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		id := way.Param(r.Context(), "id")
		if len(id) < 5 || len(id) > 32 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		planet, ok := s.ds.Planets[id]
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if planet == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		type data struct {
			Id                       string  `json:"id"`
			Orbit                    int     `json:"orbit"`
			Name                     string  `json:"name"`
			HomeWorld                bool    `json:"home_world"`
			AvailablePopulationUnits int     `json:"available_population_units"`
			EconomicEfficiency       int     `json:"economic_efficiency"`
			LSN                      int     `json:"lsn"`
			MiningDifficulty         float64 `json:"mining_difficulty"`
			ProductionPenalty        int     `json:"production_penalty"`
		}

		rsp := data{
			Id:                       planet.Id,
			Orbit:                    planet.Orbit,
			Name:                     planet.Name,
			HomeWorld:                planet.HomeWorld,
			AvailablePopulationUnits: planet.AvailablePopulationUnits,
			EconomicEfficiency:       planet.EconomicEfficiency,
			LSN:                      planet.LSN,
			MiningDifficulty:         planet.MiningDifficulty,
			ProductionPenalty:        planet.ProductionPenalty,
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetPlanets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		type data struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			HomeWorld bool   `json:"home_world"`
			Links     struct {
				Self string `json:"self"`
			} `json:"links"`
		}

		var rsp []data
		for _, planet := range s.ds.Sorted.Planets {
			d := data{
				Id:        planet.Id,
				Name:      planet.Name,
				HomeWorld: planet.HomeWorld,
			}
			d.Links.Self = fmt.Sprintf("/api/planets/%s", d.Id)
			rsp = append(rsp, d)
		}

		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		id, err := strconv.Atoi(way.Param(r.Context(), "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		system, ok := s.ds.Systems[id]
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if system == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		type coords struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		}
		type planet struct {
			Id                       string  `json:"id"`
			Orbit                    int     `json:"orbit"`
			Name                     string  `json:"name"`
			HomeWorld                bool    `json:"home_world"`
			AvailablePopulationUnits int     `json:"available_population_units"`
			EconomicEfficiency       int     `json:"economic_efficiency"`
			LSN                      int     `json:"lsn"`
			MiningDifficulty         float64 `json:"mining_difficulty"`
			ProductionPenalty        int     `json:"production_penalty"`
		}
		type data struct {
			Id      int       `json:"id"`
			Coords  coords    `json:"coords"`
			Planets []*planet `json:"planets"`
			Scanned int       `json:"scanned,omitempty"`
			Visited bool      `json:"visited"`
		}

		rsp := data{
			Id:      system.Id,
			Coords:  coords{X: system.Coords.X, Y: system.Coords.Y, Z: system.Coords.Z},
			Planets: []*planet{},
			Scanned: system.Scanned,
			Visited: system.TaggedAsVisited(),
		}
		for _, p := range system.Planets {
			rsp.Planets = append(rsp.Planets, &planet{
				Id:                       p.Id,
				Orbit:                    p.Orbit,
				Name:                     p.Name,
				HomeWorld:                p.HomeWorld,
				AvailablePopulationUnits: p.AvailablePopulationUnits,
				EconomicEfficiency:       p.EconomicEfficiency,
				LSN:                      p.LSN,
				MiningDifficulty:         p.MiningDifficulty,
				ProductionPenalty:        p.ProductionPenalty,
			})
		}

		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		type coords struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		}
		type planet struct {
			Id                       string  `json:"id"`
			Orbit                    int     `json:"orbit"`
			Name                     string  `json:"name"`
			HomeWorld                bool    `json:"home_world"`
			AvailablePopulationUnits int     `json:"available_population_units"`
			EconomicEfficiency       int     `json:"economic_efficiency"`
			LSN                      int     `json:"lsn"`
			MiningDifficulty         float64 `json:"mining_difficulty"`
			ProductionPenalty        int     `json:"production_penalty"`
		}
		type data struct {
			Id      int    `json:"id"`
			Coords  coords `json:"coords"`
			Planets []*planet `json:"planets"`
			Scanned int    `json:"scanned"`
			Visited bool   `json:"visited"`
			Links   struct {
				Self string `json:"self"`
			} `json:"links"`
		}

		var rsp []data
		for _, system := range s.ds.Sorted.Systems {
			d := data{
				Id:      system.Id,
				Coords:  coords{X: system.Coords.X, Y: system.Coords.Y, Z: system.Coords.Z},
				Planets: []*planet{},
				Scanned: system.Scanned,
				Visited: system.TaggedAsVisited(),
			}
			d.Links.Self = fmt.Sprintf("/api/systems/%d", d.Id)
			for _, p := range system.Planets {
				d.Planets = append(d.Planets, &planet{
					Id:                       p.Id,
					Orbit:                    p.Orbit,
					Name:                     p.Name,
					HomeWorld:                p.HomeWorld,
					AvailablePopulationUnits: p.AvailablePopulationUnits,
					EconomicEfficiency:       p.EconomicEfficiency,
					LSN:                      p.LSN,
					MiningDifficulty:         p.MiningDifficulty,
					ProductionPenalty:        p.ProductionPenalty,
				})
			}
			rsp = append(rsp, d)
		}

		jsonOk(w, r, rsp)
	}
}

type jCoords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
type jPlanet struct {
	Id                       string  `json:"id"`
	Orbit                    int     `json:"orbit"`
	Name                     string  `json:"name"`
	HomeWorld                bool    `json:"home_world"`
	AvailablePopulationUnits int     `json:"available_population_units"`
	EconomicEfficiency       int     `json:"economic_efficiency"`
	LSN                      int     `json:"lsn"`
	MiningDifficulty         float64 `json:"mining_difficulty"`
	ProductionPenalty        int     `json:"production_penalty"`
}
type jSystem struct {
	Id      int    `json:"id"`
	Coords  jCoords `json:"coords"`
	Planets []*jPlanet `json:"planets"`
	Scanned int    `json:"scanned"`
	Visited bool   `json:"visited"`
	Links   struct {
		Self string `json:"self"`
	} `json:"links"`
}

func (s *Server) helperGetSystems() []*jSystem {
	var allSystems []*System
	if s.ds != nil {
		allSystems = s.ds.Sorted.Systems
	}

	systems := []*jSystem{}
	for _, system := range allSystems {
		d := &jSystem{
			Id:      system.Id,
			Coords:  jCoords{X: system.Coords.X, Y: system.Coords.Y, Z: system.Coords.Z},
			Planets: []*jPlanet{},
			Scanned: system.Scanned,
			Visited: system.TaggedAsVisited(),
		}
		d.Links.Self = fmt.Sprintf("/api/systems/%d", d.Id)
		for _, p := range system.Planets {
			d.Planets = append(d.Planets, &jPlanet{
				Id:                       p.Id,
				Orbit:                    p.Orbit,
				Name:                     p.Name,
				HomeWorld:                p.HomeWorld,
				AvailablePopulationUnits: p.AvailablePopulationUnits,
				EconomicEfficiency:       p.EconomicEfficiency,
				LSN:                      p.LSN,
				MiningDifficulty:         p.MiningDifficulty,
				ProductionPenalty:        p.ProductionPenalty,
			})
		}
		systems = append(systems, d)
	}
	return systems
}