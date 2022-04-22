在 Go 项目开发中，绝大部分情况下，我们是在写能提供某种功能的后端服务，这些功能以 RPC API 接口或者 RESTful API 接口的形式对外提供，能提供这两种 API 接口的服务也统称为 Web 服务。今天这一讲，我就通过介绍 RESTful API 风格的 Web 服务，来给你介绍下如何实现 Web 服务的核心功能。

那今天我们就来看下，Web 服务的核心功能有哪些，以及如何开发这些功能。

**Web 服务的核心功能**

Web 服务有很多功能，为了便于你理解，我将这些功能分成了基础功能和高级功能两大类，并总结在了下面这张图中：

![img](https://static001.geekbang.org/resource/image/1a/2e/1a6d38450cdd0e115e505ab30113602e.jpg?wh=2248x1835)

下面，我就按图中的顺序，来串讲下这些功能。

要实现一个 Web 服务，首先我们要选择通信协议和通信格式。在 Go 项目开发中，有 HTTP+JSON 和 gRPC+Protobuf 两种组合可选。因为 iam-apiserver 主要提供的是 REST 风格的 API 接口，所以选择的是 HTTP+JSON 组合。



**Web 服务最核心的功能是路由匹配。**路由匹配其实就是根据(HTTP方法, 请求路径)匹配到处理这个请求的函数，最终由该函数处理这次请求，并返回结果，过程如下图所示：

![img](https://static001.geekbang.org/resource/image/1f/9d/1f5yydeffb32732e7d0e23a0a9cd369d.jpg?wh=2248x975)



一次 HTTP 请求经过路由匹配，最终将请求交由Delete(c *gin.Context)函数来处理。变量c中存放了这次请求的参数，在 Delete 函数中，我们可以进行参数解析、参数校验、逻辑处理，最终返回结果。

对于大型系统，可能会有很多个 API 接口，API 接口随着需求的更新迭代，可能会有多个版本，为了便于管理，我们需要对**路由进行分组**。

有时候，我们需要在一个服务进程中，同时开启 HTTP 服务的 80 端口和 HTTPS 的 443 端口，这样我们就可以做到：对内的服务，访问 80 端口，简化服务访问复杂度；对外的服务，访问更为安全的 HTTPS 服务。显然，我们没必要为相同功能启动多个服务进程，所以这时候就需要 Web 服务能够支持**一进程多服务**的功能。



我们开发 Web 服务最核心的诉求是：输入一些参数，校验通过后，进行业务逻辑处理，然后返回结果。所以 Web 服务还应该能够进行**参数解析、参数校验、逻辑处理、返回结果**。这些都是 Web 服务的业务处理功能。

上面这些是 Web 服务的基本功能，此外，我们还需要支持一些高级功能。

在进行 HTTP 请求时，经常需要针对每一次请求都设置一些通用的操作，比如添加 Header、添加 RequestID、统计请求次数等，这就要求我们的 Web 服务能够支持中间件特性。

为了保证系统安全，对于每一个请求，我们都需要进行认证。Web 服务中，通常有两种认证方式，一种是基于用户名和密码，一种是基于 Token。认证通过之后，就可以继续处理请求了。

为了方便定位和跟踪某一次请求，需要支持 RequestID，定位和跟踪 RequestID 主要是为了排障。

最后，当前的软件架构中，很多采用了前后端分离的架构。在前后端分离的架构中，前端访问地址和后端访问地址往往是不同的，浏览器为了安全，会针对这种情况设置跨域请求，所以 Web 服务需要能够处理浏览器的跨域请求。

到这里，我就把 Web 服务的基础功能和高级功能串讲了一遍。当然，上面只介绍了 Web 服务的核心功能，还有很多其他的功能，你可以通过学习Gin 的官方文档来了解。

你可以看到，Web 服务有很多核心功能，这些功能我们可以基于 net/http 包自己封装。但在实际的项目开发中， 我们更多会选择使用基于 net/http 包进行封装的优秀开源 Web 框架。本实战项目选择了 Gin 框架。

接下来，我们主要看下 Gin 框架是如何实现以上核心功能的，这些功能我们在实际的开发中可以直接拿来使用。

**为什么选择 Gin 框架？**

优秀的 Web 框架有很多，我们为什么要选择 Gin 呢？在回答这个问题之前，我们先来看下选择 Web 框架时的关注点。

在选择 Web 框架时，我们可以关注如下几点：

- 路由功能；
- 是否具备 middleware/filter 能力；
- HTTP 参数（path、query、form、header、body）解析和返回；
- 性能和稳定性；
- 使用复杂度；
- 社区活跃度。

按 GitHub Star 数来排名，当前比较火的 Go Web 框架有 Gin、Beego、Echo、Revel 、Martini。经过调研，我从中选择了 Gin 框架，原因是 Gin 具有如下特性：

轻量级，代码质量高，性能比较高；

项目目前很活跃，并有很多可用的 Middleware；

作为一个 Web 框架，功能齐全，使用起来简单。

那接下来，我就先详细介绍下 Gin 框架。

Gin是用 Go 语言编写的 Web 框架，功能完善，使用简单，性能很高。Gin 核心的路由功能是通过一个定制版的HttpRouter来实现的，具有很高的路由性能。

Gin 有很多功能，这里我给你列出了它的一些核心功能：

- 支持 HTTP 方法：GET、POST、PUT、PATCH、DELETE、OPTIONS。
- 支持不同位置的 HTTP 参数：路径参数（path）、查询字符串参数（query）、表单参数（form）、HTTP 头参数（header）、消息体参数（body）。
- 支持 HTTP 路由和路由分组。
- 支持 middleware 和自定义 middleware。
- 支持自定义 Log。
- 支持 binding 和 validation，支持自定义 validator。可以 bind 如下参数：query、path、body、header、form。
- 支持重定向。
- 支持 basic auth middleware。
- 支持自定义 HTTP 配置。
- 支持优雅关闭。
- 支持 HTTP2。
- 支持设置和获取 cookie。

**Gin 是如何支持 Web 服务基础功能的？**

接下来，我们先通过一个具体的例子，看下 Gin 是如何支持 Web 服务基础功能的，后面再详细介绍这些功能的用法。

我们创建一个 webfeature 目录，用来存放示例代码。因为要演示 HTTPS 的用法，所以需要创建证书文件。具体可以分为两步。

第一步，执行以下命令创建证书：

```
cat << 'EOF' > ca.pem
-----BEGIN CERTIFICATE-----
MIICSjCCAbOgAwIBAgIJAJHGGR4dGioHMA0GCSqGSIb3DQEBCwUAMFYxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQxDzANBgNVBAMTBnRlc3RjYTAeFw0xNDExMTEyMjMxMjla
Fw0yNDExMDgyMjMxMjlaMFYxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0
YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxDzANBgNVBAMT
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwEDfBV5MYdlHVHJ7
+L4nxrZy7mBfAVXpOc5vMYztssUI7mL2/iYujiIXM+weZYNTEpLdjyJdu7R5gGUu
g1jSVK/EPHfc74O7AyZU34PNIP4Sh33N+/A5YexrNgJlPY+E3GdVYi4ldWJjgkAd
Qah2PH5ACLrIIC6tRka9hcaBlIECAwEAAaMgMB4wDAYDVR0TBAUwAwEB/zAOBgNV
HQ8BAf8EBAMCAgQwDQYJKoZIhvcNAQELBQADgYEAHzC7jdYlzAVmddi/gdAeKPau
sPBG/C2HCWqHzpCUHcKuvMzDVkY/MP2o6JIW2DBbY64bO/FceExhjcykgaYtCH/m
oIU63+CFOTtR7otyQAWHqXa7q4SbCDlG7DyRFxqG0txPtGvy12lgldA2+RgcigQG
Dfcog5wrJytaQ6UA0wE=
-----END CERTIFICATE-----
EOF

cat << 'EOF' > server.key
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOHDFScoLCVJpYDD
M4HYtIdV6Ake/sMNaaKdODjDMsux/4tDydlumN+fm+AjPEK5GHhGn1BgzkWF+slf
3BxhrA/8dNsnunstVA7ZBgA/5qQxMfGAq4wHNVX77fBZOgp9VlSMVfyd9N8YwbBY
AckOeUQadTi2X1S6OgJXgQ0m3MWhAgMBAAECgYAn7qGnM2vbjJNBm0VZCkOkTIWm
V10okw7EPJrdL2mkre9NasghNXbE1y5zDshx5Nt3KsazKOxTT8d0Jwh/3KbaN+YY
tTCbKGW0pXDRBhwUHRcuRzScjli8Rih5UOCiZkhefUTcRb6xIhZJuQy71tjaSy0p
dHZRmYyBYO2YEQ8xoQJBAPrJPhMBkzmEYFtyIEqAxQ/o/A6E+E4w8i+KM7nQCK7q
K4JXzyXVAjLfyBZWHGM2uro/fjqPggGD6QH1qXCkI4MCQQDmdKeb2TrKRh5BY1LR
81aJGKcJ2XbcDu6wMZK4oqWbTX2KiYn9GB0woM6nSr/Y6iy1u145YzYxEV/iMwff
DJULAkB8B2MnyzOg0pNFJqBJuH29bKCcHa8gHJzqXhNO5lAlEbMK95p/P2Wi+4Hd
aiEIAF1BF326QJcvYKmwSmrORp85AkAlSNxRJ50OWrfMZnBgzVjDx3xG6KsFQVk2
ol6VhqL6dFgKUORFUWBvnKSyhjJxurlPEahV6oo6+A+mPhFY8eUvAkAZQyTdupP3
XEFQKctGz+9+gKkemDp7LBBMEMBXrGTLPhpEfcjv/7KPdnFHYmhYeBTBnuVmTVWe
F98XJ7tIFfJq
-----END PRIVATE KEY-----
EOF

cat << 'EOF' > server.pem
-----BEGIN CERTIFICATE-----
MIICnDCCAgWgAwIBAgIBBzANBgkqhkiG9w0BAQsFADBWMQswCQYDVQQGEwJBVTET
MBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0cyBQ
dHkgTHRkMQ8wDQYDVQQDEwZ0ZXN0Y2EwHhcNMTUxMTA0MDIyMDI0WhcNMjUxMTAx
MDIyMDI0WjBlMQswCQYDVQQGEwJVUzERMA8GA1UECBMISWxsaW5vaXMxEDAOBgNV
BAcTB0NoaWNhZ28xFTATBgNVBAoTDEV4YW1wbGUsIENvLjEaMBgGA1UEAxQRKi50
ZXN0Lmdvb2dsZS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAOHDFSco
LCVJpYDDM4HYtIdV6Ake/sMNaaKdODjDMsux/4tDydlumN+fm+AjPEK5GHhGn1Bg
zkWF+slf3BxhrA/8dNsnunstVA7ZBgA/5qQxMfGAq4wHNVX77fBZOgp9VlSMVfyd
9N8YwbBYAckOeUQadTi2X1S6OgJXgQ0m3MWhAgMBAAGjazBpMAkGA1UdEwQCMAAw
CwYDVR0PBAQDAgXgME8GA1UdEQRIMEaCECoudGVzdC5nb29nbGUuZnKCGHdhdGVy
em9vaS50ZXN0Lmdvb2dsZS5iZYISKi50ZXN0LnlvdXR1YmUuY29thwTAqAEDMA0G
CSqGSIb3DQEBCwUAA4GBAJFXVifQNub1LUP4JlnX5lXNlo8FxZ2a12AFQs+bzoJ6
hM044EDjqyxUqSbVePK0ni3w1fHQB5rY9yYC5f8G7aqqTY1QOhoUk8ZTSTRpnkTh
y4jjdvTZeLDVBlueZUTDRmy2feY5aZIU18vFDK08dTG0A87pppuv1LNIR3loveU8
-----END CERTIFICATE-----
EOF
```

第二步，创建 main.go 文件：

```
package main

import (
  "fmt"
  "log"
  "net/http"
  "sync"
  "time"

  "github.com/gin-gonic/gin"
  "golang.org/x/sync/errgroup"
)

type Product struct {
  Username    string    `json:"username" binding:"required"`
  Name        string    `json:"name" binding:"required"`
  Category    string    `json:"category" binding:"required"`
  Price       int       `json:"price" binding:"gte=0"`
  Description string    `json:"description"`
  CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
  sync.RWMutex
  products map[string]Product
}

func newProductHandler() *productHandler {
  return &productHandler{
    products: make(map[string]Product),
  }
}

func (u *productHandler) Create(c *gin.Context) {
  u.Lock()
  defer u.Unlock()

  // 1. 参数解析
  var product Product
  if err := c.ShouldBindJSON(&product); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // 2. 参数校验
  if _, ok := u.products[product.Name]; ok {
    c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
    return
  }
  product.CreatedAt = time.Now()

  // 3. 逻辑处理
  u.products[product.Name] = product
  log.Printf("Register product %s success", product.Name)

  // 4. 返回结果
  c.JSON(http.StatusOK, product)
}

func (u *productHandler) Get(c *gin.Context) {
  u.Lock()
  defer u.Unlock()

  product, ok := u.products[c.Param("name")]
  if !ok {
    c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not found product %s", c.Param("name"))})
    return
  }

  c.JSON(http.StatusOK, product)
}

func router() http.Handler {
  router := gin.Default()
  productHandler := newProductHandler()
  // 路由分组、中间件、认证
  v1 := router.Group("/v1")
  {
    productv1 := v1.Group("/products")
    {
      // 路由匹配
      productv1.POST("", productHandler.Create)
      productv1.GET(":name", productHandler.Get)
    }
  }

  return router
}

func main() {
  var eg errgroup.Group

  // 一进程多端口
  insecureServer := &http.Server{
    Addr:         ":8080",
    Handler:      router(),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  secureServer := &http.Server{
    Addr:         ":8443",
    Handler:      router(),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  eg.Go(func() error {
    err := insecureServer.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
      log.Fatal(err)
    }
    return err
  })

  eg.Go(func() error {
    err := secureServer.ListenAndServeTLS("server.pem", "server.key")
    if err != nil && err != http.ErrServerClosed {
      log.Fatal(err)
    }
    return err
  })

  if err := eg.Wait(); err != nil {
    log.Fatal(err)
  }
}
```



运行以上代码：

```
$ go run main.go
```

打开另外一个终端，请求 HTTP 接口：

```
# 创建产品
$ curl -XPOST -H"Content-Type: application/json" -d'{"username":"colin","name":"iphone12","category":"phone","price":8000,"description":"cannot afford"}' http://127.0.0.1:8080/v1/products
{"username":"colin","name":"iphone12","category":"phone","price":8000,"description":"cannot afford","createdAt":"2021-06-20T11:17:03.818065988+08:00"}

# 获取产品信息
$ curl -XGET http://127.0.0.1:8080/v1/products/iphone12
{"username":"colin","name":"iphone12","category":"phone","price":8000,"description":"cannot afford","createdAt":"2021-06-20T11:17:03.818065988+08:00"}
```

示例代码存放地址为webfeature。

另外，Gin 项目仓库中也包含了很多使用示例，如果你想详细了解，可以参考 gin examples。

下面，我来详细介绍下 Gin 是如何支持 Web 服务基础功能的。

**HTTP/HTTPS 支持**

因为 Gin 是基于 net/http 包封装的一个 Web 框架，所以它天然就支持 HTTP/HTTPS。在上述代码中，通过以下方式开启一个 HTTP 服务：

```
insecureServer := &http.Server{
  Addr:         ":8080",
  Handler:      router(),
  ReadTimeout:  5 * time.Second,
  WriteTimeout: 10 * time.Second,
}
...
err := insecureServer.ListenAndServe()
```

通过以下方式开启一个 HTTPS 服务：

```
secureServer := &http.Server{
  Addr:         ":8443",
  Handler:      router(),
  ReadTimeout:  5 * time.Second,
  WriteTimeout: 10 * time.Second,
}
...
err := secureServer.ListenAndServeTLS("server.pem", "server.key")
```

**JSON 数据格式支持**

Gin 支持多种数据通信格式，例如 application/json、application/xml。可以通过c.ShouldBindJSON函数，将 Body 中的 JSON 格式数据解析到指定的 Struct 中，通过c.JSON函数返回 JSON 格式的数据。

**路由匹配**

第一种匹配规则是精确匹配。例如，路由为 /products/:name，匹配情况如下表所示：

![img](https://static001.geekbang.org/resource/image/11/df/11be05d7fe7f935e01725e2635f315df.jpg?wh=2248x1418)

第二种匹配规则是模糊匹配。例如，路由为 /products/*name，匹配情况如下表所示：

![img](https://static001.geekbang.org/resource/image/b5/7b/b5ccd9924e53dd90a64af6002967b67b.jpg?wh=2248x1636)

**路由分组**

Gin 通过 Group 函数实现了路由分组的功能。路由分组是一个非常常用的功能，可以将相同版本的路由分为一组，也可以将相同 RESTful 资源的路由分为一组。例如：

```
v1 := router.Group("/v1", gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"}))
{
    productv1 := v1.Group("/products")
    {
        // 路由匹配
        productv1.POST("", productHandler.Create)
        productv1.GET(":name", productHandler.Get)
    }

    orderv1 := v1.Group("/orders")
    {
        // 路由匹配
        orderv1.POST("", orderHandler.Create)
        orderv1.GET(":name", orderHandler.Get)
    }
}

v2 := router.Group("/v2", gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"}))
{
    productv2 := v2.Group("/products")
    {
        // 路由匹配
        productv2.POST("", productHandler.Create)
        productv2.GET(":name", productHandler.Get)
    }
}
```

通过将路由分组，可以对相同分组的路由做统一处理。比如上面那个例子，我们可以通过代码

```
v1 := router.Group("/v1", gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"}))
```

给所有属于 v1 分组的路由都添加 gin.BasicAuth 中间件，以实现认证功能。中间件和认证，这里你先不用深究，下面讲高级功能的时候会介绍到。

**一进程多服务**

我们可以通过以下方式实现一进程多服务：

```go
var eg errgroup.Group
insecureServer := &http.Server{...}
secureServer := &http.Server{...}

eg.Go(func() error {
  err := insecureServer.ListenAndServe()
  if err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
  }
  return err
})
eg.Go(func() error {
  err := secureServer.ListenAndServeTLS("server.pem", "server.key")
  if err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
  }
  return err
}

if err := eg.Wait(); err != nil {
  log.Fatal(err)
})
```

上述代码实现了两个相同的服务，分别监听在不同的端口。这里需要注意的是，为了不阻塞启动第二个服务，我们需要把 ListenAndServe 函数放在 goroutine 中执行，并且调用 eg.Wait() 来阻塞程序进程，从而让两个 HTTP 服务在 goroutine 中持续监听端口，并提供服务。

**参数解析、参数校验、逻辑处理、返回结果**

此外，Web 服务还应该具有参数解析、参数校验、逻辑处理、返回结果 4 类功能，因为这些功能联系紧密，我们放在一起来说。

在 productHandler 的 Create 方法中，我们通过c.ShouldBindJSON来解析参数，接下来自己编写校验代码，然后将 product 信息保存在内存中（也就是业务逻辑处理），最后通过c.JSON返回创建的 product 信息。代码如下：

```
func (u *productHandler) Create(c *gin.Context) {
  u.Lock()
  defer u.Unlock()

  // 1. 参数解析
  var product Product
  if err := c.ShouldBindJSON(&product); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // 2. 参数校验
  if _, ok := u.products[product.Name]; ok {
    c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
    return
  }
  product.CreatedAt = time.Now()

  // 3. 逻辑处理
  u.products[product.Name] = product
  log.Printf("Register product %s success", product.Name)

  // 4. 返回结果
  c.JSON(http.StatusOK, product)
}
```

那这个时候，你可能会问：HTTP 的请求参数可以存在不同的位置，Gin 是如何解析的呢？这里，我们先来看下 HTTP 有哪些参数类型。HTTP 具有以下 5 种参数类型：

路径参数（path）。例如gin.Default().GET("/user/:name", nil)， name 就是路径参数。

查询字符串参数（query）。例如/welcome?firstname=Lingfei&lastname=Kong，firstname 和 lastname 就是查询字符串参数。

表单参数（form）。例如curl -X POST -F 'username=colin' -F 'password=colin1234' http://mydomain.com/login，username 和 password 就是表单参数。

HTTP 头参数（header）。例如curl -X POST -H 'Content-Type: application/json' -d '{"username":"colin","password":"colin1234"}' http://mydomain.com/login，Content-Type 就是 HTTP 头参数。

消息体参数（body）。例如curl -X POST -H 'Content-Type: application/json' -d '{"username":"colin","password":"colin1234"}' http://mydomain.com/login，username 和 password 就是消息体参数。

Gin 提供了一些函数，来分别读取这些 HTTP 参数，每种类别会提供两种函数，一种函数可以直接读取某个参数的值，另外一种函数会把同类 HTTP 参数绑定到一个 Go 结构体中。比如，有如下路径参数：

```
gin.Default().GET("/:name/:id", nil)
```

我们可以直接读取每个参数：

```
name := c.Param("name")
action := c.Param("action")
```

也可以将所有的路径参数，绑定到结构体中：

```
type Person struct {
    ID string `uri:"id" binding:"required,uuid"`
    Name string `uri:"name" binding:"required"`
}

if err := c.ShouldBindUri(&person); err != nil {
    // normal code
    return
}
```

Gin 在绑定参数时，是通过结构体的 tag 来判断要绑定哪类参数到结构体中的。这里要注意，不同的 HTTP 参数有不同的结构体 tag。

- 路径参数：uri。
- 查询字符串参数：form。
- 表单参数：form。
- HTTP 头参数：header。
- 消息体参数：会根据 Content-Type，自动选择使用 json 或者 xml，也可以调用 ShouldBindJSON 或者 ShouldBindXML 直接指定使用哪个 tag。

路径参数：uri。查询字符串参数：form。表单参数：form。HTTP 头参数：header。消息体参数：会根据 Content-Type，自动选择使用 json 或者 xml，也可以调用 ShouldBindJSON 或者 ShouldBindXML 直接指定使用哪个 tag。

针对每种参数类型，Gin 都有对应的函数来获取和绑定这些参数。这些函数都是基于如下两个函数进行封装的：

ShouldBindWith(obj interface{}, b binding.Binding) error

非常重要的一个函数，很多 ShouldBindXXX 函数底层都是调用 ShouldBindWith 函数来完成参数绑定的。该函数会根据传入的绑定引擎，将参数绑定到传入的结构体指针中，如果绑定失败，只返回错误内容，但不终止 HTTP 请求。ShouldBindWith 支持多种绑定引擎，例如 binding.JSON、binding.Query、binding.Uri、binding.Header 等，更详细的信息你可以参考 binding.go。

MustBindWith(obj interface{}, b binding.Binding) error

这是另一个非常重要的函数，很多 BindXXX 函数底层都是调用 MustBindWith 函数来完成参数绑定的。该函数会根据传入的绑定引擎，将参数绑定到传入的结构体指针中，如果绑定失败，返回错误并终止请求，返回 HTTP 400 错误。MustBindWith 所支持的绑定引擎跟 ShouldBindWith 函数一样。

路径参数：ShouldBindUri、BindUri；

查询字符串参数：ShouldBindQuery、BindQuery；

表单参数：ShouldBind；

HTTP 头参数：ShouldBindHeader、BindHeader；

消息体参数：ShouldBindJSON、BindJSON 等。

每个类别的 Bind 函数，详细信息你可以参考Gin 提供的 Bind 函数。

这里要注意，Gin 并没有提供类似 ShouldBindForm、BindForm 这类函数来绑定表单参数，但我们可以通过 ShouldBind 来绑定表单参数。当 HTTP 方法为 GET 时，ShouldBind 只绑定 Query 类型的参数；当 HTTP 方法为 POST 时，会先检查 content-type 是否是 json 或者 xml，如果不是，则绑定 Form 类型的参数。

所以，ShouldBind 可以绑定 Form 类型的参数，但前提是 HTTP 方法是 POST，并且 content-type 不是 application/json、application/xml。

在 Go 项目开发中，我建议使用 ShouldBindXXX，这样可以确保我们设置的 HTTP Chain（Chain 可以理解为一个 HTTP 请求的一系列处理插件）能够继续被执行。

**Gin 是如何支持 Web 服务高级功能的？**

**中间件**

Gin 支持中间件，HTTP 请求在转发到实际的处理函数之前，会被一系列加载的中间件进行处理。在中间件中，可以解析 HTTP 请求做一些逻辑处理，例如：跨域处理或者生成 X-Request-ID 并保存在 context 中，以便追踪某个请求。处理完之后，可以选择中断并返回这次请求，也可以选择将请求继续转交给下一个中间件处理。当所有的中间件都处理完之后，请求才会转给路由函数进行处理。具体流程如下图：

![img](https://static001.geekbang.org/resource/image/f0/80/f0783cb9ee8cffa969f846ebe8eae880.jpg?wh=2248x1655)

通过中间件，可以实现对所有请求都做统一的处理，提高开发效率，并使我们的代码更简洁。但是，因为所有的请求都需要经过中间件的处理，可能会增加请求延时。对于中间件特性，我有如下建议：

中间件做成可加载的，通过配置文件指定程序启动时加载哪些中间件。

只将一些通用的、必要的功能做成中间件。

在编写中间件时，一定要保证中间件的代码质量和性能。

在 Gin 中，可以通过 gin.Engine 的 Use 方法来加载中间件。中间件可以加载到不同的位置上，而且不同的位置作用范围也不同，例如：

```
router := gin.New()
router.Use(gin.Logger(), gin.Recovery()) // 中间件作用于所有的HTTP请求
v1 := router.Group("/v1").Use(gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"})) // 中间件作用于v1 group
v1.POST("/login", Login).Use(gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"})) //中间件只作用于/v1/login API接口
```

Gin 框架本身支持了一些中间件。

gin.Logger()：Logger 中间件会将日志写到 gin.DefaultWriter，gin.DefaultWriter 默认为 os.Stdout。

gin.Recovery()：Recovery 中间件可以从任何 panic 恢复，并且写入一个 500 状态码。

gin.CustomRecovery(handle gin.RecoveryFunc)：类似 Recovery 中间件，但是在恢复时还会调用传入的 handle 方法进行处理。

gin.BasicAuth()：HTTP 请求基本认证（使用用户名和密码进行认证）。

另外，Gin 还支持自定义中间件。中间件其实是一个函数，函数类型为 gin.HandlerFunc，HandlerFunc 底层类型为 func(*Context)。如下是一个 Logger 中间件的实现：

```
package main

import (
  "log"
  "time"

  "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
  return func(c *gin.Context) {
    t := time.Now()

    // 设置变量example
    c.Set("example", "12345")

    // 请求之前

    c.Next()

    // 请求之后
    latency := time.Since(t)
    log.Print(latency)

    // 访问我们发送的状态
    status := c.Writer.Status()
    log.Println(status)
  }
}

func main() {
  r := gin.New()
  r.Use(Logger())

  r.GET("/test", func(c *gin.Context) {
    example := c.MustGet("example").(string)

    // it would print: "12345"
    log.Println(example)
  })

  // Listen and serve on 0.0.0.0:8080
  r.Run(":8080")
}
```

另外，还有很多开源的中间件可供我们选择，我把一些常用的总结在了表格里：

![img](https://static001.geekbang.org/resource/image/67/10/67137697a09d9f37bd87a81bf322f510.jpg?wh=1832x1521)

**认证、RequestID、跨域**

认证、RequestID、跨域这三个高级功能，都可以通过 Gin 的中间件来实现，例如：

```
router := gin.New()

// 认证
router.Use(gin.BasicAuth(gin.Accounts{"foo": "bar", "colin": "colin404"}))

// RequestID
router.Use(requestid.New(requestid.Config{
    Generator: func() string {
        return "test"
    },
}))

// 跨域
// CORS for https://foo.com and https://github.com origins, allowing:
// - PUT and PATCH methods
// - Origin header
// - Credentials share
// - Preflight requests cached for 12 hours
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://foo.com"},
    AllowMethods:     []string{"PUT", "PATCH"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
        return origin == "https://github.com"
    },
    MaxAge: 12 * time.Hour,
}))
```

优雅关停

Go 项目上线后，我们还需要不断迭代来丰富项目功能、修复 Bug 等，这也就意味着，我们要不断地重启 Go 服务。对于 HTTP 服务来说，如果访问量大，重启服务的时候可能还有很多连接没有断开，请求没有完成。如果这时候直接关闭服务，这些连接会直接断掉，请求异常终止，这就会对用户体验和产品口碑造成很大影响。因此，这种关闭方式不是一种优雅的关闭方式。

这时候，我们期望 HTTP 服务可以在处理完所有请求后，正常地关闭这些连接，也就是优雅地关闭服务。我们有两种方法来优雅关闭 HTTP 服务，分别是借助第三方的 Go 包和自己编码实现。

方法一：借助第三方的 Go 包

如果使用第三方的 Go 包来实现优雅关闭，目前用得比较多的包是fvbock/endless。我们可以使用 fvbock/endless 来替换掉 net/http 的 ListenAndServe 方法，例如：

```
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

方法二：编码实现

借助第三方包的好处是可以稍微减少一些编码工作量，但缺点是引入了一个新的依赖包，因此我更倾向于自己编码实现。Go 1.8 版本或者更新的版本，http.Server 内置的 Shutdown 方法，已经实现了优雅关闭。下面是一个示例：

```go
// +build go1.8

package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.GET("/", func(c *gin.Context) {
    time.Sleep(5 * time.Second)
    c.String(http.StatusOK, "Welcome Gin Server")
  })

  srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
  }

  // Initializing the server in a goroutine so that
  // it won't block the graceful shutdown handling below
  go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      log.Fatalf("listen: %s\n", err)
    }
  }()

  // Wait for interrupt signal to gracefully shutdown the server with
  // a timeout of 5 seconds.
  quit := make(chan os.Signal, 1)
  // kill (no param) default send syscall.SIGTERM
  // kill -2 is syscall.SIGINT
  // kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  <-quit
  log.Println("Shutting down server...")

  // The context is used to inform the server it has 5 seconds to finish
  // the request it is currently handling
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if err := srv.Shutdown(ctx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
  }

  log.Println("Server exiting")
}
```

上面的示例中，需要把 srv.ListenAndServe 放在 goroutine 中执行，这样才不会阻塞到 srv.Shutdown 函数。因为我们把 srv.ListenAndServe 放在了 goroutine 中，所以需要一种可以让整个进程常驻的机制。

这里，我们借助了有缓冲 channel，并且调用 signal.Notify 函数将该 channel 绑定到 SIGINT、SIGTERM 信号上。这样，收到 SIGINT、SIGTERM 信号后，quilt 通道会被写入值，从而结束阻塞状态，程序继续运行，执行 srv.Shutdown(ctx)，优雅关停 HTTP 服务。

**总结**

今天我们主要学习了 Web 服务的核心功能，以及如何开发这些功能。在实际的项目开发中， 我们一般会使用基于 net/http 包进行封装的优秀开源 Web 框架。

当前比较火的 Go Web 框架有 Gin、Beego、Echo、Revel、Martini。你可以根据需要进行选择。我比较推荐 Gin，Gin 也是目前比较受欢迎的 Web 框架。Gin Web 框架支持 Web 服务的很多基础功能，例如 HTTP/HTTPS、JSON 格式的数据、路由分组和匹配、一进程多服务等。

另外，Gin 还支持 Web 服务的一些高级功能，例如 中间件、认证、RequestID、跨域和优雅关停等。

课后练习

使用 Gin 框架编写一个简单的 Web 服务，要求该 Web 服务可以解析参数、校验参数，并进行一些简单的业务逻辑处理，最终返回处理结果。欢迎在留言区分享你的成果，或者遇到的问题。

思考下，如何给 iam-apiserver 的 /healthz 接口添加一个限流中间件，用来限制请求 /healthz 的频率。

