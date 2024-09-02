package handlers

import (
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
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
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
