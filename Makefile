.PHONY: build clean deploy gomodgen

build-user: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/user_service ./functions/user_apigateway


build-event: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/on_task_created ./functions/task_sqs
	
build-task: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/task_service ./functions/task_apigateway

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x ./gomod.sh
	./gomod.sh
