package fileops

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/redisops"
)

func Write(configFilePath string, fileName string, content string, c context.Context) error {

	log := logagent.Inst(c)
	dirPath := configFilePath + "/"
	var err error
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0777)
		if err != nil {
			log.Printf("make dir error: %v\n", err)
			return err
		}
	}

	filePath := configFilePath + "/" + fileName
	// _, err := os.Stat(filePath) //os.Stat获取文件信息
	var f *os.File
	// var err error
	if checkFileExist(filePath) {
		f, err = os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY, 0644) //打开文件
		if err != nil {
			fmt.Println("file open fail", err)
			return err
		}
	} else {
		f, err = os.Create(filePath)
		if err != nil {
			log.Panic(err)
			return err
		}
	}

	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil

}

func checkFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func Read(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", err
	} else {
		return string(b), err
	}
}

func ReadFrom(file io.Reader, c context.Context) (string, error) {
	b, err := ioutil.ReadAll(file) //.ReadFile(filePath)
	log := logagent.Inst(c)
	if err != nil {

		log.Panic(err.Error())
		return "", err
	}
	str := string(b)
	return str, err
}

func Writeover(filepath string, content string, c context.Context) {
	log := logagent.Inst(c)
	pos := strings.LastIndex(filepath, "/")
	dirPath := filepath[:pos+1]
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	// var err error
	// if _, err = os.Stat(dirPath); os..IsNotExist(err) {
	// 	err = os.Mkdir(dirPath, 0777)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// }

	var f *os.File
	f, err = os.Create(filepath)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		log.Panic(err)
	}
	_, err = fmt.Fprintf(f, "%s", content)

	if err != nil {
		log.Panic(err)
	}
}

func WriteToPath(path string, configFileContent map[string]interface{}, env string, c context.Context) (string, error) {
	lastIndex := strings.LastIndex(path, "/")
	log := logagent.Inst(c)
	configFilePath := path[0:lastIndex]
	// fileName := path[lastIndex+1:]

	// dotIndex := strings.LastIndex(fileName, ".")
	// var filenamestr = fileName[0:dotIndex] + "-" + env + ".yml"
	var filenamestr = "application-" + env + ".yml"

	str, err := Read(configFilePath + string(os.PathSeparator) + filenamestr)
	log.Print(filenamestr)
	var writeContent string

	if err != nil {

		writeContent = convertops.ConvertStrMapToYaml(&configFileContent, c)
	} else {
		m := convertops.ConvertYamlToMap(str, c)
		// if _, ok := m["af-arch"]; ok {
		// 	m["af-arch"].(map[interface{}]interface{})["resource"] = configFileContent["af-arch"].(map[string]interface{})["resource"]
		// } else {
		// 	m["af-arch"] = configFileContent["af-arch"]
		// }
		m["af-arch"] = configFileContent["af-arch"]

		writeContent = convertops.ConvertMapToYaml(&m, c)
	}

	err = Write(configFilePath, filenamestr, writeContent, c)

	return writeContent, err
}

func WriteToAppPath(path, appname string, configFileContent map[string]interface{}, env string, c context.Context) (string, error) {

	log := logagent.Inst(c)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// configFilePath := path
	var filenamestr = "application-" + env + ".yml"

	log.Print(filenamestr)

	writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	_, err := rediscli.Do("HSET", "confsolver-"+appname, filenamestr, writeContent)
	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)

	return writeContent, err
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, c context.Context) []string {
	log := logagent.Inst(c)
	var files []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		log.Panic(err)
	}

	// PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if !fi.IsDir() { // 目录, 递归遍历
			// 过滤指定格式
			// ok := strings.HasSuffix(fi.Name(), ".go")
			files = append(files, fi.Name())
			// fileinfo := strings.Split(fi.Name(), ".")
			// if len(fileinfo) > 1 {

			// 	switch fileinfo[0] {
			// 	case "defaultconfig":
			// 	case "Dockerfile":
			// 	case "jenkins":
			// 	case "values":

			// 	}
			// }
		}
	}

	return files
}
