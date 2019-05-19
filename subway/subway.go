package subway

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

// Subway client for fetching GCP PubSub Events
type Subway struct {
	env    string
	client *pubsub.Client
}

// New returns the Subway Client
func New(env string) Subway {
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

	return Subway{client: client, env: env}
}

// Start will read subscriptions and pull the events payload to the stdout
func (s Subway) Start() {
	subs, err := s.listSubscriptionsFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	s.eventsFromSubscriptions(subs)
}

func (s Subway) eventsFromSubscriptions(subscriptions []*pubsub.Subscription) {
	events := make(chan []byte)

	for _, sub := range subscriptions {
		go s.readEvents(sub.ID(), events)
	}

	for e := range events {
		fmt.Printf("Event received:\n %s\n", string(e))
	}
}

func (s Subway) listSubscriptionsFromEnvironment() ([]*pubsub.Subscription, error) {
	fmt.Println("Listing all subscriptions from the project:")
	ctx := context.Background()
	var subs []*pubsub.Subscription
	it := s.client.Subscriptions(ctx)

	for {
		sub, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		if strings.Contains(sub.ID(), s.env) {
			fmt.Printf("subscription found %s\n", sub.ID())
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

func (s Subway) readEvents(subscriptionName string, events chan<- []byte) {
	ctx := context.Background()
	var mu sync.Mutex
	sub := s.client.Subscription(subscriptionName)
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		events <- msg.Data
		defer mu.Unlock()
	})

	if err != nil {
		log.Fatal(err)
	}
}
