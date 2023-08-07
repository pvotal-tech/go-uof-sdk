package uof

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNameWithSpecifier(t *testing.T) {
	players := map[int]Player{
		1234: Player{
			FullName: "John Rodriquez",
		},
	}
	fixture := Fixture{
		Name: "Euro2016",
		Competitors: []Competitor{
			{
				Name: "France",
			},
			{
				Name: "Germany",
			},
		},
	}

	testCases := []struct {
		description string
		name        string
		specifiers  map[string]string
		expected    string
	}{
		{
			description: "Replace {X} with value of specifier",
			name:        "Race to {pointnr} points",
			specifiers:  map[string]string{"pointnr": "3"},
			expected:    "Race to 3 points",
		},
		{
			description: "Replace {!X} with the ordinal value of specifier X",
			name:        "{!periodnr} period - total",
			specifiers:  map[string]string{"periodnr": "2"},
			expected:    "2nd period - total",
		},
		{
			description: "(-) Replace {X+/-c} with the value of the specifier X + or - the number c",
			name:        "Score difference {pointnr-3}",
			specifiers:  map[string]string{"pointnr": "10"},
			expected:    "Score difference 7",
		},
		{
			description: "(+) Replace {X+/-c} with the value of the specifier X + or - the number c",
			name:        "Score difference {pointnr+3}",
			specifiers:  map[string]string{"pointnr": "10"},
			expected:    "Score difference 13",
		},
		{
			description: "(-) Replace {!X+c} with the ordinal value of specifier X + c",
			name:        "{!inningnr-1} inning",
			specifiers:  map[string]string{"inningnr": "2"},
			expected:    "1st inning",
		},
		{
			description: "Replace {!X+c} with the ordinal value of specifier X + c",
			name:        "{!inningnr+1} inning",
			specifiers:  map[string]string{"inningnr": "2"},
			expected:    "3rd inning",
		},
		{
			description: "Replace {+X} with the value of specifier X with a +/- sign in front",
			name:        "Goal Diff {+goals} goals",
			specifiers:  map[string]string{"goals": "2"},
			expected:    "Goal Diff +2 goals",
		},
		{
			description: "Replace {-X} with the negated value of the specifier with a +/- sign in front",
			name:        "Goal Diff {-goals} goals",
			specifiers:  map[string]string{"goals": "2"},
			expected:    "Goal Diff -2 goals",
		},
		{
			description: "Name with 2 normal specifiers",
			name:        "Holes {from} to {to} - head2head (1x2) groups",
			specifiers:  map[string]string{"from": "1", "to": "18"}, // {"from":"1", "to":"18"},
			expected:    "Holes 1 to 18 - head2head (1x2) groups",
		},
		{
			description: "Name with 1 normal, 1 ordinal specifier",
			name:        "{!periodnr} period - {pointnr+3} points",
			specifiers:  map[string]string{"periodnr": "2", "pointnr": "10"},
			expected:    "2nd period - 13 points",
		},
		{
			description: "Name with 1 ordinal, 1 +/- specifier",
			name:        "{!half} half, {+goals} goals",
			specifiers:  map[string]string{"half": "1", "goals": "2"},
			expected:    "1st half, +2 goals",
		},
		{
			description: "Replace {%player} with name of specifier",
			name:        "{%player} total dismissals",
			specifiers:  map[string]string{"player": "sr:player:1234"},
			expected:    "John Rodriquez total dismissals",
		},
		{
			description: "Player with 1 normal, 1 ordinal specifier",
			name:        "{!half} half - {%player} {goals} goals",
			specifiers:  map[string]string{"half": "1", "player": "sr:player:1234", "goals": "2"},
			expected:    "1st half - John Rodriquez 2 goals",
		},
		{
			description: "Replace {$event} with the name of the event",
			name:        "Winner of {$event}",
			specifiers:  map[string]string{"event": "sr:tournament:1"},
			expected:    "Winner of Euro2016",
		},
		{
			description: "Replace {$competitorN} with the Nth competitor in the event",
			name:        "Winner is {$competitor2}",
			specifiers:  map[string]string{"competitor1": "sr:competitor:2"},
			expected:    "Winner is Germany",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := ParseSpecifier(tc.name, tc.specifiers, players, fixture)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
