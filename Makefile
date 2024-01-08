.PHONY: build clean deploy gomodgen

build-user: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/user_function ./functions/user_apigateway

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x ./gomod.sh
	./gomod.sh
