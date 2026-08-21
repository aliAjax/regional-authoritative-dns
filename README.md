# 11-regional-authoritative-dns

可运行的纯Go权威DNS与DNSSEC多区域解析服务。控制面提供Zone创建、版本更新、发布、冻结、回滚、记录查询和DNSSEC验证；数据面提供UDP/TCP权威DNS解析。默认使用内存仓储，PostgreSQL迁移文件和接口已预留生产适配器。

## 启动

```sh
go run ./cmd/server
```

HTTP监听`:8080`，DNS监听UDP/TCP`:5353`。除健康检查、就绪检查和指标外的HTTP接口需要`X-API-Key: dev-api-key`。

## 快速验证

```sh
curl -s -X POST localhost:8080/v1/zones -H 'X-API-Key: dev-api-key' -H 'Content-Type: application/json' -d @examples/zone.json
curl -s localhost:8080/healthz
dig @127.0.0.1 -p 5353 www.example.com A
dig @127.0.0.1 -p 5353 missing.example.com A
```

创建Zone后，将响应中的`id`用于`/v1/zones/{id}/publish`和`/v1/zones/{id}/records`。`/v1/dnssec/verify`返回签名/NSEC元数据验证摘要；真实生产部署应注入经过审查的密钥存储实现。

## 结构

`internal`按zone、record、resolver、transfer、dnssec、healthcheck、publication、auth、worker领域拆分，各领域保留domain/application/adapter/infrastructure边界。配置支持YAML和`DNS_*`环境变量覆盖；服务实现请求ID、API Key认证、请求体限制、超时、恢复、优雅停机、健康/就绪/Prometheus指标端点。

## 数据库与容器

`migrations/001_init.sql`提供Zone、记录、发布事件和传输运行表。使用`docker compose -f deploy/docker-compose.yml up --build`启动容器。
