package main

import (
	"fmt"
	"github.com/heather92115/translator/internal/database"
)

func main() {
	fmt.Println("Starting the fixer")

	err := database.CreatePool()
	if err != nil {
		fmt.Printf("Failed DB connections, %v\n", err)
		return
	}

	vocab, err := database.FindVocabByID(29919)
	if err != nil {
		fmt.Printf("Error looking for vocab, %v\n\n", err)
		return
	}

	fmt.Printf("Found a vocab, %v\n\n", vocab)

}
