#!/bin/bash

export DB=postgres 
export DATABASE_URL="postgres://app:secret@localhost:5432/pintapp?sslmode=disable" 
go run ./cmd/