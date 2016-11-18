# workengine

Experimenting with generic, async workers and tracking their status over
websockets.

# Usage

Start up the server:

`go build && ./workengine`

Fire away and watch the `FooMachines` work!:

`for i in {1..10}; do curl -X POST http://localhost:3000/run; done`
