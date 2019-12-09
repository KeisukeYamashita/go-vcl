# VCL

## Usage

### Decode

Let's say you have a VCL file.

Define a go struct how you what to retrieve the VCL expressions and attributes.

```golang
type Root struct {
    ACLs []ACL `vcl:"acl,block"`
}

type ACL struct {
    Type string `vcl:"type,label`
}    
```

Then decode your 

```golang
var r Root
_ = vcl.Decode(b, &r)
fmt.Println(r.Type)
fmt.Println(r.ACLs)
```

```console
$ go run main.go
=> "local"
```

### Encode

* TODO: Please contribute:]

## License

go-vcl is released under the MIT license.  
© 2019 KeisukeYamashita.

## Author

* [KeisukeYamashita](https://github.com/KeisukeYamashita)
