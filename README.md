# LC-3 Virtual Machine

Based on: [Write your Own Virtual Machine](https://justinmeiners.github.io/lc3-vm/) article.

## Usage

### Build

```bash
go build .
```

### Test

```bash
go test -v ./...
```

### Run

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