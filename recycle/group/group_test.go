package group

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroupReproduction(t *testing.T) {
	assert := assert.New(t)

	grp := NewGroup("GroupA", 10000, GroupReproductionTime(time.Second))

	assert.Equal(grp.Count(), 10000)
	time.Sleep(time.Second * 2)
	assert.Equal(grp.Count(), 10000)
	fmt.Println(grp.Count())
}
