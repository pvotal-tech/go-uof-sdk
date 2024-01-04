package pipe

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pvotal-tech/go-uof-sdk"
	"github.com/stretchr/testify/assert"
)

type marketsAPIMock struct {
	requests map[string]struct{}
	sync.Mutex
}

func (m *marketsAPIMock) Markets(lang uof.Lang) (uof.MarketDescriptions, error) {
	return nil, nil
}

func (m *marketsAPIMock) MarketVariant(lang uof.Lang, marketID int, variant string) (uof.MarketDescriptions, error) {
	m.Lock()
	defer m.Unlock()
	m.requests[fmt.Sprintf("%s %d %s", lang, marketID, variant)] = struct{}{}
	return nil, nil
}

func TestMarketsPipe(t *testing.T) {
	a := &marketsAPIMock{requests: make(map[string]struct{})}
	ms := Markets(a, []uof.Lang{uof.LangEN})
	assert.NotNil(t, ms)

	in := make(chan *uof.Message, 1)
	out, _ := ms(in)

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
	assert.Equal(t, 2, cnt)

	_, found := a.requests["en 145 sr:point_range:76+"]
	assert.True(t, found)
}
