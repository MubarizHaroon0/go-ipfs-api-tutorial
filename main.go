package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	ipfs "github.com/ipfs/go-ipfs-api"
)

// Storage represents the decentralized storage system
type Storage struct {
	sh *ipfs.Shell
	mu sync.Mutex
}

// NewStorage creates a new instance of Storage
func NewStorage() *Storage {
	return &Storage{
		sh: ipfs.NewShell("localhost:5001"),
	}
}

// StoreData stores data in the decentralized storage system
func (s *Storage) StoreData(data interface{}) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Add data to IPFS
	cid, err := s.sh.Add(bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	return cid, nil
}

// RetrieveData retrieves data from the decentralized storage system
func (s *Storage) RetrieveData(cid string, out interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve data from IPFS
	res, err := s.sh.Cat(cid)
	if err != nil {
		return err
	}
	defer res.Close()

	// Read the content of res into a []byte
	jsonData, err := ioutil.ReadAll(res)
	if err != nil {
		return err
	}

	// Deserialize data from JSON
	if err := json.Unmarshal(jsonData, out); err != nil {
		return err
	}

	return nil
}

func main() {
	// Create a new instance of Storage
	storage := NewStorage()

	// Sample data to store
	data := map[string]interface{}{
		"key":   "value",
		"array": []int{1, 2, 3},
	}

	// Store data
	cid, err := storage.StoreData(data)
	if err != nil {
		log.Fatal("Error storing data:", err)
	}
	fmt.Println("Data stored with CID:", cid)

	// Retrieve data
	var retrievedData map[string]interface{}
	err = storage.RetrieveData(cid, &retrievedData)
	if err != nil {
		log.Fatal("Error retrieving data:", err)
	}
	fmt.Println("Retrieved data:", retrievedData)
}
