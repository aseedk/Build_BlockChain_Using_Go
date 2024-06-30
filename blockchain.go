package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	fromAddressField = "from"
	toAddressField   = "to"
	amountField      = "amount"
)

// Blockchain struct which represents a single blockchain
type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// CreateBlockchain is a function which creates a new blockchain with a genesis block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// addBlock is a method of the Blockchain struct which adds a new block to the blockchain
func (b *Blockchain) addBlock(from, to string, amount float64) {
	// Initialize the data of the block with the from, to and amount fields
	blockData := map[string]interface{}{
		fromAddressField: from,
		toAddressField:   to,
		amountField:      amount,
	}

	// Get the last block in the blockchain
	lastBlock := b.chain[len(b.chain)-1]

	// Create a new block with the data, previous hash and timestamp
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}

	// Mine the new block
	newBlock.mine(b.difficulty)

	// Append the new block to the blockchain
	b.chain = append(b.chain, newBlock)
}

// isValid is a method of the Blockchain struct which checks if the blockchain is valid
func (b *Blockchain) isValid() bool {
	// Initialize the previous and current block variables
	var previousBlock, currentBlock Block

	// Iterate over the blocks in the blockchain and check if the hash of the block is valid
	for i := range b.chain[1:] {
		// Get the previous and current block in the blockchain
		previousBlock = b.chain[i]
		currentBlock = b.chain[i+1]

		// Check if the hash of the block is valid by comparing the hash of the block with the calculated hash
		// and the previous hash of the block with the hash of the previous block
		// If the hash of the block is not valid, return false
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}

	// If all the blocks in the blockchain are valid, return true
	return true
}

// Block struct which represents a single block in the blockchain
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

// calculateHash is a method of the Block struct which calculates the hash of the block
func (b *Block) calculateHash() string {
	// convert the data to a byte array
	data, err := json.Marshal(b.data)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// calculate the hash of the block by concatenating the previous hash, data, timestamp and proof of work
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)

	// Generate the hash of the block using SHA256
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine is a method of the Block struct which mines the block by calculating the hash of the block
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// main function
func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
	for _, block := range blockchain.chain {
		fmt.Println("--------------------")
		fmt.Println("Data:", block.data)
		fmt.Println("Hash:", block.hash)
		fmt.Println("Previous Hash:", block.previousHash)
		fmt.Println("Timestamp:", block.timestamp)
		fmt.Println("Proof of Work:", block.pow)
	}
}
