.PHONY: server build 2test lint help docker clean 2async asset test

run:
	go run . server

server:
	go run . server

build:
	go build -race -ldflags "-s -w -X 'demo/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X demo/api/controller/v1.BuildVersion=1.0.1" -tags=jsoniter -o demo .

docker:
	docker build -t demo .

clean:
	go mod tidy

test:
	go test $(go list ./... | grep -v /vendor/) -v -count=1 -coverpkg=./...

outcov:
	go test  -v -count=1 -coverpkg=./... -test.short -coverprofile=coverage.out -timeout=10s `go list ./... | grep -v /vendor/` -json > report.json

sonar: outcov
	sonar-scanner \
	  -Dsonar.projectKey=demo \
	  -Dsonar.sources=. \
	  -Dsonar.host.url=http://localhost:9000 \
	  -Dsonar.login=24652fad5f8bafa4b06bd92d5d79017fc28a07eb \
	  -Dsonar.sources.inclusions='**/*.go' \
	  -Dsonar.exclusions='doc/**,**/*_test.go,**/vendor/**,.git/**,.glide/**,asset/**,internal/asset/**' \
	  -Dsonar.tests=. -Dsonar.test.inclusions='**/*_test.go' -Dsonar.test.exclusions='**/vendor/**' \
	  -Dsonar.go.tests.reportPaths=report.json  -Dsonar.go.coverage.reportPaths=coverage.out

asset:
	go-bindata -pkg asset -o internal/asset/bindata.go asset

help:
	@echo "make: compile packages and dependencies"
	@echo "  make run: go run at server"
	@echo "  make server: go run at server"
	@echo "  make build: go build"
	@echo "  make lint: golint ./..."
	@echo "  make clean: remove object files and cached files"