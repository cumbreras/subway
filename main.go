package main

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

func main() {
	ctx := context.Background()
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	fmt.Println("Listing all subscriptions from the project:")
	subs, err := list(client)
	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan string, 100)
	for _, sub := range subs {
		go pullMsgs(client, sub.ID(), messages)
	}

	for m := range messages {
		fmt.Println("Messages %v", m)
	}
}

func list(client *pubsub.Client) ([]*pubsub.Subscription, error) {
	ctx := context.Background()
	var subs []*pubsub.Subscription
	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.Contains(s.ID(), "staging") {
			fmt.Println("subscription found %v", s.ID())
			subs = append(subs, s)
		}
	}
	return subs, nil
}

func pullMsgs(client *pubsub.Client, subName string, messages <-chan string) error {
	ctx := context.Background()
	var mu sync.Mutex
	sub := client.Subscription(subName)
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
	})
	if err != nil {
		return err
	}
	return nil
}
