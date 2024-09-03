package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

/*
--	===================	--
--	--	server vars		--
--	===================	--
*/

var wsChan = make(chan WsJSONPayload)

var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

/*
--	===================	--
--	--	server 	funcs	--
--	===================	--
*/

// Home renders the home page view
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)

	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJSONResponse defines the response sent back from websocket
type WsJSONResponse struct {
	Action 		string `json:"action"`
	Message 	string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsJSONPayload struct {
	Action 		string 				`json:"action"`
	Username 	string				`json:"username"`
	Message 	string 				`json:"message"`
	Conn 		WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJSONResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConnection{Conn: ws}
	// adds client ws connection to "clients" Map
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// starts infinite/self-healing goroutine for anything sent over established client-ws
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {

	defer func() {
		// r == "recover" interface syntax
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsJSONPayload

	// shorthand for infinite loop (anytime we get a request with a payload (JSON Post or JSON file from JS client))
	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			// pass request payload (JSON) to our channel
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJSONResponse

	for {
		event := <- wsChan

		response.Action = "Got to WsChannel Listener"
		response.Message = fmt.Sprintf("Arbitrary msg where action was: %s", event.Action)
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJSONResponse) {

	// for each client...
	for client := range clients {
		// write JSON response to client over ws... ?
		err := client.WriteJSON(response)

		// if no client or client err (means client left) close client-ws
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}

}

func renderPage(w http.ResponseWriter, template string, data jet.VarMap) error {
	view, err := views.GetTemplate(template)

	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
