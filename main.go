package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/cumbreras/subway/subway"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:1337", "http service address")

var upgrader = websocket.Upgrader{}

func events(w http.ResponseWriter, r *http.Request) {
	s := subway.New()
	subs, err := s.ListSubscriptionsFromEnvironment()

	for sub, _ := range subs {
		fmt.Printf("%s", sub)
	}

	messages := make(chan subway.Message)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()
	for {
		mt, subscription, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		go s.MessagesFromSubscription(string(subscription), messages)

		for msg := range messages {
			msg.Render()
			err = c.WriteMessage(mt, msg.Data)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

		log.Printf("recv: %s", subscription)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/events")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/events", events)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    (function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("Loading Messages");
        }
        ws.onclose = function(evt) {
            print("Closing Connection");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("Message: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    })();

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="close">Stop Messages</button>
<div id="env-fetch">
<label>Include the subscription you want to read from</label>
<p><input id="input" type="text" value="staging">
<button id="send">Send</button>
</div>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
