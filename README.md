 <div align="center">
 
 <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%" />
 
 <img src="https://github.com/lily0ng.png" width="180" height="180" alt="lily0ng" />
 
 <h1>
   <img src="https://readme-typing-svg.demolab.com?font=Fira+Code&weight=800&size=34&duration=2500&pause=800&color=667EEA&center=true&vCenter=true&width=700&lines=RootProxy;Terminal-Based+Proxy+Management;Bubble+Tea+TUI+%2B+REST+API;HTB-Inspired+Ops+Console" alt="Typing SVG" />
 </h1>
 
 <p>
 Advanced proxy management with a GUI-like terminal interface. Built for offensive operations, lab workflows, and repeatable proxy hygiene.
 </p>
 
 <p>
   <a href="https://github.com/lily0ng/RootProxy">
     <img src="https://img.shields.io/badge/Repo-GitHub-0D1117?style=for-the-badge&logo=github&logoColor=white" alt="Repo" />
   </a>
   <img src="https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
   <img src="https://img.shields.io/badge/TUI-Bubble%20Tea-7C3AED?style=for-the-badge" alt="TUI" />
   <img src="https://img.shields.io/badge/API-REST-22C55E?style=for-the-badge" alt="API" />
 </p>
 
 <p>
   <code>git clone https://github.com/lily0ng/RootProxy.git</code>
 </p>
 
 <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%" />
 
 </div>
 
 ## Overview

RootProxy is an HTB-inspired operations console for managing proxies, profiles, routing rules, chains, and operational utilities from a single TUI.

Current build status: runnable and provides:

- Proxy store (add/update/remove, active selection)
- Proxy connectivity testing (TCP dial + latency)
- Proxy import/export (JSON + text)
- Profiles store (in-memory) with per-profile chain lists
- Chain store (up to 5 hops)
- Routing rules store (domain glob/suffix + CIDR)
- Rotation (round-robin/random) via API
- Certificate store + self-signed certificate generation utilities
- Monitoring metrics store (records proxy test results)
- Minimal integrations helpers (Burp env export, proxychains.conf generator)
- Optional REST API server for tool integrations

 ## Core Features

- [x] Proxy management (HTTP/HTTPS/SOCKS4/SOCKS5)
- [x] Bulk import/export (JSON + text)
- [x] Proxy testing (latency/connectivity)
- [x] Auto-rotation (API-triggered: round-robin/random)
- [x] Certificate manager (import/self-signed generation)
- [x] Profiles (save/switch proxy chains)
- [x] Routing rules store + API
- [x] Proxy chains store + API
- [x] Monitoring metrics store + API
- [x] Security settings store + API
- [x] Integrations helpers (Burp env, proxychains.conf)

 ## Quick Start
 
 ### Requirements
 
 - Go 1.21+
 
 ### Run (TUI)
 
 ```bash
 go run ./cmd
 ```
 
 Run with a specific profile:
 
 ```bash
 go run ./cmd --profile htb-pentest
 ```
 
 ### Run (TUI + REST API)
 
 ```bash
 go run ./cmd --api 127.0.0.1:8081
 ```
 
 ### Run (REST API only / headless)
 
 ```bash
 go run ./cmd --api 127.0.0.1:8081 --headless
 ```
 
 ## TUI Hotkeys

- `1..0` switch screens
- `Ctrl+P` proxy dashboard
- `Ctrl+C` certificate manager
- `Ctrl+R` routing rules
- `Ctrl+M` monitoring
- `Ctrl+S` settings
- `F1` help
- `F4` test active proxy (connectivity + latency)
- `F10` / `q` / `Esc` exit
 
 ## API (v1)

When started with `--api`, RootProxy exposes a minimal REST surface intended for integrations.

Note: `--api` runs alongside the TUI by default. If you want an API-only process (recommended for scripting/curl), use `--headless`.

- `GET /api/v1/status`
- `GET /api/v1/proxy/list`
- `GET /api/v1/proxy/active`
- `POST /api/v1/proxy/active`
- `POST /api/v1/proxy/add`
- `POST /api/v1/proxy/update/{id}`
- `DELETE /api/v1/proxy/remove/{id}`
- `POST /api/v1/proxy/test?name=<proxy>&timeout_ms=<ms>`
- `GET /api/v1/proxy/export?format=json|text`
- `POST /api/v1/proxy/import?format=json|text`
- `POST /api/v1/profile/switch`
- `GET /api/v1/profile/list`
- `POST /api/v1/profile/upsert`
- `GET /api/v1/chain/list`
- `POST /api/v1/chain/upsert`
- `DELETE /api/v1/chain/remove/{name}`
- `GET /api/v1/routing/list`
- `POST /api/v1/routing/upsert`
- `DELETE /api/v1/routing/remove/{id}`
- `POST /api/v1/rotation/rotate`
- `GET /api/v1/cert/list`
- `POST /api/v1/cert/add`
- `POST /api/v1/cert/generate_self_signed`
- `GET /api/v1/security/get`
- `POST /api/v1/security/set`
- `GET /api/v1/monitoring/metrics`
- `GET /api/v1/monitoring/started`
- `GET /api/v1/integrations/burp/env`
- `GET /api/v1/integrations/proxychains/conf?profile=<name>`
 
 Example:
 
 ```bash
 curl -s http://127.0.0.1:8081/api/v1/status
 ```
 
 More examples:
 
 ```bash
 # Add a proxy
 curl -s -X POST http://127.0.0.1:8081/api/v1/proxy/add \
   -H 'Content-Type: application/json' \
   -d '{"name":"Local-Burp","type":"http","host":"127.0.0.1","port":8080,"auth":"none"}'

 # Set active proxy
 curl -s -X POST http://127.0.0.1:8081/api/v1/proxy/active \
   -H 'Content-Type: application/json' \
   -d '{"name":"Local-Burp"}'

 # Test a proxy (or omit name= to test current active)
 curl -s -X POST 'http://127.0.0.1:8081/api/v1/proxy/test?name=Local-Burp&timeout_ms=3000'
 ```
 
 ## Project Structure
 
 ```
 RootProxy/
 ├── cmd/
 │   └── main.go
 ├── internal/
 │   ├── cert/
 │   ├── config/
 │   ├── proxy/
 │   ├── rootproxy/
 │   └── tui/
 └── pkg/
     ├── api/
     ├── integrations/
     └── scripting/
 ```
 
 ## Author
 
 - GitHub: https://github.com/lily0ng
 
 <div align="center">
 
 <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%" />
 
 </div>
 
 ## License
 
 TBD
