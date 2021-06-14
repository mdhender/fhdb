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

package handlers

import (
	"context"
	"github.com/mdhender/fhdb/jwt"
	"log"
	"net/http"
)

// fhContextKey is the context key type for storing parameters in context.Context.
type fhContextKey string

type Session struct {
	Authenticated bool
	SpeciesId int
	Roles map[string]bool
}

func GetSession(r *http.Request) *Session {
	if sess, ok := r.Context().Value(fhContextKey("session")).(*Session); ok {
		return sess
	}
	return nil
}

func Authenticate(h http.HandlerFunc, f jwt.Factory) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, err := jwt.FromHeader(r)
		if err == nil {
			err = f.Validate(j)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		s := &Session{
			Authenticated: true,
			SpeciesId: j.Data().Id,
			Roles: make(map[string]bool),
		}
		for _, rr := range j.Data().Roles {
			s.Roles[rr] = true
		}

		// valid session, so inject ourselves into the context.
		ctx := context.WithValue(r.Context(), fhContextKey("session"), s)

		log.Println("mwAuthenticate passed")
		log.Println(h)

		// Propagate the incoming context.
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
