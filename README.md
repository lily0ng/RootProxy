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
 
 RootProxy is an HTB-inspired operations console for managing proxies, profiles, and (eventually) routing/chaining/monitoring from a single TUI.
 
 Current build status: scaffold is runnable and provides:
 
 - Proxy list + active proxy selection placeholder
 - Connectivity test (TCP dial + latency)
 - Profiles store (in-memory)
 - Certificate store + self-signed certificate generation utilities
 - Optional REST API skeleton for tool integrations

 ## Core Features

 - [x] Proxy management (HTTP/HTTPS/SOCKS4/SOCKS5)
 - [ ] Bulk import/export
 - [x] Proxy testing (latency/connectivity)
 - [ ] Auto-rotation (planned)
 - [x] Certificate manager (import/self-signed generation)
 - [x] Profiles (save/switch proxy chains)
 - [ ] Routing rules (planned)
 - [ ] Proxy chains (planned)
 - [ ] Monitoring & analytics (planned)
 - [ ] Security features (DoH/DoT, leak protection, kill switch) (planned)
 - [ ] Integrations (Burp/Nmap/Metasploit) (planned)

 ## Planned Features (Details)

 ### Routing Rules (planned)

 - **Domain-based routing**
   - `*.corp.local`, `*.hackthebox.eu` style glob matching
   - Optional regex rules (advanced)
 - **IP/Country-based routing**
   - GeoIP-driven decisions (country/ASN) for egress control
 - **Application-specific routing**
   - Route traffic based on process/app (Windows/Linux strategies)
 - **Actions**
   - `direct`, `proxy`, `chain`, `profile`
 - **Failover-aware routing**
   - Automatic fallback when a proxy/chain fails health checks

 ### Proxy Chains (planned)

 - **Multi-hop chains**
   - Up to 5 hops
 - **Chain validation**
   - Hop-by-hop connectivity test
   - End-to-end test and latency scoring
 - **Visual chain builder (TUI)**
   - Create/edit/reorder hops interactively
 - **Random chain generator**
   - Build randomized chains from tagged/filtered proxy pools

 ### Monitoring & Analytics (planned)

 - **Per-proxy health metrics**
   - Success/failure rates, last-seen, rolling latency
 - **Bandwidth tracking**
   - Traffic totals per proxy/chain/profile
 - **Logs with search**
   - Structured events (connect, failover, rotate, errors)
 - **TUI dashboards**
   - Live status panels and trend views

 ### Security Features (planned)

 - **DNS security**
   - DNS-over-HTTPS (DoH) / DNS-over-TLS (DoT) options
 - **Leak protection**
   - DNS leak checks and rule-based mitigation
   - Optional “deny-by-default” egress rules
 - **Kill switch**
   - If proxy/chain drops, block outbound traffic (platform-dependent)
 - **Certificate / MITM support**
   - Better workflows for CA install/validation for intercept proxies

 ### Integrations (planned)

 - **Burp Suite**
   - Quick profile generation (HTTP proxy + CA handling)
 - **Nmap**
   - Helpers to run scans via proxy (where supported) and manage configs
 - **Metasploit**
   - Routing/chain presets and listener-friendly profiles
 
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
 
 - `GET /api/v1/status`
 - `GET /api/v1/proxy/list`
 - `POST /api/v1/proxy/add`
 - `POST /api/v1/profile/switch`
 
 Example:
 
 ```bash
 curl -s http://127.0.0.1:8081/api/v1/status
 ```
 
 ## Features (Roadmap)
 
 - Proxy management: add/remove/edit, import/export, validation
 - Certificate management: CA import/export, self-signed, MITM support
 - Profiles: save/switch/share (encrypted), schedule switching
 - Advanced routing: domain/country/app-based, load balancing, failover
 - Proxy chains: up to 5 hops, visual chain builder, random chain generator
 - Monitoring: success/failure rate, bandwidth usage, log search
 - Security: DoH/DoT, leak protection, kill switch
 - Integrations: Burp, Nmap, Metasploit routing
 
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
