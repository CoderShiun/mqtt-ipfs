package ipfs

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
)

func CatIPFS(hash string) string {
	sh = shell.NewShell("localhost:5001")

	read, err := sh.Cat(hash)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(read)

	return string(body)
}