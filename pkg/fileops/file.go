package fileops

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/redisops"
)

func Write(configFilePath string, fileName string, content string, c context.Context) error {

	log := logagent.InstArch(c)
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
	b, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	} else {
		return string(b), err
	}
}

func ReadFrom(file io.Reader, c context.Context) (string, error) {
	b, err := io.ReadAll(file) //.ReadFile(filePath)
	log := logagent.InstArch(c)
	if err != nil {

		log.Panic(err.Error())
		return "", err
	}
	str := string(b)
	return str, err
}

func Writeover(filepath string, content string, c context.Context) {
	log := logagent.InstArch(c)
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

func WriteToPathWithFilename(path string, configFileContent map[string]interface{}, filenameenv string, c context.Context) (string, error) {
	lastIndex := strings.LastIndex(path, "/")
	log := logagent.InstArch(c)
	configFilePath := path[0:lastIndex]
	// fileName := path[lastIndex+1:]

	// dotIndex := strings.LastIndex(fileName, ".")
	// var filenamestr = fileName[0:dotIndex] + "-" + env + ".yml"
	// var filenamestr = "application-" + env + ".yml"

	str, err := Read(configFilePath + string(os.PathSeparator) + filenameenv)
	log.Print(filenameenv)
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

	err = Write(configFilePath, filenameenv, writeContent, c)

	return writeContent, err
}

func WriteToPath(path, filename string, configFileContent map[string]interface{}, env string, c context.Context) (string, error) {
	lastIndex := strings.LastIndex(path, "/")
	log := logagent.InstArch(c)
	configFilePath := path[0:lastIndex]
	// fileName := path[lastIndex+1:]

	// dotIndex := strings.LastIndex(fileName, ".")
	// var filenamestr = fileName[0:dotIndex] + "-" + env + ".yml"
	// var filenamestr = "application-" + env + ".yml"

	str, err := Read(configFilePath + string(os.PathSeparator) + filename)
	log.Print(filename)
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

	err = Write(configFilePath, filename, writeContent, c)

	return writeContent, err
}

func WriteToRepo(path, filename, team, appname, writeContent string, envdc, version, region string, c context.Context) (string, error) {

	log := logagent.InstArch(c)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// configFilePath := path
	// var filenamestr = "application-" + env + ".yml"

	log.Print(filename)

	// writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	//FIXME: write config to nexus repo
	// bodystring, filename, team, appname string)
	WriteToNexus(writeContent, filename, team, appname, envdc, version, region, c)
	_, err := rediscli.Do("HSET", "confsolver-"+appname, filename, writeContent)
	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)

	return writeContent, err
}

func GetFromRepo(team, appname, envdc, version, region, filename string, c context.Context) string {
	log := logagent.InstArch(c)

	rediscli := redisops.Pool().Get()

	defer rediscli.Close()
	value, err := redis.String(rediscli.Do("HGET", "confsolver-"+appname, filename))
	if err != nil || value == "" {
		log.Print("get from repo")
		url := "http://af-nexus.kube.com/repository/fls-aflm" + GetRepoPath(team, appname, envdc, version, region) + filename
		log.Print("url: " + url)
		username := "admin"
		password := "Paic,1234"

		// Create a HTTP client with basic authentication
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Failed to create HTTP request:", err)
			os.Exit(1)
		}
		req.SetBasicAuth(username, password)

		// Perform the HTTP request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to perform HTTP request:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		resbody, err := io.ReadAll(resp.Body)

		value = string(resbody)
		if err != nil || resp.StatusCode != 200 {
			log.Panic("consolver gen :" + value)
		}
		// Print the response status code and body
		fmt.Println("Response Status Code:", resp.StatusCode)
		fmt.Println("Response Body:", value)
		_, err = rediscli.Do("HSET", "confsolver-"+appname, filename, value)
		if err != nil {
			log.Error("redis set error")
		}

		// rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)

		// return string(resbody)
	} else {
		log.Print("get from redis")
	}

	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)
	return value
}

func GetRepoPath(team, appname, envdc, version, region string) string {
	return "/consolver/" + team + "/" + appname + "/" + version + "/" + envdc + "/" + region + "/"
}

func WriteToNexus(bodystring, filename, team, appname, envdc, version, region string, c context.Context) {
	// Open the file to be uploaded
	// file, err := os.Open("1.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://af-nexus.kube.com/service/rest/v1/components?repository=fls-aflm", nil)
	if err != nil {
		panic(err)
	}

	// Set basic authentication header
	req.SetBasicAuth("admin", "Paic,1234")

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the directory field to the form
	err = writer.WriteField("raw.directory", GetRepoPath(team, appname, envdc, version, region))
	if err != nil {
		panic(err)
	}

	err = writer.WriteField("raw.asset1.filename", filename)
	if err != nil {
		panic(err)
	}

	// Create a file part for the file field
	filePart, err := writer.CreateFormFile("raw.asset1", filename)
	if err != nil {
		panic(err)
	}
	// filePart.Write(bodybytes)
	// Copy the file content to the file part
	stringreader := strings.NewReader(bodystring)
	_, err = io.Copy(filePart, stringreader)
	if err != nil {
		panic(err)
	}

	// Close the multipart writer to finalize the form data
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Set the request body as the form data
	req.Body = io.NopCloser(body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and make the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")
	io.Copy(os.Stdout, resp.Body)
}

func WriteToAppPath(path, filename, appname string, configFileContent map[string]interface{}, env string, c context.Context) (string, error) {

	log := logagent.InstArch(c)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// configFilePath := path
	// var filenamestr = "application-" + env + ".yml"

	log.Print(filename)

	writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	_, err := rediscli.Do("HSET", "confsolver-"+appname, filename, writeContent)
	rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)

	return writeContent, err
}

// 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, c context.Context) []string {
	log := logagent.InstArch(c)
	var files []string
	dir, err := os.ReadDir(dirPth)
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
