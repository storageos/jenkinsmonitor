package main

import (
	"encoding/json"
	"fmt"
	"github.com/storageos/jenkinsmonitor/relaydriver"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type RESTRelay struct {
	*relaydriver.Driver
}

var inlineHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Relay trigger</title>
</head>

<body>
<script>

function relayOn(r) {
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/"+r, true);
	xhr.send("1")
}

function relayOff(r) {
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/"+r, true);
	xhr.send("0")
}

</script>
<h1>Relay triggers</h1>
<ul>
<li>green:<button type="button" onclick="relayOn('green')">on</button>:<button type="button" onclick="relayOff('green')">off</button></li>
<li>yellow:<button type="button" onclick="relayOn('yellow')">on</button>:<button type="button" onclick="relayOff('yellow')">off</button></li>
<li>red:<button type="button" onclick="relayOn('red')">on</button>:<button type="button" onclick="relayOff('red')">off</button></li>
<li>alarm:<button type="button" onclick="relayOn('alarm')">on</button>:<button type="button" onclick="relayOff('alarm')">off</button></li>
</ul>
</body>

</html>`

func (r *RESTRelay) JSONState() ([]byte, error) {
	errs := make([]string, 0)
	m := struct {
		R1 bool `json:"red"`
		R2 bool `json:"yellow"`
		R3 bool `json:"green"`
		R4 bool `json:"alarm"`
	}{}

	var err error
	m.R1, err = r.Driver.GetState(relaydriver.Relay1)
	if err != nil {
		errs = append(errs, err.Error())
	}

	m.R1, err = r.Driver.GetState(relaydriver.Relay1)
	if err != nil {
		errs = append(errs, err.Error())
	}

	m.R1, err = r.Driver.GetState(relaydriver.Relay1)
	if err != nil {
		errs = append(errs, err.Error())
	}

	m.R1, err = r.Driver.GetState(relaydriver.Relay1)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors reading state from driver. [%s]", strings.Join(errs, ", "))
	}

	return json.Marshal(&m)
}

func (r *RESTRelay) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("Request %s, %s\n", req.Method, req.URL.Path)

	if req.URL.Path == "/" {
		resp.Header().Add("Content-Type", "text/html")
		fmt.Fprintln(resp, inlineHTML)
		return
	}

	if req.Method == "GET" {
		buf, err := r.JSONState()
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, err.Error())
			return
		}

		fmt.Fprintf(resp, string(buf))
		return
	}

	if req.Method == "POST" {
		var relaySelected int
		switch req.URL.Path {
		case "/red":
			relaySelected = relaydriver.Relay1
		case "/yellow":
			relaySelected = relaydriver.Relay2
		case "/green":
			relaySelected = relaydriver.Relay3
		case "/alarm":
			relaySelected = relaydriver.Relay4
		default:
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch string(buf) {
		case "0":
			r.SetLow(relaySelected)
		case "1":
			r.SetHigh(relaySelected)
		default:
			resp.WriteHeader(http.StatusBadRequest)
		}
		return
	}
}

func (r *RESTRelay) Serve() error {
	return http.ListenAndServe(":6543", r)
}

func NewRESTRelay() (*RESTRelay, error) {
	rd, err := relaydriver.NewDriver()
	if err != nil {
		return nil, err
	}

	// Disconnect the relay board on shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c) // All signals fatal, probably not right
		<-c
		rd.Shutdown()

		os.Exit(0)
	}()

	return &RESTRelay{rd}, nil
	return &RESTRelay{nil}, nil
}

func main() {
	rr, err := NewRESTRelay()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = rr.Serve()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
