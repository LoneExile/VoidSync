package db

import (
	"log"
	"voidsync/config"

	"github.com/surrealdb/surrealdb.go"
)

type User struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func Init(cfg *config.Config) {
	db, err := surrealdb.New(cfg.SurrealDBEndpoint)
	if err != nil {
		panic(err)
	}
	log.Println("🟢 Connected to SurrealDB")

	// Sign in
	if _, err = db.Signin(map[string]string{
		"user": "root",
		"pass": "root",
	}); err != nil {
		panic(err)
	}
	log.Println("🟢 Signed in to SurrealDB")

	// Select namespace and database
	if _, err = db.Use("test", "test"); err != nil {
		panic(err)
	}
	log.Println("🟢 Selected namespace and database")

	user := map[string]string{
		"name":    "John",
		"surname": "Doe",
	}

	// Insert user
	data, err := db.Create("user", user)
	if err != nil {
		panic(err)
	}

	// Unmarshal data
	createdUser := make([]User, 1)
	err = surrealdb.Unmarshal(data, &createdUser)
	if err != nil {
		panic(err)
	}

	// Get user by ID
	data, err = db.Select(createdUser[0].ID)
	if err != nil {
		panic(err)
	}
	log.Println("🟢 Selected user:", data)

	// Unmarshal data
	selectedUser := new(User)
	err = surrealdb.Unmarshal(data, &selectedUser)
	if err != nil {
		panic(err)
	}

	// Change part/parts of user
	changes := map[string]string{"name": "Jane"}
	if _, err = db.Change(selectedUser.ID, changes); err != nil {
		panic(err)
	}

	// Update user
	if _, err = db.Update(selectedUser.ID, changes); err != nil {
		panic(err)
	}

	if _, err = db.Query("SELECT * FROM $record", map[string]interface{}{
		"record": createdUser[0].ID,
	}); err != nil {
		panic(err)
	}

	// Delete user by ID
	if _, err = db.Delete(selectedUser.ID); err != nil {
		panic(err)
	}
}
