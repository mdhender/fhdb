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

package jwt

import (
	"encoding/json"
	"net/http"
	"strings"
)

// pull the bearer token from a request header.
func FromHeader(r *http.Request) (*JWT, error) {
	headerAuthText := r.Header.Get("Authorization")
	if headerAuthText == "" {
		return nil, ErrMissingAuthHeader
	}
	authTokens := strings.SplitN(headerAuthText, " ", 2)
	if len(authTokens) != 2 {
		return nil, ErrBadRequest
	}
	authType, authToken := authTokens[0], strings.TrimSpace(authTokens[1])
	if authType != "Bearer" {
		return nil, ErrNotBearer
	}

	sections := strings.Split(authToken, ".")
	if len(sections) != 3 || len(sections[0]) == 0 || len(sections[1]) == 0 || len(sections[2]) == 0 {
		return nil, ErrNotJWT
	}

	var j JWT
	j.h.b64 = sections[0]
	j.p.b64 = sections[1]
	j.s = sections[2]

	// decode and extract the header from the token
	if rawHeader, err := decode(j.h.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawHeader, &j.h); err != nil {
		return nil, err
	} else if j.h.Algorithm == "" || j.h.Algorithm == "none" {
		return nil, ErrUnauthorized
	}

	// decode and extract the payload from the token
	if rawPayload, err := decode(j.p.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawPayload, &j.p); err != nil {
		return nil, err
	} else if j.h.TokenType != j.p.Private.TokenType {
		return nil, ErrUnauthorized
	} else if j.h.Algorithm != j.p.Private.Algorithm {
		return nil, ErrUnauthorized
	}

	return &j, nil
}


