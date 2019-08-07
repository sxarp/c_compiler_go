dc=docker-compose
dcex=$(dc) exec main
go=$(dcex) go
test=$(go) test
lint=golangci-lint run ./... -E golint

# start dev env
start:
	$(dc) up -d

# stop dev env
stop:
	$(dc) down

# exec test when source files are edited
watch:
	 find ./ -name '*.go' | entr make check

check:
	make test; make lint

# run tests without using cache
test:
	$(test) -cover -count=1 ./...

lint:
	$(dcex) $(lint)

lint-for-ci:
	$(lint)

# test with verbose outputs
testv:
	$(test) -cover -v -count=1 ./...

# unit test: $ make unit f=TestFuncName d=app/hoge
unit:
	$(test) -v -cover -count=1 -run $(f) ./$(d)

# attach into the container
attach:
	docker exec -it $$(docker ps -f name=c_compiler -q) /bin/bash

# compile source file at $SRC_PATH and execute
exec:
	$(dcex) /bin/bash -c "go run ./src < ${SRC_PATH} > ./tmp/out.s; gcc -o ./tmp/out.o ./tmp/out.s; ./tmp/out.o"

exec-qsort:
	SRC_PATH=./examples/quick_sort.c make exec
