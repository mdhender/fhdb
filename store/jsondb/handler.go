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

import (
	"fmt"
	"github.com/mdhender/fhdb/ports"
)

func (ds *Store) GetSystem(id string, spId int) (*ports.SystemResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	speciesId := fmt.Sprintf("SP%02d", spId)
	for _, v := range ds.Systems {
		if id == v.Key {
			system := ports.SystemResponse{
				Id: v.Id,
				Coords: ports.Coords{v.Coords.X, v.Coords.Y, v.Coords.Z},
			}
			for _, sp := range v.VisitedBy {
				if speciesId == sp {
					system.Visited = true
					break
				}
			}
			return &system, nil
		}
	}
	return nil, ports.ErrNotFound
}

func (ds *Store) GetSystems(spId int) ([]*ports.SystemsResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	speciesId := fmt.Sprintf("SP%02d", spId)
	var systems []*ports.SystemsResponse
	for _, v := range ds.Systems {
		system := ports.SystemsResponse{
			Id: v.Id,
			Coords: ports.Coords{v.Coords.X, v.Coords.Y, v.Coords.Z},
			Link: fmt.Sprintf("/api/system/%d %d %d", v.Coords.X, v.Coords.Y, v.Coords.Z),
		}
		for _, sp := range v.VisitedBy {
			if sp == speciesId {
				system.Visited = true
				break
			}
		}
		systems = append(systems, &system)
	}
	if systems == nil {
		return []*ports.SystemsResponse{nil}, nil
	}
	return systems, nil
}

func (ds *Store) GetTurnNumber() (*ports.TurnNumberResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	return &ports.TurnNumberResponse{
		TurnNumber: ds.Galaxy.TurnNumber,
	}, nil
}

func (ds *Store) GetUser(spId int) (*ports.UserResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	for _, sp := range ds.Species {
		if spId == sp.Id {
			return &ports.UserResponse{
				Id: sp.Id,
			}, nil
		}
	}
	return nil, ports.ErrNotFound
}

func (ds *Store) GetVersion() (*ports.VersionResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	return &ports.VersionResponse{
		Version: ds.Version,
	}, nil
}
