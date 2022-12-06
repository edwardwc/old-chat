package handlers

import (
	"fmt"
	"github.com/goombaio/namegenerator"
	"time"
)

func generateName() (name string) { // doesn't check if a name already exists
	name = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
	fmt.Printf("Name %v created\n", name)
	return
}

func checkName(name string) (works bool) {
	connectedClientsLock.Lock()
	if connectedClients[name] == nil {
		works = true
	} else {
		works = false
	}
	connectedClientsLock.Unlock()
	return
}
