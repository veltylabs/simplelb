# SimpleLB

Simple LB is the simplest Load Balancer ever created.

It uses a Round Robin algorithm to send requests to a set of backends and supports retries.

It also performs active health checks to remove unhealthy backends from the pool and passive recovery to add them back.

## How to use

### Installation

To install `simplelb`, use the following command:

```bash
go install github.com/kasvith/simplelb/cmd/simplelb@latest
```

### Running

Once installed, you can run `simplelb` from your terminal:

```bash
Usage:
  -backends string
        Load balanced backends, use commas to separate
  -port int
        Port to serve (default 3030)
```

**Example:**

To add the following backends to the load balancer:
- http://localhost:3031
- http://localhost:3032
- http://localhost:3033
- http://localhost:3034

```bash
simplelb -backends="http://localhost:3031,http://localhost:3032,http://localhost:3033,http://localhost:3034"
```

## How to test

To run the unit tests, use the following command:

```bash
go test -cover
```
