<div align="center">
<img src="LOGO.svg" alt="hack-browser-data logo" />
</div>


# HackBrowserData

`HackBrowserData` 是一个浏览器数据（密码|历史记录|Cookie|书签|信用卡|下载记录|localStorage|浏览器插件）的导出工具，支持全平台主流浏览器。

二开: 增加cookie处理

## 各平台浏览器支持情况

### Windows

| 浏览器        | 密码 | Cookie | 书签 | 历史记录 |
| :------- | :------: | :----: | :------: | :-----: |
| Google Chrome|    ✅     |   ✅    |    ✅     |    ✅    |
| Google Chrome Beta|    ✅    |   ✅   |    ✅    |    ✅    |
| Chromium |    ✅    |    ✅    |    ✅    |    ✅    |
| Microsoft Edge|    ✅     |   ✅    |    ✅     |    ✅    |
| 360 极速浏览器    |    ✅     |   ✅    |    ✅     |    ✅    |
| QQ |    ✅     |   ✅    |    ✅     |    ✅    |
| Brave  |    ✅    |   ✅   |    ✅    |    ✅    |
| Opera  |    ✅    |    ✅    |    ✅    |    ✅    |
| OperaGX  |    ✅    |    ✅    |    ✅    |    ✅    |
| Vivaldi  |    ✅    |    ✅    |    ✅    |    ✅    |
| Yandex |    ✅    |    ✅    |    ✅    |    ✅    |
| CocCoc |    ✅    |    ✅    |    ✅    |    ✅    |
| Firefox |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Beta |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Dev |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox ESR |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Nightly |    ✅    |   ✅   |    ✅    |    ✅    |
| IE 浏览器        |    ❌     |   ❌    |    ❌     |    ❌    |

### MacOS

由于 MacOS 的安全性设置，基于 `Chromium` 内核浏览器解密时**需要当前用户密码**

| 浏览器    | 密码 | Cookie | 书签 | 历史记录 |
| :--- | :------: | :----: | :------: | :-----: |
| Google Chrome  |    ✅     |   ✅    |    ✅     |    ✅    |
| Google Chrome Beta |    ✅    |   ✅   |    ✅    |    ✅    |
| Chromium |    ✅    |    ✅    |    ✅    |    ✅    |
| Microsoft Edge |    ✅     |   ✅    |    ✅     |    ✅    |
| Brave |    ✅    |   ✅   |    ✅    |    ✅    |
| Opera |    ✅    |    ✅    |    ✅    |    ✅    |
| OperaGX |    ✅    |    ✅    |    ✅    |    ✅    |
| Vivaldi |    ✅    |    ✅    |    ✅    |    ✅    |
| Yandex |    ✅    |    ✅    |    ✅    |    ✅    |
| CocCoc |    ✅    |    ✅    |    ✅    |    ✅    |
| Firefox |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Beta |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Dev |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox ESR |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Nightly |    ✅    |   ✅   |    ✅    |    ✅    |
| Safari   |    ❌     |   ❌    |    ❌     |    ❌|

### Linux

| 浏览器    | 密码 | Cookie | 书签 | 历史记录 |
| :----- | :------: | :----: | :------: | :-----: |
| Google Chrome |    ✅     |   ✅    |    ✅     |    ✅    |
| Google Chrome Beta |    ✅    |   ✅   |    ✅    |    ✅    |
| Chromium |    ✅    |    ✅    |    ✅    |    ✅    |
| Microsoft Edge |    ✅    |   ✅   |    ✅    |    ✅    |
| Brave |    ✅    |   ✅   |    ✅    |    ✅    |
| Opera |    ✅    |    ✅    |    ✅    |    ✅    |
| Vivaldi |    ✅    |    ✅    |    ✅    |    ✅    |
| Chromium |    ✅     |   ✅    |    ✅     |    ✅    |
| Firefox |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Beta |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Dev |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox ESR |    ✅    |   ✅   |    ✅    |    ✅    |
| Firefox Nightly |    ✅    |   ✅   |    ✅    |    ✅    |

## 安装运行
### 安装

可下载已编译好，可直接运行的 [二进制文件](https://github.com/moonD4rk/HackBrowserData/releases)

> 某些情况下，这款安全工具会被 Windows Defender 或其他杀毒软件当作病毒导致无法执行。代码已经全部开源，可自行编译。

### 从源码编译

仅支持 `go 1.18+` 以后版本，一些函数使用到了泛型

``` bash
$ git clone https://github.com/moonD4rk/HackBrowserData

$ cd HackBrowserData/cmd/hack-browser-data

$ CGO_ENABLED=1 go build
```

### 跨平台编译

由于用到了 `go-sqlite3` 库，在跨平台编译时需提前安装支持目标平台的 `GCC` 工具，下面以 `MacOS` 下分别编译 `Windows` 和 `Linux` 程序为例：

#### Windows

``` shell
brew install mingw-w64

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build
```

#### Linux

``` shell
brew install FiloSottile/musl-cross/musl-cross

CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static"
```

### 运行
双击直接运行，也可以使用命令行调用相应的命令。

```
PS C:\test> .\hack-browser-data.exe -h
NAME:
   hack-browser-data - Export passwords/cookies/history/bookmarks from browser

USAGE:
   [hack-browser-data -b chrome -f json -dir results -cc]
   Export all browingdata(password/cookie/history/bookmark) from browser
   Github Link: https://github.com/moonD4rk/HackBrowserData

GLOBAL OPTIONS:
   --verbose, --vv                   verbose (default: false)
   --compress, --zip                 compress result to zip (default: false)
   --browser value, -b value         available browsers: all|chrome|opera-gx|vivaldi|coccoc|brave|edge|chromium|chrome-beta|opera|yandex|firefox (default: "all")
   --results-dir value, --dir value  export dir (default: "results")
   --format value, -f value          file name csv|json (default: "csv")
   --profile-path value, -p value    custom profile dir path, get with chrome://version
   --help, -h                        show help (default: false)
   --version, -v                     print the version (default: false)


PS C:\test>  .\hack-browser-data.exe -b all -f json --dir results -zip
[NOTICE] [browser.go:46,pickChromium] find browser Chrome success  
[NOTICE] [browser.go:46,pickChromium] find browser Microsoft Edge success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_download.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_password.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_creditcard.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_bookmark.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_cookie.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/microsoft_edge_history.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_history.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_download.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_password.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_creditcard.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_bookmark.json success  
[NOTICE] [browsingdata.go:59,Output] output to file results/chrome_cookie.json success  

```
