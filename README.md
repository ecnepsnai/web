# API

API is a library to create quick and easy JSON-based REST APIs in Golang

## Example

```golang
server = api.New("127.0.0.1:8080")
if err := server.Start(); err != nil {
	panic(err)
}

handle := func(request api.Request) (interface{}, *Error) {
	return time.Now.Unix(), nil
}
options := api.HandleOptions{}
server.APIGET("/time", handle, options)
```