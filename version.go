// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// The version package implements version parsing.
package version

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// Number represents a version Number.
type Number struct {
	Major int
	Minor int
	Tag   string
	Patch int
	Build int
}

// Zero is occasionally convenient and readable.
// Please don't change its value.
var Zero = Number{}

// Binary specifies a binary version of juju.
type Binary struct {
	Number
	Series string
	Arch   string
}

func (v Binary) String() string {
	return fmt.Sprintf("%v-%s-%s", v.Number, v.Series, v.Arch)
}

// GetBSON turns v into a bson.Getter so it can be saved directly
// on a MongoDB database with mgo.
func (v Binary) GetBSON() (interface{}, error) {
	return v.String(), nil
}

// SetBSON turns v into a bson.Setter so it can be loaded directly
// from a MongoDB database with mgo.
func (vp *Binary) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return err
	}
	v, err := ParseBinary(s)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

func (v Binary) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (vp *Binary) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := ParseBinary(s)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

// MarshalYAML implements yaml.v2.Marshaller interface
func (v Binary) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface
func (vp *Binary) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var vstr string
	err := unmarshal(&vstr)
	if err != nil {
		return err
	}
	v, err := ParseBinary(vstr)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

var (
	binaryPat = regexp.MustCompile(`^(\d{1,9})\.(\d{1,9})(\.|-(\w+))(\d{1,9})(\.\d{1,9})?-([^-]+)-([^-]+)$`)
	numberPat = regexp.MustCompile(`^(\d{1,9})\.(\d{1,9})(\.|-(\w+))(\d{1,9})(\.\d{1,9})?$`)
)

// MustParse parses a version and panics if it does
// not parse correctly.
func MustParse(s string) Number {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// MustParseBinary parses a binary version and panics if it does
// not parse correctly.
func MustParseBinary(s string) Binary {
	v, err := ParseBinary(s)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseBinary parses a binary version of the form "1.2.3-series-arch".
func ParseBinary(s string) (Binary, error) {
	m := binaryPat.FindStringSubmatch(s)
	if m == nil {
		return Binary{}, fmt.Errorf("invalid binary version %q", s)
	}
	var v Binary
	v.Major = atoi(m[1])
	v.Minor = atoi(m[2])
	v.Tag = m[4]
	v.Patch = atoi(m[5])
	if m[6] != "" {
		v.Build = atoi(m[6][1:])
	}
	v.Series = m[7]
	v.Arch = m[8]
	return v, nil
}

// Parse parses the version, which is of the form 1.2.3
// giving the major, minor and release versions
// respectively.
func Parse(s string) (Number, error) {
	m := numberPat.FindStringSubmatch(s)
	if m == nil {
		return Number{}, fmt.Errorf("invalid version %q", s)
	}
	var v Number
	v.Major = atoi(m[1])
	v.Minor = atoi(m[2])
	v.Tag = m[4]
	v.Patch = atoi(m[5])
	if m[6] != "" {
		v.Build = atoi(m[6][1:])
	}
	return v, nil
}

// atoi is the same as strconv.Atoi but assumes that
// the string has been verified to be a valid integer.
func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func (v Number) String() string {
	var s string
	if v.Tag == "" {
		s = fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else {
		s = fmt.Sprintf("%d.%d-%s%d", v.Major, v.Minor, v.Tag, v.Patch)
	}
	if v.Build > 0 {
		s += fmt.Sprintf(".%d", v.Build)
	}
	return s
}

// Compare returns -1, 0 or 1 depending on whether
// v is less than, equal to or greater than w.
func (v Number) Compare(w Number) int {
	if v == w {
		return 0
	}
	less := false
	switch {
	case v.Major != w.Major:
		less = v.Major < w.Major
	case v.Minor != w.Minor:
		less = v.Minor < w.Minor
	case v.Tag != w.Tag:
		switch {
		case v.Tag == "":
			less = false
		case w.Tag == "":
			less = true
		default:
			less = v.Tag < w.Tag
		}
	case v.Patch != w.Patch:
		less = v.Patch < w.Patch
	case v.Build != w.Build:
		less = v.Build < w.Build
	}
	if less {
		return -1
	}
	return 1
}

// GetBSON turns v into a bson.Getter so it can be saved directly
// on a MongoDB database with mgo.
func (v Number) GetBSON() (interface{}, error) {
	return v.String(), nil
}

// SetBSON turns v into a bson.Setter so it can be loaded directly
// from a MongoDB database with mgo.
func (vp *Number) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return err
	}
	v, err := Parse(s)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

func (v Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (vp *Number) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := Parse(s)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

// MarshalYAML implements yaml.v2.Marshaller interface
func (v Number) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface
func (vp *Number) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var vstr string
	err := unmarshal(&vstr)
	if err != nil {
		return err
	}
	v, err := Parse(vstr)
	if err != nil {
		return err
	}
	*vp = v
	return nil
}

// ParseMajorMinor takes an argument of the form "major.minor" and returns ints major and minor.
func ParseMajorMinor(vers string) (int, int, error) {
	parts := strings.Split(vers, ".")
	major, err := strconv.Atoi(parts[0])
	minor := -1
	if err != nil {
		return -1, -1, fmt.Errorf("invalid major version number %s: %v", parts[0], err)
	}
	if len(parts) == 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return -1, -1, fmt.Errorf("invalid minor version number %s: %v", parts[1], err)
		}
	} else if len(parts) > 2 {
		return -1, -1, fmt.Errorf("invalid major.minor version number %s", vers)
	}
	return major, minor, nil
}
