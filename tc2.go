package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/patbcole117/testC2/comms"
)

var (
	apiIp 	= "127.0.0.1"
	apiPort = "8000"
	apiVer	= "v1"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		os.Exit(0)
	}

	tx, err := comms.NewCommsPackageTX("http")
	if err != nil {
		panic(err)
	}

	var (
		pstr_List 		= flag.String("L", "", "-L [beacons | jobs | nodes] [-id <id>]")
		pstr_New 		= flag.String("N", "", "-N [node | job -id <nid> -x '<command>']")
		pstr_Delete 	= flag.String("D", "", "-D [node -id <nid> | job -id <jid>]")
		pstr_Update 	= flag.String("U", "", "-U [node -id <nid> | job -id <jid>] [-ip <ip] [-p <port>] [-n <name>]")
		pbool_Start		= flag.Bool("START", false, "-START -id <nid> ")
		pbool_Stop		= flag.Bool("STOP", false, "-STOP -id <nid>")
		pbool_Check		= flag.Bool("CHECK", false, "-CHECK")
		pstr_Id			= flag.String("id", "", "the id of the thing to create, delete, display ot update.")
		pstr_Ip			= flag.String("a", "", "a new ip to replace the old.")
		pstr_Port		= flag.String("p", "", "a new port to replace the old.")
		pstr_Name		= flag.String("n", "", "a new name to replace the old.")
	)
	flag.Parse()

	if *pbool_Check {
		if err := Check(tx); err != nil {panic(err)}
	}
	
	switch *pstr_List{
	case "nodes":
		if *pstr_Id != "" {
			if err := ListNode(tx, *pstr_Id); err != nil {panic(err)}
		} else {
			if err := ListNodes(tx); err != nil {panic(err)}
		}
	}

	switch *pstr_New{
	case "node":
		if err := NewNode(tx); err != nil{panic(err)}
	}

	switch *pstr_Delete{
	case "node":
		if err := ValidateInt(*pstr_Id); err != nil{panic(err)}
		if err := DeleteNode(tx, *pstr_Id); err != nil {panic(err)}
	}

	switch *pstr_Update{
	case "node":
		if err := ValidateInt(*pstr_Id); err != nil{panic(err)}
		if err := UpdateNode(tx, *pstr_Id, *pstr_Name, *pstr_Ip, *pstr_Port); err != nil {panic(err)}
	} 

	if *pbool_Start {
		if err := StartNode(tx, *pstr_Id); err != nil{panic(err)}
	}

	if *pbool_Stop {
		if err := StopNode(tx, *pstr_Id); err != nil{panic(err)}
	}
}

func Check(c comms.CommsPackageTX) error {
	url := "http://" + apiIp + ":" + apiPort + "/"
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteNode(c comms.CommsPackageTX, id string) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/" + id + "/x/"
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListNodes(c comms.CommsPackageTX) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/"
	
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListNode(c comms.CommsPackageTX, id string) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/" + id + "/"
	
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNode(c comms.CommsPackageTX) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/new/"
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartNode(c comms.CommsPackageTX, id string) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/" + id + "/start/"
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func StopNode(c comms.CommsPackageTX, id string) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/" + id + "/stop/"
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateNode(c comms.CommsPackageTX, id, fName, fIp, fPort string) error {
	url := "http://" + apiIp + ":" + apiPort + "/" + apiVer + "/l/" + id + "/"
	
	body := map[string]string {
		"name": fName,
		"ip": fIp,
		"port": fPort,
	}
	
	resp, err := c.SendJSON(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
			panic(err)
		}
		fmt.Println(prettyJson.String())

	} else {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateInt(integer string) error {
	_, err := strconv.Atoi(integer)
	if err != nil {
		return err
	}
	return nil
}

func ValidatePort(port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if p > 655535 || p < 1 {
		return errors.New("invalid port")
	}
	return nil
}
