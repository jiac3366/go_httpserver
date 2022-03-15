# Apiserver deploy in kubernetes cluster

## Hi!What is this?

The project realizes the functions of mainstream httpserver based on golang / gin, including elegant start 
and termination of services, service activation and QoS service quality assurance mechanism, message parsing, 
TLS encrypted communication, configuration and code separation, etc:

## Quick Start

```shell
git clone git@github.com:jiac3366/go_httpserver.git
cd go_httpserver/k8s/
k create -f httpserver-secret.yaml
k create -f web-deployment.yaml
k create -f web-service.yaml
k get svc # get the NodePort, the example is 30658.
```

```shell
# visit service NodePort: [k8sIP]:[NodePort]/url

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
```

```shell
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



## Features
- v3.0: [基于Istio实现流量管理，灰度发布，混沌测试](readme_docs/v3.0.md)

- v2.1: 与Prometheus结合实现监控

- v1.0: [优雅终止 / 优雅启动 / 探活 / QoS / 服务发布 ](readme_docs/v1.0.md)
 

## Todo:
- JWT
- CI/CD
