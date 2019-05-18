package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

type PubSubClient struct {
	env    string
	client *pubsub.Client
}

func main() {
	env := flag.String("env", "staging", "Environment to pull events from")
	flag.Parse()
	pc := newPubSubClient(*env)
	pc.start()

}

func (pc *PubSubClient) start() {
	subs, err := pc.listSubscriptionsFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	pc.printEventsFromSubscriptions(subs)
}

func (pc *PubSubClient) printEventsFromSubscriptions(subscriptions []*pubsub.Subscription) {
	events := make(chan []byte)

	for _, sub := range subscriptions {
		go pc.readEvents(sub.ID(), events)
	}

	for e := range events {
		fmt.Printf("Event received:\n %s\n", string(e))
	}
}

func newPubSubClient(env string) PubSubClient {
	ctx := context.Background()
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	client, err := pubsub.NewClient(ctx, proj)

	if err != nil {
		log.Fatal(err)
	}

	return PubSubClient{client: client, env: env}
}

func (pc *PubSubClient) listSubscriptionsFromEnvironment() ([]*pubsub.Subscription, error) {
	fmt.Println("Listing all subscriptions from the project:")
	ctx := context.Background()
	var subs []*pubsub.Subscription
	it := pc.client.Subscriptions(ctx)

	for {
		s, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		if strings.Contains(s.ID(), pc.env) {
			fmt.Printf("subscription found %s\n", s.ID())
			subs = append(subs, s)
		}
	}

	return subs, nil
}

func (pc *PubSubClient) readEvents(subscriptionName string, events chan<- []byte) {
	ctx := context.Background()
	var mu sync.Mutex
	sub := pc.client.Subscription(subscriptionName)
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		events <- msg.Data
		defer mu.Unlock()
	})

	if err != nil {
		panic(err)
	}
}
