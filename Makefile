dc=docker-compose
dcex=$(dc) exec main
go=$(dcex) go
test=$(go) test
ensure=$(dcex) dep ensure

test:
	$(ensure)
	$(test) -cover -count=1 ./... # run tests without using cache

testv:
	$(ensure)
	$(test) -cover -v -count=1 ./... # run tests without using cache

unit:
	$(ensure)
	$(test) -v -cover -count=1 -run $(f) ./$(d) # $ make unit f=TestFuncName d=app/hoge

start:
	$(dc) up --build -d

stop:
	$(dc) down

attach:
	docker exec -it $$(docker ps -f name=c_compiler -q) /bin/bash
