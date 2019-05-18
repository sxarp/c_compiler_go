dc=docker-compose
dcex=$(dc) exec main
go=$(dcex) go
test=$(go) test
lint=golangci-lint run ./... -E golint

test:
	$(test) -cover -count=1 ./... # run tests without using cache

lint:
	$(dcex) $(lint)

lint-ci:
	$(lint)

testv:
	$(test) -cover -v -count=1 ./... # run tests without using cache

unit:
	$(test) -v -cover -count=1 -run $(f) ./$(d) # $ make unit f=TestFuncName d=app/hoge

start:
	$(dc) up --build -d

stop:
	$(dc) down

attach:
	docker exec -it $$(docker ps -f name=c_compiler -q) /bin/bash
