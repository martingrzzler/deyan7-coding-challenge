import os
import pdfplumber
from prompt import build_prompt
from openai import OpenAI
import json

oai_client = OpenAI()
dataset_dir = os.path.join(os.path.dirname(__file__), '../../dataset')
out = os.path.join(os.path.dirname(__file__), '../../dataset.jsonl')

outFile = open(out, 'w')

for filename in os.listdir(dataset_dir):
    with pdfplumber.open(os.path.join(dataset_dir, filename)) as pdf:
        pdf_text = ''
        for page in pdf.pages:
            pdf_text += "\n" + page.extract_text()

        p = build_prompt(pdf_text)

        completion = oai_client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": p}],
            response_format={"type": "json_object"}
        )

        message = completion.choices[0].message
        if message.content:
            obj = json.loads(message.content)
            outFile.write(json.dumps(obj) + '\n')
        else:
            raise Exception("No response from OpenAI at file " + filename)
        print("Processed file " + filename)

outFile.close()
