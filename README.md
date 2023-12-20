# p2p-network

# P2P File Sharing Application

This P2P file sharing application enables users to share and download files in a distributed network. It's built using Go and utilizes a simple tracker to facilitate peer discovery and file transfer in chunks.

## Features

- **File Chunking**: Efficient file sharing by breaking down files into manageable chunks.
- **Tracker-based Peer Discovery**: Utilizes a tracker for managing and discovering peers with desired files.
- **Concurrent Chunk Download**: Supports downloading different chunks of the same file concurrently from multiple peers.
- **Basic Error Handling**: Handles common network errors and file I/O issues.

## Getting Started

### Prerequisites

- Go (version 1.13 or later)
- Network access to connect with peers and the tracker

### Directory explanation

Our directory is a folder called "TestingFiles" with the following files: 'poem1.txt', 'poem2.txt', 'poem3.txt', 'poem4.txt', 'poem5.txt', 'poem6.txt', 'poem7.txt', 'poem8.txt', 'poem9.txt', 'poem10.txt','plants.json'. Total 11 files. Each file has a short poem about wisdom in it. The one 'plant.json' file has a list of plants in it.


### Running the Tracker

1. **Start the tracker:**
   ```
   go run tracker.go
   ```
   - The tracker will start on `localhost` and default port `8000`.

### Running the Peer(client/server)

1. **Start the peer:**
   ```
   go run peer.go
   ```

2. **Follow the on-screen prompts to:**
   - Enter your server port address.
   - Connect to the tracker.
   - Select or specify files for sharing.
   - Download files from peers.

## Usage Example

1. **Start the tracker:**
   ```
   go run tracker.go
   ```

2. **On a different terminal, start the peer:**
   ```
   go run peer.go
   ```

3. **Interact with the peer through the CLI to share and download files.**


## Authors

Malahim Tariq, Claire Cao, Ariel Moncrief 

Project Link: [https://github.com/cc459/p2p-network.git]


## Debugging

- **peer connection issues**: If the port is already in use you will get a "Error starting server: address already in use" rerun run the peer.go file.


### Common Issues
- **Connection Issues**: Check if the tracker and peers are accessible over the network.
- **File Handling**: Ensure file paths and permissions are correctly set for reading and writing files.
- **Concurrency**: Look out for issues related to concurrent file access and data races.
- **Input**: Please enter the file name carefully. We don't have cases handling giving error messages for incorrect input. 


### Integration Testing
- Test the application components working together - such as the interaction between the peer and the tracker.
- Simulate different network conditions to ensure the application remains stable and efficient.

### End-to-End Testing
- Test the complete workflow of the application from starting the tracker, connecting peers, to sharing and downloading files.


