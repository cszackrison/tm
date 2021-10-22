# tm
tm is a minimalist task manager. It is extremely simple. It uses go, html, sqlite. It runs in the browser at port 8080.

## Usage:

### Production
Build the binary, then run:
```
$ ./tm -db=./tasks.db -port=:8080
```

#### options:
```
-db=:memory:
-port=:8080
```

### Test
Run this command and open `localhost:8080` in a browser.
```
$ go run main.go
```

### Build
```
$ go build
```
