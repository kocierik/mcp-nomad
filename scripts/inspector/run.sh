#!/bin/bash

cd scripts/inspector
npm install

npx @modelcontextprotocol/inspector go run  ../../main.go "$@"