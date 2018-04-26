package crosspath

import (
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
)

func TestJoinAutoConverts(t *testing.T) {
	cases := []struct {
		p1     string
		p2     string
		p1OS   TargetOS
		p2OS   TargetOS
		result string
	}{
		{
			p1:     `/home/user`,
			p2:     `~\data`,
			p1OS:   Unix,
			p2OS:   Windows,
			result: `/home/user/data`,
		},
		{
			p1:     `c:\users\user`,
			p2:     `~/data`,
			p1OS:   Windows,
			p2OS:   Unix,
			result: `c:\users\user\data`,
		},
		{
			p1:     `/home/user`,
			p2:     `data\file`,
			p1OS:   Unix,
			p2OS:   Windows,
			result: `/home/user/data/file`,
		},
		{
			p1:     `c:\users\user`,
			p2:     `data/file`,
			p1OS:   Windows,
			p2OS:   Unix,
			result: `c:\users\user\data\file`,
		},
	}
	for _, c := range cases {
		p1, err := ParsePathWithDefaults(c.p1)
		assert.NilError(t, err)
		p2, err := ParsePathWithDefaults(c.p2)
		assert.NilError(t, err)
		assert.Equal(t, c.p1OS, p1.TargetOS())
		assert.Equal(t, c.p2OS, p2.TargetOS())
		res, err := p1.Join(p2)
		assert.NilError(t, err)
		assert.Equal(t, c.result, res.String())
	}
}
