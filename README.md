goSSEClient
===========

Helper for Server-Sent Events

Returns a channel of SSEvents.
```go
type SSEvent struct {
	Id   string
	Data []byte
}
```

## Example
```go
sseEvent, err := goSSEClient.OpenSSEUrl(url)
if err != nil {
	panic(err)
}

for ev := range sseEvent {
	// do what you want with the event
}

```
