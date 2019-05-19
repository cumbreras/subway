# subway
GCP PubSub Event Debugger: Read and debug gcp pubsub events easily

# Motivation
Production debugging sometimes can be hard, not knowing what event was sent or
if the structure of the message differs from what we expect. Having a visualization tool where we can see what is coming
is great to investigate what errors are happening.

This will provide a client to read the messages on the browser as they come. You could send the environment that you want to fetch from or the specific subscription that you are interesting in.

# Prerequisites
Go 1.12.X, GoogleCloud JSON Auth file and the name of the project that you are targeting within your GCP Org.

`$ export GOOGLE_CLOUD_PROJECT="awesome-project"`
`$ export GOOGLE_APPLICATION_CREDENTIALS="~/path/to/your/gcp-auth.json"`

# Run it locally
You can run the project with `go run main.go` or just build the project and run it normally.
