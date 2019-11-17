# cloudrain

Web MUD based on websocket proxy to telnet

## Build Setup

### Frontend

``` bash
# install dependencies
npm install

# serve with hot reload at localhost:7171
npm run dev

# build for production with minification
npm run build
```

### Server

```
go build
```

#### hot reload

```
curl -fLo /usr/bin/air \
    https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air
chmod +x /usr/bin/air
air
```
