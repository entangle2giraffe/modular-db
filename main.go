package main

import (
	"fmt"
	"os"
	"strconv"

	"modular-db/mydatabase"
  "github.com/google/uuid"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  \033[31mcreate <name> <age>")
	fmt.Println("  \033[32mupdate <id> <name> <age>")
	fmt.Println("  \033[33mdelete <id>")
	fmt.Println("  \033[34mget <id>")
	fmt.Println("  \033[35mlist")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	db, err := mydatabase.NewDatabaseFromEnv()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	switch command {
	case "create":
		if len(os.Args) != 4 {
			fmt.Println("Usage: create <name> <age>")
			return
		}
		name := os.Args[2]
		age, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid age")
			return
		}
		err = db.Create(name, age)
		if err != nil {
			fmt.Printf("Failed to create record: %v\n", err)
			return
		}
		fmt.Println("Record created successfully")
	case "update":
		if len(os.Args) != 5 {
			fmt.Println("Usage: update <id> <name> <age>")
			return
		}
		id, err := uuid.Parse(os.Args[2])
		if err != nil {
			fmt.Println("Invalid id")
			return
		}
		name := os.Args[3]
		age, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println("Invalid age")
			return
		}
		err = db.Update(id, name, age)
		if err != nil {
			fmt.Printf("Failed to update record: %v\n", err)
			return
		}
		fmt.Println("Record updated successfully")
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: delete <id>")
			return
		}
		id, err := uuid.Parse(os.Args[2])
		if err != nil {
			fmt.Println("Invalid id")
			return
		}
		err = db.Delete(id)
		if err != nil {
			fmt.Printf("Failed to delete record: %v\n", err)
			return
		}
		fmt.Println("Record deleted successfully")
	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: get <id>")
			return
		}
		id, err := uuid.Parse(os.Args[2])
		if err != nil {
			fmt.Println("Invalid id")
			return
		}
		record, err := db.GetById(id)
		if err != nil {
			fmt.Printf("Failed to get record: %v\n", err)
			return
		}
		fmt.Printf("Record: id=%s, name=%s, age=%d\n", record.Id.String(), record.Name, record.Age)
	case "list":
		records, err := db.GetAll()
		if err != nil {
			fmt.Printf("Failed to get records: %v\n", err)
			return
		}
		fmt.Println("Records:")
		for _, record := range records {
			fmt.Printf("  id=%s, name=%s, age=%d\n", record.Id.String(), record.Name, record.Age)
		}
	}
}
