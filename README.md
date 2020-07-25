# Narcissus

[![PkgGoDev](https://pkg.go.dev/badge/github.com/yhakbar/narcissus)](https://pkg.go.dev/github.com/yhakbar/narcissus)

Narcissus updates a Golang struct with fields that have been tagged with `ssm:"Parameter"` according to the corresponding value in SSM Parameter Store using reflection.

## Installation

```bash
go get github.com/yhakbar/narcissus
```

## Example Usage

```golang
import "github.com/yhakbar/narcissus"
type Name struct {
    FirstName string `ssm:"Name/FirstName"`
    LastName  string `ssm:"Name/LastName"`
}

type Contact struct {
    Email  string `ssm:"Contact/Email"`
    Number string `ssm:"Contact/Number"`
}

type Person struct {
    Name                       Name
    Contact                    Contact
    FavoriteNumber             int     `ssm:"FavoriteNumber"`
    FavoriteInconvenientNumber float64 `ssm:"FavoriteInconvenientNumber"`
}
person := Person{}
ssmPath := "/path/to/parameters/"
// You can get this wrapper like so: wrapper := narcissus.Wrapper{Client: client}
_ = narcissus.UpdateBySSM(&person, &ssmPath)
// If you want to reuse an SSM client, do so like this:
// wrapper := narcissus.Wrapper{Client: client}
// wrapper.UpdateBySSM(&person, &ssmPath)
```
