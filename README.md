# errsgroup
support multiple error handling synchronously with goroutine

## Installation
This package can be installed with the go get command:
```
$ go get github.com/vvatanabe/errsgroup
```

## Usage

### Basic

``` go
func main() {
	g := errsgroup.NewGroup()
	var urls = []string{
		"https://backlog.com/",
		"https://cacoo.com/",
		"https://www.typetalk.com/",
	}
	for _, url := range urls {
		url := url
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	if errs := g.Wait(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Successfully")
	}
}
```

### WithContext

``` go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errsgroup.WithContext(ctx)
	client := http.DefaultClient
	var urls = []string{
		"https://backlog.com/",
		"https://cacoo.com/",
		"https://www.typetalk.com/",
	}
	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		req = req.WithContext(ctx)
		g.Go(func() error {
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	cancel()
	if errs := g.Wait(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Successfully")
	}
}
```

### Functional Options
Set the maximum size of parallels. (default:1)
```go
func LimitSize(size int) Option
```
Set the maximum size of internal error chanel. (default:0)
```go
func ErrorChanelSize(size int) Option
```