# Lazuli Golang Package
![CI Workflow](https://github.com/augustoasilva/go-lazuli/actions/workflows/ci-workflow.yml/badge.svg)
![coverage](https://raw.githubusercontent.com/augustoasilva/go-lazuli/badges/.badges/main/coverage.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/augustoasilva/go-lazuli.svg)](https://pkg.go.dev/github.com/augustoasilva/go-lazuli)

Lazuli is a Golang package that aims to help to work with AT Protocol and Bluesky Social Network

The lazuli package is under development, so before the first release of 1.x, the code can change causing incompatibilities
against each early versions.

## Quickstart

To quickstart using go lazuli, first you need to add to your project by using:

```
go get github.com/augustoasilva/go-lazuli
```

After that you need to create a session, like:

```go
ctx := context.Background()
xrpcURL := os.Getenv("XRPC_URL")
wsURL := os.Getenv("WS_URL")

client := lazuli.NewClient(xrpcURL, wsURL)

identifier := os.Getenv("IDENTIFIER")
password := os.Getenv("PASSWORD")
sess, err := client.CreateSession(ctx, identifier, password)
if err != nil {
    slog.Error("error creating session", "error", err)
    panic(err)
}
slog.Info("session created", "session", sess)
```

And then you can start using it! You can access examples of how to use it at [example folder](https://github.com/augustoasilva/go-lazuli/tree/main/example).

## Contribution

If you want to contribute, you can do it by doing one (or more) of each options:
- Create an issue
  - If you found a bug, or if you have an idea of improvements, feel free to open an issue! But be aware that, for now
there is one developer maintaining it, so I will do my best to look at the issues and try to fix the bugs or implement your
your suggestion
- Open a pull request
  - You can look at the open issue(s), select one, comment in it that you have interest in doing it and code! So, fork the
project, develop the feature or the fix, and submit it through a pull-request.

## License

This project is licensed under MIT  ([LICENSE-MIT](LICENSE) or http://opensource.org/licenses/MIT).