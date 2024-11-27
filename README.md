<div align="center">
  <img src="assets/logov2.png" alt="GoFuzz Logo" />
  <h1><strong>GoFuzz</strong></h1>
  <p>⚡ The lightweight, fast, and concurrent fuzzing tool for web application testing ⚡</p>
</div>


GoFuzz is a simple yet powerful fuzzing tool written in Go, designed for web application security testing. With its concurrent execution model, GoFuzz can fuzz endpoints, parameters, and headers efficiently, making it a perfect companion for penetration testers and bug hunters.

# ✨ Features:

* Blazing Fast: Leverages Go's concurrency to run multiple requests in parallel.
* Customizable: Supports custom headers, wordlists, and request bodies.
* Flexible Filtering: Filter results by status codes or response contents.
* Lightweight: Minimal dependencies and easy to use.

# 🔧 Installation
To get started, install GoFuzz using the following commands:

## Clone the repository
```git clone https://github.com/CharlesTheGreat77/GoFuzz```

## Navigate to the project directory
```cd GoFuzz```

## Build and install
```bash
go mod init gofuzz
go mod tidy
go build -o gofuzz main.go
sudo mv gofuzz /usr/local/bin
```

## Verify installation:

```gofuzz -h```

## 📖 How to Use

Below is a summary of GoFuzz's options. Run gofuzz -h for details.

```
╰─ GoFuzz is a simple yet powerful fuzzing tool written in Go,
  designed for web application security testing. With its concurrent execution model,
  GoFuzz can fuzz endpoints, parameters, and headers efficiently,
  making it a perfect companion for penetration testers and bug hunters.

Usage:
  GoFuzz [flags]
  GoFuzz [command]

Available Commands:
  completion    Generate the autocompletion script for the specified shell
  fuzz          fuzz URL parameters, HTTP request bodies, and headers
  openredirect  openredirect is used to fuzz parameters in URL(s)/Headers for open redirects
  pathtraversal Fuzzes URLs to detect path traversal vulnerabilities using a wordlist.

Flags:
  -h, --help   help for GoFuzz

Use "GoFuzz [command] --help" for more information about a command.
```


# 🚀 Quick Start
1. Fuzz Paths
```
╰─ gofuzz fuzz --help
The fuzz command is used to FUZZ parameters or points in an HTTP request.
  Fuzz can be used by specifing:

  method              -> specify a POST or GET method [Default: GET]
  u [url]             -> specify the target url
  H [headers]         -> specify path to custom headers file (seperate by line)
                         [accepts burpsuite request dumps]
  body                -> specify the body to FUZZ for a given POST request
  w [wordlist]        -> specify a wordlist used to FUZZ the given parameters
  c [status-codes] -> specify status code(s) [seperated by comma: 200,403]
  s [search]          -> specify a string to search/filter for in the response body
  N [NoSearch]        -> enable to output responses that do NOT contain the search string
  timeout             -> specify the time for timeout for each request
  t [threads]         -> specify the number of threads

  example(s):

  gofuzz fuzz -w common.txt -u https://example.com/FUZZ -sc 200,403
  gofuzz fuzz -w commont.txt -u https://ex.com/#id=FUZZ -sc 200
  gofuzz fuzz -w passw.txt -method POST -u https://example.com/ -s "login failed" -N -body body.json
  gofuzz fuzz -w list.txt -u https://example.com/FUZZ -s "passw" -sc 200 -timeout 5 -t 10

Usage:
  GoFuzz fuzz [flags]

Flags:
  -b, --body string            specify the request body to fuzz
  -H, --header strings         specify headers (comma-separated
  -h, --help                   help for fuzz
  -m, --method string          specify the HTTP method to use (default "GET")
  -N, --no-search              output responses that do NOT contain the search string
  -p, --proxies string         specify proxy(ies) (line-separated) [accepts file]
  -R, --request string         specify path to custom headers file (line-separated) [accepts burpsuite request dumps]
  -s, --search string          specify the string to search/filter in response body
  -c, --status-codes strings   specify status codes to filter by (comma-separated)
  -t, --threads int            specify the number of threads to use (default 5)
  -T, --timeout int            specify a timeout per request (default 10)
  -u, --url string             specify the target URL to fuzz
  -w, --wordlist string        specify path to a wordlist used for fuzzing
```
* the fuzz command can fuzz directories, logins, etc.



2. Open Redirect
```
╰─ gofuzz openredirect --help
The openredirect command is used to FUZZ parameters in URL(s)/Headers.
  openredirect can be used by specifing:

  method              -> specify a POST or GET method [Default: GET]
  u [url]             -> specify the target url
  H [headers]         -> specify path to custom headers file (seperate by line)
  body                -> specify the body to FUZZ for a given POST request
  w [wordlist]        -> specify a wordlist used to FUZZ the given parameters
  sc [status-code(s)] -> specify status code(s) [seperated by comma: 200,403]
  s [search]          -> specify a string to search/filter for in the response body
  N [NoSearch]        -> enable to output responses that do NOT contain the search string
  timeout             -> specify the time for timeout for each request
  t [threads]         -> specify the number of threads


  example(s):
  gofuzz openredirect -u https://example.com/path?redirect=FUZZ -w urls.txt
  gofuzz openredirect -burp requestDump.txt -w urls.txt
  gofuzz openredirect -u https://example.com/path?return=FUZZ -w urls.txt -sc 302,200
```
* openredirect tries to focus on the response code for redirection and location


3. Path Traversal
```
╰─ gofuzz pathtraversal --help
The "pathtraversal" command is used to test endpoints for potential path traversal vulnerabilities by injecting and incrementing "../" sequences into the URL.

Path traversal attacks attempt to access files and directories stored outside the web root directory by manipulating variables referencing file paths. This command automates the process of crafting such payloads and sending requests to identify misconfigurations or security flaws in the target application.

Key features include:
  - Automatic generation and injection of path traversal payloads (e.g., "../../", "../../../").
  - Support for appending or replacing parts of the URL to craft dynamic attacks.

Example usage:

  gofuzz pathtraversal -u https://example.com/path/to/file?parameter=FUZZ -w paths.txt
  gofuzz pathtraversal -u https://example.com/api/v1/files/FUZZ -w sensitive_files.txt -s "root:x:"
  gofuzz pathtraversal -u https://example.com/#parameter=FUZZ -sc 200 -timeout 5 -t 10
```
* pathtraversal is used to seemlessly increments '../' on a given word


# 🎥 Demo Video
[Demo](https://github.com/user-attachments/assets/473517da-8419-4fd5-ad5d-2e7d613c05a4)


## ☕ Support

Like what you see? Your support keeps projects like GoFuzz going!
Buy Me a Coffee ☕

## 💻 Contributing

Contributions are always welcome! Whether it's reporting issues, suggesting features, or submitting pull requests, every bit helps improve GoFuzz.

## 📜 License

GoFuzz is licensed under the MIT License. See the LICENSE file for more information.

## ❤️ Acknowledgements

Inspired by the amazing work of <a href="https://github.com/xmendez/wfuzz">wfuzz</a>
