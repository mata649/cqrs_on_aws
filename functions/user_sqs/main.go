package main

import (
	"log"

	"github.com/mata649/cqrs_on_aws/functions/user_sqs/bootstrap"
)

func main() {
	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
