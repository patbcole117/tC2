package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	ERROR    	= "-1"
	STOPPED     = "0"
	LISTENING   = "1"

	SERVER_DEFAULT_IP      		string        	= "127.0.0.1"
	SERVER_DEFAULT_PORT    		string          = "80"
	SERVER_DEFAULT_TIMEOUT 		time.Duration 	= 3 * time.Second
	SERVER_DEFAULT_NAME_SIZE 	int 			= 10
)

type Node struct {
	Id 		int				`bson:"id" json:id`
	Name   	string			`json:"name"`
	Ip     	string   		`json:"ip"`
	Port   	string      	`json:"port"`
	Status 	string      	`json:"status"`
    Dob     time.Time		`json:"dob"`
	Hello 	time.Time		`json:"hello"`
	Server 	*http.Server    `json:"server"`
}

func NewNode() Node {
	n := Node {
	Ip:		SERVER_DEFAULT_IP,
	Port:	SERVER_DEFAULT_PORT,
	Status:	STOPPED,
	}
    n.Dob = time.Now()
	n.Hello = time.Now()
	n.initName(SERVER_DEFAULT_NAME_SIZE)
	return n
}

func (n *Node) initName(sz int) {
	abc := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, sz)
	for i := range b {
		b[i] = abc[rand.Intn(len(abc))]
	}
	n.Name = string(b)
}

func (n *Node) SrvStart() error {
	addy := n.Ip + ":" + n.Port
	m := http.NewServeMux()
	s := http.Server{
		Addr:         addy,
		Handler:      m,
		ReadTimeout:  SERVER_DEFAULT_TIMEOUT,
		WriteTimeout: SERVER_DEFAULT_TIMEOUT,
	}
	m.HandleFunc("/", n.urlRoot)

	n.Server = &s
	go s.ListenAndServe()
	time.Sleep(SERVER_DEFAULT_TIMEOUT)

	n.Status = LISTENING
	return nil
}

func (n *Node) SrvStop() error {
	if n.Server == nil {
		return errors.New("server is nil.")
	}
	if err := n.Server.Close(); err != nil {
		return err
	}
	time.Sleep(SERVER_DEFAULT_TIMEOUT)
	n.Status = STOPPED
	return nil
}

func (n *Node) UnmarshalJSON(j []byte) error {
    type Alias Node
    aux := &struct {
        Dob 	string	`json:"dob"`
		Hello 	string	`json:"hello"`	
        *Alias
    }{
        Alias:  (*Alias)(n),
    }

    if err := json.Unmarshal(j, &aux); err != nil {
        return err
    }
   
    t, err := time.Parse(time.RFC3339, aux.Dob)
    if err != nil {
        return err
    }
	n.Dob = t
    
	t, err = time.Parse(time.RFC3339, aux.Hello)
    if err != nil {
        return err
    }
	n.Hello = t

    return nil
}

func (n *Node) MarshalJSON() ([]byte, error) {
    type Alias Node
    return json.Marshal(&struct {
        Dob 	string  `json:"dob"`
		Hello 	string  `json:"hello"`
        *Alias
    }{
        Dob: 	n.Dob.Format(time.RFC3339),
		Hello: 	n.Hello.Format(time.RFC3339),
        Alias:  (*Alias)(n),
    })
}

func (n *Node) ToJsonPretty() (string) {
	b, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		return string(make([]byte, 0))
	}
	return string(b)
}

func (n *Node) urlRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Date", time.Now().Format(time.UnixDate))
	msg := fmt.Sprintf("[+] Hello from %s.\n", n.Name)
	io.WriteString(w, msg)
}