# RootProxy

```
[HTB]  ██████╗░░█████╗░░█████╗░████████╗
      ██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝
      ██████╔╝██║░░██║██║░░██║░░░██║░░░
      ██╔══██╗██║░░██║██║░░██║░░░██║░░░
      ██║░░██║╚█████╔╝╚█████╔╝░░░██║░░░
      ╚═╝░░╚═╝░╚════╝░░╚════╝░░░░╚═╝░░░

RootProxy - Terminal-Based Proxy Management System
```

A terminal-first proxy management system for penetration testers, built in Go with a Bubble Tea TUI.

Repository: https://github.com/lily0ng/RootProxy.git

## Features (Roadmap)

- Proxy management (HTTP/HTTPS/SOCKS4/SOCKS5)
- Bulk import/export
- Proxy testing (latency/connectivity)
- Auto-rotation (planned)
- Certificate manager (import/self-signed generation)
- Profiles (save/switch proxy chains)
- Routing rules (planned)
- Proxy chains (planned)
- Monitoring & analytics (planned)
- Security features (DoH/DoT, leak protection, kill switch) (planned)
- Integrations (Burp/Nmap/Metasploit) (planned)

## Quick Start

### Requirements

- Go 1.21+

### Run (TUI)

From the project root:

```bash
go run ./cmd
```

Run with a specific profile:

```bash
go run ./cmd --profile htb-pentest
```

### Run (TUI + REST API)

Start the optional API server:

```bash
go run ./cmd --api 127.0.0.1:8081
```

## TUI Keys

- `1..0` switch screens
- `Ctrl+P` proxy dashboard
- `Ctrl+C` certificate manager
- `Ctrl+R` routing rules
- `Ctrl+M` monitoring
- `Ctrl+S` settings
- `F1` help
- `F4` test active proxy (TCP connectivity + latency)
- `F10` / `q` / `Esc` exit

## API (v1)

When started with `--api`, RootProxy exposes a small REST surface intended for integrations.

- `GET /api/v1/status`
- `GET /api/v1/proxy/list`
- `POST /api/v1/proxy/add`
- `POST /api/v1/profile/switch`

Example:

```bash
curl -s http://127.0.0.1:8081/api/v1/status
```

## Project Structure

```
RootProxy/
├── cmd/
│   └── main.go
├── internal/
│   ├── cert/
│   │   ├── generator.go
│   │   └── manager.go
│   ├── config/
│   │   ├── profiles.go
│   │   └── settings.go
│   ├── proxy/
│   │   ├── chain.go
│   │   ├── manager.go
│   │   ├── types.go
│   │   └── validator.go
│   ├── rootproxy/
│   │   └── app.go
│   └── tui/
│       ├── model.go
│       ├── screens.go
│       └── theme.go
└── pkg/
    ├── api/
    │   ├── routes.go
    │   └── server.go
    ├── integrations/
    │   └── integrations.go
    └── scripting/
        └── scripting.go
```

## Notes

- This repository currently provides a working scaffold (TUI + API skeleton) with core domain primitives.
- The next milestones are persistence (profiles/proxies on disk), routing rules, chain execution, and a richer proxy editing UI.

## License

TBD
