#!/bin/bash

header="Authorization: bearer $id_token"
endpoint="$gateway_uri/interpret"

response=$(curl -s -H "$header" -d '{"diary": "1/2 345 abc, def, ghi;"}' $endpoint)

expected_output="Diary Output: Token Type: 2, Lexeme: 1/2, Token Type: 3, Lexeme: 345, Token Type: 4, Lexeme: abc, Token Type: 6, Lexeme: ,, Token Type: 4, Lexeme: def, Token Type: 6, Lexeme: ,, Token Type: 4, Lexeme: ghi, Token Type: 1, Lexeme: ;"
if [ "$response" != "$expected_output" ]; then
    echo -e "API test failed.\\nGot: $response\nexpected: $expected_output"
    exit 1
fi
echo "API test succeeded.\\nGot: $response"
