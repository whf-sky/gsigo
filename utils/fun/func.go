package fun

import "os"

//文件/文件夹是否存在
//filename 文件路径
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return  err == nil ||os.IsExist(err)
}

//文件/文件夹是否存在
//filename 文件路径
func  DirExists(path string) bool {
	return FileExists(path)
}

//是否为文件夹
//path 文件路径
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return  fileInfo.IsDir()
}