package model

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestImageRatio(t *testing.T) {
	assert := assert.New(t)

	i := Image{
		Width:  5760,
		Height: 3840,
	}

	assert.Equal(1.5, i.Ratio())
}

func TestImageScale(t *testing.T) {
	assert := assert.New(t)

	i := Image{
		Width:  5760,
		Height: 3840,
	}

	assert.Equal(1024, i.Scale(1024).Dx())
	assert.Equal(682, i.Scale(1024).Dy())
}
