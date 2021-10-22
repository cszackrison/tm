# tm
tm is a minimalist task manager. It is extremely simple and lightweight.

## Usage:

### Production
Build the binary, then run:
```
$ ./tm -db=tasks.db -port=:8080
```
#### Command-line Parameters
Parameters can be given using the following format: `KEY=VALUE`.

| Name | Key | Value | Description |
|------|-----|----------|-------------|
| Database Path | `db` | `string` | The path to the sqlite database.<br/>default is `:memory:` |
| Port | `port` | `string` | The port to use.<br/>default is `:8080` |

#### URL Parameters
Parameters can be given using the following format: `?KEY=VALUE`.

| Name | Key | Value | Description |
|------|-----|----------|-------------|
| Board ID | `boardId` | `string` | The ID of the board you want tasks bound to. |


## Development
Run this command and open `localhost:8080?boardId=hello` in a browser.
```
$ go run main.go
```

## Build
```
$ go build
```
