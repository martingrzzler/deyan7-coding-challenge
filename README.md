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
python3 -m venv venv
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
go run ./cmd/rag -question "Welche Leuchte hat SCIP Nummer dd2ddf15-037b-4473-8156-97498e721fb3?
```
