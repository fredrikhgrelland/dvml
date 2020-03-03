.PHONY: run

run:
	go run main.go -dir conf

debug:
	go run main.go -dir conf -debug