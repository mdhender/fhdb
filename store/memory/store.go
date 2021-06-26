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

import (
	"fmt"
	"github.com/mdhender/fhdb/ports"
)

type Store struct {
	Version    string
	TurnNumber int

	Systems map[string]*System // indexed by systemId
	Planets []*Planet          // indexed by planetId
	Species []*Species         // indexed by spId

	Colonies map[string]*Colony // colonies are "named planets"
	Ships    map[string]*Ship   // key for ship is spId / shipId
}

func (ds *Store) GetKnownSpecies(id int, roles map[string]bool) ([]*ports.KnownSpeciesResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	} else if len(roles) == 0 {
		return nil, ports.ErrUnauthorized
	} else if id < 1 || !(id < len(ds.Species)) {
		return nil, ports.ErrUnauthorized
	}
	var results []*ports.KnownSpeciesResponse
	for _, v := range ds.Species {
		if id == v.Id { // don't report on self
			continue
		} else if _, ok := roles[fmt.Sprintf("SP%02d", v.Id)]; !ok {
			continue
		}
		sp := ports.KnownSpeciesResponse{
			Id: v.Id,
		}
		results = append(results, &sp)
	}
	if results == nil {
		return []*ports.KnownSpeciesResponse{}, nil
	}
	return results, nil
}

func (ds *Store) GetSpecies(id int, roles map[string]bool) (*ports.SpeciesResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	} else if len(roles) == 0 {
		return nil, ports.ErrUnauthorized
	} else if _, ok := roles[fmt.Sprintf("SP%02d", id)]; !ok {
		return nil, ports.ErrUnauthorized
	} else if id < 1 || !(id < len(ds.Species)) {
		return nil, ports.ErrUnauthorized
	}
	sp := ds.Species[id]
	if sp == nil {
		return nil, ports.ErrNotFound
	}
	rsp := ports.SpeciesResponse{
		Id: sp.Id,
	}
	return &rsp, nil
}

func (ds *Store) GetSystem(id string, spId int) (*ports.SystemResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	system, ok := ds.Systems[id]
	if !ok {
		return nil, ports.ErrNotFound
	}
	rsp := ports.SystemResponse{
		Id:     system.Id,
		Coords: ports.Coords{X: system.Coords.X, Y: system.Coords.Y, Z: system.Coords.Z},
	}
	return &rsp, nil
}

func (ds *Store) GetSystems(spId int) ([]*ports.SystemsResponse, error) {
	if ds == nil {
		return nil, ports.ErrInternalError
	}
	//speciesId := fmt.Sprintf("SP%02d", spId)
	var systems []*ports.SystemsResponse
	for _, v := range ds.Systems {
		system := ports.SystemsResponse{
			Id:     v.Id,
			Coords: ports.Coords{v.Coords.X, v.Coords.Y, v.Coords.Z},
			Link:   fmt.Sprintf("/api/system/%d %d %d", v.Coords.X, v.Coords.Y, v.Coords.Z),
		}
		//for _, sp := range v.VisitedBy {
		//	if sp == speciesId {
		//		system.Visited = true
		//		break
		//	}
		//}
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
		TurnNumber: ds.TurnNumber,
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
