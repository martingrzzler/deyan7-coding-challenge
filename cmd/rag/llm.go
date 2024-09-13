package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIClient struct {
	Model  string
	APIKey string
	URL    string
}

func NewGPT4OMiniClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		Model:  "gpt-4o-mini",
		APIKey: apiKey,
		URL:    "https://api.openai.com/v1/chat/completions",
	}
}

func (c *OpenAIClient) AnswerQuestion(q string, dbResult string) (string, error) {
	opts := OpenAIRequestOptions{
		Model:          c.Model,
		ResponseFormat: OpenAIResponseFormat{Type: "text"},
		Messages: []OpenAIMessage{
			{Role: "system", Content: dbRes2AnswerBasePrompt},
			{Role: "user", Content: fmt.Sprintf("User question: %s\nDB Result: %s\nAnswer:", q, dbResult)},
		},
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return "", fmt.Errorf("could not marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, b)
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) QuestionToQuery(q string) (Query, error) {
	opts := OpenAIRequestOptions{
		Model:          c.Model,
		ResponseFormat: OpenAIResponseFormat{Type: "json_object"},
		Messages: []OpenAIMessage{
			{Role: "system", Content: question2QueryBasePrompt},
			{Role: "user", Content: fmt.Sprintf("User question: %s", q)},
		},
	}

	body, err := json.Marshal(opts)
	if err != nil {
		return Query{}, fmt.Errorf("could not marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(body))
	if err != nil {
		return Query{}, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Query{}, fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return Query{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, b)
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return Query{}, fmt.Errorf("could not decode response: %w", err)
	}

	var m map[string]Query
	err = json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &m)
	if err != nil {
		return Query{}, fmt.Errorf("could not unmarshal query: %w", err)
	}

	return m["query"], nil
}

type OpenAIRequestOptions struct {
	Model          string               `json:"model"`
	ResponseFormat OpenAIResponseFormat `json:"response_format"`
	Messages       []OpenAIMessage      `json:"messages"`
}

type OpenAIResponseFormat struct {
	Type string `json:"type"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

const dbRes2AnswerBasePrompt string = `
You are a helpful assistent that maps a database result returned as JSON to an answer of a previous user question.
Your answer should be in plain text.
`

const question2QueryBasePrompt string = `
You are a helpful assistent that maps a user querstion to a database query.
The api accepts a json object with the following structure:
{
  "query": {
    "type": "many" | "one",
    "where": [
      {
        "field": string,
        "value": float | string,
        "op": "eq" | "gte" | "lte" | "gt" | "lt"
      }
    ],
    "return_fields": [string]
  }
}
You can choose between the following fields of a bulb products within the db.
name VARCHAR,
 anwendungs_gebiete JSONB,
 -- list of strings
 vorteile JSONB,
 -- list of strings
 eigenschaften JSONB,
 -- list of strings
 nenn_strom_a INTEGER,
 -- in Ampere
 strom_steuer_a_min INTEGER,
 -- in Ampere
 stroem_steuer_a_max INTEGER,
 -- in Ampere
 nenn_leistung_w INTEGER,
 -- in Watt
 nenn_spannung_v INTEGER,
 -- in Volt
 durchmesser_mm INTEGER,
 -- in mm
 laenge_mm INTEGER,
 -- in mm
 laenge_mit_sockel_mm INTEGER,
 -- Länge mit Sockel jedoch ohne Sockelstift
 lcl_mm INTEGER,
 -- Abstand Lichtschwerpunkt (LCL)
 kabel_laenge_mm INTEGER,
 -- Kabel-/Leitungslänge, Eingangsseite
 elekroden_abstand_mm INTEGER,
 -- Elektrodenabstand kalt
 produkt_gewicht_g INTEGER,
 -- in Gramm
 max_umgebungsgtemperatur_c INTEGER,
 -- in Grad Celsius
 lebensdauer_h INTEGER,
 -- in Stunden
 sockel_anode VARCHAR,
 -- Socker Anode (Normbezeichnung)
 sockel_kathode VARCHAR,
 -- Sockel Kathode (Normbezeichnung)
 kuehlung VARCHAR,
 -- Kühlung enum
 brennstellung VARCHAR,
 -- Brennstellung
 deklarations_datum DATE,
 -- Datum der Deklaration
 erzeugniss_nummern JSONB,
 -- Primäre Erzeugnisnummer, can be multiple split by |
 stoff VARCHAR,
 -- Stoff der Kandidatenliste
 stoff_cas_nummer VARCHAR,
 -- CAS-Nummer des Stoffes
 scip_nummern JSONB,
 -- SCIP Deklarationsnummer, can be multiple split by |
 ean VARCHAR,
 -- EAN
 metel_code VARCHAR,
 -- METEL-Code
 seg_no VARCHAR,
 -- SEG-No.
 stk_nummer VARCHAR,
 -- STK-Nummer
 uk_org VARCHAR,
 -- UK-Org.

 Make sure to always include the field "name" in the return_fields.
`
