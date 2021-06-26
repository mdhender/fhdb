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

type Species struct {
	Id                  int
	Name                string
	AutoOrders          bool
	BankedEconomicUnits int
	FleetCost           int
	FleetPercentCost    float64
	Government          struct {
		Name string
		Type string
	}
	Homeworld struct {
		Colony       *Colony
		OriginalBase int
	}
	Relationships map[int]Relationship
	Tech          map[string]*Tech
}

type Relationship int

const (
	None Relationship = iota
	Ally
	Enemy
	Neutral
)

type Tech struct {
	Level     int
	Init      int
	Knowledge int
	BankedXp  int
}

//func (ds *Store) SpeciesMap(id string) []*Species {
//	var result []*Species
//	for _, s := range ds.Sorted.Species {
//		if id == "*" || id == s.Name || id == fmt.Sprintf("%d", s.Id) || id == fmt.Sprintf("SP%d", s.Id) {
//			result = append(result, s)
//		}
//	}
//	if result == nil {
//		return []*Species{}
//	}
//	return result
//}
