dc=docker-compose
dcex=$(dc) exec main
go=$(dcex) go
test=$(go) test
lint=golangci-lint run ./... -E golint

check:
	make test; make lint

test:
	$(test) -cover -count=1 ./... # run tests without using cache

watch:
	 find ./ -name '*.go' | entr make check

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

exec:
	$(dcex) /bin/bash -c "go run ./src < ${SRC_PATH} > ./tmp/out.s; gcc -o ./tmp/out.o ./tmp/out.s; ./tmp/out.o"
