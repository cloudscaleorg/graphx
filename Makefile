# Local Development #

.PHONY: etcd-up
etcd-up:
	docker-compose up -d graphx-etcd
	docker exec graphx-etcd /bin/sh -c "/usr/local/bin/etcdctl version"

.PHONY: etcd-down
etcd-down:
	docker-compose down graphx-etcd

.PHONY: swagger-up
swagger-up:
	docker-compose up -d graphx-swagger-ui

.PHONY: swagger-down
swagger-down:
	docker-compose down -d graphx-swagger-ui

.PHONY: graphx-up
graphx-up:
	docker-compose up -d graphx-node

.PHONY: graphx-down
graphx-down:
	docker-compose down -d graphx-node

.PHONY: local-dev-up
local-dev-up:
	make etcd-up
	make swagger-up
	make graphx-up

.PHONY: local-dev-down
local-dev-down:
	docker-compose down

# Testing #

.PHONY: unit-verbose
unit-verbose:
	go test -v -count=1 -race ./...

.PHONY: automated-integration-etcd
automated-integration-etcd:
	make etcd-down
	make etcd-up
	go test -count=1 -race -tags etcdintegration ./...

.PHONY: verbose-integration-etcd
verbose-integration-etcd:
	go test -v -count=1 -race -tags etcdintegration ./...

