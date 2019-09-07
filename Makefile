.PHONY: etcd-up
etcd-up:
	docker run -d -p 2379:2379 -p 2380:2380 --name etcd-gcr-v3.4.0-graphx gcr.io/etcd-development/etcd:v3.4.0 \
		/usr/local/bin/etcd \
		--name s1 \
		--listen-client-urls http://0.0.0.0:2379 \
		--advertise-client-urls http://0.0.0.0:2379 \
		--listen-peer-urls http://0.0.0.0:2380 \
		--initial-advertise-peer-urls http://0.0.0.0:2380 \
		--initial-cluster s1=http://0.0.0.0:2380 \
		--initial-cluster-token tkn \
		--initial-cluster-state new \
		--log-level debug \
		--logger zap \
		--log-outputs stderr
	docker exec etcd-gcr-v3.4.0-graphx /bin/sh -c "/usr/local/bin/etcdctl version"

.PHONY: etcd-down
etcd-down:
	-docker kill etcd-gcr-v3.4.0-graphx
	-docker rm etcd-gcr-v3.4.0-graphx

.PHONY: unit-verbose
unit-verbose:
	go test -v -count=1 -race ./...

.PHONY: automated-integration-etcd
automated-integration-etcd:
	make etcd-down
	make etcd-up
	go test -count=1 -race -tags integration-etcd ./...

.PHONY: verbose-integration-etcd
verbose-integration-etcd:
	go test -v -count=1 -race -tags integration-etcd ./...
