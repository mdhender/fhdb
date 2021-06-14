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

// Package jwt implements naive JSON Web Tokens.
// Don't use this for anything other than testing.
package jwt

import (
	"time"
)

func (j *JWT) Data() Data {
	return Data{
		Id:       j.p.Private.Id,
		Username: j.p.Private.Username,
		Email:    j.p.Private.Email,
		Roles:    j.p.Private.Roles,
	}
}

func (j *JWT) IsValid() bool {
	now := time.Now().UTC()
	if j == nil {
		//log.Printf("jwt is nil\n")
		return false
	} else if !j.isSigned || j.h.Algorithm != j.p.Private.Algorithm || j.h.TokenType != j.p.Private.TokenType {
		//log.Printf("alg %q typ %q signed %v borked\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if j.p.NotBefore != 0 && !now.Before(time.Unix(j.p.NotBefore, 0)) {
		//log.Printf("alg %q typ %q signed %v !now.Before(notBefore)\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if j.p.IssuedAt == 0 {
		//log.Printf("alg %q typ %q signed %v no issue timestamp\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if !now.After(time.Unix(j.p.IssuedAt, 0)) {
		//log.Printf("alg %q typ %q signed %v !now.After(issuedAt) %s %s\n", j.h.Algorithm, j.h.TokenType, j.isSigned, now.Format("2006-01-02T15:04:05.99999999Z"), time.Unix(j.p.IssuedAt, 0).Format("2006-01-02T15:04:05.99999999Z"))
		return false
	} else if j.p.ExpirationTime == 0 {
		//log.Printf("alg %q typ %q signed %v no expiration timestamp\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if !time.Unix(j.p.ExpirationTime, 0).After(now) {
		//log.Printf("alg %q typ %q signed %v !expiresAt.After(now)\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	}
	return true
}

type JWT struct {
	h struct {
		Algorithm   string `json:"alg,omitempty"` // message authentication code algorithm
		TokenType   string `json:"typ,omitempty"`
		ContentType string `json:"cty,omitempty"`
		KeyID       string `json:"kid,omitempty"` // optional identifier used to sign. doesn't work.
		b64         string // header marshalled to JSON and then base-64 encoded
	}
	p struct {
		// The principal that issued the JWT.
		Issuer string `json:"iss,omitempty"`
		// The subject of the JWT.
		Subject string `json:"sub,omitempty"`
		// The recipients that the JWT is intended for.
		// Each principal intended to process the JWT must identify itself with a value in the audience claim.
		// If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present,
		// then the JWT must be rejected.
		Audience []string `json:"aud,omitempty"`
		// The expiration time on and after which the JWT must not be accepted for processing.
		// The value must be a NumericDate:[9] either an integer or decimal, representing seconds past 1970-01-01 00:00:00Z.
		ExpirationTime int64 `json:"exp,omitempty"`
		// The time on which the JWT will start to be accepted for processing.
		// The value must be a NumericDate.
		NotBefore int64 `json:"nbf,omitempty"`
		// The time at which the JWT was issued.
		// The value must be a NumericDate.
		IssuedAt int64 `json:"iat,omitempty"`
		// Case sensitive unique identifier of the token even among different issuers.
		JWTID string `json:"jti,omitempty"`
		// Private data for use by the application.
		Private struct {
			Algorithm string   `json:"alg"`
			TokenType string   `json:"typ"`
			Id        int      `json:"id,omitempty"`
			Username  string   `json:"username,omitempty"`
			Email     string   `json:"email,omitempty"`
			Roles     []string `json:"roles,omitempty"`
		} `json:"private"`
		b64 string // payload marshalled to JSON and then base-64 encoded
	}
	s        string // signature base-64 encoded
	isSigned bool   // true only if the signature has been verified
}

type Data struct {
	Id       int
	Username string
	Email    string
	Roles    []string
}

