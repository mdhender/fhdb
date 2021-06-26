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

type Ship struct {
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
}

type ShipData struct {
	Class                 string `json:"class"`
	MinManufacturingLevel int    `json:"min_ma"`
	CostFtl               int    `json:"cost_ftl"`
	CostSublight          int    `json:"cost_sublight"`
	Tonnage               int    `json:"tonnage"`
	CarryingCapacity      int    `json:"carrying_capacity"`
}
