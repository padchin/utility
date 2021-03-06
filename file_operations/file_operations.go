package file_operations

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetFilesByPath возвращает список файлов (не каталогов). Если указан абсолютный путь, то будет возвращен список файлов
// с абсолютными путями.
func GetFilesByPath(path string, recursive bool, extensions ...string) ([]string, error) {
	var asFilesList []string
	var err error

	// если указанный путь является абсолютным, то используется он для построения списка
	sCurrentPath := ""

	if !filepath.IsAbs(path) {
		sCurrentPath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	sPath, err := filepath.Abs(filepath.Dir(path))

	if err != nil {
		return nil, err
	}

	if recursive {
		err = filepath.Walk(sPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					if len(extensions) > 0 {
						for _, ext := range extensions {
							if strings.HasSuffix(info.Name(), ext) {
								asFilesList = append(asFilesList, path)
							}
						}
					} else {
						asFilesList = append(asFilesList, path)
					}
				}

				return nil
			})
		if err != nil {
			return nil, err
		}
	} else {
		files, err := os.ReadDir(filepath.Join(sCurrentPath, path))

		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if len(extensions) > 0 {
				for _, ext := range extensions {
					if strings.HasSuffix(f.Name(), ext) {
						if len(sCurrentPath) > 0 {
							asFilesList = append(asFilesList, filepath.Join(sCurrentPath, path, f.Name()))
						} else {
							asFilesList = append(asFilesList, filepath.Join(path, f.Name()))
						}
					}
				}
			} else {
				if len(sCurrentPath) > 0 {
					asFilesList = append(asFilesList, filepath.Join(sCurrentPath, path, f.Name()))
				} else {
					asFilesList = append(asFilesList, filepath.Join(path, f.Name()))
				}
			}
		}
	}

	return asFilesList, nil
}

// GetFilesListByMask получение списка файлов по маске
func GetFilesListByMask(mask string) ([]string, error) {
	return filepath.Glob(mask)
}

// CopyFile copies the contents of the file named src to the file named by dst. The file will be created if it does not
// already exist. If the destination file exists, all it's contents will be replaced by the contents of the source file.
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
		errCopy := out.Close()
		if err == nil {
			err = errCopy
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
	_ = CopyFile(file, file+".bak")
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, jsonBytes, 0644)

	if err != nil {
		return err
	}

	return nil
}

func JSONLoad(obj interface{}, file string) error {
	jsonBytes, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, obj)

	if err != nil {
		return err
	}

	return nil
}
