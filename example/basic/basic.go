package main

import (
	"fmt"
	"net/http"

	"github.com/vvatanabe/errsgroup"
)

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
