# Video Chat App using WebRTC and Golang
### A simple video chat application developed using Golang and React.
### Frontend

The `client` is written in React and uses [Vite](https://vitejs.dev/) for the dev server. Run the following commands in the `client` directory

* `npm i` to install all the dependencies
* `npm run dev` to start the local dev server

### Backend

Written in Go. A simple WebSocket server for signalling implemented using 
[gorilla/websocket](https://github.com/gorilla/websocket)

* `go build` to compile and build the binary
* `./video-chat-app` to run the backend server on `:8000`

<br>

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-javascript.svg)](https://forthebadge.com)
