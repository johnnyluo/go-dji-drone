package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		item, err := r.ReadString('\n')
		if nil != err {
			fmt.Println("err", err)
			return
		}
		result := bytes.TrimSuffix([]byte(item), []byte{'\n'})
		strResult := string(result)
		fmt.Println(strResult)
	}

}
