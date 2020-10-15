package enrollment

import (
	"testing"

	"github.com/northwesternmutual/grammes"
	p "github.com/northwesternmutual/grammes/query/predicate"
	t "github.com/northwesternmutual/grammes/query/traversal"
	"github.com/stretchr/testify/assert"
)

func Test_BuildQuery_fromZIP(te *testing.T) {
	g := grammes.Traversal()
	expected := g.V().Has("zip", "name", "78704").
		In("lives").
		HasLabel("provider").As("p").
		OutE("provides").
		And(t.NewTraversal().Has("max_rate", p.LessThanOrEqual(50)).Raw(),
			t.NewTraversal().Has("min_rate", p.GreaterThanOrEqual(0)).Raw(),
			t.NewTraversal().Has("service", "childCare").Raw()).
		InV().
		Select("p").
		//Order().By("sitter_id").
		Range(0, 10).
		Properties().HasKey("sitter_id").Value()

	query, err := BuildQuery(&GRPCModel{
		PostalCode: "78704",
		CareType:   "childCare",
		HourlyRate: &HourlyRateGRPCModel{Min: 0, Max: 50},
		PageSize:   10,
	})
	assert.NoError(te, err)

	assert.Equal(te, expected.String(), query.String())
}

func Test_BuildQuery_fromService(te *testing.T) {
	g := grammes.Traversal()
	expected := g.V().
		Has("service", "childCare").
		InE("provides").
		And(t.NewTraversal().Has("max_rate", p.LessThanOrEqual(100)).Raw(),
			t.NewTraversal().Has("min_rate", p.GreaterThanOrEqual(0)).Raw()).
		OutV().
		HasLabel("provider").
		//Order().By("sitter_id").
		Range(20, 30).
		Properties().HasKey("sitter_id").Value()

	query, err := BuildQuery(&GRPCModel{
		CareType:   "childCare",
		HourlyRate: &HourlyRateGRPCModel{Min: 0, Max: 100},
		PageSize:   10,
		PageToken:  "2",
	})
	assert.NoError(te, err)

	assert.Equal(te, expected.String(), query.String())
}
