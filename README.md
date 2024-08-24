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
sudo mv gofuzz /usr/local/bin
gofuzz -h
```

# Usage 🧠
```
Usage of gofuzz:
  -body string
        specify POST request body
  -burp string
        specify path to burp request
  -custom-headers string
        specify the file that contains headers [separated by line]
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
gofuzz -url https://example.com/FUZZ -wordlist list.txt -custom-headers headers.txt
```

Fuzz parameters in URL:
```bash
gofuzz -url https://example.com/api/search=FUZZ -timeout 3 -threads 10 -wordlist list.txt
```

Fuzz with POST Requests:
```bash
gofuzz -url https://example.com/upload/file=FUZZ -method POST -body '{"test": "123456"}' -custom-headers headers.txt -wordlist list.txt
```
* To FUZZ the body of the post:
    ```bash
    gofuzz -url https://example.com/api/upload -method POST -body '{"payload": "FUZZ"}' -custom-headers headers.txt -wordlist list.txt
    ```

Filter by status code(s):
```bash
gofuzz -url https://example.com/FUZZ -wordlist list.txt | grep -E " 200 | 400 " -A 2
```
* **-A** *after-content*, gives the 2 lines after the match. *(Response Length, Response Body)*

BurpSuite Requests:
One can copy and paste http requests to fuzz/intruder your hacking adventures.
```bash
gofuzz -burp request.txt -wordlist list.txt -threads 10
```
* Burpsuite uses HTTP/2, change to 1.1 as should be.
```
POST /product/stock HTTP/1.1
Host: 0a7e00d204a2c01c80b71227002c00e3.web-security-academy.net
Cookie: session=XV7TWpUFWZsOMadSpqUmllHnsyX7aozv
Content-Length: 96
Sec-Ch-Ua: "Chromium";v="127", "Not)A;Brand";v="99"
Content-Type: application/x-www-form-urlencoded
Accept-Language: en-US
Sec-Ch-Ua-Mobile: ?0
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.89 Safari/537.36
Sec-Ch-Ua-Platform: "macOS"
Accept: */*
Origin: https://0a7e00d204a2c01c80b71227002c00e3.web-security-academy.net
Sec-Fetch-Site: same-origin
Sec-Fetch-Mode: cors
Sec-Fetch-Dest: empty
Referer: https://0a7e00d204a2c01c80b71227002c00e3.web-security-academy.net/product?productId=1
Accept-Encoding: gzip, deflate, br
Priority: u=1, i

stockApi=http%3A%2F%2F192.168.0.FUZZ%3A8080%2Fadmin
```
* Specifying FUZZ will be the position for which we try the wordlist on.

## Video Example
[Recording](https://github.com/user-attachments/assets/4d053735-9290-45e8-963c-14eb9f9221ec)




# Coffee ☕️
If you enjoy this project or my other projects, It wouldn't hurt to grab me a <a href="https://buymeacoffee.com/doobthegoober">coffee</a>! 🙏
