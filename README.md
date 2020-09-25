# baas-demo BaaS应用开发示例

>  区块链合约应用调用端开发示例。

基于迅雷链 SDK ，开发者可在部署合约之后开发合约专有的应用服务。

接下来将以一个K-V存储合约及其应用的 demo 为例，详细说明如何**基于迅雷链BaaS平台和SDK开发合约应用**以调用该存储合约。

需要注意的是，在 demo 中 SDK 以官方示例 [Server](https://github.com/XunleiBlockchain/baas-sdk-go) 的形式存在，开发者可自行集成。

## 1 链上合约代码

该 K-V 存储的链上 Solidity 合约代码为：

```solidity
pragma solidity >=0.4.22 <0.6.0;

contract KVDB {
    mapping(bytes32 => bytes) private _data;

    function set(bytes32 key, bytes memory value) public payable returns (bool) {
        _data[key]=value;
        return true;
    }

    function get(bytes32 key) public view returns (bytes memory) {
        return _data[key];
    }
}
```

显然，该合约提供两个接口 ：

1. **set设置map内的 key、value**
2. **get获取map内的 key、value**

## 2 应用开发框架

demo 使用 beego 实现路由框架，并提供 RESTful 风格的、面向用户的合约调用服务。

### 1.1 文件结构

```sh
.
├── main.go            # 程序入口
├── Makefile           # 编译
├── start.sh           # 启动脚本
├── README.md          # README
├── conf               # demo系统配置目录 主要配置SDK-Server的host地址
│   └── app.conf
└── routers            # beego路由框架的路由
│   └── router.go
├── controllers        # controller层
│   ├── contract.go
│   └── default.go
├── models             # model层
│   ├── error.go
│   ├── req_call.go
│   ├── req_execute.go
│   ├── req_gettx.go
│   └── resp.go
└── contract           # 合约目录 目录下的所有合约都应注册至contract/contract.go中的合约中心
    ├── contract.go
    ├── ERC20          # ERC20测试合约
    │   └── ERC20.go
    └── KVDB           # K-V存储测试合约
        └── KVDB.go
```

其中，conf/app.conf的配置内容如下：

```
appname = baas-demo                     // 应用名
httpport = 8081                         // 监听端口

copyrequestbody = true                  // 拷贝请求体

# backend sdk-server
sdk-server = http://localhost:8080      // sdk-server的host
```

在文件目录下执行 `go run main.go` 即可运行demo示例。

### 1.2 框架析构

demo提供的http服务所使用的数据结构定义于models包。

对于开发者来说，需要特别关注与用户调用密切相关的以下三个结构：

```go
// 执行 对应contract/execute接口
type ReqExecute struct {
    Account  string        `json:"account"`    // 调用合约的用户账户地址
    Contract string        `json:"contract"`   // SDK视图下的合约名
    Addr     string        `json:"addr"`       // 合约地址
    Method   string        `json:"method"`     // 合约方法
    Params   []interface{} `json:"params"`     // 合约方法的参数 以列表形式给出
}

// 调用 对应contract/call接口
type ReqCall struct {
    Account  string        `json:"account"`    // 调用合约的用户账户地址
    Contract string        `json:"contract"`   // SDK视图下的合约名
    Addr     string        `json:"addr"`       // 合约地址
    Method   string        `json:"method"`     // 合约方法
    Params   []interface{} `json:"params"`     // 合约方法的参数 以列表形式给出
}

// 查询 对应contract/getTx接口
type ReqGetTx struct {
    Account string `json:"account"`            // 调用查询接口的用户账户地址
    Hash    string `json:"hash"`               // 交易Hash
}
```

HTTP服务将从用户的JSON请求中解析得到用户账户地址（Account）、合约名（Contract）、合约方法（Method）、合约参数（Params）、交易hash（Hash）等信息，依据不同接口执行不同流程。

1. 对于 Execute、Call，框架将对其Params字段进行ABI编码：

```go
    // 1. 首先依据合约名找到合约实例：someContract
    // 2. 合约示例调用其Data接口实现对合约参数的ABI编码
    cdata, err := someContract.Data(method, params)
    // 3. 向SDK发送编码后的数据请求
    backReqParam := []interface{}{
        map[string]string{
            "from": account,              // 用户钱包账户地址
            "to":   someContract.Addr(),  // 合约地址
            "data": cdata,                // ABI编码结果
        },
    }
    backCall("sendContractTransaction", backReqParam)
```

1. 对于getTx，直接使用其hash值：

```go
    // 1. 直接根据用户请求携带的`用户账户地址`和`交易hash`向SDK发送数据请求
    backReqParam := []string{
        account,    // 用户账户地址
        hash,       // 交易hash值
    }
    return backCall("getTransactionReceipt", backReqParam)
```

容易看到，应用对外提供的接口最终都以 `方法method` 和 `入参params` 的形式进入了backCall。

接下来来看backCall做了哪些事情：

```go
func backCall(method string, params interface{}) (interface{}, error) {
    // 1. 设置请求体参数 包括id、协议、方法和参数
    backReqParams := make(map[string]interface{})
    backReqParams["id"] = 0
    backReqParams["jsonrpc"] = "2.0"
    backReqParams["method"] = method
    backReqParams["params"] = params
    // 2. 创建请求并填入参数 设置连接超时和读写超时时间
    backReq := httplib.Post("sdk-server-address"))
    backReq.SetTimeout(60*time.Second, 60*time.Second)
    reqBody, err := json.Marshal(backReqParams)
    if err != nil {
        return nil, err
    }
    backReq.Body(reqBody)

    // 3. 解码SDK-Server返回数据
    var ret backResp
    err = backReq.ToJSON(&ret)
    if err != nil {
        return nil, err
    }
    if ret.Errcode != 0 {
        return nil, fmt.Errorf("errcode: %d errmsg: %s", ret.Errcode, ret.Errmsg)
    }
    return ret.Result, nil
}
```

### 1.3 接口测试

在现有应用开发示例的基础上，开发者可以使用 `curl` 工具可以简单快速地对示例程序进行接口测试。

**在框架析构部分对用户请求的数据结构做出了说明。开发者在模拟用户请求进行测试时，应使用数据结构规定的JSON字段。**

以下示例将在 K-V合约中存储以“this is a test”并校验执行结果 ：

```json
// set
curl -H 'Content-Type: application/json' -d '{"account":"0x54fb1c7d0f011dd63b08f85ed7b518ab82028100","contract":"kvdb","addr":"0xd1abccccbd0e74e3427ab4f487e4b8ddbb8082a3","method":"set","params":["1111111111111111111111111111111111111111111111111111111111111111","this is a test"]}' http://127.0.0.1:8088/v1/contract/execute
=> 返回交易hash
// 请求参数说明：
// 1. account:    用户账户地址
// 2. contract:   调用合约名,此处为kvdb
// 3. addr:       合约地址
// 4. method:     合约方法,此处为set：设置键值对
// 5. params:     合约参数,对于kvdb的set方法,其参数列表共有两个元素：
//                    (1) bytes32类型的key(详情参考solidity合约部分),即此处的64个'1'
//                        每两个'1'构成一个16进制数'0x11',共计32个16进制数
//                    (2) byte类型的value
// 开发者如需自行构造key 只需满足5(1)的8位拓展ASCII编码条件，即：使用数字'0'~'9'、字符'a'~'f'拼成总长度为64的字符串即可。

// 查询合约交易-执行结果
curl -H 'Content-Type: application/json' -d '{"account":"0x54fb1c7d0f011dd63b08f85ed7b518ab82028100","hash":"0x81f68810182a40f25416c2db3efa5a98b2c7e1dfe75078b27e9c608cbf215604"}' http://127.0.0.1:8088/v1/contract/getTx
=>
=> {
  "code": 0,
  "msg": "",
  "data": {
    "blockHash": "0x62d79b7421a28e9b47b10a2e963d0c4c2ae15fe7ff2c59dfe06e31d311533807",
    "blockNumber": "0xf52c",
    "contractAddress": null,
    "cumulativeGasUsed": "0x8cc1",
    "from": "0x54fb1c7d0f011dd63b08f85ed7b518ab82028100",
    "gasUsed": "0x8cc1",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1",
    "to": "0xd1abccccbd0e74e3427ab4f487e4b8ddbb8082a3",
    "tokenAddress": "0x0000000000000000000000000000000000000000",
    "transactionHash": "0x81f68810182a40f25416c2db3efa5a98b2c7e1dfe75078b27e9c608cbf215604",
    "transactionIndex": "0x0"
  }
}
// 请求参数说明：
// 1. account:    用户账户地址
// 2. hash:       合约交易hash

// get
curl -H 'Content-Type: application/json' -d '{"account":"0x54fb1c7d0f011dd63b08f85ed7b518ab82028100","contract":"kvdb","addr":"0xd1abccccbd0e74e3427ab4f487e4b8ddbb8082a3","method":"get","params":["1111111111111111111111111111111111111111111111111111111111111111"]}' http://127.0.0.1:8088/v1/contract/call
=>
=> {
  "code": 0,
  "msg": "",
  "data": "this is a test"
}
// 请求参数说明：
// 1. account:    用户账户地址
// 2. contract:   调用合约名,此处为kvdb
// 3. addr:       合约地址
// 4. method:     合约方法,此处为get：设置特定键对应的值
// 5. params:     合约参数,对于kvdb的get方法,其参数列表只有一个元素：
//                    (1) bytes32类型的key(详情参考solidity合约部分),即此处的64个'1'
```

## 3 应用开发规范

合约应用开发需遵循以下规范：

1. **为合约应用生成合约ABI \*[(点此了解合约ABI)](https://github.com/openethereum/ethabi)\***
2. **合约应用实现Contract接口**

```go
type Contract interface {
    Def() string                                   // 合约定义
    Addr() string                                  // 合约地址
    Data(string, []interface{}) (string, error)    // 请求数据ABI编码
    Result(string, string) (interface{}, error)    // 返回结果ABI解码
}
```

1. **将合约注册至合约中心以提供http服务**

以K-V存储合约为例：

**1) 新建合约结构时，依赖合约定义生成ABI：**

```go
func newKVDB() *KVDB {
    kvdb := &KVDB{
        def: `[{"constant":true,"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"}],"name":"get","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"payable":false,"stateMutability":"view","type":"function"},
               {"constant":false,"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes","name":"value","type":"bytes"}],"name":"set","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"}]`,
        addr: "0x54fb1c7d0f011dd63b08f85ed7b518ab82028101",
    }
    abi, err := abi.JSON(strings.NewReader(kvdb.def))
    if err != nil {
        panic(fmt.Errorf("newKVDB abi.JSON err: %v", err))
    }
    kvdb.abi = abi
    return kvdb
}
```

**2) 实现 Contract 接口：**

```go
func (self *KVDB) Def() string {
    return self.def
}

func (self *KVDB) Addr() string {
    return self.addr
}

// 在 K-V 存储中，链上合约提供了 set存储 和 get读取 两个接口
// 相应的，这里调用不同的方法，在方法中对参数进行ABI编码
func (self *KVDB) Data(method string, params []interface{}) (string, error) {
    switch method {
    case "set":
        return self.Set(params)    // 内部实现为ABI打包：abi.Pack
    case "get":
        return self.Get(params)    // 内部实现为 abi.Pack
    }
}

// ABI解码获取返回数据
func (self *KVDB) Result(method string, ret string) (interface{}, error) {
    err := self.abi.Unpack(&val, method, common.FromHex(ret))
    // ...
    return string(val.([]byte)), nil
}
```

**3) 将合约注册至合约中心：**

```go
func init() {
    contract.Register(newKVDB())
}
```

## 4 新增合约应用

demo中对合约逻辑和端口服务逻辑进行了解耦。开发者无需关心合约如何提供端口服务。

基于该demo和开发规范，开发者依次完成以下步骤，即可新增合约应用：

1. 在框架的 contract/ 路径下实现新的合约代码：包括合约定义、结构及方法，init注册函数等。
2. 在 main.go 中引入合约代码所在的Go Package以执行合约注册。*(在示例中使用"_ baas-demo/contract/KVDB"的方式引入合约包是为了执行定义于其中的init函数)*

