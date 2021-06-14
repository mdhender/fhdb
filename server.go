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
	"github.com/mdhender/fhdb/handlers"
	"github.com/mdhender/fhdb/ports"
	"github.com/mdhender/fhdb/store"
	"github.com/mdhender/fhdb/store/jsondb"
	"github.com/mdhender/fhdb/way"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	http.Server
	Data      string // path to read and write data files
	DtFmt     string // format string for timestamps in responses
	Router    *way.Router
	Templates struct {
		Root string
	}
	debug bool
	ds    *store.Store
	jdb   *jsondb.Store
}

func (s *Server) handleCalcMishap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fields := strings.Split(way.Param(r.Context(), "from"), " ")
		if len(fields) != 3 {
			http.Error(w, "invalid from", http.StatusBadRequest)
		}
		fromX, err := strconv.Atoi(fields[0])
		if err != nil {
			http.Error(w, "invalid from.x", http.StatusBadRequest)
			return
		}
		fromY, err := strconv.Atoi(fields[1])
		if err != nil {
			http.Error(w, "invalid from.y", http.StatusBadRequest)
			return
		}
		fromZ, err := strconv.Atoi(fields[2])
		if err != nil {
			http.Error(w, "invalid from.z", http.StatusBadRequest)
			return
		}
		fields = strings.Split(way.Param(r.Context(), "to"), " ")
		if len(fields) != 3 {
			http.Error(w, "invalid to", http.StatusBadRequest)
		}
		toX, err := strconv.Atoi(fields[0])
		if err != nil {
			http.Error(w, "invalid to.x", http.StatusBadRequest)
			return
		}
		toY, err := strconv.Atoi(fields[1])
		if err != nil {
			http.Error(w, "invalid to.y", http.StatusBadRequest)
			return
		}
		toZ, err := strconv.Atoi(fields[2])
		if err != nil {
			http.Error(w, "invalid to.z", http.StatusBadRequest)
			return
		}
		mishapAge, err := strconv.Atoi(way.Param(r.Context(), "age"))
		if err != nil {
			http.Error(w, "invalid age", http.StatusBadRequest)
			return
		} else if mishapAge < 0 {
			http.Error(w, "age must be a non-negative integer", http.StatusBadRequest)
			return
		}
		mishapGV, err := strconv.Atoi(way.Param(r.Context(), "gv"))
		if err != nil {
			http.Error(w, "invalid gv", http.StatusBadRequest)
			return
		} else if mishapGV < 1 {
			http.Error(w, "gv must be a positive integer", http.StatusBadRequest)
			return
		}
		deltaX := fromX - toX
		deltaY := fromY - toY
		deltaZ := fromZ - toZ
		mishapChance := 100 * (deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ) / mishapGV

		if mishapChance > 10000 {
			mishapChance = 10000
		} else if mishapAge > 0 {
			// Add aging effect
			successChance := 10000 - mishapChance
			successChance -= (2 * mishapAge * successChance) / 100
			if successChance < 0 {
				successChance = 0
			}
			mishapChance = 10000 - successChance
		}
		rsp := ports.MishapResponse{
			Age:          mishapAge,
			GV:           mishapGV,
			MishapChance: float64(mishapChance)/100,
		}
		rsp.From.X = fromX
		rsp.From.Y = fromY
		rsp.From.Z = fromZ
		rsp.To.X = toX
		rsp.To.Y = toY
		rsp.To.Z = toZ
		jsonOk(w, r, rsp)
	}
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

func (s *Server) handleGetSpecies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		id := way.Param(r.Context(), "id")
		if id == "" {
			id = "*"
		}
		type species struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}
		var rsp []*species
		for _, e := range s.ds.SpeciesMap(id) {
			rsp = append(rsp, &species{
				Id:   e.Id,
				Name: e.Name,
			})
		}
		if rsp == nil {
			rsp = []*species{}
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := handlers.GetSession(r)
		if sess == nil || !sess.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		id := way.Param(r.Context(), "id")
		rsp, err := s.jdb.GetSystem(id, sess.SpeciesId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := handlers.GetSession(r)
		if sess == nil || !sess.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rsp, err := s.jdb.GetSystems(sess.SpeciesId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetTurn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := handlers.GetSession(r)
		if sess == nil || !sess.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rsp, err := s.jdb.GetTurnNumber()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := handlers.GetSession(r)
		if sess == nil || !sess.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rsp, err := s.jdb.GetUser(sess.SpeciesId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleGetVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rsp, err := s.jdb.GetVersion()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonOk(w, r, rsp)
	}
}

func (s *Server) handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.ds == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		s.ds.Write(s.Data)
		w.WriteHeader(http.StatusOK)
	}
}

type jCoords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
type jItem struct {
	Code     string `json:"code"`
	Location string `json:"location"`
	Quantity int    `json:"qty"`
}
type jShip struct {
	Id                 string   `json:"id"`
	Capacity           int      `json:"capacity"`
	DeepSpace          bool     `json:"deep_space"`
	ForcedJump         bool     `json:"forced_jump"`
	FTL                bool     `json:"ftl"`
	Hiding             bool     `json:"hiding"`
	Landed             bool     `json:"landed"`
	Location           string   `json:"location"`
	MaintenanceCost    int      `json:"maintenance_cost"`
	Orbiting           bool     `json:"orbiting"`
	WithdrewFromCombat bool     `json:"withdrew_from_combat"`
	Inventory          []*jItem `json:"inventory"`
}
type jPlanet struct {
	Id                       string   `json:"id"`
	Orbit                    int      `json:"orbit"`
	Name                     string   `json:"name"`
	HomeWorld                bool     `json:"home_world"`
	AvailablePopulationUnits int      `json:"available_population_units"`
	EconomicEfficiency       int      `json:"economic_efficiency"`
	Inventory                []*jItem `json:"inventory"`
	LSN                      int      `json:"lsn"`
	MiningDifficulty         float64  `json:"mining_difficulty"`
	ProductionPenalty        int      `json:"production_penalty"`
	Ships                    []*jShip `json:"ships"`
	Shipyards                int      `json:"shipyards"`
}
type jSystem struct {
	Id      string     `json:"id"`
	X       int        `json:"x"`
	Y       int        `json:"Y"`
	Z       int        `json:"z"`
	Planets []*jPlanet `json:"planets"`
	Scanned int        `json:"scanned"`
	Ships   []*jShip   `json:"ships"`
	Visited bool       `json:"visited"`
	Links   struct {
		Self string `json:"self"`
	} `json:"links"`
}

func (s *Server) helperGetSystems(all []*store.System) []*jSystem {
	var systems []*jSystem
	for _, system := range all {
		d := &jSystem{
			Id:      system.Id,
			X:       system.X,
			Y:       system.Y,
			Z:       system.Z,
			Planets: []*jPlanet{},
			Scanned: system.Scanned,
			Ships:   []*jShip{},
			Visited: system.TaggedAsVisited(),
		}
		d.Links.Self = fmt.Sprintf("/api/systems/%s", d.Id)
		for _, p := range system.Planets {
			jp := &jPlanet{
				Id:                       p.Id,
				Orbit:                    p.Orbit,
				Name:                     p.Name,
				HomeWorld:                p.HomeWorld,
				AvailablePopulationUnits: p.AvailablePopulationUnits,
				EconomicEfficiency:       p.EconomicEfficiency,
				Inventory:                []*jItem{},
				LSN:                      p.LSN,
				MiningDifficulty:         p.MiningDifficulty,
				ProductionPenalty:        p.ProductionPenalty,
				Ships:                    []*jShip{},
				Shipyards:                p.Shipyards,
			}
			for code, inv := range p.Inventory {
				if inv.Quantity == 0 {
					continue
				}
				jp.Inventory = append(jp.Inventory, &jItem{
					Code:     code,
					Location: inv.Location,
					Quantity: inv.Quantity,
				})
			}
			for _, ship := range p.Ships {
				js := &jShip{
					Id:        ship.Id,
					Inventory: []*jItem{},
				}
				jp.Ships = append(jp.Ships, js)
			}
			d.Planets = append(d.Planets, jp)
		}
		systems = append(systems, d)
	}
	if systems == nil {
		return []*jSystem{}
	}
	return systems
}
