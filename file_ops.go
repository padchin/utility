package utility

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetFilesByPath возвращает список файлов рекурсивно или нет с указанными расширениями или все не являющиеся каталогами
func GetFilesByPath(path string, recursive bool, extensions ...string) (files_list []string, err error) {
	//todo возвращать полное имя файла вместе с каталогами
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
		current_path := ""
		files, err2 := os.ReadDir(s_path)
		// проверяется если указанный путь является абсолютным, то используется он для построения списка
		if !filepath.IsAbs(s_path) {
			current_path, err = os.Getwd()
			if err != nil {
				return nil, err
			}
		}
		if err2 != nil {
			log.Println(err)
			return nil, err2
		}
		for _, f := range files {
			if len(extensions) > 0 {
				for _, ext := range extensions {
					if strings.HasSuffix(f.Name(), ext) {
						if len(current_path) > 0 {
							files_list = append(files_list, filepath.Join(current_path, f.Name()))
						} else {
							files_list = append(files_list, filepath.Join(path, f.Name()))
						}
					}
				}
			} else {
				if len(current_path) > 0 {
					files_list = append(files_list, filepath.Join(current_path, f.Name()))
				} else {
					files_list = append(files_list, filepath.Join(path, f.Name()))
				}
			}
		}
	}
	return
}

// GetFilesListByMask получение списка файлов по маске
func GetFilesListByMask(mask string) ([]string, error) {
	return filepath.Glob(mask)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		err_copy := out.Close()
		if err == nil {
			err = err_copy
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	return nil
}

func JSONDump(obj interface{}, file string) error {
	json_bytes, err_json := json.MarshalIndent(obj, "", "    ")
	if err_json != nil {
		return err_json
	}
	err_write := ioutil.WriteFile(file, json_bytes, 0644)
	if err_write != nil {
		return err_write
	}
	return nil
}

//JSONLoad
func JSONLoad(obj interface{}, file string) error {
	json_bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(json_bytes, obj)
	if err != nil {
		return err
	}
	return nil
}
