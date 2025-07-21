package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	todo := context.TODO()
	value, ok := ctx.Value("key").(string)
	if !ok {
		fmt.Println("not ok")
	}
	fmt.Println(value)
	fmt.Println(ctx)
	fmt.Println(todo)
}
