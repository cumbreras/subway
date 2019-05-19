package main

import (
	"flag"
	"github.com/cumbreras/subway/pubsubclient"
)

func main() {
	env := flag.String("env", "staging", "Environment to pull events from")
	flag.Parse()
	pc := pubsubclient.NewPubSubClient(*env)
	pc.Start()
}
