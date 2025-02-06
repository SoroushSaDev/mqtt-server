package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ACL struct {
	Username string  `json:"username"`
	Topics   []Topic `json:"topics"`
}

type Topic struct {
	Topic      string `json:"topic"`
	Permission string `json:"permission"`
}

var (
	mu        sync.Mutex
	authRules = &auth.Ledger{
		Auth: auth.AuthRules{},
		ACL:  auth.ACLRules{},
	}
)

func fetchDataFromAPI() {
	apiURL := "http://127.0.0.1:8000/mqtt/auth"

	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("Error fetching data:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println("Error reading response body:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var data struct {
			Users []User `json:"users"`
			ACLs  []ACL  `json:"acls"`
		}

		if err := json.Unmarshal(body, &data); err != nil {
			log.Println("Error unmarshalling data:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		mu.Lock()
		authRules.Auth = auth.AuthRules{}
		authRules.ACL = auth.ACLRules{}

		for _, user := range data.Users {
			authRules.Auth = append(authRules.Auth, auth.AuthRule{
				Username: auth.RString(user.Username),
				Password: auth.RString(user.Password),
				Allow:    true,
			})
		}

		for _, acl := range data.ACLs {
			aclRule := auth.ACLRule{Username: auth.RString(acl.Username), Filters: auth.Filters{}}
			for _, topic := range acl.Topics {
				var permission auth.Access
				switch auth.RString(topic.Permission) {
				case "rw":
					permission = auth.ReadWrite
					break
				case "r":
					permission = auth.ReadOnly
					break
				case "w":
					permission = auth.WriteOnly
					break
				default:
					permission = auth.Deny
					break
				}
				aclRule.Filters[auth.RString(topic.Topic)] = permission
			}
			authRules.ACL = append(authRules.ACL, aclRule)
		}

		mu.Unlock()

		log.Println("Updated users and ACLs from API")
		time.Sleep(10 * time.Second) // Adjust polling interval as needed
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	go fetchDataFromAPI()

	server := mqtt.New(nil)
	err := server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: authRules,
	})
	if err != nil {
		log.Fatal(err)
	}

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: ":1883",
	})
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	// server.Log.Info("main.go finished")
}
