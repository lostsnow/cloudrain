# cloudrain

Web MUD based on websocket proxy to telnet

## Build Setup

### Web

```
git submodule update --recursive --remote
```

build

``` bash
cd web

# install dependencies
npm install

# serve with hot reload at localhost:7171
npm run dev

# build for development
npm run build-dev

# build for production with minification
npm run build
```

### Server

```
go build -v
```

#### hot reload

```
# binary will be $(go env GOPATH)/bin/air
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

air
```

### Docker

```shell
docker-compose build --build-arg GOPROXY="https://goproxy.cn,direct" \
  --build-arg VUE_APP_WEBSOCKET_URL=ws://localhost:7071/ws
docker-compose up -d
```