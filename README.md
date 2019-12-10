# VCL

> VCL parser written in Go

[![GitHub Actions][github-actions-badge]][github-actions]
[![GoDoc][godoc-badge]][godoc]
[![Go Report Card][go-report-card-badge]][go-report-card]
[![GolangCI][golangci-badge]][golangci]

[![License][license-badge]][license]
[![Dependabot][dependabot-badge]][dependabot]

[![DeepSource][deepsource-badge]][deepsource]

## Usage

### Decode

Let's say you have a VCL file.

```vcl
acl purge_ip {
    "localhost";
    "127.0.0.1";
}
```

Define a go struct how you what to retrieve the VCL expressions and attributes.

```golang
type Root struct {
    ACLs []*ACL `vcl:"acl,block"`
}

type ACL struct {
    Type      string `vcl:"type,label`
    Endpoints string `vcl:"endpoints,flat"`
}  
```

Then decode your like following.

```golang
var r Root
err := vcl.Decode(b, &r)
fmt.Println(r.Type)
fmt.Println(r.ACLs)
```

```console
$ go run main.go
=> "local"
=> []string{"localhost","127.0.0.1"}
```

## Supported tags

I am not a VCL master so there may be not supported features.

There are struct tags you can use for you input.

* `block`: Represents a unit of your block like `acl`, `sub`, etc...
* `label`: The label of your block.
* `flat`: Represents a expression field
* `attr`: (Default) Attribute of your block

## Releases

Release tag will be based on [Semantic Versioning 2.0.0](https://semver.org/).  
See the [CHANGELOGS.md](./CHANGELOGS.md)

## How to Contribute

I am always welcome for any contributions.

* Raise a Issue.
* Create a PR.

Simple:)

## License

go-vcl is released under the MIT license.  
© 2019 KeisukeYamashita.

## Author

* [KeisukeYamashita](https://github.com/KeisukeYamashita)

<!-- badge links -->

[dependabot]: https://dependabot.com 
[dependabot-badge]: https://badgen.net/badge/icon/Dependabot?icon=dependabot&label&color=blue

[license]: LICENSE
[license-badge]: https://img.shields.io/badge/license-Apache%202.0-%23E93424

[godoc]: https://godoc.org/github.com/KeisukeYamashita/go-vcl
[godoc-badge]: https://img.shields.io/badge/godoc.org-reference-blue.svg

[go-report-card]: https://goreportcard.com/report/github.com/KeisukeYamashita/go-vcl
[go-report-card-badge]: https://goreportcard.com/badge/github.com/KeisukeYamashita/go-vcl

[deepsource]: https://deepsource.io/gh/KeisukeYamashita/go-vcl/?ref=repository-badge
[deepsource-badge]: https://static.deepsource.io/deepsource-badge-light.svg


[github-actions]: https://github.com/KeisukeYamashita/go-vcl/actions
[github-actions-badge]: https://github.com/KeisukeYamashita/go-vcl/workflows/Main%20Workflow/badge.svg

[golangci]: https://golangci.com/r/github.com/KeisukeYamashita/go-vcl
[golangci-badge]: https://golangci.com/badges/github.com/KeisukeYamashita/go-vcl.svg
