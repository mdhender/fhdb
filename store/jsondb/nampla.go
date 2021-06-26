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

type NamedPlanet struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Coords   Coords `json:"coords"`
	Orbit    int    `json:"orbit"`
	Status   struct {
		HomePlanet      bool `json:"home_planet"`
		Colony          bool `json:"colony"`
		Populated       bool `json:"populated"`
		MiningColony    bool `json:"mining_colony"`
		ResortColony    bool `json:"resort_colony"`
		DisbandedColony bool `json:"disbanded_colony"`
	} `json:"status"`
	Hiding       bool           `json:"hiding"`
	Hidden       int            `json:"hidden"`
	PlanetIndex  int            `json:"planet_index"`
	SiegeEff     int            `json:"siege_eff"`
	Shipyards    int            `json:"shipyards"`
	IUsNeeded    int            `json:"IUs_needed"`
	AUsNeeded    int            `json:"AUs_needed"`
	AutoIUs      int            `json:"auto_IUs"`
	AutoAUs      int            `json:"auto_AUs"`
	IUsToInstall int            `json:"IUs_to_install"`
	AUsToInstall int            `json:"AUs_to_install"`
	MiBase       int            `json:"mi_base"`
	MaBase       int            `json:"ma_base"`
	PopUnits     int            `json:"pop_units"`
	UseOnAmbush  int            `json:"use_on_ambush"`
	Message      int            `json:"message"`
	Inventory    map[string]int `json:"inventory"`
}
