# GoFuzz
<div align="center">
  <img src="assets/logo.png" alt="GoFuzz Logo" width="200" />
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
go build -o gofuzz main.go
sudo mv gofuzz /usr/local/bin
```

## Verify installation:

```gofuzz -h```

## 📖 How to Use

Below is a summary of GoFuzz's options. Run gofuzz -h for details.

```
Usage of gofuzz:
  -H string
    	specify the file that contains headers [separated by line]
  -body string
    	specify POST request body (or file containing the body)
  -burp string
    	specify path to burp request
  -h	show usage
  -method string
    	specify the request method [POST, GET] (default "GET")
  -s string
    	specify a string to search for in response body 'Login Successful'
  -sc string
    	specify a status code(s) to output
  -t int
    	specify thread count [default: 3] (default 3)
  -timeout int
    	specify timeout in seconds [default 5] (default 5)
  -u string
    	specify the host url
  -w string
    	specify a wordlist used to fuzz
```

# 🚀 Quick Start
1. Fuzz Paths

```gofuzz -u https://example.com/FUZZ -w paths.txt```

2. Fuzz Query Parameters

```gofuzz -u https://example.com/api?param=FUZZ -w params.txt```

3. Fuzz POST Body

```gofuzz -method POST -u https://example.com -body '{"key":"FUZZ"}' -w payloads.txt```

3. Custom Headers	

```gofuzz -u https://example.com/FUZZ -H headers.txt```

4. Filter Status Codes	

```gofuzz -u https://example.com/FUZZ -sc 200,403```

5. "Grep" a string in respones body

```bash
gofuzz -u https://example.com -sc -body body.json -s 'Login Successful'
```

# 🎥 Demo Video
[Recording](https://github.com/user-attachments/assets/4d053735-9290-45e8-963c-14eb9f9221ec)
## ☕ Support

Like what you see? Your support keeps projects like GoFuzz going!
Buy Me a Coffee ☕

## 💻 Contributing

Contributions are always welcome! Whether it's reporting issues, suggesting features, or submitting pull requests, every bit helps improve GoFuzz.

## 📜 License

GoFuzz is licensed under the MIT License. See the LICENSE file for more information.

## ❤️ Acknowledgements

Inspired by the amazing work of <a href="https://github.com/xmendez/wfuzz">wfuzz</a>
