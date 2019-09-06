package uof

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	assert.Equal(t, "Ctrl", Producer(3).String())
	assert.Equal(t, "Ctrl", Producer(3).Name())
	assert.Equal(t, "Betradar Ctrl", Producer(3).Description())
	assert.Equal(t, InvalidName, Producer(-1).String())
}

func TestURN(t *testing.T) {
	u := URN("sr:match:123")
	assert.Equal(t, 123, u.ID())
	assert.Equal(t, URNTypeMatch, u.Type())
}

func TestLanguage(t *testing.T) {
	var l Lang
	l.Parse("hr")
	assert.Equal(t, LangHR, l)
	assert.Equal(t, "hr", l.Code())
	assert.Equal(t, "Croatian", l.Name())

	ls := Languages("hr,en,de")
	assert.Len(t, ls, 3)
	assert.Equal(t, LangHR, ls[0])
	assert.Equal(t, LangEN, ls[1])
	assert.Equal(t, LangDE, ls[2])
}
