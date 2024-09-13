package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/martingrzzler/deyan7challenge/internal/persist"
)

var openAIAPIKey = os.Getenv("OPENAI_API_KEY")

func main() {
	if openAIAPIKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is required")
		os.Exit(1)
	}

	userQuestion := flag.String("question", "", "The question to answer about the lightbulbs")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *userQuestion == "" {
		fmt.Println("question flag is required")
		os.Exit(1)
	}

	// userQuestion := "Gebe mir alle Leuchtmittel mit mindestens 1500W und einer Lebensdauer von mehr als 3000 Stunden?"

	fmt.Println("Question:", *userQuestion)
	fmt.Println("Thinking...")

	openaiClient := NewGPT4OMiniClient(openAIAPIKey)
	db, err := persist.Connect()
	if err != nil {
		fmt.Println(fmt.Errorf("could not connect to database: %w", err))
		os.Exit(1)
	}

	q, err := openaiClient.QuestionToQuery(*userQuestion)
	if err != nil {
		fmt.Println(fmt.Errorf("could not convert question to query: %w", err))
		os.Exit(1)
	}

	if *debug {
		fmt.Println("Query:")
		d, _ := json.MarshalIndent(q, "", "  ")
		if err == nil {
			fmt.Println(string(d))
		}
	}

	var dbRes string
	switch q.Type {
	case "one":
		result, err := QueryOne(db, q)
		if err != nil {
			fmt.Println(fmt.Errorf("could not query database: %w", err))
			os.Exit(1)
		}

		data, err := json.Marshal(result)
		if err != nil {
			fmt.Println(fmt.Errorf("could not marshal database result: %w", err))
			os.Exit(1)
		}

		dbRes = string(data)
	case "many":
		results, err := QueryMany(db, q)
		if err != nil {
			fmt.Println(fmt.Errorf("could not query database: %w", err))
			os.Exit(1)
		}

		data, err := json.Marshal(results)
		if err != nil {
			fmt.Println(fmt.Errorf("could not marshal database results: %w", err))
			os.Exit(1)
		}

		dbRes = string(data)
	default:
		fmt.Println(fmt.Errorf("unexpected query type: %s", q.Type))
		os.Exit(1)
	}

	if *debug {
		fmt.Println("Database Result:")
		fmt.Println(dbRes)
	}

	answer, err := openaiClient.AnswerQuestion(*userQuestion, dbRes)
	if err != nil {
		fmt.Println(fmt.Errorf("could not answer question: %w", err))
		os.Exit(1)
	}

	fmt.Println("Answer:")
	fmt.Println(answer)
}
