logger
======

`logger` is not well named. it is an event bus that can write events in a variety of formats.

## Requirements
- Enable or disable event types by flag
- Should be able to be configured by a config object that can be parsed from json or yaml
- Show or hide output for event types by flag
  - the implication is some events are used only for eventing, some are used for tracing
- Add one or many listeners for events by flag
  - each listener should be identifyable so they can be enabled or disabled, or removed.
- Output can be to one or more writers
  - each writer just needs to satisfy an interface to allow messages to be passed to it.
- Default supported output formats are text and json, but more can be added by users.
- Should support a number of message types out of the box:
  - Informational (string messages)
  - Error (errors or exceptions)
  - HTTP Request (using the primitives from the stdlib)

## Design
- Messages are represented by strongly typed "events"
- An "event" is composed of
  - Timestamp : can be set by the caller, but generally when the event was generated.
  - Flag : the "type" of the event; determines if the event is enabled or disabled.
- Most basic messages (info, debug, warning) are just a flag+message pair, represented by a `MessageEvent`.
- Most error events (error, fatal) are just a flag+error pair, represented by an `ErrorEvent`.
- Writing to JSON or writing to Text should look the same to the caller, and be determined by the configured writer(s).

Example Syntax:

Creating a logger based on environment configuration:
```go
log := logger.NewFromEnv()
```

Creating a logger using a config object:
```go
log := logger.NewFromConfig(&logger.Config{
  Flags: []string{"info", "error", "fatal"},
})
```

Logging an informational message:
```go
log.Infof("%s foo", bar) 
```

Listening for events with a given flag:
```go
log.Listen(log.Info, "datadog", func(e logger.Event) {
  //...
  // sniff the type of the message and do something specific with it
})
```

## Development notes:
We have a couple standing guidelines with go-logger:
- The library should be as general as possible; limit external dependencies.
  - Permissible external dependencies:
    - go-assert
    - go-exception
    - go-util
    - go-util/env
- Keep the existing API as stable as possible; add new methods, leave old ones.
- If you have to remove a method, have a clear migration path.
- Tests shouldn't write to stdout / stderr.
- The test suite should complete in < 1 second.
