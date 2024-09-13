# Deyan7 Coding Challenge

## Setup

### Requirements
- Python
- Docker
- Golang

### Running the project

1. Recreate the dataset (optional) `data.jsonl`
```bash
export OPENAI_API_KEY=<your_openai_api_key>
python3 -m venv ./cmd/preprocess/venv
source ./cmd/preprocess/venv/bin/activate
pip install -r ./cmd/preprocess/requirements.txt
python ./cmd/preprocess/main.py
```

2. Start a docker postgres instance
```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d --rm postgres
```

3. Migrate the data (`dataset.jsonl`)
```bash
go run ./cmd/insert/main.go -input ./dataset.jsonl -migrate
```

4. Ask a question
```bash
export OPENAI_API_KEY=<your_openai_api_key>
go run ./cmd/rag -question "Welche Leuchte hat SCIP Nummer dd2ddf15-037b-4473-8156-97498e721fb3?"
go run ./cmd/rag -question "Gebe mir alle Leuchtmittel mit mindestens 1500W und einer Lebensdauer von mehr als 3000 Stunden?"
go run ./cmd/rag -question "Welche Leuchte hat die Erzeugnissnummer 4008321299963?"
go run ./cmd/rag -question "Wie viel wiegt XBO 4000 W/HS XL OFR?"
go run ./cmd/rag -question "Welche Leuchten haben einen Durchmesser kleiner als 50mm und welchen Durchmesser haben sie?"
```
Add -debug to see more information about the query.

## Limitations

This RAG pipeline relies on a api that queries the database. It only supports queries of one or more rows filtered by a WHERE clause (combined with AND if multiple).
Therefore questions like "Welche Leuchte hat die höchste Lebensdauer?" can't be answered yet. This feature can be added by extending the api.
