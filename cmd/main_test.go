package main

import (
	"testing"
)

func TestMain(m *testing.T) {
	// просто проверка, что main не падает
	go main()
}
