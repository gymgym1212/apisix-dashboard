package dag2lua

import (
	"errors"
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println(errors.New("no script config param!"))
		return
	}
	script := os.Args[1]

	fmt.Println(script)
}
