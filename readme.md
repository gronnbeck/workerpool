# workerpool

`workerpool` is a simple worker pool library written for go.

## Examples

### Wait for all work to complete

```go
func main() {
  pool := workerpool.New(1000)
  pool.Start()
  work := 10
  res := make(chan error)

  doWork := func() error {
    return nil
  }

  go func() {
    for i:=0; i < work; i++ {
      ctx := context.Background()
      pool.Queue(ctx, func() {
        res <- doWork() // some function that does work
      })
    }
  }()


  for i := i < work; i ++ {
    select {
    case err := <-res:
      if err != nil {
        // handle error here or
        panic(err)
      }
    }
  }
}
```
