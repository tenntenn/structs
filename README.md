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
Merge(A{N:100}, B{S:"sample"})
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
Of("N", 100, "S", "sample")

/*
struct{
     N int	`json:"n"`
     S string	`json:"s"`
}{
	N: 100,
	S: "sample",
}
*/
Of("N", 100, `json:"n"`, "S", "sample", `json:"s"`)
```


<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/tenntenn/structs
[gopkg-badge]: https://pkg.go.dev/badge/github.com/tenntenn/structs?status.svg
