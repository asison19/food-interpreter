#!/bin/bash

input='{"diary": "1/2 345 abc, def, ghi;"}'
expected_output="Diary Output: Token Type: 2, Lexeme: 1/2, Token Type: 3, Lexeme: 345, Token Type: 4, Lexeme: abc, Token Type: 6, Lexeme: ,, Token Type: 4, Lexeme: def, Token Type: 6, Lexeme: ,, Token Type: 4, Lexeme: ghi, Token Type: 1, Lexeme: ;"

response=$(curl -s -H "Authorization: bearer $(gcloud auth print-identity-token)" -d $input "$gateway_uri/interpret")

if [$response != $expected_output]; then
    echo "Got: $response; expected: $expected_output"
    exit 1
fi
