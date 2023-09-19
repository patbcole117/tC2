package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/patbcole117/tC2/node"
)

var db	            dbManager   = NewDBManager()
var masterNodes 	[]node.Node
var url             string      = "127.0.0.1:8000"

func Run() {
    r := chi.NewRouter()
    r.Get("/",                  Check)
    r.Get("/v1/l",              GetNodesFromDB)
    r.Get("/v1/l/new",          NewNode)
    r.Get("/v1/l/{id}",         GetNodeFromDB)
    r.Post("/v1/l/{id}",        UpdateNode)
    r.Get("/v1/l/{id}/start",   StartNode)
    r.Get("/v1/l/{id}/stop",    StopNode)
    r.Get("/v1/l/{id}/x",       DeleteNode)
    
    http.ListenAndServe(url, r)
}

func Check (w http.ResponseWriter, r *http.Request) {
    fmt.Println(masterNodes)
    _, bmsg := FgoodMsg("check.")
    w.Write(bmsg)
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("DeleteNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    _, err = db.DeleteNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("DeleteNode->db.DeleteNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, ix := GetNodeFromMaster(id)
    if  n == nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg(fmt.Sprintf("node %d does not exist in master.", id))
        w.Write(bmsg)
        return
    }

    if err = n.SrvStop(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
}

    masterNodes = append(masterNodes[:ix], masterNodes[ix+1:]...)
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("deleted %d.", id))
    w.Write(bmsg)
}

func NewNode (w http.ResponseWriter, r *http.Request) {
    n := node.NewNode()
    id, err := db.GetNextNodeID()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("NewNode->db.GetNextNodeID->"+err.Error())
        w.Write(bmsg)
        return
    }
    n.Id = id

    result, err := db.InsertNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("NewNode->db.InsertNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    masterNodes = append(masterNodes, n)
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("inserted %s.", result.InsertedID))
    w.Write(bmsg)
}

func GetNodeFromDB(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodeFromDB->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, err := db.GetNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodeFromDB->db.GetNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    jnode, err := json.Marshal(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodeFromDB->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jnode)
}

func GetNodesFromDB(w http.ResponseWriter, r *http.Request) {
    nodes, err := db.GetNodes()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodesFromDB->db.GetNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    jnodes, err := json.Marshal(nodes)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodesFromDB->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jnodes)
}

func StartNode(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StartNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, _ := GetNodeFromMaster(id)
    if  n == nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg(fmt.Sprintf("node %d does not exist in master.", id))
        w.Write(bmsg)
        return
    }

    n.SrvStart()

    res, err := db.UpdateNode(*n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StartNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("started %d.", n.Id))
    w.Write(bmsg)
}

func StopNode(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, _ := GetNodeFromMaster(id)
    if  n == nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg(fmt.Sprintf("node %d does not exist in master.", id))
        w.Write(bmsg)
        return
    }

    if err = n.SrvStop(); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            _, bmsg := FbadMsg("StopNode->db.UpdateNode->"+err.Error())
            w.Write(bmsg)
            return
    }

    res, err := db.UpdateNode(*n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("stopped %d.", n.Id))
    w.Write(bmsg)

}

func UpdateNode (w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->io.ReadAll->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, _:= GetNodeFromMaster(id);
    if  n == nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg(fmt.Sprintf("node %d does not exist in master.", id))
        w.Write(bmsg)
        return
    }
    
    var b map[string]string
    if err = json.Unmarshal(body, &b); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    n.Id        = id
    if b["name"] != ""{
        n.Name      = b["name"]
    }
    if b["ip"] != ""{
        n.Ip        = b["ip"]
    }
    if b["port"] != ""{
        n.Port        = b["Port"]
    }

    res, err := db.UpdateNode(*n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("updated %d.", n.Id))
    w.Write(bmsg)
}

//Helpers

func GetNodeFromMaster(id int) (*node.Node, int) {
    for i := range masterNodes {
        if masterNodes[i].Id == id {
            return &masterNodes[i], i
        }
    }
    return nil, 0
}

func FdebugMsg (msg string) (string, []byte) {
	s := fmt.Sprintf(`{"type": "debug", "msg": %s}`, msg)
	return s, []byte(s)
}
func FbadMsg (msg string) (string, []byte) {
	s := fmt.Sprintf(`{"type": "bad", "msg": %s}`, msg)
	return s, []byte(s)
}
func FgoodMsg (msg string) (string, []byte) {
	s := fmt.Sprintf(`{"type": "good", "msg": %s}`, msg)
	return s, []byte(s)
}