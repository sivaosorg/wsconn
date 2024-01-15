# wsconn

![GitHub contributors](https://img.shields.io/github/contributors/sivaosorg/gocell)
![GitHub followers](https://img.shields.io/github/followers/sivaosorg)
![GitHub User's stars](https://img.shields.io/github/stars/pnguyen215)

Implement the WebSocket library using the Gin framework and the Gorilla WebSocket library in Go.

## Table of Contents

- [wsconn](#wsconn)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Modules](#modules)
    - [Running Tests](#running-tests)
    - [Tidying up Modules](#tidying-up-modules)
    - [Upgrading Dependencies](#upgrading-dependencies)
    - [Cleaning Dependency Cache](#cleaning-dependency-cache)
  - [Usage](#usage)
    - [WebSocket Connection](#websocket-connection)
    - [Registering Topics](#registering-topics)

## Introduction

This repository provides a simple WebSocket implementation using the Gin framework and Gorilla WebSocket library in Go. It allows real-time communication between clients and a server through WebSocket connections.

## Features

- Topic-based Subscription: Clients can subscribe to specific topics, and the server will broadcast messages to all subscribers of that topic.
- Dynamic Topic Registration: Topics can be dynamically registered by clients, allowing for flexible and dynamic communication channels.
- Concurrency Handling: The implementation uses Gorilla WebSocket and supports concurrent connections and message broadcasting.
- Closure Handling: Optionally, the server can be configured to handle closure events, such as detecting when a client connection is closed.

## Prerequisites

Golang version v1.20

## Installation

- Latest version

```bash
go get -u github.com/sivaosorg/wsconn@latest
```

- Use a specific version (tag)

```bash
go get github.com/sivaosorg/wsconn@v1.0.6
```

## Modules

Explain how users can interact with the various modules.

### Running Tests

To run tests for all modules, use the following command:

```bash
make test
```

### Tidying up Modules

To tidy up the project's Go modules, use the following command:

```bash
make tidy
```

### Upgrading Dependencies

To upgrade project dependencies, use the following command:

```bash
make deps-upgrade
```

### Cleaning Dependency Cache

To clean the Go module cache, use the following command:

```bash
make deps-clean-cache
```

## Usage

### WebSocket Connection

Connect to the WebSocket server using a WebSocket client. For example, in a browser, you can use JavaScript or tools like WebSocket.org's WebSocket Tester.

```javascript
const socket = new WebSocket("ws://localhost:8080/subscribe");

// Handle connection open event
socket.addEventListener("open", (event) => {
  console.log("WebSocket connection opened:", event);

  // Subscribe to a topic
  const subscription = {
    topic: "your-topic-name",
    content: "sample",
    userId: "user123",
    isPersistent: true,
    // Add any additional subscription parameters as needed
  };

  socket.send(JSON.stringify(subscription));
});

// Handle incoming messages
socket.addEventListener("message", (event) => {
  const message = JSON.parse(event.data);
  console.log("Received message:", message);
});

// Handle connection close event
socket.addEventListener("close", (event) => {
  console.log("WebSocket connection closed:", event);
});

// Handle connection error event
socket.addEventListener("error", (event) => {
  console.error("WebSocket error:", event);
});
```

### Registering Topics

You can dynamically register topics using a RESTful API endpoint:

```bash
curl -X POST http://localhost:8080/register -d '{"topic": "your-topic"}'
```
