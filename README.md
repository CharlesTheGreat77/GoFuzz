# GoFuzz

<div align="center">

  <img src="assets/logo.png" alt="logo" width="auto" height="auto" />
  <h1>GoFuzz</h1>
</div>
GoFuzz is a concurrent fuzzer written in go to assist in fuzzing parameters for web application testing. This project is inspired by <a href="https://github.com/xmendez/wfuzz">wfuzz</a>. It is by all means no replacement, but great for lightweight fuzzing.

# Install 🚀
| Prerequisite | Version |
|--------------|---------|
| Go           |  <=1.22 |

```bash
git clone https://github.com/CharlesTheGreat77/GoFuzz
cd GoFuzz
go mod init gofuzz
go mod tidy
go build -o gofuzz main.go
```

# Usage 🧠
```
Usage of ./gofuzz:
  -body string
        specify POST request body
  -custom-headers string
        specify the file that contains headers [seperated by line]
  -h    show usage
  -method string
        specify the request method [POST, GET] (default "GET")
  -threads int
        specify thread count [default: 3] (default 3)
  -timeout int
        specify timeout in seconds [default 5] (default 5)
  -url string
        specify the host url
  -wordlist string
        specify a wordlist used to fuzz
```

# Examples 🦫
Fuzz for Paths:
```bash
./gofuzz -url https://example.com/FUZZ -wordlist list.txt -custom-headers headers.txt
```

Fuzz parameters in URL:
```bash
./gofuzz -url https://example.com/api/search=FUZZ -timeout 3 -threads 10 -wordlist list.txt
```

Fuzz with POST Requests:
```bash
./gofuzz -url https://example.com/upload/file=FUZZ -method POST -body '{"test": "123456"}' -custom-headers headers.txt -timeout 6
```
* To FUZZ the body of the post, we can just use bash:
    ```bash
    cat payloads.txt | while read payload; do ./gofuzz -url https://example.com/api/upload -method POST -body $payload -custom-headers headers.txt -timeout 6; done
    ```
## Video Example
[Recording](https://github.com/user-attachments/assets/4d053735-9290-45e8-963c-14eb9f9221ec)




# Coffee ☕️
If you enjoy this project or my other projects, It wouldn't hurt to grab me a <a href="https://buymeacoffee.com/doobthegoober">coffee</a>! 🙏
