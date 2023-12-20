package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// Tracker represents a simple peer-to-peer tracker.
// It maintains a map of peers and the files they have.
type Tracker struct {
	peers map[string][]string // Map of peer addresses to their files
	lock  sync.Mutex          // Mutex for safe concurrent access to the peers map
}

// NewTracker creates and returns a new Tracker instance.
// It initializes the peers map and the mutex lock.
func NewTracker() *Tracker {
	return &Tracker{
		peers: make(map[string][]string),
		lock:  sync.Mutex{},
	}
}

// handleConnection manages a single peer connection.
// It processes incoming messages from peers.
func (t *Tracker) handleConnection(conn net.Conn) {
	defer conn.Close()                     // Ensure the connection is closed after the function returns
	peerAddr := conn.RemoteAddr().String() // Get the address of the connected peer

	buffer := make([]byte, 1024) // Buffer to store incoming data
	n, err := conn.Read(buffer)  // Read data from the connection
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	message := string(buffer[:n])        // Convert the buffer bytes into a string
	parts := strings.Split(message, ":") // Split the message into parts using ":" as the delimiter

	// Handle message based on its type
	if parts[0] == "REGISTER" && len(parts) == 3 {
		fileName := parts[1] // Extract the file name
		peerPort := parts[2] // Extract the peer's server port

		// Store the peer's IP address and port as a single string
		peerIP := strings.Split(peerAddr, ":")[0]
		peerInfo := peerIP + ":" + peerPort

		// Update the tracker's peers map with the new information
		t.lock.Lock()
		t.peers[peerInfo] = append(t.peers[peerInfo], fileName)
		t.lock.Unlock()

		// Log the new registration
		fmt.Println("Peer", peerInfo, "has file:", fileName)

		// Send a response back to the peer
		conn.Write([]byte("OK"))
	}

	if parts[0] == "REQUEST_FILE" && len(parts) == 2 {
		fileName := parts[1]                     // Extract the file name
		peerList := t.getPeersWithFile(fileName) // Get a list of peers that have the file
		if len(peerList) > 0 {
			rand.Seed(time.Now().UnixNano())          // Seed the random number generator
			randomIndex := rand.Intn(len(peerList))   // Select a random index from the peerList
			conn.Write([]byte(peerList[randomIndex])) // Send the randomly selected peer's info
		} else {
			conn.Write([]byte("NO_PEER")) // Send a response indicating no peer has the file
		}
	}

	if parts[0] == "EXIT" {
		// Handle peer exit
		t.lock.Lock()
		delete(t.peers, peerAddr) // Remove the peer from the tracker's map
		t.lock.Unlock()

		// Log the peer's exit
		fmt.Println("Peer", peerAddr, "has exited")
	}
}

// getPeersWithFile returns a slice of peers that have the specified file.
func (t *Tracker) getPeersWithFile(fileName string) []string {
	t.lock.Lock() // Ensure exclusive access to the peers map
	defer t.lock.Unlock()

	var peerList []string // Initialize an empty slice for peers with the file
	for peer, files := range t.peers {
		for _, f := range files {
			if f == fileName {
				peerList = append(peerList, peer) // Add the peer to the list if they have the file
				break
			}
		}
	}
	return peerList
}

// Start begins the tracker server on the specified host and port.
// It listens for incoming connections and handles them.
func (t *Tracker) Start(host string, port string) {
	listener, err := net.Listen("tcp", host+":"+port) // Start listening on the specified host and port
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close() // Ensure the listener is closed when the function returns

	// Log that the tracker is running
	fmt.Println("Tracker running on " + host + ":" + port)

	for {
		conn, err := listener.Accept() // Accept new connections
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go t.handleConnection(conn) // Handle the connection
	}
}

func main() {
	tracker := NewTracker() // Create a new instance of Tracker

	// Feel free to change the tracker IP and tracker port based on your machine
	trackerIP := "localhost"
	trackerPort := "29392"

	// Start the tracker on the local machine ("localhost") on port "20000"
	tracker.Start(trackerIP, trackerPort)
}
