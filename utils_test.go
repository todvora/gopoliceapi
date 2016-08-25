package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslateKey(t *testing.T) {
	assert.Equal(t, "manufacturer", TranslateKey("vyrobce"))
	assert.Equal(t, "engine", TranslateKey("motor"))
	assert.Equal(t, "unknownXYZ", TranslateKey("unknownXYZ"))
}

func TestStandardizeDate(t *testing.T) {
	assert.Equal(t, "1.4.2016", StandardizeDate("1. dubna 2016"))
}
