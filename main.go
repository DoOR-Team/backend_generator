package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/DoOR-Team/goutils/command"
	"github.com/DoOR-Team/goutils/log"
)

var appName = flag.String("name", "", "app名称")

var Urls = []string{
	"http://xuelang-algo-test.oss-cn-hangzhou.aliyuncs.com/test-be.tar.gz",
	"https://github.com/jpbirdy/backend-generator/raw/master/test-be.tar.gz",
}

func downloadBaseFile() error {
	for _, url := range Urls {
		res, err := http.Get(url)
		if err != nil {
			log.Printf("url: %s, 下载失败%s，更换下载源", url, err.Error())
			continue
		}
		f, err := os.Create(fmt.Sprintf("test-be.tar.gz"))
		if err != nil {
			log.Printf("url: %s, 下载失败%s，更换下载源", url, err.Error())
			continue
		}
		io.Copy(f, res.Body)
		return nil
	}
	return errors.New("下载失败")
}

//压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压 tar.gz
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		if file != nil {
			io.Copy(file, tr)
		}
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), os.ModePerm)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(name, "/") {
		return nil, nil
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
}

func strFirstToUpper(str string) string {
	str = strings.Replace(str, "-", "_", -1)
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		if temp[y] == "be" {
			continue
		}
		vv := []rune(temp[y])
		//if y != 0 {
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				if vv[i] > rune('Z') {
					vv[i] -= 32
				}
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
		//}
	}
	return upperStr
}

func replaceFileString(fileName string, old string, new string) {
	in, _ := os.Open(fileName)
	out, _ := os.OpenFile("tmp", os.O_RDWR|os.O_CREATE, 0766)

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err:", err)
			os.Exit(-1)
		}
		newLine := strings.Replace(string(line), old, new, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			fmt.Println("write to file fail:", err)
			os.Exit(-1)
		}
		index++
	}
	os.Remove(fileName)
	os.Rename("tmp", fileName)
	fmt.Printf("%sFINISH!\n", fileName)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	flag.Parse()

	// 参数非空校验
	if *appName == "" {
		log.Println("请输入想创建的app名称\nbackend_generator --name APPNAME")
		return
	}

	initForbiddenChar()

	if !checkName(*appName) {
		log.Println("app名称不合法")
		log.Println(`名称中不应当出现- @ ! # $ % ^ & * () [] {} | \ ; : \ ' ’ ， 。 《 》 < > · ~ 。`)
		log.Println(`连接符请使用_（下划线）`)
		return
	}

	var err error
	//
	//err := os.Mkdir(*appName, os.ModePerm)
	//if err != nil {
	//	log.Printf("文件夹创建错误，%s\n", err.Error())
	//	return
	//}

	if !PathExists("test-be") {
		err = downloadBaseFile()
		if err != nil {
			log.Printf("下载错误，%s\n", err.Error())
			return
		}

		err = DeCompress("test-be.tar.gz", "./")
		if err != nil {
			log.Printf("解压缩错误，%s\n", err.Error())
			panic(err)
		}
		defer os.Remove("test-be.tar.gz")

		err = os.Rename("test-be", *appName)
		if err != nil {
			panic(err)
		}

	} else {
		log.Println("test-be存在，直接copy")
		err, logStr := command.Shellout("cp -r test-be " + *appName)
		if err != nil {
			panic(err)
		}
		log.Println(logStr)
	}

	k8sName := strings.ReplaceAll(*appName, "_", "-")

	className := strFirstToUpper(*appName)
	replaceFileString(fmt.Sprintf("%s/protos/service.proto", *appName), "Demo", className)
	replaceFileString(fmt.Sprintf("%s/service/service.go", *appName), "Demo", className)
	replaceFileString(fmt.Sprintf("%s/service/service.go", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/service/service_test.go", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/tables/table_test.go", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/Dockerfile", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/main.go", *appName), "Demo", className)
	replaceFileString(fmt.Sprintf("%s/main.go", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/.gitignore", *appName), "test-be", *appName)
	replaceFileString(fmt.Sprintf("%s/configs/daily/app.yaml", *appName), "test-be", k8sName)
	replaceFileString(fmt.Sprintf("%s/configs/daily/app.yaml", *appName), "test_be", *appName)
	replaceFileString(fmt.Sprintf("%s/configs/production/app.yaml", *appName), "test-be", k8sName)
	replaceFileString(fmt.Sprintf("%s/configs/production/app.yaml", *appName), "test_be", *appName)

	log.Println("Generator success")
}

var forbiddenStrings = `-@!#$%^&*()[]{}|\;:/'’<>·~.?"`

var forbiddenCharMap map[uint8]bool

func initForbiddenChar() {
	forbiddenCharMap = make(map[uint8]bool)
	for i := 0; i < len(forbiddenStrings); i++ {
		forbiddenCharMap[forbiddenStrings[i]] = true
	}
}

func checkName(s string) bool {
	for i := 0; i < len(s); i++ {
		if _, ok := forbiddenCharMap[s[i]]; ok {
			return false
		}
	}
	return true
}
