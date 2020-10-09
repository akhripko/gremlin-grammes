package main

import (
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/akhripko/gremlin-grammes/src/options"
	"github.com/northwesternmutual/grammes"
	p "github.com/northwesternmutual/grammes/query/predicate"
	t "github.com/northwesternmutual/grammes/query/traversal"
)

func main() {
	config := options.ReadEnv()
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix(">>")
	log.Println(config.GremlinAddr)

	// Load CA cert
	//caCert, err := ioutil.ReadFile("SFSRootCAG2.pem")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)
	//
	//tlsConfig := &tls.Config{
	//	Certificates: []tls.Certificate{},
	//	RootCAs:      caCertPool,
	//}
	//log.Println(tlsConfig)

	// Creates a new client with the localhost IP.
	client, err := grammes.DialWithWebSocket("wss://127.0.0.1:8182",
		grammes.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatalf("Error while creating client: %s\n", err.Error())
	}

	// Create a graph traversal.
	g := grammes.Traversal()

	//g.V().has('zip', 'name', '78704').
	//	in('lives').
	//	hasLabel('provider').as('p').
	//	outE('provides').
	//	and(has('max_rate', lte(50)), has('min_rate', gte(0))).
	//	inV().
	//	hasLabel('service').
	//	select('p').
	//	properties().hasKey('sitter_id').dedup().value()

	query := g.V().Has("zip", "name", "78704").
		In("lives").
		HasLabel("provider").As("p").
		OutE("provides").
		And(t.NewTraversal().Has("max_rate", p.LessThanOrEqual(50)).Raw(),
			t.NewTraversal().Has("min_rate", p.GreaterThanOrEqual(0)).Raw()).
		InV().
		HasLabel("service").
		Select("p").
		Properties().HasKey("sitter_id").Dedup().Value()

	log.Printf("\nquery:\n===========\n%s\n===========\n", query.String())

	res, err := client.ExecuteQuery(query)
	if err != nil {
		log.Fatalf("Querying error: %s\n", err.Error())
	}

	data, err := UnmarshalInt32List(res)
	if err != nil {
		log.Fatalf("Querying error: %s\n", err.Error())
	}
	log.Printf("\ncount:%d\n%v\n", len(data), data)
}

func UnmarshalInt32List(recs [][]byte) ([]int32, error) {
	var (
		data  GremlinInt32List
		items []GremlinInt32Value
		err   error
	)
	if len(recs) == 0 {
		return nil, nil
	}
	for _, r := range recs {
		err = json.Unmarshal(r, &data)
		if err != nil {
			return nil, err
		}
		items = append(items, data.Items...)
	}
	res := make([]int32, 0, len(items))
	for _, v := range items {
		res = append(res, v.Value)
	}
	return res, nil
}

type GremlinInt32List struct {
	Items []GremlinInt32Value `json:"@value"`
}

type GremlinInt32Value struct {
	Value int32 `json:"@value"`
}
