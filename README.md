# [![PongHub](static/band.png)](https://health.ch3nyang.top)

<div align="center">

🌏 [Live Demo](https://health.ch3nyang.top) | 📖 [简体中文](README_CN.md)

</div>

## Introduction

PongHub is an open-source service status monitoring website designed to help users track and verify service availability. It supports:

- Non-intrusive monitoring
- One-click CI/CD deployment using GitHub Actions and GitHub Pages
- Multi-port checks for individual services
- Status code matching and response body regex matching
- Custom request bodies
- Customizable configurations such as check intervals, retry attempts, timeout durations, etc.

## Quick Start

1. Star and Fork [PongHub](https://github.com/WCY-dt/ponghub)

2. Modify the [`config.yaml`](config.yaml) file in the root directory to configure your service checks.

3. Modify the [`CNAME`](CNAME) file in the root directory to set your custom domain name.

   > [!NOTE]
   > If you do not need a custom domain, you can delete the `CNAME` file.

4. Commit and push your changes to your repository. GitHub Actions will automatically run and deploy to GitHub Pages and require no intervention.

> [!IMPORTANT]
> If GitHub Actions does not trigger automatically, you can manually trigger it once.

## Configuration Guide

The `config.yaml` file follows this format:

| Field                     | Type   | Description                                      | Required |
|---------------------------|--------|--------------------------------------------------|----------|
| `timeout`                 | Integer| Timeout for each request in seconds              | No       |
| `retry`                   | Integer| Number of retry attempts on request failure      | No       |
| `max_log_days`            | Integer| Number of days to retain logs; logs older than this will be deleted | No       |
| `services`                | Array  | List of services to monitor                      | Yes      |
| `services.name`           | String | Name of the service                              | Yes      |
| `services.health`         | Array  | Health check configurations for the service      | No       |
| `services.health.url`     | String | URL to check                                     | Yes      |
| `services.health.method`  | String | HTTP method (`GET`/`POST`/`PUT`)                 | No       |
| `services.health.status_code` | Integer | Expected HTTP status code (default `200`)       | No       |
| `services.health.response_regex` | String | Regex to match response body content            | No       |
| `services.health.body`    | String | Request body content, used only for `POST` requests | No       |
| `services.api`            | Array  | API check configurations, same format as above   | No       |

Here is an example configuration file:

```yaml
timeout: 5
retry: 2
max_log_days: 30
services:
  - name: "GitHub API"
    health:
      - url: "https://api.github.com"
    api:
      - url: "https://api.github.com/repos/wcy-dt/ponghub"
        method: "GET"
        status_code: 200
        response_regex: "full_name"
  - name: "Ch3nyang's  Websites"
    health:
      - url: "https://example.com/health"
        response_regex: "status"
      - url: "https://example.com/status"
        method: "POST"
        body: '{"key": "value"}'
```

> [!TIP]
> The `health` and `api` sections must have at least one entry. They are processed similarly, with this distinction made for future expansion.

## Disclaimer

[PongHub](https://github.com/WCY-dt/ponghub) is intended for personal learning and research only. The developers are not responsible for its usage or outcomes. Do not use it for commercial purposes or illegal activities.
