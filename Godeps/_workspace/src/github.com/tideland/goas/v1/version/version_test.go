// Tideland Go Application Support - Version - Unit Tests
//
// Copyright (C) 2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package version_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"github.com/tideland/goas/v1/version"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestVersion tests the creation of a new versions and their
// accessor ethods.
func TestVersion(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	tests := []struct {
		id         string
		vsn        version.Version
		major      int
		minor      int
		patch      int
		preRelease string
		metadata   string
	}{
		{
			id:         "v1.2.3",
			vsn:        version.New(1, 2, 3),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "v1.0.3",
			vsn:        version.New(1, -2, 3),
			major:      1,
			minor:      0,
			patch:      3,
			preRelease: "",
			metadata:   "",
		}, {
			id:         "v1.2.3-alpha.2014-08-03",
			vsn:        version.New(1, 2, 3, "alpha", "2014-08-03"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha.2014-08-03",
			metadata:   "",
		}, {
			id:         "v1.2.3-alphabeta.7.11",
			vsn:        version.New(1, 2, 3, "alpha beta", "007", "1+1"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alphabeta.7.11",
			metadata:   "",
		}, {
			id:         "v1.2.3+007.a",
			vsn:        version.New(1, 2, 3, version.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "",
			metadata:   "007.a",
		}, {
			id:         "v1.2.3-alpha+007.a",
			vsn:        version.New(1, 2, 3, "alpha", version.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "alpha",
			metadata:   "007.a",
		}, {
			id:         "v1.2.3-ALPHA+007.a",
			vsn:        version.New(1, 2, 3, "ALPHA", version.Metadata, "007", "a"),
			major:      1,
			minor:      2,
			patch:      3,
			preRelease: "ALPHA",
			metadata:   "007.a",
		},
	}

	for i, test := range tests {
		assert.Logf("test #%d: %q", i, test.id)
		assert.Equal(test.vsn.Major(), test.major)
		assert.Equal(test.vsn.Minor(), test.minor)
		assert.Equal(test.vsn.Patch(), test.patch)
		assert.Equal(test.vsn.PreRelease(), test.preRelease)
		assert.Equal(test.vsn.Metadata(), test.metadata)
		assert.Equal(test.vsn.String(), test.id)
	}
}

// TestLess tests the comparing of two versions.
func TestLess(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	tests := []struct {
		vsnA version.Version
		vsnB version.Version
		less bool
	}{
		{
			vsnA: version.New(1, 2, 3),
			vsnB: version.New(1, 2, 3),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3),
			vsnB: version.New(1, 2, 4),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3),
			vsnB: version.New(1, 3, 3),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3),
			vsnB: version.New(2, 2, 3),
			less: true,
		}, {
			vsnA: version.New(3, 2, 1),
			vsnB: version.New(1, 2, 3),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3, "alpha"),
			vsnB: version.New(1, 2, 3),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "1"),
			vsnB: version.New(1, 2, 3, "alpha"),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "1"),
			vsnB: version.New(1, 2, 3, "alpha", "2"),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "4711"),
			vsnB: version.New(1, 2, 3, "alpha", "471"),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "48"),
			vsnB: version.New(1, 2, 3, "alpha", "4711"),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3, version.Metadata, "alpha", "1"),
			vsnB: version.New(1, 2, 3, version.Metadata, "alpha", "2"),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3, version.Metadata, "alpha", "2"),
			vsnB: version.New(1, 2, 3, version.Metadata, "alpha", "1"),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", version.Metadata, "alpha", "2"),
			vsnB: version.New(1, 2, 3, "alpha", version.Metadata, "alpha", "1"),
			less: false,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "48", version.Metadata, "alpha", "2"),
			vsnB: version.New(1, 2, 3, "alpha", "4711", version.Metadata, "alpha", "1"),
			less: true,
		}, {
			vsnA: version.New(1, 2, 3, "alpha", "2"),
			vsnB: version.New(1, 2, 3, "alpha", "1b"),
			less: false,
		},
	}

	for i, test := range tests {
		assert.Logf("test #%d: %q < %q", i, test.vsnA, test.vsnB)
		assert.Equal(test.vsnA.Less(test.vsnB), test.less)
	}
}

// EOF
