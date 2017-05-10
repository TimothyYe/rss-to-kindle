package kindle

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"rss-to-kindle/utils"
)

//GenerateMobi ...
func GenerateMobi(dir string) string {
	cmd := exec.Command("/usr/local/bin/kindlegen", filepath.Join(dir, "main.opf"), "-c2", "-gif")
	err := cmd.Run()

	if err != nil {
		fmt.Println("Failed to invoke kindlegen:", err.Error())
	}

	mobiPath := filepath.Join(dir, "main.mobi")
	_, err = ioutil.ReadFile(mobiPath)
	utils.ExitIfErr(err)

	return mobiPath
}
