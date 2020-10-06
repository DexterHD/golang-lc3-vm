# LC-3 Virtual Machine

This is my implementation of "Little Computer 3" virtual machine. I was made it to get better
understanding about what is going on inside a computer and better understand how programming languages work.

This project based on: [Write your Own Virtual Machine](https://justinmeiners.github.io/lc3-vm/) article.
You can also find some basics about LC3 on [Wikipedia](https://en.wikipedia.org/wiki/Little_Computer_3).

In `lc3-isa.pdf` file within this repository you can find The Instruction Set Architecture (ISA) of the LC-3.
Inside `apps` directory you can find some applications which can be run by this Virtual Machine.

## Build

- To run linter `make lint`
- To build project just run `make build`
- To run tests `make test`
- To remove build artifacts `make clean`

## Run

It's easy to run any application. Just specify application object file as a first argument to a binary.

```bash
./golang-lc3-vm ./apps/hello-world.obj
./golang-lc3-vm ./apps/rogue.obj
```

## Alternative "Go" implementations

- https://github.com/ziggy42/gLC3
- https://github.com/robmorgan/go-lc3-vm

## Useful links

- [Wikipedia: Little Computer 3](https://en.wikipedia.org/wiki/Little_Computer_3)
- [LC3 Web simulator](https://wchargin.github.io/lc3web/)
