# httpserver deploy in kubernetes cluster

## Hi!What is this?

The project realizes the functions of mainstream httpserver based on golang / gin, including elegant start 
and termination of services, service activation and QoS service quality assurance mechanism, message parsing, 
TLS encrypted communication, configuration and code separation, etc:

## Usage

```shell
git clone git@github.com:jiac3366/go_httpserver.git
cd go_httpserver/k8s/
k create -f httpserver-secret.yaml
k create -f web-deployment.yaml
k create -f web-service.yaml
k get svc # 取得NodePort
```

```shell
# Example

curl --location --request POST 'http://192.168.34.3:30658/api/orders' \
--header 'Authorization: Basic amlhYzozMzY2' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "id":666,
    "partner": {
        "name": "jiac",
        "age": 22,
        "email": "463045792@qq.com"
    }
}'

> {
    "message": "OK"
}


curl --location --request GET 'http://192.168.34.3:30658/api/orders' \
--header 'Authorization: Basic amlhYzozMzY2'

>[
    {
        "id": 666,
        "partner": {
            "name": "jiac",
            "age": 22,
            "company": "",
            "email": "463045792@qq.com"
        },
        "description": ""
    }
]
```



## Feature

- 优雅终止

  ```go
  <-quit
  log.Println("Shutting down server...")
  // The context is used to inform the server it has 5 seconds to finish
  // the request it is currently handling
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  //defer cancel()
  defer func() {
      // extra handling here
      cancel()
  }()
  
  if err := productRouter.Shutdown(ctx); err != nil {
      log.Fatalf("productRouter forced to shutdown:%+v", err)
  }
  log.Println("productRouter Exited Properly")
  
  if err := healthyServer.Shutdown(ctx); err != nil {
      log.Fatalf("healthyServer forced to shutdown:%+v", err)
  }
  log.Println("healthyServer Exited Properly")
  ```

- 优雅启动 / 探活 / QoS 保证

  ```yaml
            livenessProbe: # 存活检查，检查容器是否正常，不正常则重启实例
              httpGet: # HTTP请求检查方法
                path: /healthz # 请求路径
                port: 8088 # 检查端口
                scheme: HTTP # 检查协议
              initialDelaySeconds: 5 # 启动延时，容器延时启动健康检查的时间
              periodSeconds: 10 # 间隔时间，进行健康检查的时间间隔
              successThreshold: 1 # 健康阈值，表示后端容器从失败到成功的连续健康检查成功次数
              failureThreshold: 1 # 不健康阈值，表示后端容器从成功到失败的连续健康检查成功次数
              timeoutSeconds: 3 # 响应超时，每次健康检查响应的最大超时时间
            readinessProbe: # 就绪检查，检查容器是否就绪，不就绪则停止转发流量到当前实例
              httpGet:
                path: /healthz
                port: 8088
                scheme: HTTP
              initialDelaySeconds: 5
              periodSeconds: 10
              successThreshold: 1
              failureThreshold: 1
              timeoutSeconds: 3
            startupProbe: # 启动探针，可以知道应用程序容器什么时候启动了
              failureThreshold: 10
              httpGet:
                path: /healthz
                port: 8088
                scheme: HTTP
              initialDelaySeconds: 1  #
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 3
  ```

  

- 资源需求

  ```yaml
            resources: # 资源需求
              limits: # limits用于设置容器使用资源的最大上限,避免异常情况下节点资源消耗过多
                cpu: "1" # 设置cpu limit，1核心 = 1000m
                memory: 1Gi # 设置memory limit，1G = 1024Mi
              requests: # requests用于预分配资源,当集群中的节点没有request所要求的资源数量时,容器会创建失败
                cpu: 250m # 设置cpu request
                memory: 500Mi # 设置memory request
  ```

  

- 日志等级 / 代码配置分离(从YAML中传入debug参数指定日志等级)
  ![image-20211130080451332](https://cdn.jsdelivr.net/gh/jiac3366/image-host@master/httpserver/2299500dacebaf4028b0015266fb924.pmtjq7y8hq8.png)
  
- 身份授权（基于Secret的basic auth）
  ![image-20211130080451337](https://cdn.jsdelivr.net/gh/jiac3366/image-host@master/httpserver/image-20211130080451337.26njmyu47mlc.png)

- 对内外发布

  - Service
    - 基于service ClusterIP
      ![image-20211130075828190](https://cdn.jsdelivr.net/gh/jiac3366/image-host@master/httpserver/959419e376223ca57ce31df0e69ee03.5ecrpo120rc0.png)
    - 基于service NodePort
      ![image-20211130081828163](https://cdn.jsdelivr.net/gh/jiac3366/image-host@master/httpserver/6eec6296b4e14f507965220791edbd2.2nf87vzu3b00.png)
  - Ingress

  

======待完善=======

- HTTPS
- 更完善的授权机制（JWT）
- CICD
