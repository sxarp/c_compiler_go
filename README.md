# Start development & run example

Start development environment.

```sh
$ make start
docker-compose up -d
Creating network "c_compiler_go_default" with the default driver
Creating c_compiler_go_main_1 ... done
```

Run tests.

```sh
$ make test
docker-compose exec main go test -cover -count=1 ./...
ok  	github.com/sxarp/c_compiler_go/src	2.659s	coverage: 35.7% of statements
ok  	github.com/sxarp/c_compiler_go/src/asm	0.009s	coverage: 88.4% of statements
ok  	github.com/sxarp/c_compiler_go/src/ast	0.005s	coverage: 16.7% of statements
?   	github.com/sxarp/c_compiler_go/src/em	[no test files]
ok  	github.com/sxarp/c_compiler_go/src/gen	1.824s	coverage: 96.0% of statements
?   	github.com/sxarp/c_compiler_go/src/h	[no test files]
ok  	github.com/sxarp/c_compiler_go/src/psr	0.005s	coverage: 93.3% of statements
?   	github.com/sxarp/c_compiler_go/src/str	[no test files]
ok  	github.com/sxarp/c_compiler_go/src/tok	0.008s	coverage: 86.3% of statements
ok  	github.com/sxarp/c_compiler_go/src/tp	0.004s	coverage: 83.3% of statements
```

Run quick sort sample.

```sh
$ make exec-qsort
SRC_PATH=./examples/quick_sort.c make exec
docker-compose exec main /bin/bash -c "go run ./src < ./examples/quick_sort.c > ./tmp/out.s; gcc -o ./tmp/out.o ./tmp/out.s; ./tmp/out.o"
12 17 20 11 28 19 0 23 16 9 6 25 28 13 2 5 16 13 16 1 24 29 20 27 26 5 4 21 14 27
0 1 2 4 5 5 6 9 11 12 13 13 14 16 16 16 17 19 20 20 21 23 24 25 26 27 27 28 28 29
```


# Organization of directories

* [src](https://github.com/sxarp/c_compiler_go/tree/master/src)

Entry point of compiler where `main.go` resides.

* [src/asm](https://github.com/sxarp/c_compiler_go/tree/master/src/asm)

Definition of the DSL for generating x86-64 assembly.

* [src/ast](https://github.com/sxarp/c_compiler_go/tree/master/src/ast)

Definition for the struct of AST.

* [src/em](https://github.com/sxarp/c_compiler_go/tree/master/src/em)

Helpers to show error messages when parsers fail.

* **[src/gen](https://github.com/sxarp/c_compiler_go/tree/master/src/gen)**

**The Compiler is defined here using the DSLs and with the supports of other packages.**

* [src/h](https://github.com/sxarp/c_compiler_go/tree/master/src/h)

Helpers for testing.

* [src/psr](https://github.com/sxarp/c_compiler_go/tree/master/src/psr)

Definition of the parsers.

* [src/str](https://github.com/sxarp/c_compiler_go/tree/master/src/str)


Utilities to write assembly codes.

* [src/tok](https://github.com/sxarp/c_compiler_go/tree/master/src/tok)


Definition of the tokenizer.

* [src/tp](https://github.com/sxarp/c_compiler_go/tree/master/src/tp)

Definition of the type system.
