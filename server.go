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
		jsonOk(w, r, s.helperGetSystems([]*System{system})[0])
	}
}

func (s *Server) handleGetSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		rsp := s.helperGetSystems(s.ds.Sorted.Systems)
		jsonOk(w, r, rsp)
	}
}

type jCoords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
type jItem struct {
	Code     string `json:"code"`
	Quantity int    `json:"qty"`
	Location string `json:"location,omitempty"`
}
type jShip struct {
	Id        string   `json:"id"`
	Landed    bool     `json:"landed"`
	Orbiting  bool     `json:"orbiting"`
	DeepSpace bool     `json:"deep_space"`
	Hiding    bool     `json:"hiding"`
	Inventory []*jItem `json:"inventory,omitempty"`
}
type jPlanet struct {
	Id                       string   `json:"id"`
	Orbit                    int      `json:"orbit"`
	Name                     string   `json:"name"`
	HomeWorld                bool     `json:"home_world"`
	AvailablePopulationUnits int      `json:"available_population_units"`
	EconomicEfficiency       int      `json:"economic_efficiency"`
	LSN                      int      `json:"lsn"`
	MiningDifficulty         float64  `json:"mining_difficulty"`
	ProductionPenalty        int      `json:"production_penalty"`
	Ships                    []*jShip `json:"ships,omitempty"`
}
type jSystem struct {
	Id      int        `json:"id"`
	Coords  jCoords    `json:"coords"`
	Planets []*jPlanet `json:"planets"`
	Scanned int        `json:"scanned"`
	Ships   []*jShip   `json:"ships,omitempty"`
	Visited bool       `json:"visited"`
	Links   struct {
		Self string `json:"self"`
	} `json:"links"`
}

func (s *Server) helperGetSystems(all []*System) []*jSystem {
	systems := []*jSystem{}
	for _, system := range all {
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
