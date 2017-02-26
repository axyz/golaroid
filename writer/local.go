package writer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func Local(variant string, path string, data []byte, opts map[string]interface{}) error {
	folder := opts["folder"].(string)
	fmt.Println("TODO: writing to " + path)
	dst := filepath.Join(folder, variant, path)
	dirError := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if dirError != nil {
		fmt.Println(dirError.Error())
	}
	f, fileError := os.Create(dst)
	if fileError != nil {
		fmt.Println(fileError.Error())
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err := w.Write(data)
	w.Flush()
	return err
}
