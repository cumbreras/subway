# subway
GCP PubSub Event Debugger: Read and debug gcp pubsub events easily

# Motivation
Production debugging sometimes can be hard, not knowing what event was sent or
if the structure of the message differs from what we expect. Having a visualization tool where we can see what is coming
is great to investigate what errors are happening.

# Prerequisites
Go 1.1/1.2, GoogleCloud JSON Auth file and the name of the project that you are targeting within your GCP Org.

# Installation
Not ready yet, but you can run the project with `go run main.go -env=XXX` or just build the project and run it normally.

# Flags
[] environment flag i.e. `-env=staging`
[] subscription flag i.e. `-subscription=payment-created`
