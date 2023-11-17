# actorsystem-go

`actorsystem-go` is a Go library that provides an Actor class for connecting to AMQP, creating and/or binding a queue using a routing key, and consuming messages.

## Installation

```bash
go get github.com/actorsystem/actorsystem-go
```

## Usage

1. Create a `.env` file in your project directory with the following content:

```
AMQP_URL=your_amqp_url_here
```

2. Install the required packages:

```bash
go get github.com/joho/godotenv
go get github.com/actorsystem/actorsystem-go
```

3. Use the library in your Go code:

```go
package main

import (
	"fmt"
	"github.com/your-username/actorsystem-go"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the AMQP URL from the environment
	amqpURI := os.Getenv("AMQP_URL")
	if amqpURI == "" {
		log.Fatal("AMQP_URL not found in .env file")
	}

	actor := &actorsystem.Actor{
		Exchange:   "your_exchange",
		RoutingKey: "your_routing_key",
		Queue:      "your_queue",
		Handler: func(message []byte) error {
			// Your message handling logic here
			fmt.Printf("Received message: %s\n", message)
			return nil
		},
	}

	err = actor.Start(amqpURI)
	if err != nil {
		fmt.Printf("Error starting actor: %v", err)
	}
}
```

## Lush Lounge Undertones

This library is designed to provide a smooth and relaxing experience while working with AMQP messaging in your Go applications. Kick back, enjoy the elegance of the `Actor` class, and let the messages flow through your system effortlessly.

Cheers to a stress-free messaging experience with `actorsystem-go`!
```

Make sure to replace `"your_amqp_url_here"` with your actual AMQP URL in the `.env` file.
