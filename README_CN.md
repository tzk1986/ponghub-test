[![PongHub](static/band.png)](https://health.ch3nyang.top)

## 简介

PongHub 是一个开源的服务状态检查网站，旨在帮助用户监控和验证服务的可用性。它支持

- 利用 GitHub Actions 和 GitHub Pages 一键 CI/CD 部署
- 支持单个服务的多端口检查
- 支持状态码匹配和响应体内容正则表达式匹配
- 支持自定义请求体
- 支持自定义检查间隔、重试次数、超时时间等配置

[Live Demo](https://health.ch3nyang.top)

## 快速开始

1. Star 并 Fork [PongHub](https://github.com/WCY-dt/ponghub)

2. 在你 Fork 的仓库中，进入 `Settings` -> `Secrets and variables` -> `Actions` 页面，添加以下 Secrets：

    - `DEPLOY_TOKEN`: 你的 GitHub Personal Access Token，至少需要 `repo` 和 `workflow` 权限（[获取位置](https://github.com/settings/tokens)）
    - `DEPLOY_CNAME`: 你希望使用的自定义域名，例如 `health.ch3nyang.top`
    - `GH_USERNAME`: 你的 GitHub 用户名，例如 `WCY-dt`

    ![设置 secrets](static/step-secret.png)

3. 修改根目录下的 [`config.yaml`](config.yaml) 文件，配置你的服务检查项。

4. 提交修改并推送到你的仓库，GitHub Actions 将自动运行并部署到 GitHub Pages。后期将自动更新，无需干预。

## 配置说明

配置文件 `config.yaml` 的格式如下：

```yaml
timeout: 5
retry: 2
max_log_days: 30
services:
  - name: "GitHub API"
    health:
      - url: "https://api.github.com"
        method: "GET"
        status_code: 200
    api:
      - url: "https://api.github.com/repos/wcy-dt/ponghub"
        method: "GET"
        status_code: 200
        response_regex: "full_name"
  - name: "Ch3nyang's  Websites"
    health:
      - url: "https://example.com/health"
        method: "GET"
        status_code: 200
        response_regex: "status"
      - url: "https://example.com/status"
        method: "POST"
        body: '{"key": "value"}'
```

- `timeout`: 每次请求的超时时间，单位为秒
- `retry`: 请求失败时的重试次数
- `max_log_days`: 日志保留天数，超过此天数的日志将被删除
- `services`: **[可选]** 服务列表
  - `name`: 服务名称
  - `health`: **[可选]** 健康检查配置列表
    - `url`: 检查的 URL
    - `method`: HTTP 方法（GET、POST 等）
    - `status_code`: **[可选]** 期望的 HTTP 状态码
    - `response_regex`: **[可选]** 响应体内容的正则表达式匹配
    - `body`: **[可选]** 请求体内容，仅在 POST 请求时使用
  - `api`: **[可选]** API 检查配置列表，格式同上。

## 免责声明

[PongHub](https://github.com/WCY-dt/ponghub) 仅用于个人学习和研究，不对程序的使用行为或结果负责。请勿将其用于商业用途或非法活动。
