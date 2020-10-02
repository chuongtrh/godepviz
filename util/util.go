package util

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq func
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RenderImage func
func RenderImage(str string, isSaveImage bool, path string) ([]byte, error) {

	id := fmt.Sprintf("%x", time.Now().Unix())
	dotFile := fmt.Sprintf("%s/%s.dot", path, id)
	pngFile := fmt.Sprintf("%s/%s.png", path, id)

	ioutil.WriteFile(dotFile, []byte(str), 0755)

	png, err := exec.Command("dot", "-Tpng", dotFile).Output()
	if err != nil {
		return []byte{}, err
	}
	if isSaveImage {
		ioutil.WriteFile(pngFile, png, 0755)
	}
	return png, err
}
