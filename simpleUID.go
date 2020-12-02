// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//An another one simple random UID's generator
//The UID length, chars and format can be customized
//Math.rand is using as a random generator, seed is initialized with time.Now().UnixNano() by default

/*
	usage example:

		g := simpleUID.NewGenerator(&simpleUID.Cfg{
			Alfa:      "1234567890",
			Format:    "XXX-XXXXXX-XXX",
			Validator: "[0-9]{3}-[0-9]{6}-[0-9]{3}",
		})
		uid1 := g.New()
		uid2 := g.New()

		g.Validate("111-222222-333")

*/
package simpleUID

import (
	"errors"
	"math/rand"
	"regexp"
	"time"
)

type UID interface {
	New() string
	Validate(string) (string, error)
	Validator() string
}

//simple short uid generator
type simpleRandomUID struct {
	alfa            string
	format          string
	validator       string
	validatorRgxp   *regexp.Regexp
	randomGenerator *rand.Rand
}

//Configuration descriptor for UID generator
type Cfg struct {
	Alfa      string //The chars used in the uid generation, for example "1234567890"
	Format    string //uid format, every X is replaced with a random generated char, for example "XXX-XXXXXX-XXX"
	Validator string //uid validation regexp, for example "[0-9]{3}-[0-9]{6}-[0-9]{3}"
	Seed      *int64 //Random seed generator, if null, the time.Now().UnixNano() is used
}

func NewGenerator(c *Cfg) UID {

	//default UID format
	if c == nil {
		c = &Cfg{
			Alfa:      "1234567890",
			Format:    "XXX-XXXXXX-XXX",
			Validator: "[0-9]{3}-[0-9]{6}-[0-9]{3}",
		}
	}

	if c.Seed == nil {
		seed := int64(time.Now().UnixNano())
		c.Seed = &seed
	}
	return &simpleRandomUID{
		c.Alfa,
		c.Format,
		c.Validator,
		regexp.MustCompile("(" + c.Validator + ")"),
		rand.New(rand.NewSource(*c.Seed))}
}

// generates a new uid according the format
func (s *simpleRandomUID) New() string {
	size := len(s.format)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		if s.format[i] == 'X' {
			buf[i] = s.alfa[s.randomGenerator.Intn(len(s.alfa))]
		} else {
			buf[i] = s.format[i]
		}
	}
	return string(buf)
}

//looking for the uid in the string
//return uid or error if not found
func (s *simpleRandomUID) Validate(str string) (string, error) {
	matches := s.validatorRgxp.FindStringSubmatchIndex(str)
	if len(matches) > 0 {
		return string(s.validatorRgxp.ExpandString(nil, "$1", str, matches)), nil
	}
	return "", errors.New("uid isn't found in the string")
}

func (s *simpleRandomUID) Validator() string {
	return s.validator
}
