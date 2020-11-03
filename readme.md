## 后端工程初始化

### 使用方法

1. 下载代码，将代码下载至$GOPATH/src/github.com/DoOR-Team/中
    ```
    mkdir -p $GOPATH/src/github.com/DoOR-Team/
    cd $GOPATH/src/github.com/DoOR-Team/
    git clone git@github.com:DoOR-Team/backend_generator.git
    ```
2.  编译该项目
    ```
    cd $GOPATH/src/github.com/DoOR-Team/backend_generator
    git pull && go build && go install
    ```
3. 生成一个项目，注意该项目只能存在于$GOPATH/src/github.com/DoOR-Team/，如果需要在其他目录，自行修改对应文件。APPNAME自行替换为自己的名字
    ```
    cd $GOPATH/src/github.com/DoOR-Team/
    backend_generator --name APPNAME
    ```
4. 进入生成的项目，并生成protobuf，APPNAME自行替换为自己的名字
    ```
    cd $GOPATH/src/github.com/DoOR-Team/APPNAME
    sh genprotos.sh
    go build
    ```

5. 如果希望该项目能够对接自动部署，需要在该项目下添加.deploy文件，文件内容参考如下：
   ```yaml
    deploy:
      # 如果不是前端服务，不设置subdomain
      # subdomain: test_deploy #部署后的subdomain，daily环境会添加_daily后缀
      # 企业微信机器人的webhook，对具体项目，创建特定的机器人，否则运维群会爆炸。
      wxboturl: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=bb18509f-2324-47dc-8672-64736b11ab8f
      deploy_type: k8s
    ```
    下面的已经废弃
    ```yaml
    deploy:
      port: 0 #线上部署服务端口号，0为随机
      subdomain: test_deploy #部署后的subdomain，daily环境会添加_daily后缀
      daily_port: 0 #日常测试环境服务的端口号，0为随机
      # 企业微信机器人的webhook，对具体项目，创建特定的机器人，否则运维群会爆炸。
      wxboturl: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=bb18509f-2324-47dc-8672-64736b11ab8f
    ```