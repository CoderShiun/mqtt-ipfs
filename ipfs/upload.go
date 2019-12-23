package ipfs

import (
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

func UploadIPFS(str string) string {
	sh = shell.NewShell("localhost:5001")

	hash, err := sh.Add(bytes.NewBufferString(str))
	if err != nil {
		fmt.Println("uploading err：", err)
	}

	/*fmt.Println(hash)

	dirsh, err := sh.AddDir("/home/shiun/forIPFS")
	if err != nil {
		fmt.Println("uploading err：", err)
	}

	fmt.Println(dirsh)*/

	return hash
}