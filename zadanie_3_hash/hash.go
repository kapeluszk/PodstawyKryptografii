package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"os"
	"time"
)

// Generate a random string
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}

// Hash functions
func hashMD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func hashSHA1(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func hashSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func hashSHA512(input string) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func hashSHA3_256(input string) string {
	hash := sha3.New256()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func hashSHA3_512(input string) string {
	hash := sha3.New512()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

// Extract the first 12 bits of the hash
func getFirst12Bits(hash string) string {
	decoded, _ := hex.DecodeString(hash)
	firstByte := decoded[0]
	secondByte := decoded[1]

	// Get the first 4 bits of the second byte
	first12Bits := (int(firstByte) << 4) | (int(secondByte) >> 4)
	return fmt.Sprintf("%03x", first12Bits) // Return as 3 hexadecimal characters
}

// Compare the speed and length of hash functions
func compareHashes(input string) {
	// Define the list of hash functions
	hashes := []struct {
		name     string
		hashFunc func(string) string
	}{
		{"MD5", hashMD5},
		{"SHA-1", hashSHA1},
		{"SHA-256", hashSHA256},
		{"SHA-512", hashSHA512},
		{"SHA3-256", hashSHA3_256},
		{"SHA3-512", hashSHA3_512},
	}

	// Iterate over each hash function and display results
	for _, hash := range hashes {
		// Generate the hash value and measure the time taken
		var hashValue string
		time_table := make([]time.Duration, 0)
		for _ = range 1000000 {
			// Measure the time taken to generate the hash
			start := time.Now()
			hashValue = hash.hashFunc(input)
			duration := time.Since(start)
			time_table = append(time_table, duration)
		}
		// Calculate the average time taken
		var sum time.Duration
		for _, duration := range time_table {
			sum += duration
		}
		duration := sum / time.Duration(len(time_table))

		// Print the results in a readable format
		fmt.Printf("Hash Function: %-10s\n", hash.name)
		fmt.Printf("Hash Value   : %s\n", hashValue)
		fmt.Printf("Average Time Taken   : %v\n", duration)
		fmt.Println("---------------------------------------------------")
	}
}

// Find collisions on the first 12 bits
func findCollision() {
	seen := make(map[string]string) // Map to store first 12 bits and corresponding input
	for {
		input := generateRandomString(10) // Generate random input
		hash := hashSHA256(input)         // Hash the input using SHA-256
		first12Bits := getFirst12Bits(hash)

		if existingInput, found := seen[first12Bits]; found {
			// Collision found
			fmt.Println("Collision found!")
			fmt.Printf("Input 1: %s, Input 2: %s\n", existingInput, input)
			fmt.Printf("First 12 bits: %s\n", first12Bits)
			return
		}
		seen[first12Bits] = input // Store the first 12 bits and input
	}
}

func main() {
	// Part 1: Compare the speed and length of different hash functions
	fmt.Println("Enter a string to hash:")
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	if scanner.Scan() {
		input = scanner.Text() // Get the full input string
		fmt.Println("You entered:", input)
	} else {
		fmt.Println("Error reading input:", scanner.Err())
	}
	fmt.Println("Comparing hash functions on the input string:", input)
	compareHashes(input)

	// Part 2: Find collisions on the first 12 bits of SHA-256
	fmt.Println("Finding collisions on the first 12 bits of SHA-256...")
	findCollision()
}
