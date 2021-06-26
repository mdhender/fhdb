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

type Planet struct {
	Id                       int
	System                   *System
	Coords                   Coords
	Diameter                 int            `json:"diameter"`
	EconEfficiency           float64        `json:"econ_efficiency"`
	Gases                    map[string]int `json:"gases"`
	Gravity                  float64        `json:"gravity"`
	Message                  int            `json:"message"`
	MiningDifficulty         float64        `json:"mining_difficulty"`
	MiningDifficultyIncrease float64        `json:"md_increase"`
	PressureClass            int            `json:"pressure_class"`
	TemperatureClass         int            `json:"temperature_class"`
	Colonies                 []*Colony
	Ships                    []*Ship
}

// Less is a helper for sorting
func (p *Planet) Less(p2 *Planet) bool {
	return p.Coords.Less(p2.Coords)
}

type GasType int

const (
	GTNone GasType = iota
)
