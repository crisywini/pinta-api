#!/bin/bash


migrate -path migrations/postgres \
  -database "postgres://app:secret@localhost:5432/pintapp?sslmode=disable" up
