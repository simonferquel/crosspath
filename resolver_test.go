package crosspath

import (
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
)

func TestResolverWithDefaults(t *testing.T) {
	cases := []struct {
		path string
		os   TargetOS
	}{
		{
			path: "/var/data",
			os:   Unix,
		},
		{
			path: "var/data",
			os:   Unix,
		},
		{
			path: `\var\data`,
			os:   Windows,
		},
		{
			path: `var\data`,
			os:   Windows,
		},
		{
			path: `c:\var\data`,
			os:   Windows,
		},
		{
			path: `c:/var/data`,
			os:   Windows,
		},
		{
			path: `\\unc\path`,
			os:   Windows,
		},
		{
			path: `//unix/path/with/double/slash`,
			os:   Unix,
		},
		{
			path: `//?/UNC/unix/path/with/double/slash`,
			os:   Windows,
		},
		{
			path: `//./pipe/docker`,
			os:   Windows,
		},
		{
			path: `data\?with\invalidpaths\for\windows`,
			os:   Unix,
		},
	}

	for _, c := range cases {
		p, err := ParsePathWithDefaults(c.path)
		assert.NilError(t, err)
		assert.Equal(t, c.os, p.TargetOS())
	}
}
