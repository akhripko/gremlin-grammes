.PHONY: run_gremlin start_deps stop_deps mod generate

run_gremlin:
	docker run --name gremlin_local -p 8182:8182 -d tinkerpop/gremlin-server

start_deps:
	docker start gremlin_local

stop_deps:
	docker stop gremlin_local

mod:
	go mod download
	go mod vendor

generate:
	go generate ./ent

