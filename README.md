# tenntenn/structs

[![pkg.go.dev][gopkg-badge]][gopkg]

```go
type A { N int `json:"n"`}
type B { S string `json:"s"`}

/*
struct{
     N int	`json:"n"`
     S string	`json:"s"`
}{
	N: 100,
	S: "sample",
}
*/
structs.Merge(A{N:100}, B{S:"sample"})
```

```go
/*
struct{
     N int
     S string
}{
	N: 100,
	S: "sample",
}
*/
structs.Of(structs.F("N", 100), structs.F("S", "sample"))

/*
struct{
     N int	`json:"n"`
     S string	`json:"s"`
}{
	N: 100,
	S: "sample",
}
*/
structs.Of(structs.F("N", 100, structs.Tag(`json:"n"`)), structs.F("S", "sample", structs.Tag(`json:"s"`)))
```


<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/tenntenn/structs
[gopkg-badge]: https://pkg.go.dev/badge/github.com/tenntenn/structs?status.svg
