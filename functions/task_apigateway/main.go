package main

import (
	"log"

	"github.com/mata649/cqrs_on_aws/functions/task_apigateway/bootstrap"
)

func main() {
	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}

}
