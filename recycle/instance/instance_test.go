package instance

import (
	"testing"
	"time"

	"github.com/septemhill/test/recycle"
	"github.com/septemhill/test/recycle/group"
	"github.com/stretchr/testify/assert"
)

func TestInstanceGroup(t *testing.T) {
	assert := assert.New(t)

	grpA := group.NewGroup("GroupA", 1000, time.Second*2, recycle.GenderStrategy{})

	inst := instImpl{
		grp:  grpA,
		unit: 100,
	}

	assert.Equal("GroupA", inst.grp.Name())
}
