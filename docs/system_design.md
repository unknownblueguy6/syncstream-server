# syncstream System Design

### Programming Languages

- **Golang**: In addition to wanting to learn a new language, the memory safe, simple syntax, and concurrency model (goroutines) of the language makes it suitable for the backend development of our server

- **TypeScript**: TypeScript is a superset of JavaScript that introduces static typing, which means that variables, function parameters, and return values have specified types for easy debugging. 

## Libraries & Technologies
- **Gorilla/Websockets**: for the Golang implementation of the Websocket protocol. It provides a simple and easy-to-use API for creating WebSocket connections between clients and servers.
- **Google/UUID**: To generate/parse UUIDs on the Go server. The Google/UUID package provides functions to generate UUIDs based on different versions (e.g., UUIDv4) and to parse UUIDs from strings.
- **Redis Database**: To be used for ephemeral token storage. Ephemeral tokens are short-lived tokens used for authentication or authorization purposes. Redis's fast in-memory operations make it suitable for storing and retrieving these tokens efficiently.

## Server Endpoints

### HTTP/S
#### `/init`
- On installation, our extension sends an initializing ping in the form of a POST request. The server responds with a UUID and a bearer Authentication token, which is used for all future requests.
- **Request**: `{}`
- **Response**: `{"id": "string", "bearerAuth": "string"}`

#### `/create`
- Send a POST request to the web server to create a room. The server returns a 6 character alphabetic code, which can be used to join the room.
- **Request**: `{ "id": "string", "url": "string", "streamState": { "currentTime": "float64", "paused": "bool", "playbackRate": "float32" }, "timestamp": "uint64" }`
- **Response**: `{ "code": "string" }`

#### `/joinToken`
- Send a POST request to the web server to get the token to join a room. The client sends their UUID and the room code they want to join, and the server returns an ephemeral token back.
- **Request**: `{ "id": "string", "code": "string" }`
- **Response**: `{ "token": "string" }`

#### `/delete`
- Send a POST request to the web server to delete a room, and kick out anyone in the room. Only the room creator can delete a room.
- **Request**: `{ "id": "string", "code": "string" }`
- **Response**: `{ "success": "bool" }`

## Websocket
### `/join?{token}`
- Send a GET request to the web server, which will upgrade the connection to the Websocket protocol if the token is valid, and connect you to the room. The client will receive "messages" from the server, which can be either binary or text.

### Message Event Format
- **Message**: `{ "sourceID": "string", "timestamp": "uint64", "type": "uint8", "data": {...} }`  # data contains an arbitrary JSON object

## Client Side
- The client side includes a web browser running on the user's device.
- The browser displays the user interface of the web application, in the form of a browser extension.
- The browser extension establishes a Web Socket connection with the server.

## Server Side
- The server side includes a web server that handles incoming requests from clients.
- The server also establishes Web Socket connections with the clients.

