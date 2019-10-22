package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var globalRandString string

func init() {
	rand.Seed(time.Now().UnixNano())
	globalRandString = randString(16)
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	r := &resolver{cache: make(map[string]bool)}

	for sc.Scan() {
		name := strings.ToLower(sc.Text())

		if !r.containsWildcard(name) {
			fmt.Println(name)
		}
	}
}

type resolver struct {
	cache map[string]bool
}

func (r *resolver) isWildcard(name string) bool {
	if r.cache[name] {
		return true
	}

	check := fmt.Sprintf("%s.%s", globalRandString, name)
	_, err := net.LookupHost(check)
	r.cache[name] = err == nil
	return err == nil
}

func (r *resolver) containsWildcard(name string) bool {

	parts := strings.Split(name, ".")

	// Given one.two.target.com
	// Start from the right hand side, check:
	//   $rand.target.com
	//   $rand.two.target.com
	// Do not check $rand.one.two.target.com otherwise
	// we'd just be resolving every line of input which
	// is what we're trying to avoid doing.
	for i := len(parts) - 2; i > 0; i-- {
		candidate := strings.Join(parts[i:len(parts)], ".")
		if r.isWildcard(candidate) {
			return true
		}
	}

	return false
}

func randString(length int) string {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	out := bytes.Buffer{}

	for i := 0; i < length; i++ {
		out.WriteByte(chars[rand.Intn(len(chars))])
	}

	return out.String()
}
