package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, firstName, lastName, password string) *Account {
	acc, err := NewAccount(firstName, lastName, password)
	if err != nil {
		log.Fatal(err)
	}

	id, err := store.CreateAccount(acc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New account ID =>", id)

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Anthony", "Neinn", "hunter7979")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database...")
		seedAccounts(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
