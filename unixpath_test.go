package crosspath

import (
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	"github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestUnixPathStringer(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{
			name: "absolute",
			path: "/var/data",
		},
		{
			name: "absolute_trailingslash",
			path: "/var/data/",
		},
		{
			name: "absolute repeated slash prefix",
			path: "//var/data",
		},
		{
			name: "relative",
			path: "data",
		},
		{
			name: "relative multi segment",
			path: "data/1",
		},
		{
			name: "relative repeated slash",
			path: "data//1",
		},
		{
			name: "relative with dots and double dots",
			path: "data/./../1",
		},
		{
			name: "home routed",
			path: "~/data",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, err := NewUnixPath(c.path)
			assert.NilError(t, err)
			assert.Equal(t, c.path, p.Raw())
		})
	}
}

func TestUnixPathKind(t *testing.T) {
	cases := []struct {
		name string
		path string
		kind Kind
	}{
		{
			name: "absolute",
			path: "/var/data",
			kind: Absolute,
		},
		{
			name: "absolute_trailingslash",
			path: "/var/data/",
			kind: Absolute,
		},
		{
			name: "absolute repeated slash prefix",
			path: "//var/data",
			kind: Absolute,
		},
		{
			name: "relative",
			path: "data",
			kind: Relative,
		},
		{
			name: "relative multi segment",
			path: "data/1",
			kind: Relative,
		},
		{
			name: "relative repeated slash",
			path: "data//1",
			kind: Relative,
		},
		{
			name: "relative with dots and double dots",
			path: "data/./../1",
			kind: Relative,
		},
		{
			name: "home routed",
			path: "~/data",
			kind: HomeRooted,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, err := NewUnixPath(c.path)
			assert.NilError(t, err)
			assert.Equal(t, c.kind, p.Kind())
		})
	}
}

func TestUnixNormalize(t *testing.T) {
	cases := []struct {
		name       string
		original   string
		normalized string
	}{
		{
			name:       "absolute",
			original:   "/var/data",
			normalized: "/var/data",
		},
		{
			name:       "absolute_trailingslash",
			original:   "/var/data/",
			normalized: "/var/data",
		},
		{
			name:       "absolute repeated slash prefix",
			original:   "//var/data",
			normalized: "/var/data",
		},
		{
			name:       "relative",
			original:   "data",
			normalized: "data",
		},
		{
			name:       "relative_trailing slash",
			original:   "data/",
			normalized: "data",
		},
		{
			name:       "relative multi segment",
			original:   "data/1",
			normalized: "data/1",
		},
		{
			name:       "relative repeated slash",
			original:   "data//1",
			normalized: "data/1",
		},
		{
			name:       "relative with dots and double dots",
			original:   "data/./../1",
			normalized: "1",
		},
		{
			name:       "home routed",
			original:   "~/data",
			normalized: "~/data",
		},
		{
			name:       "back",
			original:   "..",
			normalized: "..",
		},
		{
			name:       "relative-with-manybacks",
			original:   "data/../../test",
			normalized: "../test",
		},
		{
			name:       "absolute with backs",
			original:   "/var/../data",
			normalized: "/data",
		},
		{
			name:       "absolute with many backs",
			original:   "/var/../../data",
			normalized: "/data",
		},
		{
			name:       "home with many backs",
			original:   "~/var/../../data",
			normalized: "~/../data",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, err := NewUnixPath(c.original)
			assert.NilError(t, err)
			assert.Equal(t, c.normalized, p.String())
		})
	}
}

func TestUnixPathJoin(t *testing.T) {
	path1 := "/var/data"
	path2 := "../hello"
	p1, _ := NewUnixPath(path1)
	p2, _ := NewUnixPath(path2)
	res, err := p1.Join(p2)
	assert.NilError(t, err)
	assert.Equal(t, "/var/data/../hello", res.Raw())
}

func TestUnixPathJoinHomeRooted(t *testing.T) {
	path1 := "/home/user"
	path2 := "~/data"
	p1, _ := NewUnixPath(path1)
	p2, _ := NewUnixPath(path2)
	res, err := p1.Join(p2)
	assert.NilError(t, err)
	assert.Equal(t, "/home/user/data", res.String())
}

func TestUnixPathConvert(t *testing.T) {
	cases := []struct {
		source        string
		expected      string
		expectedError string
	}{
		{
			source:   "var/data",
			expected: `var\data`,
		},
		{
			source:   "~/var/data",
			expected: `~\var\data`,
		},
		{
			source:        "/var/data",
			expectedError: "only relative and home rooted paths can be converted",
		},
	}

	for _, c := range cases {
		src, _ := NewUnixPath(c.source)
		res, err := src.Convert(Windows)
		if c.expectedError != "" {
			assert.Check(t, cmp.ErrorContains(err, c.expectedError))
		} else {
			assert.NilError(t, err)
			assert.Equal(t, c.expected, res.String())
		}
	}
}
