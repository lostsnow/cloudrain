# cloudrain

Web MUD based on websocket proxy to telnet

## Build Setup

### Web

git subtree (for develop)

```
git clone git@github.com:lostsnow/cloudrain.git
git remote add -f web git@github.com:lostsnow/cloudrain-web.git

git subtree add --prefix=web web master --squash
git subtree pull --prefix=web web master --squash
git subtree push --prefix=web web master
```

build

``` bash
cd web

# install dependencies
npm install

# serve with hot reload at localhost:7171
npm run dev

# build for production with minification
npm run build
```

### Server

```
go build -v
```

#### hot reload

```
curl -fLo /usr/bin/air \
    https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air
chmod +x /usr/bin/air
air
```
