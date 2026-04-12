.PHONY: test

test:
	go test -v ./tests/...

run: 
	docker-compose up -d 

seedData: 
	./seed_data_for_testing_ui.sh 

tearDown: 
	docker-compose down
