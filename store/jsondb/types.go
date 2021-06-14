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

type Store struct {
	Version string `json:"version"`

	Galaxy  struct {
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
	} `json:"galaxy"`

	Systems []*struct {
		Id         int      `json:"id"`
		Key        string   `json:"key"`
		Coords     Coords   `json:"coords"`
		Type       string   `json:"type"`
		Color      string   `json:"color"`
		Size       int      `json:"size"`
		HomeSystem bool     `json:"home_system"`
		Planets    []int    `json:"planets"`
		VisitedBy  []string `json:"visited_by"`
		Message    int      `json:"message"`
	} `json:"systems"`

	Planets []*struct {
		Id               int            `json:"id"`
		TemperatureClass int            `json:"temperature_class"`
		PressureClass    int            `json:"pressure_class"`
		Gases            map[string]int `json:"gases"`
		Diameter         int            `json:"diameter"`
		Gravity          int            `json:"gravity"`
		MiningDifficulty int            `json:"mining_difficulty"`
		EconEfficiency   int            `json:"econ_efficiency"`
		MdIncrease       int            `json:"md_increase"`
		Message          int            `json:"message"`
	} `json:"planets"`

	Species map[string]*struct {
		Id         int    `json:"id"`
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
		EconUnits        int `json:"econ_units"`
		HpOriginalBase   int `json:"hp_original_base"`
		FleetCost        int `json:"fleet_cost"`
		FleetPercentCost int `json:"fleet_percent_cost"`
		NamedPlanets     map[string]*struct {
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
		} `json:"namplas"`
		Ships map[string]*struct {
			Id             int            `json:"id"`
			Name           string         `json:"name"`
			Location       string         `json:"location"`
			Coords         Coords         `json:"coords"`
			Orbit          int            `json:"orbit"`
			Class          string         `json:"class"`
			Type           string         `json:"type"`
			Tonnage        int            `json:"tonnage"`
			Age            int            `json:"age"`
			Status         string         `json:"status"`
			Dest           Coords         `json:"dest"`
			LoadingPoint   int            `json:"loading_point"`
			UnloadingPoint int            `json:"unloading_point"`
			RemainingCost  int            `json:"remaining_cost"`
			Message        int            `json:"message"`
			Inventory      map[string]int `json:"inventory"`
		} `json:"ships"`
	} `json:"species"`

	Commands map[string]string `json:"commands"`

	Items map[string]*struct {
		Name      string         `json:"name"`
		Cost      int            `json:"cost"`
		Tech      map[string]int `json:"tech"`
		CarryCost int            `json:"carry_cost"`
	} `json:"items"`

	Ships map[string]*struct {
		Class               string `json:"class"`
		MinManufactingLevel int    `json:"min_ma"`
		CostFtl             int    `json:"cost_ftl"`
		CostSublight        int    `json:"cost_sublight"`
		Tonnage             int    `json:"tonnage"`
		CarryingCapacity    int    `json:"carrying_capacity"`
	} `json:"ships"`

	Tech map[string]string `json:"tech"`
}

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type GasMinMax struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Technology struct {
	Level     int `json:"level"`
	Init      int `json:"init"`
	Knowledge int `json:"knowledge"`
	Xp        int `json:"xp"`
}
