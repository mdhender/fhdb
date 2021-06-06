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

type Store struct {
	Planets map[string]*Planet
	Species map[string]*Species
	Systems map[string]*System
	Sorted  struct {
		Planets []*Planet
		Species []*Species
		Systems []*System
	}
}

type System struct {
	Id        string
	X         int
	Y         int
	Z         int
	Empty     bool
	Inventory map[string]*Item
	Planets   []*Planet
	Scanned   int
	Shipyards int
	Ships     map[string]*Ship
	Visited   bool
}

type Planet struct {
	Id                       string
	System                   *System
	Orbit                    int
	Name                     string
	HomeWorld                bool
	AvailablePopulationUnits int
	EconomicEfficiency       int
	Inventory                map[string]*Item
	LSN                      int
	MiningDifficulty         float64
	ProductionPenalty        int
	Shipyards                int
	Ships                    map[string]*Ship
}

type Ship struct {
	Id                 string
	Code string
	Age                int
	Location           string
	Capacity           int
	DeepSpace          bool
	ForcedJump         bool
	FTL                bool
	Hiding             bool
	Landed             bool
	MaintenanceCost    int
	Orbiting           bool
	MALevel int
	WithdrewFromCombat bool
	Inventory          map[string]*Item
}

type Item struct {
	Code         string
	Location     string
	Quantity     int
}

type Species struct {
	Id            int
	Name          string
	EconomicUnits int
	TechLevels    map[string]*TechLevel
}

type TechLevel struct {
	Value int
}


// Less is a helper for sorting
func (s *System) Less(s2 *System) bool {
	if s.X < s2.X {
		return true
	} else if s.X == s2.X {
		if s.Y < s2.Y {
			return true
		} else if s.Y == s2.Y {
			return s.Z < s2.Z
		}
	}
	return false
}

// TaggedAsVisited checks to see if the player manually tagged the system as visited.
func (s *System) TaggedAsVisited() bool {
	return s.Visited && s.Scanned == 0
}

// Less is a helper for sorting
func (p *Planet) Less(p2 *Planet) bool {
	if p.System.X < p2.System.X {
		return true
	} else if p.System.X == p2.System.X {
		if p.System.Y < p2.System.Y {
			return true
		} else if p.System.Y == p2.System.Y {
			if p.System.Z < p2.System.Z {
				return true
			} else if p.System.Z == p2.System.Z {
				return p.Orbit < p2.Orbit
			}
		}
	}
	return false
}
