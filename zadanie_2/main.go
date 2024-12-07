package main

import (
	"math/big"
	"math/rand/v2"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func generatePrime() int {
	p := 0
	for true {
		p = rand.IntN(10000000000-1000000000) + 1000000000
		if big.NewInt(int64(p)).ProbablyPrime(50) {
			break
		}
	}
	return p
}

func isPrimitiveRoot(g, p int) bool {
	phi := p - 1
	factors := primeFactors(phi)
	for _, factor := range factors {
		if new(big.Int).Exp(big.NewInt(int64(g)), big.NewInt(int64(phi/factor)), big.NewInt(int64(p))).Cmp(big.NewInt(1)) == 0 {
			return false
		}
	}
	return true
}

func primeFactors(n int) []int {
	factors := []int{}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func generatePrimitiveRoot(p int) int {
	for {
		g := rand.IntN(p-2000000) + 2000000
		if isPrimitiveRoot(g, p) {
			return g
		}
	}
}

func diffieHellman(p, g, id int, send, receive chan int) {
	var A int = 0
	a := rand.IntN(1000000-100000) + 100000
	A = (g ^ a) % p
	println("generated public key for goroutine no. ", id, ": ", A)
	println("generated private key for goroutine no. ", id, ": ", a)
	println("sending public key...")
	time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
	send <- A
	println("receiving public key...")
	B := <-receive
	println("received public key: ", B)
	K := (B ^ a) % p
	println("shared secret key for goroutine no. ", id, ": ", K)
}

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	//go diffie_hellman()
	p := generatePrime()
	g := generatePrimitiveRoot(p)
	println("p: ", p, " g: ", g)
	id := 1

	go diffieHellman(p, g, id, ch1, ch2)
	id += 1
	go diffieHellman(p, g, id, ch2, ch1)

	time.Sleep(10 * time.Second)

	close(ch1)
	close(ch2)
}
