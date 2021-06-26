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
	"github.com/mdhender/fhdb/config"
	"github.com/mdhender/fhdb/handlers"
	"github.com/mdhender/fhdb/jwt"
	"net/http"
)

// Routes initializes all routes exposed by the Server.
func (s *Server) Routes(cfg *config.Config) error {
	for _, route := range []struct {
		method  string
		pattern string
		handler http.HandlerFunc
	}{
		{"GET", "/api/calc/mishap/:from/:to/:age/:gv", s.handleCalcMishap()},
		{"GET", "/api/version", s.handleGetVersion()},
	} {
		s.Router.HandleFunc(route.method, route.pattern, route.handler)
	}
	for _, route := range []struct {
		method  string
		pattern string
		handler http.HandlerFunc
	}{
		{"GET", "/api/flush", s.handleSave()},
		{"GET", "/api/planet/:id", s.handleGetPlanet()},
		{"GET", "/api/planets", s.handleGetPlanets()},
		{"GET", "/api/species", s.handleGetKnownSpecies()},
		{"GET", "/api/species/:id", s.handleGetSpecies()},
		{"GET", "/api/system/:id", s.handleGetSystem()},
		{"GET", "/api/systems", s.handleGetSystems()},
		{"GET", "/api/turn", s.handleGetTurn()},
		{"GET", "/api/user", s.handleGetUser()},
	} {
		s.Router.HandleFunc(route.method, route.pattern, handlers.Authenticate(route.handler, jwt.NewFactory(cfg.Server.JWT.Key)))
	}
	//s.Router.NotFound = handlers.Static("/", cfg.Server.Web.Root, true, true)
	return nil
}

// 	s.Handler = mwCORS(handlers.Authenticate(mwVersion(s.Router, s.jdb.Version), jwt.NewFactory(cfg.Server.JWT.Key)))
