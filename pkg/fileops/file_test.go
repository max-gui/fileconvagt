package fileops

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/stretchr/testify/assert"
)

var tempdir, filename, filefullname, filecontent string

func setup() {
	tempdir = "testfolder"
	os.Mkdir(tempdir, 0777)
	// tempdir = dir
	log.Println(tempdir)
	filename = "testfile.yml"
	filefullname = tempdir + string(os.PathSeparator) + filename
}

func teardown() {
	os.RemoveAll(tempdir)
}

// func Test_Cases(t *testing.T) {
// 	// <setup code>
// 	// setup()

// 	t.Run("Read=Read", Test_Read)
// 	t.Run("Read=ReadFrom", Test_ReadFrom)
// 	t.Run("Write=Write", Test_Write)
// 	t.Run("Write=ExistedFile", Test_Write)
// 	t.Run("Write=WhiteToPath", Test_WhiteToPath)
// 	// <tear-down code>
// 	// teardown()
// }

func Test_Read(t *testing.T) {
	// tempdir = t.TempDir()
	// log.Println(tempdir)
	// filename := "testfile"
	// filefullname := tempdir + string(os.PathSeparator) + filename
	filecontent = "aaa"
	ioutil.WriteFile(filefullname, []byte(filecontent), 0644)

	str, err := Read(filefullname)
	// log.Println("fdsafsad")

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_Read result is:\n%s", str)
	// if err != nil {
	// 	t.Fatalf("read error: %s", err.Error())
	// } else if strings.Compare(str, filecontent) != 0 {
	// 	t.Fatalf("read content error! content should be:%s, get:%s", filecontent, str)
	// } else {
	// 	log.Printf("test is ok, read result is %s", str)
	// }

}

func Test_ReadFrom(t *testing.T) {
	// tempdir := t.TempDir()
	// log.Println(tempdir)
	// filename := "testfile"
	// filefullname := tempdir + string(os.PathSeparator) + filename
	c := context.Background()
	filecontent := "aaa"
	ioutil.WriteFile(filefullname, []byte(filecontent), 0644)

	f, _ := os.OpenFile(filefullname, os.O_RDONLY, 0644)
	// ReadFrom(f)
	defer f.Close()
	str, err := ReadFrom(f, c)
	// log.Println("fdsafsad")

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_ReadFrom result is:\n%s", str)

}

func Test_Write(t *testing.T) {
	// tempdir := t.TempDir()
	// log.Println(tempdir)
	// filename := "testfile"
	// filefullname := tempdir + string(os.PathSeparator) + filename
	c := context.Background()
	filecontent = "aaa"
	// test for writing new file
	err := Write(tempdir, filename, filecontent, c)
	// if err != nil {
	// 	log.Printf("error is %s", err.Error())
	// }

	assert.NoError(t, err, "Write is ok")

	var str string
	str, err = Read(filefullname)

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_Write result is:\n%s", str)
	// assert(str, err, t)

	// test for writing existed file
	// err = Write(tempdir, filename, filecontent)
	// if err != nil {
	// 	log.Printf("error is %s", err.Error())
	// }

	// str, err = Read(filefullname)
	// // log.Println("fdsafsad")
	// if err != nil {
	// 	t.Fatalf("read error: %s", err.Error())
	// } else if strings.Compare(str, filecontent) != 0 {
	// 	t.Fatalf("read content error! content should be:%s, get:%s", filecontent, str)
	// } else {
	// 	log.Printf("test is ok, read result is %s", str)
	// }

}

// func assert(str string, err error, t *testing.T) {
// 	if err != nil {
// 		t.Fatalf("read error: %s", err.Error())
// 	} else if strings.Compare(str, filecontent) != 0 {
// 		t.Fatalf("read content error! content should be:%s, get:%s", filecontent, str)
// 	} else {
// 		log.Println("test is ok, read result is:")
// 		fmt.Println(str)
// 	}
// }

func Test_WhiteToPath(t *testing.T) {
	// tempdir := t.TempDir()
	// log.Println(tempdir)
	// filename := "testfile.yml"
	// filefullname := tempdir + string(os.PathSeparator) + filename
	//test for no file
	c := context.Background()
	mapcontent := map[string]interface{}{
		"af-arch": 1,
		"b":       2,
	}
	writecontent := map[string]interface{}{
		"af-arch": 1,
		"b":       2,
	}

	env := "test"

	dotIndex := strings.LastIndex(filefullname, ".")
	var filefullnameext = filefullname[0:dotIndex] + "-" + env + ".yml"

	filecontent = convertops.ConvertStrMapToYaml(&mapcontent, c)
	// test for writing new file
	WriteToPath(filefullname, writecontent, env, c)

	var str string
	str, err := Read(filefullnameext)

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_WhiteToPath result is:\n%s", str)

	var f0 = func(archmap map[string]interface{}) {
		writecontent["af-arch"] = archmap
		writecontent["c"] = 2
		mapcontent["af-arch"] = archmap
		delete(writecontent, "b")
	}
	//test for existed with arch

	f0(map[string]interface{}{
		"resource": 1,
	})

	filecontent = convertops.ConvertStrMapToYaml(&mapcontent, c)
	WriteToPath(filefullname, writecontent, env, c)
	str, err = Read(filefullnameext)

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_WhiteToPath result is:\n%s", str)

	//test for existed with arch:resource

	f0(map[string]interface{}{
		"resource": map[string]interface{}{
			"test": 1,
		},
	})

	filecontent = convertops.ConvertStrMapToYaml(&mapcontent, c)
	WriteToPath(filefullname, writecontent, env, c)
	str, err = Read(filefullnameext)

	assert.NoError(t, err, "read is ok")
	assert.Equal(t, filecontent, str)
	log.Printf("Test_WhiteToPath result is:\n%s", str)

}

func TestMain(m *testing.M) {
	setup()
	// constset.StartupInit()
	// sendconfig2consul()
	// configgen.Getconfig = getTestConfig

	exitCode := m.Run()
	teardown()
	// // 退出
	os.Exit(exitCode)
}
