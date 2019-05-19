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

type Message struct {
	data             []byte
	subscriptionName string
}

func (m Message) render() {
	fmt.Printf("Message received for %s:\n %s\n", m.subscriptionName, string(m.data))
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
	messages := make(chan Message)

	for _, sub := range subscriptions {
		go s.pullMessages(sub.ID(), messages)
	}

	for msg := range messages {
		msg.render()
	}
}

func (s Subway) listSubscriptionsFromEnvironment() ([]*pubsub.Subscription, error) {
	fmt.Printf("Listing subscription for %s environment:\n", s.env)
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

func (s Subway) pullMessages(subscriptionName string, messages chan<- Message) {
	ctx := context.Background()
	var mu sync.Mutex
	sub := s.client.Subscription(subscriptionName)
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		messages <- Message{subscriptionName: sub.ID(), data: msg.Data}
		defer mu.Unlock()
	})

	if err != nil {
		log.Fatal(err)
	}
}
