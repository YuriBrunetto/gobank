package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("yuro", "brunatko", "hunter")
	// se o erro for nulo, passa o teste
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}
