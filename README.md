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

Our directory is a folder called "files" with the following files: 'poem1.txt', 'poem2.txt', 'poem3.txt', 'poem4.txt', 'poem5.txt', 'poem6.txt', 'poem7.txt', 'poem8.txt', 'poem9.txt', 'poem10.txt', 'plants.json'. In total 11 files. Each file has a short poem about wisdom in it. The one 'plant.json' file has a list of plants in it.

### Downloading the files

Download the tracker.go file and the peer.go file and store them in separate folders.

In addition, please download the directory before you start. Alternatively, you may create your own folder titled "files" in which to add short text files. Put this folder in the same folder as the peer.go file.

### Running the Tracker

1. **Start the tracker:**
   ```
   go run tracker.go
   ```
   - The tracker will start on `localhost` and default port `20000`.
   - Feel free to change the tracker IP and port number based on your machine.

### Running the Peer(client/server)

1. **Start the peer:**
   ```
   go run peer.go
   ```
   - Please enter an available port number for the peer server.

2. **Follow the on-screen prompts to:**
   - Connect to the tracker. (Enter the tracker IP and port).
   - Select or specify files for sharing. (For example, enter "poem1.txt" without the quotations to download poem1.txt")
   - Download files from peers. (Receive the file in chunks)

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

- **Port Number Issues**: If the port number is already in use on your peer machine you will get a "Error starting server: address already in use". To resolve this issue, rerun the peer.go file. (Similarly, rerun the tracker.go file if the tracker port number is already in use.)

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


