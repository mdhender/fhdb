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

package ports

import "errors"

var ErrInternalError = errors.New("internal error")
var ErrNotFound = errors.New("not found")

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type MishapResponse struct {
	From Coords `json:"from"`
	To Coords `json:"to"`
	Age int `json:"age"`
	GV int `json:"gv"`
	MishapChance float64 `json:"mishap_chance"`
}

type SystemResponse struct {
	Id     int `json:"id"`
	Coords Coords `json:"coords"`
	Visited bool `json:"visited"`
	Link string `json:"link"`
}

type SystemsResponse struct {
	Id     int `json:"id"`
	Coords Coords `json:"coords"`
	Visited bool `json:"visited"`
	Link string `json:"link"`
}

type TurnNumberResponse struct {
	TurnNumber int `json:"turn_number"`
}

type UserResponse struct {
	Id int `json:"id"`
}

type VersionResponse struct {
	Version string `json:"version"`
}
