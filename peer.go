// peer.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const ChunkSize = 1024 // Size of each file chunk in bytes... not fully implemented

type P2PPeer struct {
	peers          []net.Conn // Slice of network connections to other peers
	availableFiles []string   // List of available files for sharing
}

// NewP2PPeer creates and returns a new P2PPeer instance
func NewP2PPeer() *P2PPeer {
	return &P2PPeer{
		peers:          make([]net.Conn, 0),
		availableFiles: make([]string, 0),
	}
}

// pickRandomFile selects a random file from a given directory
func (c *P2PPeer) pickRandomFile(directory string) (string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	// Add file to array
	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	if len(fileNames) == 0 {
		return "", fmt.Errorf("No files found in directory")
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select a random file
	selectedFile := fileNames[rand.Intn(len(fileNames))]

	// Return randomly selected file and no error
	return selectedFile, nil
}

// startPeerServer starts a TCP server to listen for incoming connections
// This is the server aspect of the peer
func (c *P2PPeer) startPeerServer(port string) string {
	// Create a server on user input port number
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting my server:", err.Error())
		return "error"
	}
	defer listener.Close()

	fmt.Println("My server is running on port " + port)

	// Server listening for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading request:", err.Error())
				return
			}
			// Parse incoming message
			request := string(buffer[:n])
			parts := strings.Split(request, ":")

			// Read request for chunks
			if len(parts) == 3 && parts[0] == "GET_CHUNK" {
				fileName := parts[1]
				chunkIndex, _ := strconv.Atoi(parts[2])
				// Send over file chunk
				c.serveFileChunk(conn, fileName, chunkIndex)
			}
		}()
	}
}

// connectToTracker connects to a tracker server and registers the peer
func (c *P2PPeer) connectToTracker(trackerHost string, trackerPort string, fileName string, myServerPort string) {
	// Start a TCP connection with tracker
	conn, err := net.Dial("tcp", trackerHost+":"+trackerPort)
	if err != nil {
		fmt.Println("Error connecting to tracker:", err.Error())
		return
	}
	defer conn.Close()

	// Register to join the network
	infoMessage := fmt.Sprintf("REGISTER:%s:%s", fileName, myServerPort)
	_, err = conn.Write([]byte(infoMessage))
	if err != nil {
		fmt.Println("Error sending message to tracker:", err.Error())
		return
	}

	// Buffer for incoming data
	buffer := make([]byte, 1024)

	// Read data into buffer
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving message from tracker:", err.Error())
	}

	// Convert the received bytes to a string
	receivedMsg := string(buffer[:n])

	// Check if the received message is "OK"
	if receivedMsg == "OK" {
		fmt.Println("Received 'OK' from tracker")
		fmt.Println("Sucessfully registered with tracker")
	} else {
		fmt.Println("Received different message:", receivedMsg)
	}
}

// requestFileFromTracker asks the tracker for peers who have a specific file
func (c *P2PPeer) requestFileFromTracker(trackerHost string, trackerPort string, fileName string) string {
	// Start TCP connection with tracker
	conn, err := net.Dial("tcp", trackerHost+":"+trackerPort)
	if err != nil {
		fmt.Println("Error connecting to tracker:", err.Error())
		return ""
	}
	defer conn.Close()

	// Send a file request message to tracker
	requestMessage := fmt.Sprintf("REQUEST_FILE:%s", fileName)
	_, err = conn.Write([]byte(requestMessage))
	if err != nil {
		fmt.Println("Error sending request to tracker:", err.Error())
		return ""
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from tracker:", err.Error())
		return ""
	}

	// Parse response from tracker
	response := string(buffer[:n])
	if response == "NO_PEER" {
		fmt.Println("No peer has the requested file")
		return ""
	}

	// Return information for a peer that contains the requested file
	return response
}

// connectToPeer establishes a connection with another peer
func (c *P2PPeer) connectToPeer(peerHost string, peerPort string) {
	conn, err := net.Dial("tcp", peerHost+":"+peerPort)
	if err != nil {
		fmt.Println("Error connecting to peer:", err.Error())
		return
	}
	c.peers = append(c.peers, conn)
	fmt.Println("Connected to peer at " + peerHost + ":" + peerPort)
}

// downloadFile handles the downloading of a file from a peer
func (c *P2PPeer) downloadFile(peerConn net.Conn, fileName string, totalChunks int) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}
	defer outFile.Close()

	// For each chunk send a get chunk request message to the uploading peer
	for i := 0; i < totalChunks; i++ {
		requestMessage := fmt.Sprintf("GET_CHUNK:%s:%d", fileName, i)
		// Send request message
		_, err = peerConn.Write([]byte(requestMessage))
		if err != nil {
			fmt.Println("Error sending chunk request:", err.Error())
			return
		}

		// Parse the received chunk
		// This may present an error if the buffer is not large enough for the entire chunk
		buffer := make([]byte, ChunkSize)
		n, err := peerConn.Read(buffer)

		if err != nil {
			fmt.Println("Error receiving chunk:", err.Error())
			return
		}

		if n == 0 {
			fmt.Println("Received empty chunk")
			continue
		}

		// Write chunk to the file
		bytesWritten, err := outFile.Write(buffer[:n])

		if err != nil {
			fmt.Println("Error writing chunk to file:", err.Error())
			return
		}
		fmt.Printf("Chunk %d written, %d bytes\n", i, bytesWritten)
	}

	outFile.Sync() // Flush the file buffer to disk
	fmt.Println("Download complete for file:", fileName)
}

// serveFileChunk sends a requested chunk of a file to another peer
func (c *P2PPeer) serveFileChunk(conn net.Conn, fileName string, chunkIndex int) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		return
	}
	defer file.Close()

	// Seek to the start of the chunk
	_, err = file.Seek(int64(chunkIndex*ChunkSize), 0)
	if err != nil {
		fmt.Println("Error seeking file:", err.Error())
		return
	}

	// Read the chunk
	buffer := make([]byte, ChunkSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file chunk:", err.Error())
		return
	}

	// Send the chunk
	_, err = conn.Write(buffer[:n])
	if err != nil {
		fmt.Println("Error sending file chunk:", err.Error())
		return
	}

	// Logging message
	fmt.Println("Served chunk", chunkIndex, "of file", fileName)
}

func main() {

	// Creating a new P2P peer
	peer := NewP2PPeer()

	reader := bufio.NewReader(os.Stdin) // User input

	// Input user peer port
	fmt.Print("Enter my server port: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	// Start the peer server in a separate goroutine
	go peer.startPeerServer(port)

	time.Sleep(1 * time.Second) // Delay so that messages will not overlap

	fmt.Print("Enter tracker IP: ") // Prompt for tracker IP
	trackerHost, _ := reader.ReadString('\n')
	trackerHost = strings.TrimSpace(trackerHost)

	fmt.Print("Enter tracker port: ") // Prompt for tracker port
	trackerPort, _ := reader.ReadString('\n')
	trackerPort = strings.TrimSpace(trackerPort)

	// Pick a random file from a specified directory
	selectedFile, err := peer.pickRandomFile("files")
	if err != nil {
		fmt.Println("Error selecting my file:", err.Error())
		return
	}

	fmt.Println("My initial selected file:", selectedFile)
	fileName := selectedFile

	// Initiate connection to tracker
	peer.connectToTracker(trackerHost, trackerPort, fileName, port)

	// Loop to request files
	for {
		// Prompt for file request
		fmt.Print("Enter the name of the file you want to request (or type 'EXIT' to quit): ")
		requestedFile, _ := reader.ReadString('\n')
		requestedFile = strings.TrimSpace(requestedFile)

		// Check if the user wants to exit the loop
		if strings.ToUpper(requestedFile) == "EXIT" {
			fmt.Println("Sending exit message to tracker and exiting file request loop.")
			conn, err := net.Dial("tcp", trackerHost+":"+trackerPort)
			if err != nil {
				fmt.Println("Error connecting to tracker:", err.Error())
			} else {
				requestMessage := "EXIT"                    // Exit the network
				_, err = conn.Write([]byte(requestMessage)) // Inform the tracker that peer is leaving
				if err != nil {
					fmt.Println("Error sending exit message to tracker:", err.Error())
				}
				conn.Close()
			}
			break
		}

		fileName := requestedFile
		totalChunks := 10 // Number of chunks for each file... Potentially revise

		// Message tracker for information about the peer who possesses the file
		peerInfo := peer.requestFileFromTracker(trackerHost, trackerPort, fileName)
		fmt.Println("Here is the information for the peer who has the file you are requesting:", peerInfo)
		if peerInfo != "" {
			hostPort := strings.Split(peerInfo, ":")
			// Connect to peer with the file
			peerConn, err := net.Dial("tcp", hostPort[0]+":"+hostPort[1])
			if err != nil {
				fmt.Println("Error connecting to peer:", err.Error())
				continue
			}
			// Download the file from the peer
			peer.downloadFile(peerConn, fileName, totalChunks)
			peerConn.Close()
		} else {
			// Most likely file does not exist in the network
			fmt.Println("No peer information available for the requested file.")
		}
	}
}
