package pipe

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/pvotal-tech/go-uof-sdk"
	"github.com/stretchr/testify/assert"
)

type playerAPIMock struct {
	requests map[int]struct{}
	sync.Mutex
}

func (m *playerAPIMock) Player(lang uof.Lang, playerID int) (*uof.Player, error) {
	m.Lock()
	defer m.Unlock()
	m.requests[uof.UIDWithLang(playerID, lang)] = struct{}{}
	return nil, nil
}

func TestPlayerPipe(t *testing.T) {
	a := &playerAPIMock{requests: make(map[int]struct{})}
	p := Player(a, []uof.Lang{uof.LangEN}, false)
	assert.NotNil(t, p)

	in := make(chan *uof.Message)
	out, _ := p(in)

	m := oddsChangeMessage(t)
	in <- m

	close(in)
	cnt := 0
	for om := range out {
		if om.Type == uof.MessageTypeOddsChange {
			assert.Equal(t, m, om)
		} else {
			cnt++
		}
	}
	assert.Equal(t, 41, cnt)
	assert.Equal(t, 41, len(a.requests))
}

func oddsChangeMessage(t *testing.T) *uof.Message {
	buf, err := ioutil.ReadFile("../testdata/odds_change-0.xml")
	assert.NoError(t, err)
	m, err := uof.NewQueueMessage("hi.pre.-.odds_change.1.sr:match.1234.-", buf)
	assert.NoError(t, err)
	return m
}
