## 后端工程初始化

### 使用方法

1. 下载代码，将代码下载至$GOPATH/src/gitlab.xuelang.xyz/xuelang_algo/中
    ```
    mkdir -p $GOPATH/src/gitlab.hz-xuelang.xyz/xuelang_algo/
    cd $GOPATH/src/gitlab.hz-xuelang.xyz/xuelang_algo/
    git clone git@101.132.72.35:xuelang_algo/backend-generator.git
    ```
2.  编译该项目
    ```
    cd $GOPATH/src/gitlab.hz-xuelang.xyz/xuelang_algo/backend-generator
    go build && go install
    ```
3. 生成一个项目，注意该项目只能存在于$GOPATH/src/gitlab.xuelang.xyz/xuelang_algo/，如果需要在其他目录，自行修改对应文件。APPNAME自行替换为自己的名字
    ```
    cd $GOPATH/src/gitlab.hz-xuelang.xyz/xuelang_algo/
    backend-generator --name APPNAME
    ```
4. 进入生成的项目，并生成protobuf，APPNAME自行替换为自己的名字
    ```
    cd $GOPATH/src/gitlab.hz-xuelang.xyz/xuelang_algo/APPNAME
    sh genprotos.sh
    go build
    ```
