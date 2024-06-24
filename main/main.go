package main

import (
	"github.com/caesar-rocks/core"
)

type PostController struct {
	*core.BaseResourceController
}

// func (pc *PostController) Index(ctx *core.CaesarCtx) error {
// 	fmt.Println("PostController Index")
// 	return nil
// }

func main() {
	// r := core.NewRouter()
	c := &PostController{}
	c.Index(nil)
}
