
# httpserver deploy in kubernetes cluster


## What is this?
The project realizes the functions of mainstream httpserver based on golang / gin, including elegant start 
and termination of services, service activation and QoS service quality assurance mechanism, message parsing, 
TLS encrypted communication, configuration and code separation, etc:
- 优雅启动
- 优雅终止
- 探活
- 资源需求
- QoS 保证
- HTTPS
- 代码与配置分离
- 日志等级
- 身份授权（基于ConfigMap）
