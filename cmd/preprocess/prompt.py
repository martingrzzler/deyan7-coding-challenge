prompt_base = "Given the following extracted text from a pdf document that describes a product return a json object that represents the product."

prompt_json_schema = """
The json object should have the following schema:
{
  "name": "string",
  "anwendungs_gebiete": ["string1", "string2", "..."],
  "vorteile": ["string1", "string2", "..."],
  "eigenschaften": ["string1", "string2", "..."],
  "nenn_strom_a": 0,
  "strom_steuer_a_min": 0,
  "strom_steuer_a_max": 0,
  "nenn_leistung_w": 0,
  "nenn_spannung_v": 0,
  "durchmesser_mm": 0,
  "laenge_mm": 0,
  "laenge_mit_sockel_mm": 0,
  "lcl_mm": 0,
  "kabel_laenge_mm": 0,
  "elekroden_abstand_mm": 0,
  "produkt_gewicht_g": 0,
  "max_umgebungsgtemperatur_c": 0,
  "lebensdauer_h": 0,
  "sockel_anode": "string",
  "sockel_kathode": "string",
  "kuehlung": "string",
  "brennstellung": "string",
  "deklarations_datum": "YYYY-MM-DD",
  "erzeugniss_nummern": ["string1", "string2", "..."],
  "stoff": "string",
  "stoff_cas_nummer": "string",
  "scip_nummern": ["string1", "string2", "..."],
  "ean": "string",
  "metel_code": "string",
  "seg_no": "string",
  "stk_nummer": "string",
  "uk_org": "string"
}
"""

def build_prompt(pdf_text: str) -> str:
    return f"{prompt_base}\nextracted text:\n{pdf_text}\n{prompt_json_schema}"
