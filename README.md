# tcprouter

a down to earth tcp router based on traefik tcp streaming and supports multiple backends using [valkyrie](https://github.com/abronan/valkeyrie)


## Build

```bash
go get -u github.com/xmonader/tcprouter
```



## Running

configfile: router.toml
```toml
[server]
addr = "0.0.0.0"
port = 443

[server.dbbackend]
type 	 = "redis"
addr     = "127.0.0.1"
port     = 6379
refresh  = 10
```
then 
`./tcprouter router.toml`


Please notice if you are using low numbered port like 80 or 443 you can use sudo or setcap before running the binary.
- `sudo ./tcprouter router.toml`
- setcap: `sudo setcap CAP_NET_BIND_SERVICE=+eip PATH_TO_TCPROUTER`



### router.toml
We have two toml sections so far

#### [server]

```toml
[server]
addr = "0.0.0.0"
port = 443
```
in `[server]` section we define the listening interface/port the tcprouter intercepting: typically that's 443 for TLS connections.

#### [server.dbbackend]
```toml
[server.dbbackend]
type 	 = "redis"
addr     = "127.0.0.1"
port     = 6379
refresh  = 10
```
in `server.dbbackend` we define the backend kv store and its connection information `addr,port` and how often we want to reload the data from the kv store using `refresh` key in seconds.



## Data representation in KV

```
127.0.0.1:6379> KEYS *
1) "/tcprouter/services/bing"
2) "/tcprouter/services/google"
3) "/tcprouter/services/facebook"

127.0.0.1:6379> get /tcprouter/services/google
"{\"Key\":\"tcprouter/services/google\",\"Value\":\"eyJhZGRyIjoiMTcyLjIxNy4xOS40Njo0NDMiLCJzbmkiOiJ3d3cuZ29vZ2xlLmNvbSJ9\",\"LastIndex\":75292246}"

```

### Decoding data from python

```ipython

In [64]: res = r.get("/tcprouter/service/google")     

In [65]: decoded = json.loads(res)                    

In [66]: decoded                                      
Out[66]: 
{'Key': '/tcprouter/service/google',
 'Value': 'eyJhZGRyIjogIjE3Mi4yMTcuMTkuNDY6NDQzIiwgInNuaSI6ICJ3d3cuZ29vZ2xlLmNvbSIsICJuYW1lIjogImdvb2dsZSJ9'}


```
`Value` payload is base64 encoded because of how golang is marshaling.

```ipython
In [67]: base64.b64decode(decoded['Value'])           
Out[67]: b'{"addr": "172.217.19.46:443", "sni": "www.google.com", "name": "google"}'

```

## Examples

### Go

This example can be found at [examples/main.go](./examples/main.go)
```go

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"

	"github.com/abronan/valkeyrie/store/redis"
)

func init() {
	redis.Register()
}

type Service struct {
	Addr string `json:"addr"`
	SNI  string `json:"sni"`
	Name string `json:"bing"`
}

func main() {

	// Initialize a new store with redis
	kv, err := valkeyrie.NewStore(
		store.REDIS,
		[]string{"127.0.0.1:6379" },
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store redis")
	}
	google := &Service{Addr:"172.217.19.46:443", SNI:"www.google.com", Name:"google"}
	encGoogle, _ := json.Marshal(google)
	bing := &Service{Addr:"13.107.21.200:443", SNI:"www.bing.com", Name:"bing"}
	encBing, _ := json.Marshal(bing)

	kv.Put("/tcprouter/services/google", encGoogle, nil)
	kv.Put("/tcprouter/services/bing", encBing, nil)


}

```




### Python
```python3
import base64
import json
import redis

r = redis.Redis()

def create_service(name, sni, addr):
    service = {}
    service['Key'] = '/tcprouter/service/{}'.format(name)
    record = {"addr":addr, "sni":sni, "name":name}
    json_dumped_record_bytes = json.dumps(record).encode()
    b64_record = base64.b64encode(json_dumped_record_bytes).decode()
    service['Value'] = b64_record
    r.set(service['Key'], json.dumps(service))
    
create_service('facebook', "www.facebook.com", "102.132.97.35:443")
create_service('google', 'www.google.com', '172.217.19.46:443')
create_service('bing', 'www.bing.com', '13.107.21.200:443')
            

```


If you want to test that locally you can modify `/etc/hosts`

```


127.0.0.1 www.google.com
127.0.0.1 www.bing.com
127.0.0.1 www.facebook.com

```
So your browser go to your `127.0.0.1:443` on requesting google or bing.
