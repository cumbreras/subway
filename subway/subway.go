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
	client *pubsub.Client
}

type Message struct {
	Data             []byte
	SubscriptionName string
}

// Render will print to the stdout the message
func (m Message) Render() {
	fmt.Printf("Message received for %s:\n %s\n", m.SubscriptionName, string(m.Data))
}

// New returns the Subway Client
func New() Subway {
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

	return Subway{client: client}
}

func (s Subway) MessagesFromSubscription(subscription string, messages chan<- Message) {
	go s.pullMessages(subscription, messages)
}

// ListSubscriptionsFromEnvironment will return all subscriptions available for a particular ENV
func (s Subway) ListSubscriptionsFromEnvironment() ([]*pubsub.Subscription, error) {
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

		if strings.Contains(sub.ID(), "staging") {
			fmt.Printf("subscription found %s\n", sub.ID())
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

func (s Subway) pullMessages(SubscriptionName string, messages chan<- Message) {
	ctx := context.Background()
	var mu sync.Mutex
	sub := s.client.Subscription(SubscriptionName)
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		messages <- Message{SubscriptionName: sub.ID(), Data: msg.Data}
		defer mu.Unlock()
	})

	if err != nil {
		log.Fatal(err)
	}
}
