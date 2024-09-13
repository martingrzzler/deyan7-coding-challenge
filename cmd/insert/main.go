package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/martingrzzler/deyan7challenge/internal/persist"
)

func main() {
	inputFile := flag.String("input", "", "input jsonl file with products")
	migrate := flag.Bool("migrate", false, "migrate database")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input file")
		os.Exit(1)
	}

	in, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println(fmt.Errorf("could not open file: %w", err))
		os.Exit(1)
	}
	reader := bufio.NewReader(in)

	db, err := persist.Connect()
	if err != nil {
		fmt.Println(fmt.Errorf("could not connect to database: %w", err))
		os.Exit(1)
	}
	defer db.Close()

	if *migrate {
		if err := Migrate(db); err != nil {
			fmt.Println(fmt.Errorf("could not migrate database: %w", err))
			os.Exit(1)
		}
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(fmt.Errorf("could not read line: %w", err))
			os.Exit(1)
		}

		var product Product
		if err := json.Unmarshal(line, &product); err != nil {
			fmt.Println(fmt.Errorf("could not unmarshal line: %w", err))
			os.Exit(1)
		}

		if err := InsertProduct(db, product); err != nil {
			fmt.Println(fmt.Errorf("could not insert product: %w", err))
			os.Exit(1)
		}
	}
}

func InsertProduct(db *sql.DB, product Product) error {
	_, err := db.Exec(`
INSERT INTO product_data (
  name,
  anwendungs_gebiete,
  vorteile,
  eigenschaften,
  nenn_strom_a,
  strom_steuer_a_min,
  stroem_steuer_a_max,
  nenn_leistung_w,
  nenn_spannung_v,
  durchmesser_mm,
  laenge_mm,
  laenge_mit_sockel_mm,
  lcl_mm,
  kabel_laenge_mm,
  elekroden_abstand_mm,
  produkt_gewicht_g,
  max_umgebungsgtemperatur_c,
  lebensdauer_h,
  sockel_anode,
  sockel_kathode,
  kuehlung,
  brennstellung,
  deklarations_datum,
  erzeugniss_nummern,
  stoff,
  stoff_cas_nummer,
  scip_nummern,
  ean,
  metel_code,
  seg_no,
  stk_nummer,
  uk_org
  ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32
  )`,
		product.Name,
		product.AnwendungsGebiete,
		product.Vorteile,
		product.Eigenschaften,
		product.NennStromA,
		product.StromSteuerAMin,
		product.StromSteuerAMax,
		product.NennLeistungW,
		product.NennSpannungV,
		product.DurchmesserMM,
		product.LaengeMM,
		product.LaengeMitSockelMM,
		product.LCLMM,
		product.KabelLaengeMM,
		product.ElekrodenAbstandMM,
		product.ProduktGewichtG,
		product.MaxUmgebungstemperaturC,
		product.LebensdauerH,
		product.SockelAnode,
		product.SockelKathode,
		product.Kuehlung,
		product.Brennstellung,
		product.DeklarationsDatum,
		product.ErzeugnissNummern,
		product.Stoff,
		product.StoffCasNummer,
		product.ScipNummern,
		product.EAN,
		product.MetelCode,
		product.SegNo,
		product.StkNummer,
		product.UkOrg,
	)

	return err
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
DROP TABLE IF EXISTS product_data;
CREATE TABLE
  product_data (
    name VARCHAR,
    anwendungs_gebiete JSONB,
    -- list of strings
    vorteile JSONB,
    -- list of strings
    eigenschaften JSONB,
    -- list of strings
    nenn_strom_a REAL,
    -- in Ampere
    strom_steuer_a_min REAL,
    -- in Ampere
    stroem_steuer_a_max REAL,
    -- in Ampere
    nenn_leistung_w REAL,
    -- in Watt
    nenn_spannung_v REAL,
    -- in Volt
    durchmesser_mm REAL,
    -- in mm
    laenge_mm REAL,
    -- in mm
    laenge_mit_sockel_mm REAL,
    -- Länge mit Sockel jedoch ohne Sockelstift
    lcl_mm REAL,
    -- Abstand Lichtschwerpunkt (LCL)
    kabel_laenge_mm REAL,
    -- Kabel-/Leitungslänge, Eingangsseite
    elekroden_abstand_mm REAL,
    -- Elektrodenabstand kalt
    produkt_gewicht_g REAL,
    -- in Gramm
    max_umgebungsgtemperatur_c REAL,
    -- in Grad Celsius
    lebensdauer_h REAL,
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
    uk_org VARCHAR
    -- UK-Org.
  );
`)

	if err != nil {
		return fmt.Errorf("could not create table: %w", err)
	}

	return nil
}

type Product struct {
	Name                    string   `json:"name"`
	AnwendungsGebiete       []string `json:"anwendungs_gebiete"`
	Vorteile                []string `json:"vorteile"`
	Eigenschaften           []string `json:"eigenschaften"`
	NennStromA              float32  `json:"nenn_strom_a"`
	StromSteuerAMin         float32  `json:"strom_steuer_a_min"`
	StromSteuerAMax         float32  `json:"strom_steuer_a_max"`
	NennLeistungW           float32  `json:"nenn_leistung_w"`
	NennSpannungV           float32  `json:"nenn_spannung_v"`
	DurchmesserMM           float32  `json:"durchmesser_mm"`
	LaengeMM                float32  `json:"laenge_mm"`
	LaengeMitSockelMM       float32  `json:"laenge_mit_sockel_mm"`
	LCLMM                   float32  `json:"lcl_mm"`
	KabelLaengeMM           float32  `json:"kabel_laenge_mm"`
	ElekrodenAbstandMM      float32  `json:"elekroden_abstand_mm"`
	ProduktGewichtG         float32  `json:"produkt_gewicht_g"`
	MaxUmgebungstemperaturC float32  `json:"max_umgebungsgtemperatur_c"`
	LebensdauerH            float32  `json:"lebensdauer_h"`
	SockelAnode             string   `json:"sockel_anode"`
	SockelKathode           string   `json:"sockel_kathode"`
	Kuehlung                string   `json:"kuehlung"`
	Brennstellung           string   `json:"brennstellung"`
	DeklarationsDatum       string   `json:"deklarations_datum"`
	ErzeugnissNummern       []string `json:"erzeugniss_nummern"`
	Stoff                   string   `json:"stoff"`
	StoffCasNummer          string   `json:"stoff_cas_nummer"`
	ScipNummern             []string `json:"scip_nummern"`
	EAN                     string   `json:"ean"`
	MetelCode               string   `json:"metel_code"`
	SegNo                   string   `json:"seg_no"`
	StkNummer               string   `json:"stk_nummer"`
	UkOrg                   string   `json:"uk_org"`
}
