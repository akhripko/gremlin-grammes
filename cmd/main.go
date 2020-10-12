package main

import (
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/akhripko/gremlin-grammes/src/enrollment"

	"github.com/akhripko/gremlin-grammes/src/options"
	"github.com/northwesternmutual/grammes"
)

func main() {
	config := options.ReadEnv()
	//log.SetFlags(log.Lmsgprefix)
	//log.SetPrefix(">>")
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
	//g := grammes.Traversal()

	//query := g.V().Has("zip", "name", "78704").
	//	In("lives").
	//	HasLabel("provider").As("p").
	//	OutE("provides").
	//	And(t.NewTraversal().Has("max_rate", p.LessThanOrEqual(50)).Raw(),
	//		t.NewTraversal().Has("min_rate", p.GreaterThanOrEqual(0)).Raw(),
	//		t.NewTraversal().Has("service", "childCare").Raw()).
	//	InV().
	//	Select("p").Range(0, 10).
	//	Properties().HasKey("sitter_id").Value()

	//query := g.V().
	//	Has("service", "childCare").
	//	InE("provides").
	//	And(t.NewTraversal().Has("max_rate", p.LessThanOrEqual(100)).Raw(),
	//		t.NewTraversal().Has("min_rate", p.GreaterThanOrEqual(0)).Raw()).
	//	OutV().
	//	HasLabel("provider").
	//	Range(0, 10).
	//	Properties().HasKey("sitter_id").Value()

	query, err := enrollment.BuildQuery(&enrollment.GRPCModel{
		PostalCode: "78704",
		CareType:   "childCare",
		HourlyRate: &enrollment.HourlyRateGRPCModel{Min: 0, Max: 50},
		PageSize:   10,
	})

	//query, err := enrollment.BuildQuery(&enrollment.GRPCModel{
	//	CareType:   "childCare",
	//	HourlyRate: &enrollment.HourlyRateGRPCModel{Min: 0, Max: 50},
	//	PageSize:   10,
	//	PageToken: "2",
	//})

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
	log.Printf("\ncount:%d\n", len(data))
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
