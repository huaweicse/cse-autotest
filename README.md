# cse-autotest
Auto test for CSE Go-Chassis/mesher


----------
This project provides 4 docker images:

- sdkat_consumer_gosdk
- sdkat_consumer_mesher
- sdkat_provider_gosdk
- sdkat_provider_mesher

## REST API

### Consumer
- `GET /{providerApi}`
#### Query params of consumer
```
protocol    # rest|highway, default: rest
provider    # provider name, default: sdkat_provider_gosdk
times       # call times, default: 1
```
### Example
```bash
curl http://localhost:8000/svc?times=3
```

## Implement

|App|Server End|Client End|
|:---:|:---:|:---:|
|sdkat_consumer_gosdk|Go|Go-Chassis|
|sdkat_consumer_mesher|Go|Go|
|sdkat_provider_gosdk|Go-Chassis||
|sdkat_provider_mesher|Go||

## Configuration

- On VM, please change the config files in the `conf` dir.

- In docker, please mount your customized config files into `/tmp` dir in the container.

## 如何测试
用于测试治理功能的接口，其典型的返回信息如下：
```json
{
    "Result": [
        {
            "num": 1,
            "time": "2019-01-04 15:05:17.330",
            "statusCode": 200,
            "provider": {
                "micro_service": {
                    "application": "sdkat",
                    "service_name": "sdkat_provider_gosdk",
                    "version": "3.0"
                },
                "instance_name": "desktop-0008_rest_0.0.0.0_9090",
                "instance_alias": "1"
            }
        },
        {
            "num": 2,
            "time": "2019-01-04 15:05:17.330",
            "statusCode": 200,
            "provider": {
                "micro_service": {
                    "application": "sdkat",
                    "service_name": "sdkat_provider_gosdk",
                    "version": "3.0"
                },
                "instance_name": "desktop-0008_rest_0.0.0.0_9090",
                "instance_alias": "1"
            }
        }
    ]
}
```

### 负载均衡。
首先发起20次调用，收集服务端实例个数。
- 轮询、随机。发起20次调用，根据实例个数，检查实例名排列是否为轮询、随机。
- session粘滞。发起20次调用，根据实例名列表，检查是否调用同一个实例。
- 权值。调用"/delayInstance/{instanceName}/{ms}"，恶化某个实例的通信质量。再发起调用，查看调用该恶化实例的次数是否减少。

### 容错
"/failTwice/{statusCode}"会失败两次，成功一次。开启同实例重试功能后，调用该接口应当可以避免错误。
"/failInstance/{instanceName}/{statusCode}"接口会在调用到{instanceName}对应的实例时失败，开启不同实例间重试功能后，调用该接口应当可以避免错误。

## 路由、灰度
实例信息中包含了版本信息，统计各个版本比例，即可确认路由、灰度策略是否生效。

## 限流
根据时间统计tps

## 错误注入
根据状态码和错误信息

## 熔断
设定熔断阈值，调用failTwice（错误率2/3）接口检查是否生效。用delay接口可检查超时是否生效。

##黑白名单
根据错误信息