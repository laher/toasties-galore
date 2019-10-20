

run-all:
	docker-compose start postgres
	cd chillybin && nohup go run . &
	cd jafflr && nohup go run . &

test:
	CHILLYBIN_ADDR=http://localhost:7001 JAFFLR_ADDR=http://localhost:7000 go test -v ./integration-tests 
