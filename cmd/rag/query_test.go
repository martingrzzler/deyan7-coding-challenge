package main

import (
	"testing"

	"github.com/martingrzzler/deyan7challenge/internal/persist"
)

func TestQueryOne(t *testing.T) {
	db, err := persist.Connect()
	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	q := Query{
		Type: QueryTypeOne,
		Where: []Where{
			{Field: "name", Value: "XBO 2500 W/HS XL OFR", Op: OperationEqual},
		},
		ReturnFields: []string{"produkt_gewicht_g", "name"},
	}

	result, err := QueryOne(db, q)
	if err != nil {
		t.Fatalf("could not query database: %v", err)
	}

	name := result["name"].(string)
	if name != "XBO 2500 W/HS XL OFR" {
		t.Errorf("unexpected name: %v", name)
	}

	weight := result["produkt_gewicht_g"].(float64)
	if weight != 571 {
		t.Errorf("unexpected weight: %v", weight)
	}
}

func TestQueryMany(t *testing.T) {
	db, err := persist.Connect()
	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	q := Query{
		Type: QueryTypeMany,
		Where: []Where{
			{Field: "lebensdauer_h", Op: OperationGT, Value: 3000},
			{Field: "nenn_leistung_w", Op: OperationGTE, Value: 1500},
		},
		ReturnFields: []string{"name", "lebensdauer_h", "nenn_leistung_w"},
	}

	results, err := QueryMany(db, q)
	if err != nil {
		t.Fatalf("could not query database: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("unexpected number of results: %v", len(results))
	}
}
