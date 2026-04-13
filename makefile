.PHONY: test

test:
	go test -v ./tests/...

run: 
	docker-compose up -d

start: 
	docker-compose up --build -d 

seedData: 
	./seed_data_for_testing_ui.sh 

tearDown: 
	docker-compose down -v

stop: 
	docker-compose stop