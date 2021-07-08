package main

import (
	"bufio"
	"code.lyntime.com/common/goutils/command"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"code.lyntime.com/common/goutils/log"
)

var appName = flag.String("name", "", "app名称")
var template = flag.String("template", "test_be", "template")
var group = flag.String("group", "", "group")

var templateUrl = map[string]string{
	"test_be": "ssh://git@code.lyntime.com:30022/common/test_be.git",
	// "test_ae": "http://xuelang-algo-test.oss-cn-hangzhou.aliyuncs.com/test_ae.tar.gz",
}
var templatePath = map[string]string{
	"test_be": ".lyntime/code/common",
	// "test_ae": "http://xuelang-algo-test.oss-cn-hangzhou.aliyuncs.com/test_ae.tar.gz",
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
		// if y != 0 {
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
		// }
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
	log.Info(err)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type FilePath struct {
	AbsPath  string
	FileName string
	Path     string
}

func GetAllFile(pathname string, filenames *[]FilePath) error {
	rd, err := ioutil.ReadDir(pathname)
	path, _ := os.Open(pathname + "/")
	absPath, _ := filepath.Abs(filepath.Dir(path.Name()))

	path.Close()

	for _, fi := range rd {
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if fi.IsDir() {
			// fmt.Printf("[%s]\n", pathname+"/"+fi.Name())
			GetAllFile(pathname+"/"+fi.Name(), filenames)
		} else {
			// fmt.Println(absPath)
			*filenames = append(*filenames, FilePath{
				AbsPath:  absPath + "/" + fi.Name(),
				FileName: fi.Name(),
				Path:     pathname,
			})
		}
	}
	return err
}

func main() {
	flag.Parse()

	// 参数非空校验
	if *appName == "" {
		// log.Println("请输入想创建的app名称\nbackend_generator --name APPNAME --template [test_be|test_ae]")
		log.Println("请输入想创建的app名称\nbackend_generator --group GROUPNAME --name APPNAME --template [test_be|test_ae]")
		return
	}

	if *template == "" {
		*template = "test_be"
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
	// err := os.Mkdir(*appName, os.ModePerm)
	// if err != nil {
	//	log.Printf("文件夹创建错误，%s\n", err.Error())
	//	return
	// }


	user, err := user.Current()

	log.Info("目录", user.HomeDir + "/" + templatePath[*template]+"/"+*template)

	templateAbsPath := user.HomeDir + "/" + templatePath[*template]+"/"+*template


	if !PathExists(templateAbsPath) {

		log.Info("模板文件不存在，正在创建...")
		command := fmt.Sprintf(`
set +x
mkdir -p %s
cd %s
git clone %s
`, templatePath[*template], templatePath[*template], templateUrl[*template])

		log.Info("正在执行", command)
		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Error(err)
		}

	} else {
		log.Println(*template + "存在，更新，并copy")
		cmd := exec.Command("bash", "-c", "set +x\ncd " + templateAbsPath + "\ngit pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Error(err)
		}
	}

	err, logStr := command.Shellout(fmt.Sprintf(`
set +x
cp -r %s %s
cd %s
rm -rf .git
`, templateAbsPath , *appName, *appName))
	if err != nil {
		log.Error(logStr)
		panic(err)
	}
	log.Println(logStr)

	k8sName := strings.ReplaceAll(*appName, "_", "-")
	className := strFirstToUpper(*appName)

	templateK8sName := strings.ReplaceAll(*template, "_", "-")

	allFiles := make([]FilePath, 0)
	GetAllFile(*appName, &allFiles)
	for _, f := range allFiles {
		replaceFileString(f.AbsPath, "Demo", className)
		replaceFileString(f.AbsPath, *template, *appName)
		replaceFileString(f.AbsPath, templateK8sName, k8sName)
		replaceFileString(f.AbsPath, "common", *group)
	}

	replaceFileString(*appName+"/.gitignore", *template, *appName)

	os.Rename(*appName+"/protos/service.proto", *appName+"/protos/"+*appName+".proto")

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
