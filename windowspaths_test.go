package crosspath

import (
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	"github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestWindowsPathParsing(t *testing.T) {
	cases := []struct {
		convertSlashes bool
		source         string
		expectedError  string
		expectedKind   Kind
		expectedOut    string
	}{
		{
			convertSlashes: false,
			source:         `c:\data\file`,
			expectedKind:   Absolute,
			expectedOut:    `c:\data\file`,
		},
		{
			convertSlashes: false,
			source:         `data\file`,
			expectedKind:   Relative,
			expectedOut:    `data\file`,
		},
		{
			convertSlashes: false,
			source:         `..\file`,
			expectedKind:   Relative,
			expectedOut:    `..\file`,
		},
		{
			convertSlashes: false,
			source:         `~\data\file`,
			expectedKind:   HomeRooted,
			expectedOut:    `~\data\file`,
		},
		{
			convertSlashes: false,
			source:         `\data\file`,
			expectedKind:   AbsoluteFromCurrentDrive,
			expectedOut:    `\data\file`,
		},
		{
			convertSlashes: false,
			source:         `\`,
			expectedKind:   AbsoluteFromCurrentDrive,
			expectedOut:    `\`,
		},
		{
			convertSlashes: false,
			source:         `c:data\file`,
			expectedKind:   RelativeFromDriveCurrentDir,
			expectedOut:    `c:data\file`,
		},
		{
			convertSlashes: false,
			source:         `\\some\unc\path`,
			expectedKind:   UNC,
			expectedOut:    `\\some\unc\path`,
		},
		{
			convertSlashes: false,
			source:         `\\?\c:\data\path`,
			expectedKind:   Absolute,
			expectedOut:    `\\?\c:\data\path`,
		},
		{
			convertSlashes: false,
			source:         `\\?\UNC\unc\path`,
			expectedKind:   UNC,
			expectedOut:    `\\?\UNC\unc\path`,
		},
		{
			convertSlashes: false,
			source:         `\\.\pipe\docker`,
			expectedKind:   WindowsDevice,
			expectedOut:    `\\.\pipe\docker`,
		},
	}
	for _, c := range cases {
		p, err := NewWindowsPath(c.source, c.convertSlashes)
		if c.expectedError != "" {
			assert.Check(t, cmp.ErrorContains(err, c.expectedError))
		} else {
			assert.Equal(t, c.expectedOut, p.String())
			assert.Equal(t, c.expectedKind, p.Kind())
		}
	}
}

func TestWindowsPathNormalize(t *testing.T) {
	cases := []struct {
		source   string
		expected string
	}{
		{
			source:   `c:\data\..\var`,
			expected: `c:\var`,
		},
		{
			source:   `c:\data\..\.\var\.`,
			expected: `c:\var`,
		},
		{
			source:   `c:\data\..\..`,
			expected: `c:\`,
		},
		{
			source:   `\data\..\var`,
			expected: `\var`,
		},
		{
			source:   `\data\..\..`,
			expected: `\`,
		},
		{
			source:   `data\..\var`,
			expected: `var`,
		},
		{
			source:   `data\..\..`,
			expected: `..`,
		},
		{
			source:   `~\data\..\var`,
			expected: `~\var`,
		},
		{
			source:   `~\data\..\..`,
			expected: `~\..`,
		},
		{
			source:   `\\?\c:\data\..\..`,
			expected: `\\?\c:\data\..\..`,
		},
		{
			source:   `\\.\pipe\docker\..\docker`,
			expected: `\\.\pipe\docker`,
		},
		{
			source:   `\\server\data\..\var`,
			expected: `\\server\var`,
		},
		{
			source:   `\\server\data\..\..`,
			expected: `\\server`,
		},
	}
	for _, c := range cases {
		p, err := NewWindowsPath(c.source, false)
		assert.NilError(t, err)
		assert.Equal(t, c.expected, p.String())
	}
}

func TestWindowsPathJoin(t *testing.T) {
	path1 := `c:\var\data`
	path2 := `..\hello`
	p1, _ := NewWindowsPath(path1, false)
	p2, _ := NewWindowsPath(path2, false)
	res, err := p1.Join(p2)
	assert.NilError(t, err)
	assert.Equal(t, `c:\var\data\..\hello`, res.Raw())
}

func TestWindowsPathJoinHomeRooted(t *testing.T) {
	path1 := `c:\users\user`
	path2 := `~\data`
	p1, _ := NewWindowsPath(path1, false)
	p2, _ := NewWindowsPath(path2, false)
	res, err := p1.Join(p2)
	assert.NilError(t, err)
	assert.Equal(t, `c:\users\user\data`, res.String())
}

func TestWindowsPathConvert(t *testing.T) {
	cases := []struct {
		source        string
		expected      string
		expectedError string
	}{
		{
			expected: "var/data",
			source:   `var\data`,
		},
		{
			expected: "~/var/data",
			source:   `~\var\data`,
		},
		{
			source:        `\var\data`,
			expectedError: "only relative and home rooted paths can be converted",
		},
	}

	for _, c := range cases {
		src, _ := NewWindowsPath(c.source, false)
		res, err := src.Convert(Unix)
		if c.expectedError != "" {
			assert.Check(t, cmp.ErrorContains(err, c.expectedError))
		} else {
			assert.NilError(t, err)
			assert.Equal(t, c.expected, res.String())
		}
	}
}
