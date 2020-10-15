package enrollment

import (
	"errors"
	"strconv"

	"github.com/northwesternmutual/grammes"
	p "github.com/northwesternmutual/grammes/query/predicate"
	t "github.com/northwesternmutual/grammes/query/traversal"
)

const DefaultPageSize int32 = 20

type HourlyRateGRPCModel struct {
	Min float32
	Max float32
}

type GRPCModel struct {
	PostalCode string
	CareType   string
	Gender     string
	HourlyRate *HourlyRateGRPCModel
	PageSize   int32
	PageToken  string
}

func BuildQuery(req *GRPCModel) (t.String, error) {
	var err error
	g := grammes.Traversal()
	if req == nil {
		return g, errors.New("empty request")
	}
	// build core request
	switch {
	case isZIPModel(req):
		g, err = providersFromZIP(g, req)
	default:
		g, err = providersFromService(g, req)
	}
	if err != nil {
		return g, err
	}
	// order result
	//g = order(g, req)
	// pagination
	g, err = pagination(g, req)
	if err != nil {
		return g, err
	}
	// in result
	return forResult(g, req), nil
}

func isZIPModel(req *GRPCModel) bool {
	return len(req.PostalCode) > 0
}

func providersFromZIP(g t.String, req *GRPCModel) (t.String, error) {
	query := g.V().Has("zip", "name", req.PostalCode).In("lives")
	query = addProviderFilter(query, req).As("p")
	// add limits for: provider -provides(and(limits...))-> service
	query = query.OutE("provides")
	limits := make([]t.String, 0, 3)
	limits = appendRateLimits(limits, req)
	limits = appendServiceLimit(limits, req)
	if len(limits) > 0 {
		query = query.And(limits...)
	}
	// select filtered providers
	return query.InV().Select("p"), nil
}

func providersFromService(g t.String, req *GRPCModel) (t.String, error) {
	query := g.V()
	// select services
	query = addServiceFilter(query, req)
	// add limits for: service <-provides(and(limits...))- provider
	query = query.InE("provides")
	limits := make([]t.String, 0, 1)
	limits = appendRateLimits(limits, req)
	if len(limits) > 0 {
		query = query.And(limits...)
	}
	return query.OutV().HasLabel("provider"), nil
}

func addProviderFilter(g t.String, req *GRPCModel) t.String {
	if len(req.Gender) > 0 {
		return g.Has("provider", "gender", req.Gender)
	}
	return g.HasLabel("provider")
}

func addServiceFilter(g t.String, req *GRPCModel) t.String {
	if len(req.CareType) > 0 {
		return g.Has("service", req.CareType)
	}
	return g.HasLabel("service")
}

func appendServiceLimit(limits []t.String, req *GRPCModel) []t.String {
	if len(req.CareType) > 0 {
		limits = append(limits, getRawHas("service", req.CareType))
	}
	return limits
}

func appendRateLimits(limits []t.String, req *GRPCModel) []t.String {
	if req.HourlyRate == nil {
		return limits
	}
	if req.HourlyRate.Max > 0 {
		limits = append(limits, getRawHas("max_rate", p.LessThanOrEqual(req.HourlyRate.Max)))
	}
	if req.HourlyRate.Max > 0 {
		limits = append(limits, getRawHas("min_rate", p.GreaterThanOrEqual(req.HourlyRate.Min)))
	}
	return limits
}

func getRawHas(first interface{}, params ...interface{}) t.String {
	return t.NewTraversal().Has(first, params...).Raw()
}

func order(q t.String, _ *GRPCModel) t.String {
	//return q.Order().By("id", t.NewCustomTraversal("asc"))
	return q.Order().By(t.NewCustomTraversal("id"))
}

func forResult(g t.String, _ *GRPCModel) t.String {
	return g.Properties().HasKey("sitter_id").Value()
}

func pagination(q t.String, req *GRPCModel) (t.String, error) {
	pageToken, err := getPageToken(req)
	if err != nil {
		return q, err
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	from := pageToken * pageSize
	to := from + pageSize
	return q.Range(from, to), nil
}

func getPageToken(req *GRPCModel) (int32, error) {
	if len(req.PageToken) == 0 {
		return 0, nil
	}
	token, err := strconv.ParseInt(req.PageToken, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(token), nil
}
