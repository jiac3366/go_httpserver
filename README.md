
# httpserver deploy in kubernetes cluster


## What is this?
The project realizes the functions of mainstream httpserver based on golang / gin, including elegant start 
and termination of services, service activation and QoS service quality assurance mechanism, message parsing, 
TLS encrypted communication, configuration and code separation, etc:
- 优雅启动

- 优雅终止

- 探活/QoS 保证

- 资源需求 

- 代码与配置分离

- 日志等级
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
