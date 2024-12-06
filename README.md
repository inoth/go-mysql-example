#### update go mod 
```
require (
	github.com/go-mysql-org/go-mysql v0.0.0
)
replace github.com/go-mysql-org/go-mysql => ../go-mysql
```

```sh
# run server
go run ./server/.

# run client
go run ./client/.
```