package fileops

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/golang/snappy"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/fileconvagt/pkg/filesets"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/cacheops"
)

func Write(configFilePath string, fileName string, content string, c context.Context) error {

	log := logagent.RecreateWith(c)
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
	log := logagent.RecreateWith(c)
	if err != nil {

		log.Panic(err.Error())
		return "", err
	}
	str := string(b)
	return str, err
}

func Writeover(filepath string, content string, c context.Context) {
	log := logagent.RecreateWith(c)
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
	log := logagent.RecreateWith(c)
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
	log := logagent.RecreateWith(c)
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

func WriteToRepoInter(writeContent, repo, filename, appname string, getPath func() string, c context.Context) (string, error) {

	log := logagent.RecreateWith(c)
	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// configFilePath := path
	// var filenamestr = "application-" + env + ".yml"

	log.Print(filename)

	// writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	//FIXME: write config to nexus repo
	// bodystring, filename, team, appname string)
	WriteToNexusInter(writeContent, repo, filename, getPath, c)
	// _, err := rediscli.Do("HSET", "confrepo-"+appname, filename, writeContent)
	// if err != nil {
	// 	log.Error(err)
	// 	log.Error("redis set error")
	// }
	// rediscli.Do("EXPIRE", "confrepo-"+appname, 60*10)
	prefix := "confrepo-" + appname
	kvmap := cacheops.GetKVMap()
	kvmap.SetkWithPrefixDuration(prefix, filename, writeContent, 10*60, c)
	// if ok, rediscli := redisops.Health(c); ok {
	// 	_, err := rediscli.Pipelined(c, func(pipe redis.Pipeliner) error {
	// 		pipe.HSet(c, "confrepo-"+appname, filename, writeContent)
	// 		pipe.Expire(c, "confrepo-"+appname, 10*time.Minute)
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		log.Error(err)
	// 		log.Error("redis set error")
	// 	}
	// }
	return writeContent, nil
}

func WriteToRepo(path, filename, team, appname, writeContent string, envdc, version, region string, c context.Context) (string, error) {

	log := logagent.RecreateWith(c)
	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// configFilePath := path
	// var filenamestr = "application-" + env + ".yml"

	log.Print(filename)

	// writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	//FIXME: write config to nexus repo
	// bodystring, filename, team, appname string)
	WriteToNexus(writeContent, filename, team, appname, envdc, version, region, c)
	// _, err := rediscli.Do("HSET", "confrepo-"+appname, filename, writeContent)
	// if err != nil {
	// 	log.Error(err)
	// 	log.Error("redis set error")
	// }
	// rediscli.Do("EXPIRE", "confrepo-"+appname, 60*10)
	kvmap := cacheops.GetKVMap()
	prefix := "confrepo-" + appname
	kvmap.SetkWithPrefixDuration(prefix, filename, writeContent, 10*60, c)
	// if ok, rediscli := redisops.Health(c); ok {

	// 	_, err := rediscli.Pipelined(c, func(pipe redis.Pipeliner) error {
	// 		pipe.HSet(c, "confrepo-"+appname, filename, writeContent)
	// 		pipe.Expire(c, "confrepo-"+appname, 10*time.Minute)
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		log.Error(err)
	// 		log.Error("redis set error")
	// 	}
	// }

	return writeContent, nil
}

func GetIaCRepoPath(team, name, version string, c context.Context) string {
	return "/IaC/" + team + "/" + name + "/" + version + "/"
}

func GetFromRepoInter(repo, appname, filename string, getPath func() string, c context.Context) (string, error) {
	log := logagent.RecreateWith(c)

	// rediscli := redisops.Pool().Get()
	// defer rediscli.Close()
	// value, err := redis.String(rediscli.Do("HGET", "confrepo-"+appname, filename))
	var err error
	var value string
	// redishealth, rediscli := redisops.Health(c)
	prefix := "confrepo-" + appname
	kvmap := cacheops.GetKVMap()

	// key := "confrepo-" + appname
	// if redishealth {
	// 	value, err = rediscli.HGet(c, key, filename).Result()
	// } else {
	// 	err = fmt.Errorf("redis connection is not healthy")
	// }
	fullpath := repo + getPath() + filename
	if val, _, ok, _, _ := kvmap.ReadkWithPrefixType(prefix, filename, c); ok {
		value = val.(string)
		log.Print("get from redis")
		return value, nil
	}

	//get from nexus repo
	// if err != nil || value == "" {
	log.Info("get from repo")
	url := *filesets.Nexusurl + "/repository/" + fullpath
	log.Info("url: " + url)
	username := *filesets.Nexususr
	password := *filesets.Nexuspass

	// Create a HTTP client with basic authentication
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Failed to create HTTP request:", err)
		return "", err
	}
	req.SetBasicAuth(username, password)

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to do HTTP request:", err)
		return "", err
		// os.Exit(1)
	}
	defer resp.Body.Close()

	resbody, err := io.ReadAll(resp.Body)

	value = string(resbody)
	if err != nil || resp.StatusCode != 200 {
		log.Error("consolver gen :" + value)
		return "", fmt.Errorf("not found: %v", resp.StatusCode)
	}
	// Print the response status code and body
	log.Info("Response Status Code:", resp.StatusCode)
	// log.Debug("Response Body:", value)
	// _, err = rediscli.Do("HSET", "confrepo-"+appname, filename, value)
	// if err != nil {
	// 	log.Error("redis set error")
	// }
	// 	if redishealth {
	// 		err = rediscli.HSet(c, key, filename, value).Err()
	// 		if err != nil {
	// 			log.Error("redis set error")
	// 		}
	// 	}

	// 	// rediscli.Do("EXPIRE", "confrepo-"+appname, 60*10)

	// 	// return string(resbody)
	// } else {
	// 	log.Print("get from redis")
	// }

	// _, err = rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)
	// if redishealth {
	// 	err = rediscli.Expire(c, key, 10*time.Minute).Err()
	// 	if err != nil {
	// 		log.Error("redis expire error")
	// 	}
	// }

	// if len(value) < 1 {
	// 	err = fmt.Errorf("%s is empty", fullpath)
	// }

	return value, err
}

func GetFromRepo(team, appname, envdc, version, region, filename string, c context.Context) (string, error) {
	return GetFromRepoInter("fls-aflm", appname, filename, func() string {
		return GetRepoPath(team, appname, envdc, version, region)
	}, c)
}

func GetRepoPath(team, appname, envdc, version, region string) string {
	return "/consolver/" + team + "/" + appname + "/" + version + "/" + envdc + "/" + region + "/"
}

func WriteToNexusInter(bodystring, repo, filename string, getPath func() string, c context.Context) {
	// Open the file to be uploaded
	// file, err := os.Open("1.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	log := logagent.RecreateWith(c)
	// Create a new HTTP request
	req, err := http.NewRequest("POST", *filesets.Nexusurl+"/service/rest/v1/components?repository="+repo, nil)
	if err != nil {
		log.Panic("Failed to create HTTP request:", err)
	}

	// Set basic authentication header
	req.SetBasicAuth(*filesets.Nexususr, *filesets.Nexuspass)

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpath := getPath()
	// Add the directory field to the form
	err = writer.WriteField("raw.directory", fpath)
	if err != nil {
		log.Panic("Failed to WriteField raw.directory:", err)
		// panic(err)
	}

	err = writer.WriteField("raw.asset1.filename", filename)
	if err != nil {
		log.Panic("Failed to WriteField raw.asset1.filename:", err)
		// panic(err)
	}

	// Create a file part for the file field
	filePart, err := writer.CreateFormFile("raw.asset1", filename)
	if err != nil {
		log.Panic("Failed to CreateFormFile raw.asset1:", err)
	}
	// filePart.Write(bodybytes)
	// Copy the file content to the file part
	stringreader := strings.NewReader(bodystring)
	_, err = io.Copy(filePart, stringreader)
	if err != nil {
		log.Panic("Failed to copy io:", err)
	}

	// Close the multipart writer to finalize the form data
	err = writer.Close()
	if err != nil {
		log.Panic("Failed to close writer:", err)
		// panic(err)
	}

	// Set the request body as the form data
	req.Body = io.NopCloser(body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and make the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil || (resp.StatusCode != 204 && resp.StatusCode != 200) {
		log.Panic("Failed to do HTTP request:", err)
		// panic(err)
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Info("Response Status:", resp.Status)
	log.Info("Response Body:", resp.Body)
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Body:")
	// io.Copy(os.Stdout, resp.Body)
}

func WriteToNexus(bodystring, filename, team, appname, envdc, version, region string, c context.Context) {
	WriteToNexusInter(bodystring, "fls-aflm", filename, func() string {
		return GetRepoPath(team, appname, envdc, version, region)
	}, c)
	// // Open the file to be uploaded
	// // file, err := os.Open("1.txt")
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // defer file.Close()

	// log := logagent.InstArch(c)
	// // Create a new HTTP request
	// req, err := http.NewRequest("POST", *filesets.Nexusurl+"/service/rest/v1/components?repository=fls-aflm", nil)
	// if err != nil {
	// 	log.Panic("Failed to create HTTP request:", err)
	// }

	// // Set basic authentication header
	// req.SetBasicAuth(*filesets.Nexususr, *filesets.Nexuspass)

	// // Create a multipart form
	// body := &bytes.Buffer{}
	// writer := multipart.NewWriter(body)

	// // Add the directory field to the form
	// err = writer.WriteField("raw.directory", GetRepoPath(team, appname, envdc, version, region))
	// if err != nil {
	// 	log.Panic("Failed to WriteField raw.directory:", err)
	// 	// panic(err)
	// }

	// err = writer.WriteField("raw.asset1.filename", filename)
	// if err != nil {
	// 	log.Panic("Failed to WriteField raw.asset1.filename:", err)
	// 	// panic(err)
	// }

	// // Create a file part for the file field
	// filePart, err := writer.CreateFormFile("raw.asset1", filename)
	// if err != nil {
	// 	log.Panic("Failed to CreateFormFile raw.asset1:", err)
	// }
	// // filePart.Write(bodybytes)
	// // Copy the file content to the file part
	// stringreader := strings.NewReader(bodystring)
	// _, err = io.Copy(filePart, stringreader)
	// if err != nil {
	// 	log.Panic("Failed to copy io:", err)
	// }

	// // Close the multipart writer to finalize the form data
	// err = writer.Close()
	// if err != nil {
	// 	log.Panic("Failed to close writer:", err)
	// 	// panic(err)
	// }

	// // Set the request body as the form data
	// req.Body = io.NopCloser(body)
	// req.Header.Set("Content-Type", writer.FormDataContentType())

	// // Create an HTTP client and make the request
	// client := http.DefaultClient
	// resp, err := client.Do(req)
	// if err != nil || (resp.StatusCode != 204 && resp.StatusCode != 200) {
	// 	log.Panic("Failed to do HTTP request:", err)
	// 	// panic(err)
	// }
	// defer resp.Body.Close()

	// // Print the response status and body
	// log.Info("Response Status:", resp.Status)
	// log.Info("Response Body:", resp.Body)
	// // fmt.Println("Response Status:", resp.Status)
	// // fmt.Println("Response Body:")
	// // io.Copy(os.Stdout, resp.Body)
}

func WriteToAppPath(path, filename, appname string, configFileContent map[string]interface{}, env string, c context.Context) (string, error) {

	log := logagent.RecreateWith(c)
	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// configFilePath := path
	// var filenamestr = "application-" + env + ".yml"

	log.Print(filename)

	writeContent := convertops.ConvertStrMapToYaml(&configFileContent, c)

	// err := Write(configFilePath, filenamestr, writeContent)
	// _, err := rediscli.Do("HSET", "confsolver-"+appname, filename, writeContent)
	// rediscli.Do("EXPIRE", "confsolver-"+appname, 60*10)
	kvmap := cacheops.GetKVMap()
	prefix := "confsolver-" + appname
	kvmap.SetkWithPrefixDuration(prefix, filename, writeContent, 10*60, c)
	// var err error
	// if ok, rediscli := redisops.Health(c); ok {

	// 	key := "confrepo-" + appname
	// 	_, err = rediscli.Pipelined(c, func(pipe redis.Pipeliner) error {
	// 		pipe.HSet(c, key, filename, writeContent)
	// 		pipe.Expire(c, key, 10*time.Minute)
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		log.Error(err)
	// 		log.Error("redis set error")
	// 	}
	// } else {
	// 	err = fmt.Errorf("redis connection is not healthy")
	// }

	return writeContent, nil
}

// 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, c context.Context) []string {
	log := logagent.RecreateWith(c)
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

// 將任意 go 值用 gob 編碼，並透過 snappy 壓縮寫入（streaming）
func WriteGobSnappy(v interface{}) ([]byte, error) {

	var buff bytes.Buffer
	sw := snappy.NewBufferedWriter(&buff)
	defer sw.Close()
	enc := gob.NewEncoder(sw)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	err := sw.Close()
	bytes := buff.Bytes()
	return bytes, err
}

// 從 gob+snappy 檔案讀出並解碼到 v（v 必須是指標）
func ReadGobSnappy(byties []byte, out interface{}) error {
	buff := bytes.NewBuffer(byties)
	sr := snappy.NewReader(buff)
	dec := gob.NewDecoder(sr)
	return dec.Decode(out)
}
