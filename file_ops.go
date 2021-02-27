package utility

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetFilesByPath возвращает список файлов рекурсивно или нет с указанными расширениями или все не являющиеся каталогами
func GetFilesByPath(path string, recursive bool, extensions ...string) (files_list []string, err error) {
	var s_path string
	s_path, err = filepath.Abs(filepath.Dir(path))
	if err != nil {
		log.Println(err)
		return
	}
	if recursive {
		err2 := filepath.Walk(s_path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					if len(extensions) > 0 {
						for _, ext := range extensions {
							if strings.HasSuffix(info.Name(), ext) {
								files_list = append(files_list, path)
							}
						}
					} else {
						files_list = append(files_list, path)
					}
				}
				return nil
			})
		if err2 != nil {
			log.Println(err2)
			return nil, err2
		}
	} else {
		files, err2 := ioutil.ReadDir(s_path)
		if err2 != nil {
			log.Println(err)
			return nil, err2
		}
		for _, f := range files {
			if len(extensions) > 0 {
				for _, ext := range extensions {
					if strings.HasSuffix(f.Name(), ext) {
						files_list = append(files_list, f.Name())
					}
				}
			} else {
				files_list = append(files_list, f.Name())
			}
		}
	}
	return
}
