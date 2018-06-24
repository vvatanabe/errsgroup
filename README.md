# errsgroup [![Build Status](https://travis-ci.org/vvatanabe/errsgroup.svg?branch=master)](https://travis-ci.org/vvatanabe/errsgroup) [![Coverage Status](https://coveralls.io/repos/github/vvatanabe/errsgroup/badge.svg?branch=master)](https://coveralls.io/github/vvatanabe/errsgroup?branch=master)
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
func MaxParallelSize(size int) Option
```
Set the maximum size of internal error chanel. (default:0)
```go
func ErrorChanelSize(size int) Option
```

## Acknowledgments
[golang.org/x/sync/errgroup](https://github.com/golang/sync/tree/master/errgroup) really inspired me. I appreciate it.

## Bugs and Feedback
For bugs, questions and discussions please use the Github Issues.