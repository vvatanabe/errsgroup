package main

import (
	"fmt"
	"net/http"

	"context"

	"github.com/vvatanabe/errsgroup"
)

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
