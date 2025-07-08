## Use Case: Graceful Shutdown of an HTTP Server

In real-world production systems, abrupt termination of a running HTTP server (e.g., via `Ctrl+C`, system shutdown, or Docker stop) can cause active requests to be cut off, data loss, or resource leaks. A graceful shutdown ensures in-flight requests are **given time to complete**, and resources are released **cleanly**.

---

## What's Going On:

* An HTTP server is created using Go’s `net/http` package.
* Each request simulates a long-running operation (5 seconds).
* The request handler listens for cancellation via `r.Context().Done()`.
* A `context.WithCancel()` is used to propagate cancellation from the OS interrupt (like `Ctrl+C`) to the server and the request handler.
* A signal listener (`os.Signal`) captures `os.Interrupt`, and then **calls `Shutdown()`** with a timeout.
* A `context.WithTimeout()` ensures the server waits for up to 10 seconds to complete any active requests before forcefully shutting down.

---

## Why It Worked

* The handler checks for request cancellation using `r.Context()`, so it **responds early** if the request is cancelled.
* `context.WithCancel()` ties the entire lifecycle together — from the main process to the request-level handlers.
* `os.Signal` allows Go to listen to **external shutdown triggers** like `SIGINT` and initiate cleanup.
* `http.Server.Shutdown()` is designed to gracefully stop accepting new requests and wait for existing ones to finish.

> Without context and shutdown handling, your server would terminate immediately, potentially leaving in-progress requests hanging or half-written.

---

## Output

When running the server and hitting `Ctrl+C`:

```
Server listening on :8080
Handling request...
Interrupt received, shutting down...
Worker timeout or cancelled: context canceled
Server stopped cleanly.
```

---

## When To Use This Pattern

### Web Servers, APIs, and Microservices

* Gracefully handle shutdown signals (like SIGINT, SIGTERM)
* Ensure users don’t experience partial or broken responses
* Release DB connections, file handles, and clean up goroutines

### CLI Tools or Long-Running Services

* Respect user cancellation (`Ctrl+C`)
* Abort long-running operations cleanly

---
