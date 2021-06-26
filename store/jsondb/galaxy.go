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

type Galaxy struct {
	TurnNumber        int `json:"turn_number"`
	DNumSpecies       int `json:"d_num_species"`
	NumSpecies        int `json:"num_species"`
	Radius            int `json:"radius"`
	MinRadius         int `json:"min_radius"`
	MaxRadius         int `json:"max_radius"`
	StdNumStars       int `json:"std_num_stars"`
	MinStars          int `json:"min_stars"`
	MaxStars          int `json:"max_stars"`
	StdNumSpecies     int `json:"std_num_species"`
	MinSpecies        int `json:"min_species"`
	MaxSpecies        int `json:"max_species"`
	MaxItems          int `json:"max_items"`
	MaxLocations      int `json:"max_locations"`
	MaxTransactions   int `json:"max_transactions"`
	NumCommands       int `json:"num_commands"`
	NumContactWords   int `json:"num_contact_words"`
	NumShipClasses    int `json:"num_ship_classes"`
	SizeofChar        int `json:"sizeof char"`
	SizeofInt         int `json:"sizeof int"`
	SizeofLong        int `json:"sizeof long"`
	SizeofShort       int `json:"sizeof short"`
	SizeofGalaxyData  int `json:"sizeof galaxy_data"`
	SizeofStarData    int `json:"sizeof star_data"`
	SizeofPlanetData  int `json:"sizeof planet_data"`
	SizeofNamplaData  int `json:"sizeof nampla_data"`
	SizeofSpeciesData int `json:"sizeof species_data"`
	SizeofShipData    int `json:"sizeof ship_data"`
	SizeofTransData   int `json:"sizeof trans_data"`
}
