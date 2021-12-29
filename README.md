# README

```sh
go get github.com/hkloudou/nrpc
```

## sample
``` go
conn, err := nats.Connect(nats.DefaultURL)
if err != nil {
    panic(err)
}
ser := nrpc.New[wrapperspb.StringValue, string](conn)
ser.Queue("test", "q1", func(req *wrapperspb.StringValue) (*string, error) {
    return nrpc.PointerOf("test1.1"), nil
})
ser.Queue("test", "q1", func(req *wrapperspb.StringValue) (*string, error) {
    return nrpc.PointerOf("test1.2"), nil
})
ser.Queue("test", "q2", func(req *wrapperspb.StringValue) (*string, error) {
    return nrpc.PointerOf("test2"), nil
})
cli := nrpc.NewRequest[wrapperspb.StringValue, string](conn, "test")
for i := 0; i < 3; i++ {
    res, err := cli.Request(wrapperspb.String("t1"), 5*time.Second)
    log.Println("res", *res, "err", err)
}
```