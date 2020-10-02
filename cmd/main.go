package main

import (
	"log"

	"github.com/akhripko/gremlin-grammes/src/options"
	"github.com/northwesternmutual/grammes"
	//"github.com/northwesternmutual/grammes/query/traversal"
)

func main() {
	config := options.ReadEnv()
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
	client, err := grammes.DialWithWebSocket("wss://127.0.0.1:8182")
	if err != nil {
		log.Fatalf("Error while creating client: %s\n", err.Error())
	}

	// Create a graph traversal.
	g := grammes.Traversal()

	res, err := client.ExecuteQuery(g.V("state_TX").PropertyMap())
	if err != nil {
		log.Fatalf("Querying error: %s\n", err.Error())
	}

	// Log out the response.
	for _, r := range res {
		log.Println(string(r))
	}

	//
	// g.V().hasLabel('state').has('name', 'TX').inE('located').outV().hasLabel('city').inE('located').outV().inE('lives').outV().hasLabel('provider').outE('provides').inV().hasLabel('service').and(has('max_rate', lte(50)), has('min_rate', gte(0))).path().limit(10)
	//
	//g.V().hasLabel('state')
	//.has('name', 'TX')
	//.inE('located')
	//.outV().hasLabel('city')
	//.inE('located')
	//.outV()
	//.inE('lives')
	//.outV().hasLabel('provider')
	//.outE('provides')
	//.inV().hasLabel('service')
	//.and(has('max_rate', lte(50)), has('min_rate', gte(0)))
	//.path()
	//.limit(200)

	query := g.V().HasLabel("state").Has("name", "TX").InE("located").OutV().
		HasLabel("city").
		InE("located").
		OutV().
		InE("lives").
		OutV().
		HasLabel("provider").
		OutE("provides").
		InV().
		HasLabel("service").
		//And(traversal.String("has('max_rate', lte(50))"), "has('min_rate', gte(0))").
		Path().
		Limit(100)

	res, err = client.ExecuteQuery(query)
	if err != nil {
		log.Fatalf("Querying error: %s\n", err.Error())
	}

	// Log out the response.
	for _, r := range res {
		log.Printf("\n>>>>PATH:\n%s\n==============\n", string(r))
	}
}
