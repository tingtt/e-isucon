# この Write up について

## コマンド

以下の記述はそれぞれのマシンでのコマンド実行を示します。
ローカルで実行する場合やサーバー上でのユーザーが重要な場合は明記します。

- ローカルマシン（自分のPC） : MacBook Pro 2020 - Montrey 12.5.1
    ```bash
    # [localmachine ~]
    ```

- 問題環境（SSHでリモート接続するサーバー） : AWS EC2 Linux

    ```bash
    # [ec2-user@server ~]
    ```
    or
    ```bash
    # [root@server ~]
    ```
    or （何もなし）
    ```bash
    ```

> 意味
> ```
> [(ユーザー名):(マシン) (いる場所)]
> ```

---

# スコア

1349 → 12380

---

# SSH のセットアップ

> https://growi.tingtt.jp/cttit/isucon_manual_ssh

1. 秘密鍵のDL

	`~/Downloads/` に `eisucon.pem` をDL

1. 接続
	```bash
	# [localmachine ~]
	ssh ec2-user@13.208.215.218 -i ~/Downloads/eisucon.pem
	```

---

# サーバーの探索

> この WriteUp は問題作者が作成しているため探索が最短です。
> 本家のISUCONでは更に探索が必要になる場合があります。

## 外部から探索

[cttit (tingtt.jp) - jobs](https://cttit.tingtt.jp/contests/isucon/1/jobs)から `チーム環境 - 題材`にあるアドレスを見に行く。　

###  http://13.208.215.218/

- Webのフロントエンドが動いている。
	- ログインしてみる。
	- イベント一覧（`/events`） が遅い。
	- 管理画面（`/manage/user`, `/manage/event`）が遅い。
	- ブラウザのデバッグツールで遅いページの通信を見てみる。
		- `/events?embed=user&embed=documents`
			- 29 秒
			- 1473821 byte
			- `Transfer-Encoding: chunked`（エンコード方式）
		- `/users`
			- 440 ミリ秒
			- 13090 byte
			- `Transfer-Encoding: chunked`（エンコード方式）

###  http://13.208.215.218:1323/events

- WebAPI（Webのバックエンド）が`1323`ポートで動いている。
	- `/events`のデータ量が多い。

## 内部から探索

```bash
## ファイル一覧
$ ls
docker-compose.version  prc_hub_back-e-isucon
```
→  `docker` が動いてるかも
→ `docker-compose` も使われてそう

```bash
## ./prc_hub_back-e-isucon 内のファイル一覧
$ ls prc_hub_back-e-isucon/
Dockerfile  LICENSE  Makefile  README.md  application  docker-compose.yml  domain  go.mod  go.sum  main.go  presentation  utils.go
```
→ `go` で書かれた何かしらのアプリケーション
→ `docker-compose`が使われてる

```bash
## ./prc_hub_back-e-isucon ディレクトリに移動
$ cd ./prc_hub_back-e-isucon

## docker コンテナの一覧
$ docker-compose ps
NAME                              COMMAND                  SERVICE             STATUS              PORTS
prc_hub_back-e-isucon-backend-1   "go run . --log.leve…"   backend             running             0.0.0.0:1323->1323/tcp, :::1323->1323/tcp
prc_hub_back-e-isucon-mysql-1     "docker-entrypoint.s…"   mysql               running             33060/tcp, 0.0.0.0:49153->3306/tcp, :::49153->3306/tcp

## docker のログを確認
$ docker-compose logs backend
backend-backend-1  | time="2023-01-18T05:07:57.569173122Z" level=info msg="Log level: \"debug\"" func="prc_hub_back/domain/model/logrus.(*WrapLogrus).SetLevel" file="/go/src/app/domain/model/logrus/main.go:65"
backend-backend-1  | time="2023-01-18T05:07:57.869434252Z" level=info msg="cors enabled" func=prc_hub_back/presentation/echo.Start file="/go/src/app/presentation/echo/main.go:27"
backend-backend-1  | time="2023-01-18T05:07:57.869487541Z" level=debug msg="cors allow origins: [*]" func=prc_hub_back/presentation/echo.Start file="/go/src/app/presentation/echo/main.go:28"
backend-backend-1  |
backend-backend-1  |    ____    __
backend-backend-1  |   / __/___/ /  ___
backend-backend-1  |  / _// __/ _ \/ _ \
backend-backend-1  | /___/\__/_//_/\___/ v4.8.0
backend-backend-1  | High performance, minimalist Go web framework
backend-backend-1  | https://echo.labstack.com
backend-backend-1  | ____________________________________O/_______
backend-backend-1  |                                     O\
backend-backend-1  | ⇨ http server started on [::]:1323
```
→ 1323 ポートで `go` でバックエンドが動いている
→ `go`は`echo`というWebフレームワークを使用している

→ `MySQL`が`docker`コンテナで動いている

→ `go`のソースコードや`MySQL`の設定やスキーマを書き換えてチューニングできそう

```bash
## プロセス一覧
ps ax
```

一部抜粋

>  ```
>  3384 ?        Ssl    0:01 /usr/bin/node server.js
>  ```

→ `node.js` が動いている（Webサーバー）

> ```
> 4297 ?        Ssl    0:04 go run . --log.level=debug --jwt.issuer=prc_hub --jwt.secret=e0VzhtkQ --mysql.host=mysql --mysql.db=prc_hub --mysql.user=prc_hub --mysql.password=secret --migrate-sql-file=./domain/model/eisucon/migrate.sql
> ```

→ `go` で WebAPIが動いている

---

# 開発環境のセットアップ

> https://growi.tingtt.jp/cttit/isucon_manual_src_apply_cli

注意： ↑ のマニュアルとは少し違う構成で作成します。

## チューニングできそうなのアプリケーションのソースコードをDL

1. サーバー上のディレクトリを圧縮
	```bash
	## ファイル一覧
	$ ls
	docker-compose.version  prc_hub_back-e-isucon
	# 圧縮
	$ tar cf backend.tar.gz prc_hub_back-e-isucon
	## ファイル一覧
	$ ls
	backend.tar.gz  docker-compose.version  prc_hub_back-e-isucon
	```

1. 自分のマシンにDL
	```bash
	# [localmachine ~/Downloads]

	# ディレクトリ作成
	$ mkdir eisucon
	# 移動
	$ cd eisucon

	# [localmachine ~/Downloads/eisucon]

	# 探索したアプリケーションのファイルをDL
	$ scp -i eisucon.pem ec2-user@13.208.215.218:/home/ec2-user/backend.tar.gz ./
	# 解凍
	$ tar xf backend.tar.gz
	# 解凍前野ファイルを削除
	$ rm backend.tar.gz
	# ディレクトリ名を変更
	$ mv prc_hub_back-e-isucon backend

	# ファイル一覧
	$ tree -L 3
	.
	└── backend
	    ├── Dockerfile
	    ├── LICENSE
	    ├── Makefile
	    ├── README.md
	    ├── application
	    │   ├── eisucon
	    │   ├── event
	    │   └── user
	    ├── docker-compose.yml
	    ├── domain
	    │   └── model
	    ├── go.mod
	    ├── go.sum
	    ├── main.go
	    ├── presentation
	    │   ├── echo
	    │   └── oas.yml
	    └── utils.go
	
	9 directories, 10 files

	# VSCodeで開く
	$ code ./
	```

## GitHub に Push

1. GitHubにプライベートリポジトリ作成
2. DLしたディレクトリを push
	```bash
	# [localmachine ~/Downloads/eisucon]

	# 初期化
	git init
	git checkout -b main
	git add .
	git commit -m 'initial commit'

	# 作成したプライベートリポジトリを紐付け
	git remote add origin git@github.com:tingtt/e-isucon.git

	# push
	git push origin main
	```

## サーバーで clone

1. GitHubでアクセストークン作成
	[New Personal Access Token (Classic) (github.com)](https://github.com/settings/tokens/new)
2. `git`をインストール
	```bash
	# [ec2-user@server ~]

	# root ユーザーに変更
	sudo su -
	```

	```bash
	# [root@server ~]

	# git をインストール
	yum install git -y
	```
3. リポジトリを clone
	```bash
	cd /usr/local/src

	git clone \
		https://tingtt:<token>@github.com/tingtt/e-isucon.git \
		eisucon
	```

## WebAPI (`go`)と DB(`mysql`) を起動中のものと入れ替える。

1. 起動中のものを停止
	```bash
	# [root@server ~]
	
	# 移動
	cd /home/ec2-user/prc_hub_back-e-isucon
	
	# 停止
	docker-compose down --volumes --rmi local
	```
2. 先程 clone したコードで起動
	```bash
	# 移動
	cd /usr/local/src/eisucon/backend

	# 環境変数ファイルをコピー
	cp /home/ec2-user/prc_hub_back-e-isucon/.env ./

	# 起動
	docker-compose up -d

	# 確認
	docker-compose ps
	NAME                COMMAND                  SERVICE             STATUS              PORTS
	backend-backend-1   "go run . --log.leve…"   backend             running             0.0.0.0:1323->1323/tcp, :::1323->1323/tcp
	backend-mysql-1     "docker-entrypoint.s…"   mysql               running             33060/tcp, 0.0.0.0:49154->3306/tcp, :::49154->3306/tcp
	```

## Makefileで作業を楽にする。

```bash
# [localmachine ~/Downloads/eisucon]
code Makefile
```

```Makefile
.PHONY: ps
ps:
cd ./backend ; \
docker-compose ps

.PHONY: build
build:
cd ./backend ; \
docker-compose build

.PHONY: up
up:
cd ./backend ; \
docker-compose up -d

.PHONY: down
down:
cd ./backend ; \
docker-compose down

.PHONY: purge
purge:
cd ./backend ; \
docker-compose down --volumes

.PHONY: upgrade
upgrade: purge build up
```

```bash
git add Makefile
git commit -m 'ci: add `Makefile` for ci/cd'
```

---

# 探索 その2

##  `docker-compose.yml`,  `Dockerfile`, `.env`を見てみる

> `docker-compose.yml`
> ```yaml
> version: "3.7"
>
> services:
>   backend:
>     build:
>       context: .
>       dockerfile: Dockerfile
>       target: dev
>     volumes:
>       - .:/go/src/app
>     ports:
>       - ${PORT:-1323}:${PORT:-1323}
>     environment:
>       TZ: ${TZ:-UTC}
>       PORT: ${PORT:-1323}
>       MYSQL_DATABASE: ${MYSQL_DATABASE:-prc_hub}
>       MYSQL_USER: ${MYSQL_USER:-prc_hub}
>       MYSQL_PASSWORD: ${MYSQL_PASSWORD}
>     command: ${ARGS:-}
>     depends_on:
>       - mysql
>     restart: unless-stopped
>
>   mysql:
>     image: mysql:8
>     volumes:
>       - type: bind
>         source: "./.mysql/init.sql"
>         target: "/docker-entrypoint-initdb.d/init.sql"
>       - type: bind
>         source: "./.mysql/my.cnf"
>         target: "/etc/mysql/conf.d/my.cnf"
>       # - ./.mysql/log:/var/log/mysql
>       - mysql_data:/var/lib/mysql
>     ports:
>       - 3306
>     environment:
>       TZ: ${TZ:-UTC}
>       MYSQL_DATABASE: ${MYSQL_DATABASE:-prc_hub}
>       MYSQL_USER: ${MYSQL_USER:-prc_hub}
>       MYSQL_PASSWORD: ${MYSQL_PASSWORD}
>       MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
>     restart: unless-stopped
>
> volumes:
>   mysql_data:
>
> ```

> `Dockerfile`
> ```Dockerfile
> FROM golang:1.18-alpine as dev
>
> ENV ROOT=/go/src/app
> ENV CGO_ENABLED 0
> WORKDIR ${ROOT}
>
> RUN apk update && apk add git
> COPY go.mod go.sum ./
> RUN go mod download
> EXPOSE ${PORT}
>
> ENTRYPOINT ["go", "run", "."]
>
>
> FROM golang:1.18-alpine as builder
>
> ENV ROOT=/go/src/app
> WORKDIR ${ROOT}
>
> RUN apk update && apk add git
> COPY go.mod go.sum ./
> RUN go mod download
>
> COPY . ${ROOT}
> RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary
>
>
> FROM alpine as prod
>
> ENV ROOT=/go/src/app
> WORKDIR ${ROOT}
> COPY --from=builder ${ROOT}/binary ${ROOT}
>
> EXPOSE ${PORT}
> ENTRYPOINT ["/go/src/app/binary"]
> ```

> `.env`
> ```shell
> PORT=1323
>
> MYSQL_DATABASE=prc_hub
> MYSQL_USER=prc_hub
> MYSQL_PASSWORD=secret
> MYSQL_ROOT_PASSWORD=secret
>
> ARGS=--log.level=debug --jwt.issuer=prc_hub --jwt.secret=e0VzhtkQ --mysql.host=mysql --mysql.db=prc_hub --mysql.user=prc_hub --mysql.password=secret --migrate-sql-file="./domain/model/eisucon/migrate.sql"
> ```

→ コンテナは2つ

- `backend`
	- `./Dockerfile`が参照されている
	- `go`が実行されている
		- 実行されるコマンド
			```bash
			go run . \
				--log.level=debug \
				--jwt.issuer=prc_hub \
				--jwt.secret=e0VzhtkQ \
				--mysql.host=mysql \
				--mysql.db=prc_hub \
				--mysql.user=prc_hub \
				--mysql.password=secret \
				--migrate-sql-file="./domain/model/eisucon/migrate.sql"
			```
- `mysql`
	- `MySQL`サーバーが実行されている
	- `./.mysql/my.cnf`にコンフィグファイルがある
	- `./.mysql/init.sql`に初期化用のSQLがある
	- 接続コマンド
		```bash
		docker-compose exec mysql \
			mysql -uroot -psecret prc_hub
		```

---

# ベンチマーク測定時のログを集める

## `alp` とは

JSON または LTSV 形式のアクセスログからエンドポイントそれぞれの呼び出し回数や応答時間などを集計しするツール。

- [tkuchiki/alp: Access Log Profiler (github.com)](https://github.com/tkuchiki/alp)
- [alpの使い方(基本編) (zenn.dev)](https://zenn.dev/tkuchiki/articles/how-to-use-alp)

今回の題材では`Go`のフレームワークに`echo`が使用されているため、`echo`で`alp`に対応したログを出力する。

## LTSV 形式のログを実装

### ログフォーマットを設定

> `./presentation/echo/main.go:12`
> ```go
> func Start(port uint, jwtIssuer string, jwtSecret string, allowOrigins []string) {
> 	// echoサーバーのインスタンス生成
> 	e := echo.New()
> 	
> 	...
> ```

`./presentation/echo/main.go` に `echo`に関する記述があるためここに書く。

↓ 追記後

```go
func Start(port uint, jwtIssuer string, jwtSecret string, allowOrigins []string) {
	// echoサーバーのインスタンス生成
	e := echo.New()
	
	// LSTV形式のロギング
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time:${time_rfc3339}\thost:${remote_ip}\tstatus:${status}\tmethod:${method}\turi:${uri}\tsize:${bytes_out}\tua:${user_agent}\tapptime:${latency}\tvhost:${host}\treqtime_human:${latency_human}\thost:${host}\n",
		Output: os.Stdout,
	}))

	...
```

### ローカルで実行して確認

```bash
# [localmachine ~/Downloads/eisucon/backend]

# 起動
$ docker-compose up -d

# ログを出力させて確認
$ curl -X GET --head localhost:1323/events

# ログを確認
$ docker-compose logs backend | \
	sed 's/backend-backend-1  | //g' | \
	grep 'time:'
time:2023-01-19T00:23:08Z	host:172.19.0.1	status:200	method:GET	uri:/events	size:579916	ua:curl/7.79.1	apptime:1034316558	vhost:localhost:1323	reqtime_human:1.034316558s	host:localhost:1323

# alpで確認（テーブルのデータが正しく表示されればOK）
$ docker-compose logs backend | \
	sed 's/backend-backend-1  | //g' | \
	grep 'time:' | \
	alp ltsv
+-------+-----+-----+-----+-----+-----+--------+---------+----------------+----------------+----------------+----------------+----------------+----------------+----------------+--------+------------+------------+------------+------------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |   URI   |      MIN       |      MAX       |      SUM       |      AVG       |      P90       |      P95       |      P99       | STDDEV | MIN(BODY)  | MAX(BODY)  | SUM(BODY)  | AVG(BODY)  |
+-------+-----+-----+-----+-----+-----+--------+---------+----------------+----------------+----------------+----------------+----------------+----------------+----------------+--------+------------+------------+------------+------------+
| 1     | 0   | 1   | 0   | 0   | 0   | GET    | /events | 1034316558.000 | 1034316558.000 | 1034316558.000 | 1034316558.000 | 1034316558.000 | 1034316558.000 | 1034316558.000 | 0.000  | 579916.000 | 579916.000 | 579916.000 | 579916.000 |
+-------+-----+-----+-----+-----+-----+--------+---------+----------------+----------------+----------------+----------------+----------------+----------------+----------------+--------+------------+------------+------------+------------+
```

### git commit / git push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'feat(log): add ltsv logging for alp'
git push origin head
```

## サーバーのアプリケーションを更新

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

## ベンチマークを実行し、ログを保存

1. [cttit (tingtt.jp)](https://cttit.tingtt.jp/contests/isucon/1/jobs)でベンチマークを実行

2. ログの確認と保存

	```bash
	# [root@server /usr/local/src/eisucon]
	cd backend
	docker-compose logs -f backend
	```

	↓のようなログが大量に表示されていればOK
	> ```
	> backend-backend-1  | time:2023-01-19T00:51:12Z	host:60.94.171.14	status:200	method:GET	uri:/events/762?embed=user&embed=documents	size:628	ua:Go-http-client/1.1	apptime:8992019	vhost:13.208.215.218:1323	reqtime_human:8.992019ms	host:13.208.215.218:1323
	> backend-backend-1  | time:2023-01-19T00:51:12Z	host:60.94.171.14	status:200	method:GET	uri:/users/41	size:122	ua:Go-http-client/1.1	apptime:4749584	vhost:13.208.215.218:1323	reqtime_human:4.749584ms	host:13.208.215.218:1323
	> backend-backend-1  | time:2023-01-19T00:51:12Z	host:60.94.171.14	status:200	method:GET	uri:/users	size:13100	ua:Go-http-client/1.1	apptime:246191094	vhost:13.208.215.218:1323	reqtime_human:246.191094ms	host:13.208.215.218:1323
	> ```
	
	新しいログが流れてこなくなればベンチマークが終了したと判断し、ログを保存する。
	（`ctrl + C`でログの表示を終了する。）

	```bash
	# [root@server /usr/local/src/eisucon/backend]

	# ユーザー切り替え
	su ec2-user

	# ログをファイルに書き出し（sedコマンドで整形, grepコマンドで必要な部分のみ抽出）
	docker-compose logs backend | \
		sed 's/backend-backend-1  | //g' | \
		grep 'time:' \
		> /home/ec2-user/benchmark_logs/$(date +%s).log
		
	# ファイルを確認（`q`キーで表示を終了）
	less /home/ec2-user/benchmark_logs/1674090251.log
	```

## 保存したログをローカルにダウンロードし、`alp`で集計

```bash
# [localmachine ~/Downloads/eisucon]

# ダウンロード
scp -i ~/Downloads/eisucon/eisucon.pem \
	ec2-user@13.208.215.218:/home/ec2-user/benchmark_logs/1674090251.log \
	./benchmark_ltsv.log

# alp で集計
cat benchmark_ltsv.log | alp ltsv -m '/users/[0-9]+$,/users/[0-9]+/star,/events/[0-9]+$,/events/[0-9]+/documents$,/events/[0-9]+/documents/[0-9]+' -r --sort avg
```

↓のように集計される

>```
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> | COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |               URI               |      MIN      |       MAX       |       SUM        |       AVG       |       P90       |       P95       |       P99       |     STDDEV     | MIN(BODY)  |  MAX(BODY)  |  SUM(BODY)   |  AVG(BODY)  |
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> | 55    | 0   | 55  | 0   | 0   | 0   | GET    | /events                         | 852337946.000 | 14620435702.000 | 692974486820.000 | 12599536124.000 | 13996009616.000 | 14355144273.000 | 14620435702.000 | 1746267527.553 | 579916.000 | 1685069.000 | 82868431.000 | 1506698.745 |
> | 1     | 0   | 1   | 0   | 0   | 0   | POST   | /reset                          | 440917747.000 | 440917747.000   | 440917747.000    | 440917747.000   | 440917747.000   | 440917747.000   | 440917747.000   | 0.000          | 26.000     | 26.000      | 26.000       | 26.000      |
> | 304   | 0   | 304 | 0   | 0   | 0   | GET    | /users                          | 43246514.000  | 337367411.000   | 64833827274.000  | 213269168.664   | 280923577.000   | 305331844.000   | 328082232.000   | 57303671.520   | 13090.000  | 13127.000   | 3984289.000  | 13106.214   |
> | 338   | 0   | 338 | 0   | 0   | 0   | POST   | /users/sign_in                  | 71633730.000  | 167896713.000   | 33724124567.000  | 99775516.470    | 122221643.000   | 129423701.000   | 143434002.000   | 16457519.871   | 241.000    | 241.000     | 81458.000    | 241.000     |
> | 364   | 0   | 364 | 0   | 0   | 0   | POST   | /events                         | 7922853.000   | 73410888.000    | 6497303348.000   | 17849734.473    | 25614275.000    | 28200073.000    | 33725184.000    | 6569905.672    | 396.000    | 432.000     | 149283.000   | 410.118     |
> | 366   | 0   | 366 | 0   | 0   | 0   | GET    | /events/[0-9]+$                 | 5841207.000   | 79701697.000    | 4868403652.000   | 13301649.322    | 18314376.000    | 21165449.000    | 31062334.000    | 5753671.406    | 562.000    | 1023.000    | 247743.000   | 676.893     |
> | 366   | 0   | 366 | 0   | 0   | 0   | GET    | /events/[0-9]+/documents/[0-9]+ | 5506501.000   | 102161571.000   | 4848362680.000   | 13246892.568    | 18241060.000    | 22081150.000    | 34597913.000    | 6937755.430    | 88.000     | 104.000     | 34873.000    | 95.281      |
> | 366   | 0   | 366 | 0   | 0   | 0   | POST   | /users/[0-9]+/star              | 6899961.000   | 62490976.000    | 4766748845.000   | 13023903.948    | 17478514.000    | 20504411.000    | 28955388.000    | 4674604.733    | 16.000     | 17.000      | 5958.000     | 16.279      |
> | 368   | 0   | 368 | 0   | 0   | 0   | GET    | /users/[0-9]+$                  | 2146637.000   | 57282665.000    | 2062257447.000   | 5603960.454     | 8201930.000     | 9393148.000     | 14444932.000    | 3623852.566    | 117.000    | 128.000     | 44934.000    | 122.103     |
> | 740   | 0   | 740 | 0   | 0   | 0   | GET    | /events/[0-9]+/documents$       | 897742.000    | 55614768.000    | 2057932341.000   | 2780989.650     | 4373881.000     | 5272174.000     | 10869283.000    | 2579170.784    | 5.000      | 439.000     | 91981.000    | 124.299     |
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> ```

## Makefileで便利なスクリプトを用意しておく

今後は、

1. `alp`で集計したデータをもとにチューニング
2. ベンチマーク実行

を繰り返すことになるのでよく行う操作はすぐに実行できるようにしておきます。

```Makefile
# Makefileに追記

.PYONY: log-save
log-save: /home/ec2-user/benchmark_logs
	cd ./backend ; \
		docker-compose logs backend | \
		sed 's/backend-backend-1 | //g' | \
		grep 'time:' \
			> /home/ec2-user/benchmark_logs/$$(date +%s).log

.PYONY: log-dl
SSH_FILE ?= ~/Downloads/eisucon.pem
IP ?= 13.208.215.218
log-dl:
	ssh -i ${SSH_FILE} ec2-user@${IP} "ls -t /home/ec2-user/benchmark_logs/ | head -1" | \
		xargs -I SomeString scp -i ${SSH_FILE} ec2-user@${IP}:/home/ec2-user/benchmark_logs/SomeString ./benchmark_ltsv.log

.PYONY: alp
alp:
	cat benchmark_ltsv.log | \
	alp ltsv \
	-m '/users/[0-9]+$$,/users/[0-9]+/star,/events/[0-9]+$$,/events/[0-9]+/documents$$,/events/[0-9]+/documents/[0-9]+' \
	-r --sort avg
```

commit と push

```bash
git add .
git commit -m 'chore: add scripts for alp'
git push origin head
```

---

# クエリチューニング

`alp`の集計データから遅い部分を探る。

↓のグラフは応答速度の平均値で降順になっているので、上にあるものほど遅いと考えられます。

>```
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> | COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |               URI               |      MIN      |       MAX       |       SUM        |       AVG       |       P90       |       P95       |       P99       |     STDDEV     | MIN(BODY)  |  MAX(BODY)  |  SUM(BODY)   |  AVG(BODY)  |
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> | 55    | 0   | 55  | 0   | 0   | 0   | GET    | /events                         | 852337946.000 | 14620435702.000 | 692974486820.000 | 12599536124.000 | 13996009616.000 | 14355144273.000 | 14620435702.000 | 1746267527.553 | 579916.000 | 1685069.000 | 82868431.000 | 1506698.745 |
> | 1     | 0   | 1   | 0   | 0   | 0   | POST   | /reset                          | 440917747.000 | 440917747.000   | 440917747.000    | 440917747.000   | 440917747.000   | 440917747.000   | 440917747.000   | 0.000          | 26.000     | 26.000      | 26.000       | 26.000      |
> | 304   | 0   | 304 | 0   | 0   | 0   | GET    | /users                          | 43246514.000  | 337367411.000   | 64833827274.000  | 213269168.664   | 280923577.000   | 305331844.000   | 328082232.000   | 57303671.520   | 13090.000  | 13127.000   | 3984289.000  | 13106.214   |
> | 338   | 0   | 338 | 0   | 0   | 0   | POST   | /users/sign_in                  | 71633730.000  | 167896713.000   | 33724124567.000  | 99775516.470    | 122221643.000   | 129423701.000   | 143434002.000   | 16457519.871   | 241.000    | 241.000     | 81458.000    | 241.000     |
> | 364   | 0   | 364 | 0   | 0   | 0   | POST   | /events                         | 7922853.000   | 73410888.000    | 6497303348.000   | 17849734.473    | 25614275.000    | 28200073.000    | 33725184.000    | 6569905.672    | 396.000    | 432.000     | 149283.000   | 410.118     |
> | 366   | 0   | 366 | 0   | 0   | 0   | GET    | /events/[0-9]+$                 | 5841207.000   | 79701697.000    | 4868403652.000   | 13301649.322    | 18314376.000    | 21165449.000    | 31062334.000    | 5753671.406    | 562.000    | 1023.000    | 247743.000   | 676.893     |
> | 366   | 0   | 366 | 0   | 0   | 0   | GET    | /events/[0-9]+/documents/[0-9]+ | 5506501.000   | 102161571.000   | 4848362680.000   | 13246892.568    | 18241060.000    | 22081150.000    | 34597913.000    | 6937755.430    | 88.000     | 104.000     | 34873.000    | 95.281      |
> | 366   | 0   | 366 | 0   | 0   | 0   | POST   | /users/[0-9]+/star              | 6899961.000   | 62490976.000    | 4766748845.000   | 13023903.948    | 17478514.000    | 20504411.000    | 28955388.000    | 4674604.733    | 16.000     | 17.000      | 5958.000     | 16.279      |
> | 368   | 0   | 368 | 0   | 0   | 0   | GET    | /users/[0-9]+$                  | 2146637.000   | 57282665.000    | 2062257447.000   | 5603960.454     | 8201930.000     | 9393148.000     | 14444932.000    | 3623852.566    | 117.000    | 128.000     | 44934.000    | 122.103     |
> | 740   | 0   | 740 | 0   | 0   | 0   | GET    | /events/[0-9]+/documents$       | 897742.000    | 55614768.000    | 2057932341.000   | 2780989.650     | 4373881.000     | 5272174.000     | 10869283.000    | 2579170.784    | 5.000      | 439.000     | 91981.000    | 124.299     |
> +-------+-----+-----+-----+-----+-----+--------+---------------------------------+---------------+-----------------+------------------+-----------------+-----------------+-----------------+-----------------+----------------+------------+-------------+--------------+-------------+
> ```

`./presentation/echo/main.go:59` に各エンドポイントに関する記述があるのでここから探索する。

> ```go
> 	...
>  
> 	// handlerの登録
> 	var server *Server
> 
> 	// ↓ スコア測定に直接関係するエンドポイント
> 	e.GET("/events", server.GetEvents)
> 	e.POST("/events", server.PostEvents)
> 	e.GET("/events/:id", server.GetEventsId)
> 	e.GET("/events/:id/documents", server.GetEventsIdDocuments)
> 	e.POST("/events/:id/documents", server.PostEventsIdDocuments)
> 	e.GET("/events/:id/documents/:document_id", server.GetEventsIdDocumentsDocumentId)
> 	e.POST("/reset", server.PostReset)
> 	e.GET("/users", server.GetUsers)
> 	e.POST("/users/sign_in", server.PostUsersSignIn)
> 	e.GET("/users/:id", server.GetUsersId)
> 	e.POST("/users/:id/star", server.PostUsersIdStar)
> 
> 	// ↓ スコアに直接関係しないため他部分の変更による影響がない場合は原則変更しなくて構わない
> 	e.DELETE("/events/:id", server.DeleteEventsId)
> 	e.PATCH("/events/:id", server.PatchEventsId)
> 	e.DELETE("/events/:id/documents/:document_id", server.DeleteEventsIdDocumentsDocumentId)
> 	e.PATCH("/events/:id/documents/:document_id", server.PatchEventsIdDocumentsDocumentId)
> 	e.DELETE("/users", server.DeleteUsers)
> 	e.POST("/users", server.PostUsers)
> 	e.DELETE("/users/:id", server.DeleteUsersId)
> 	e.PATCH("/users/:id", server.PatchUsersId)
> 
> 	// echoサーバーの起動
> 	logger.Logger().Fatal(e.Start(fmt.Sprintf(":%d", port)))
> 	
> 	...


## `/events`

> [!WARNING] 一番難易度が高いエンドポイントなので`/users`、`/users/[id]`を先にするほうがおすすめ

### コードを追っていく。　

> [!INFO] IDEで関数を`ctrl + Click`するとソースコードを掘っていける。（macは`cmd + Click`）

> `./presentation/echo/main.go:63`
> ```go
> 	e.GET("/events", server.GetEvents)
> ```

> `./presentation/echo/handler_events.get.go:11`
> ```go
> // (GET /events)
> func (*Server) GetEvents(ctx echo.Context) error {
> 	// Get jwt claim
> 	var jwtId *int64
> 	jcc, err := jwt.Check(ctx)
> 	if err == nil {
> 		jwtId = &jcc.Id
> 	}
> 
> 	// Bind query
> 	query := new(event.GetEventListQueryParam)
> 	type Query struct {
> 		Published       *bool   `query:"published"`
> 		Name            *string `query:"name"`
> 		NameContain     *string `query:"name_contain"`
> 		Location        *string `query:"location"`
> 		LocationContain *string `query:"location_contain"`
> 	}
> 	queryTmp := new(Query)
> 	if err := ctx.Bind(queryTmp); err != nil {
> 		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
> 	}
> 	v := ctx.QueryParams()
> 	embed := v["embed"]
> 	query.Embed = &embed
> 	query.Name = queryTmp.Name
> 	query.NameContain = queryTmp.NameContain
> 	query.Location = queryTmp.Location
> 	query.LocationContain = queryTmp.LocationContain
> 
> 	// Get events
> 	events, err := event.GetEventList(*query, jwtId)
> 	if err != nil {
> 		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
> 	}
> 
> 	if events == nil {
> 		return JSONPretty(ctx, http.StatusOK, []interface{}{})
> 	}
> 	return JSONPretty(ctx, http.StatusOK, events)
> }
> ```

- `41`行目あたりが根本処理っぽい

> `./application/event/event.go:50`
> ```go
> func GetEventList(q GetEventListQueryParam, requestUserId *int64) ([]event.EventEmbed, error) {
> 	u := new(userDomain.User)
> 
> 	if requestUserId != nil {
> 		// リクエスト元のユーザーを取得
> 		var u2 userDomain.User
> 		u2, err := user.Get(*requestUserId)
> 		if err != nil {
> 			return nil, err
> 		}
> 		u = &u2
> 	} else if requestUserId == nil {
> 		// リクエストユーザーが指定されていない場合は最小権限のユーザーを仮使用
> 		u = &userDomain.User{
> 			Id:                  0,
> 			PostEventAvailabled: false,
> 			Manage:              false,
> 			Admin:               false,
> 		}
> 	}
> 
> 	return event.GetEventList(q, *u)
> }
> ```

- `71`行目で更に関数を呼び出している

> `./domain/model/event/event_get_list.go:18`
> ```go
> func GetEventList(q GetEventListQueryParam, requestUser user.User) ([]EventEmbed, error) {
> 	// MySQLサーバーに接続
> 	db, err := OpenMysql()
> 	if err != nil {
> 		return nil, err
> 	}
> 	// return時にMySQLサーバーとの接続を閉じる
> 	defer db.Close()
> 
> 	embedUser := false
> 	embedDocuments := false
> 	if q.Embed != nil {
> 		for _, e := range *q.Embed {
> 			if e == "user" {
> 				embedUser = true
> 			}
> 			if e == "documents" {
> 				embedDocuments = true
> 			}
> 		}
> 	}
> 
> 	// `Event`リストを取得
> 
> 	// 取得用変数
> 	events := []EventEmbed{}
> 
> 	// クエリを作成
> 	query := "SELECT * FROM events WHERE"
> 	queryParams := []interface{}{}
> 	if q.Name != nil {
> 		// イベント名の一致で絞り込み
> 		query += " name = ? AND"
> 		queryParams = append(queryParams, *q.Name)
> 	}
> 	if q.NameContain != nil {
> 		// イベント名に文字列が含まれるかで絞り込み
> 		query += " name LIKE ? AND"
> 		queryParams = append(queryParams, "%"+*q.NameContain+"%")
> 	}
> 	if q.Location != nil {
> 		// `Location`の一致で絞り込み
> 		query += " location = ? AND"
> 		queryParams = append(queryParams, *q.Location)
> 	}
> 	if q.LocationContain != nil {
> 		// `Location`に文字列が含まれるかで絞り込み
> 		query += " location LIKE ? AND"
> 		queryParams = append(queryParams, "%"+*q.LocationContain+"%")
> 	}
> 	if q.Published != nil {
> 		// `Published`で絞り込み
> 		query += " published = ?"
> 		queryParams = append(queryParams, *q.Published)
> 	}
> 	// 不要な末尾の句を切り取り
> 	query = strings.TrimSuffix(query, "WHERE")
> 	query = strings.TrimSuffix(query, "AND")
> 
> 	// 実行
> 	r1, err := db.Query(query, queryParams...)
> 	if err != nil {
> 		return nil, err
> 	}
> 	defer r1.Close()
> 
> 	// １行ずつ処理
> 	for r1.Next() {
> 		// 一時変数に割当
> 		var (
> 			id          int64
> 			name        string
> 			description *string
> 			location    *string
> 			published   bool
> 			completed   bool
> 			userId      int64
> 		)
> 		err = r1.Scan(&id, &name, &description, &location, &published, &completed, &userId)
> 		if err != nil {
> 			return nil, err
> 		}
> 		// 配列追加用変数
> 		event := EventEmbed{
> 			Event: Event{
> 				Id:          id,
> 				Name:        name,
> 				Description: description,
> 				Location:    location,
> 				Datetimes:   []EventDatetime{},
> 				Published:   published,
> 				Completed:   completed,
> 				UserId:      userId,
> 			},
> 		}
> 
> 		// `EventDatetime`を取得
> 		r2, err := db.Query("SELECT * FROM event_datetimes WHERE event_id = ?", id)
> 		if err != nil {
> 			return nil, err
> 		}
> 		defer r2.Close()
> 		for r2.Next() {
> 			var (
> 				eId   string
> 				start *time.Time
> 				end   *time.Time
> 			)
> 			err = r2.Scan(&eId, &start, &end)
> 			if err != nil {
> 				return nil, err
> 			}
> 			// 配列に追加
> 			event.Event.Datetimes = append(event.Event.Datetimes, EventDatetime{*start, *end})
> 		}
> 
> 		if embedUser {
> 			// `User`を取得
> 			u, err := user.Get(userId)
> 			if err != nil {
> 				return nil, err
> 			}
> 			// 変数に追加
> 			event.User = &u
> 		}
> 
> 		if embedDocuments {
> 			// `Documents`を取得
> 			ed, err := GetDocumentList(GetDocumentQueryParam{EventId: &id})
> 			if err != nil {
> 				return nil, err
> 			}
> 			event.Documents = &ed
> 		}
> 
> 		events = append(events, event)
> 	}
> 
> 	return events, nil
> }
> ```

- このファイルのコードで`SQL`が実行されている。
- クエリ例
	1. イベント
		```sql
		# 46行目からクエリ作成
		# 78行目で実行
		SELECT * FROM events;
		```
		- `*`が使用されているため無駄な読み込みが発生している可能性がある。

	1. イベント日時
		```sql
		# 115行目
		SELECT * FROM event_datetimes WHERE event_id = <イベントID>;
		```
		- イベント一覧の`for`文内で実行されているため `N+1`問題

	1. ユーザー
		```sql
		# 136行目の`user.Get()`の内部で実行
		SELECT * FROM users WHERE id = <ユーザーID>;
		```
		- イベント一覧の`for`文内で実行されているため `N+1`問題
		- `*`が使用されているため無駄な読み込みが発生している可能性がある。
		- コネクションを再利用していなさそう

	1. イベントドキュメント
		```sql
		# 146行目の`GetDocumentList()`の内部で実行
		SELECT * FROM documents WHERE event_id = <イベントID>;
		```
		- イベント一覧の`for`文内で実行されているため `N+1`問題
		- `*`が使用されているため無駄な読み込みが発生している可能性がある。
		- コネクションを再利用していなさそう

### クエリ変更

`./.mysql/init.sql`を見て照らし合わせながら書き換える。

N+1 問題（このコードでは N+N+N+1 ）だが、`JOIN`では対応が難しい（`for`文内で実行されているクエリで得られるレコードが2件以上の場合がある）ので、

- `UNION`句を使用したクエリ
	- `JOIN`では横方向の結合だが、`UNION`では縦方向の結合ができる。
	- N+N+N+1 回 → 1 回
- `Eager loading`
	- N+1 の内の N 回分のクエリを1回にし、得られたデータはインメモリで結合する。
	- N+N+N+1 回 → 4 回

が対応例としてあげられる。

今回は`UNION`を使用したクエリでのチューニングを行っています。

```sql
SELECT * FROM (
	WITH params AS ( SELECT id AS event_id FROM events )
		SELECT
			e.id, e.name, e.description, e.location, e.published, e.completed, e.user_id,
			null AS start, null AS end,
			null AS doc_id, null AS doc_name, null AS doc_url,
			u.id AS u_id, u.name AS u_name, u.email AS u_email, u.post_event_availabled, u.manage, u.admin, u.twitter_id, u.github_username
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		WHERE e.id IN (SELECT event_id FROM params)
	UNION ALL
		SELECT
			dt.event_id, null, null, null, null, null, null,
			dt.start, dt.end,
			null, null, null,
			null, null, null, null, null, null, null, null
		FROM event_datetimes dt
		WHERE dt.event_id IN (SELECT event_id FROM params)
	UNION ALL
		SELECT
			doc.event_id, null, null, null, null, null, null,
			null, null,
			doc.id, doc.name, doc.url,
			null, null, null, null, null, null, null, null
		FROM documents doc
		WHERE doc.event_id IN (SELECT event_id FROM params)
	) AS e
ORDER BY id, name IS NULL ASC
```

### `go`のコードを書き換える

省略（下の url から変更例を見れます。）

### サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: update query to get events'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: update query to get events · tingtt/e-isucon@1e20ab3 (github.com)](https://github.com/tingtt/e-isucon/commit/1e20ab36c1cc866a43dab1b49e6499d1efed7e48)

## `/users/[id]`, `/users` その1

### コードを見てみる

> `./presentation/handler_user.get.go`
> ```go
> // (GET /users/{id})
> func (*Server) GetUsersId(ctx echo.Context) error {
> 	// Get jwt claim
> 	_, err := jwt.CheckProvided(ctx)
> 	// jcc, err := jwt.CheckProvided(ctx)
> 	if err != nil {
> 		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
> 	}
> 
> 	// Bind id
> 	var id Id
> 	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
> 	if err != nil {
> 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
> 	}
> 
> 	// Get user
> 	u, err := user.Get(id)
> 	if err != nil {
> 		return JSONMessage(ctx, user.ErrToCode(err), err.Error())
> 	}
> 
> 	return JSONPretty(ctx, http.StatusOK, u)
> }
> ```

- 30行目の`user.Get(id)`が根本の処理っぽい

> `./application/user/get.go:5`
> ```go
> func Get(id int64) (_ user.User, err error) {
> 	return user.Get(id)
> }
> ```

> `./domain/model/user/get.go`
> ```go
> func Get(id int64) (User, error) {
> 	// MySQLサーバーに接続
>	db, err := OpenMysql()
>
> 	if err != nil {
> 		return User{}, err
> 	}
> 	// return時にMySQLサーバーとの接続を閉じる
> 	defer db.Close()
> 	
> 	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
> 	r, err := db.Query("SELECT * FROM users WHERE id = ?", id)
> 	if err != nil {
> 		return User{}, err
> 	}
> 	defer r.Close()
> 	
> 	if !r.Next() {
> 		// 1行もレコードが無い場合
> 		// not found
> 		err = ErrUserNotFound
> 		return User{}, err
> 	}
> 	
> 	// 一時変数に割り当て
> 	var (
>		id2 int64
>		name string
>		email string
>		password string
>		postEventAvailabled bool
>		manage bool
>		admin bool
>		twitterId *string
>		githubUsername *string
> 	)
> 	
> 	err = r.Scan(
> 		&id2, &name, &email, &password, &postEventAvailabled,
> 		&manage, &admin, &twitterId, &githubUsername,
> 	)
> 	if err != nil {
> 		return User{}, err
> 	}
> 	
> 	// スター数を取得
> 	var count uint64 = 0
> 	r2, err := db.Query("SELECT COUNT(*) FROM user_stars WHERE target_user_id = ?", id)
> 	if err != nil {
> 	return User{}, err
> 	}
> 	defer r2.Close()
> 	if !r2.Next() {
> 		return User{}, ErrConflictUserStars
> 	}
> 	err = r2.Scan(&count)
> 	if err != nil {
> 		return User{}, err
> 	}
> 	
> 	u := User{
> 	Id: id,
> 		Name: name,
> 		Email: email,
> 		Password: password,
> 		StarCount: count,
> 		PostEventAvailabled: postEventAvailabled,
> 		Manage: manage,
> 		Admin: admin,
> 		TwitterId: twitterId,
> 		GithubUsername: githubUsername,
> 	}
> 	
> 	return u, nil
> }
> ```

- クエリ例
	1. ユーザー取得
		```sql
		SELECT * FROM users WHERE id = <ユーザーID>;
		```
	1. スター数取得
		```sql
		SELECT COUNT(*) FROM user_stars WHERE target_user_id = <ユーザーID>;
		```
		- `JOIN`句で２つのクエリを１つにまとめられそう。
		- `users`、`user_start`に分割されたテーブルを結合すれば`JOIN`句も無くすことができそう。
		- `/users`の処理ではユーザー一覧の`for`文内で実行されているため `N+1`問題

### クエリの変更

####  `./.mysql/init.sql`を見て照らし合わせながら書き換える。

`users`テーブルを取得する際に`user_stars`テーブルを結合し、`COUNT`句でスター数を算出するように変更。

```sql
# 全件取得
SELECT
	u.id, u.name, u.email, u.password,
	u.post_event_availabled, u.manage, u.admin,
	u.twitter_id, u.github_username,
	COUNT(s.target_user_id) AS star_count
FROM
	users u
LEFT JOIN
	user_stars s ON u.id = s.target_user_id
GROUP BY u.id;
```

```sql
# １件取得
SELECT
	u.id, u.name, u.email, u.password,
	u.post_event_availabled, u.manage, u.admin,
	u.twitter_id, u.github_username,
	COUNT(s.target_user_id) AS star_count
FROM
	users u
LEFT JOIN
	user_stars s ON u.id = s.target_user_id
WHERE u.id = <ユーザーID>
GROUP BY u.id;
```

```sql
# 1件取得（メールアドレスで検索）
SELECT
	u.id, u.name, u.email, u.password,
	u.post_event_availabled, u.manage, u.admin,
	u.twitter_id, u.github_username,
	COUNT(s.target_user_id) AS star_count
FROM
	users u
LEFT JOIN
	user_stars s ON u.id = s.target_user_id
WHERE u.email = <メールアドレス>
GROUP BY u.id;
```

#### SQLを実行してみる。

```bash
# [localmachine ~/Downloads/eisucon/backend]
docker-compose exec mysql mysql -uroot -psecret prc_hub
```
```sql
SELECT
	u.id, u.name, u.email, u.password,
	u.post_event_availabled, u.manage, u.admin,
	u.twitter_id, u.github_username,
	COUNT(s.target_user_id) AS star_count
FROM
	users u
LEFT JOIN
	user_stars s ON u.id = s.target_user_id
GROUP BY u.id;
```

↓のような結果になればOK

```
+-----+---------------------+--------------------------------+--------------------------------------------------------------+-----------------------+--------+-------+------------+-----------------+------------+
| id  | name                | email                          | password                                                     | post_event_availabled | manage | admin | twitter_id | github_username | star_count |
+-----+---------------------+--------------------------------+--------------------------------------------------------------+-----------------------+--------+-------+------------+-----------------+------------+
|   1 | throbbing-pond      | throbbing-pond@prchub.com      | $2a$10$N6XHQHo9zOIdkv.SiqeGgeRouYBu4VpuCizK25BBmM4.pHyJHwwBi |                     1 |      1 |     1 | NULL       | NULL            |          1 |
|   2 | old-wood            | old-wood@prchub.com            | $2a$10$jiRWkv27tJqOPB8RyW/WzuzO1/d05KoqWrYhhJyew.WRInfAW9nq2 |                     1 |      1 |     1 | NULL       | NULL            |          2 |
|   3 | throbbing-moon      | throbbing-moon@prchub.com      | $2a$10$mvVpOERZ49q1YjawGSBXae1qHv.fr5vRTtpYXRAzzrnOINJccLsPa |                     1 |      1 |     1 | NULL       | NULL            |          3 |
|   4 | white-pond          | white-pond@prchub.com          | $2a$10$oR2AXxuXa8zOyOE7xYo9oOuL4Jgea3T7DnP0YFjCh6a0yLS1h25tm |                     1 |      1 |     1 | NULL       | NULL            |          4 |
|   5 | spring-darkness     | spring-darkness@prchub.com     | $2a$10$/LkM3OLtgNhMYyhmyaEz2OweIs4Z43MNTRsOxpSV60yH5O1S43Fs6 |                     1 |      1 |     1 | NULL       | NULL            |          5 |
|   6 | nameless-wildflower | nameless-wildflower@prchub.com | $2a$10$OVR3Qh//EoXdu68mBbxlauEPJUc2ZBr.FIYaKWifCajjGbYCwOKQS |                     1 |      1 |     1 | NULL       | NULL            |          6 |
|   7 | lingering-meadow    | lingering-meadow@prchub.com    | $2a$10$OrRRXE.RGlj/Vdt9bL9.EuX3GISas8y1rEBT/UKZ8ssn9U43vJ5.C |                     1 |      1 |     1 | NULL       | NULL            |          7 |

...
```

#### `go`のコードを書き換える。

##### `domain/model/user/get.go`

- クエリを１つにまとめる
- 不要な変数定義を削除

```go
func Get(id int64) (User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return User{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// usersテーブルからidが一致する行を取得し、変数eに代入する
	r, err := db.Query(
		`SELECT
			u.id, u.name, u.email, u.password,
			u.post_event_availabled, u.manage, u.admin,
			u.twitter_id, u.github_username,
			COUNT(s.target_user_id) AS star_count
		FROM
			users u
		LEFT JOIN
			user_stars s ON u.id = s.target_user_id
		WHERE u.id = ?
		GROUP BY u.id`,
		id,
	)
	if err != nil {
		return User{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		return User{}, ErrUserNotFound
	}

	// 変数に割り当て
	u := User{}
	err = r.Scan(
		&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
		&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetByEmail(email string) (User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return User{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
	r, err := db.Query(
		`SELECT
			u.id, u.name, u.email, u.password,
			u.post_event_availabled, u.manage, u.admin,
			u.twitter_id, u.github_username,
			COUNT(s.target_user_id) AS star_count
		FROM
			users u
		LEFT JOIN
			user_stars s ON u.id = s.target_user_id
		WHERE u.email = ?
		GROUP BY u.id`,
		email,
	)
	if err != nil {
		return User{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		return User{}, ErrUserNotFound
	}

	// 変数に割り当て
	u := User{}
	err = r.Scan(
		&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
		&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

```

##### `domain/model/user/get_list.go`

- クエリを１つにまとめる
- 不要な変数定義を削除

```go
func GetList(q GetUserListQueryParam) ([]User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return nil, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// クエリを作成
	query :=
		`SELECT
			u.id, u.name, u.email, u.password,
			u.post_event_availabled, u.manage, u.admin,
			u.twitter_id, u.github_username,
			COUNT(s.target_user_id) AS star_count
		FROM
			users u
		LEFT JOIN
			user_stars s ON u.id = s.target_user_id
		WHERE`
	queryParams := []interface{}{}
	if q.PostEventAvailabled != nil {
		// 権限で絞り込み
		query += " u.post_event_availabled = ? AND"
		queryParams = append(queryParams, *q.PostEventAvailabled)
	}
	if q.Manage != nil {
		// 権限で絞り込み
		query += " u.manage = ? AND"
		queryParams = append(queryParams, *q.Manage)
	}
	if q.Admin != nil {
		// 権限で絞り込み
		query += " u.admin = ? AND"
		queryParams = append(queryParams, *q.Admin)
	}
	if q.Name != nil {
		// ドキュメント名の一致で絞り込み
		query += " u.name = ? AND"
		queryParams = append(queryParams, *q.Name)
	}
	if q.NameContain != nil {
		// ドキュメント名に文字列が含まれるかで絞り込み
		query += " u.name LIKE ?"
		queryParams = append(queryParams, "%"+*q.NameContain+"%")
	}
	// 不要な末尾の句を切り取り
	query = strings.TrimSuffix(query, "WHERE")
	query = strings.TrimSuffix(query, " AND")

	// `users`テーブルからを取得
	r, err := db.Query(query+" GROUP BY u.id", queryParams...)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 取得したテーブルを１行ずつ処理
	// 配列`users`に代入する
	var users []User
	for r.Next() {
		// 変数に割り当て
		u := User{}
		err = r.Scan(
			&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
			&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
		)
		if err != nil {
			return nil, err
		}

		// 配列に追加
		users = append(users, u)
	}

	return users, nil
}
```

### サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: update query to get users'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: update query to get users · tingtt/e-isucon@1df112b (github.com)](https://github.com/tingtt/e-isucon/commit/1df112b652b9ece2956fee8010a608e1d095f9aa)

## `/users/[id]`, `/users` その2

「その１」では`JOIN`句を使ってクエリをまとめました。
`JOIN`句を使用すると`N+1`の状態よりは高速化できますが、`JOIN`句も使わなくて良い構造にするのが理想です。

- `JOIN`をしているため遅い
- `COUNT`, `GROUP BY`をしているため遅い

### テーブル構造を調べる

> `./.mysql/init.sql`
> ```sql
> --
> -- Table structure for table `users`
> --
> CREATE TABLE `users` (
>   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
>   `name` varchar(255) NOT NULL,
>   `email` varchar(255) NOT NULL UNIQUE,
>   `password` varchar(255) NOT NULL,
>   `post_event_availabled` tinyint(1) NOT NULL DEFAULT 0,
>   `manage` tinyint(1) NOT NULL DEFAULT 0,
>   `admin` tinyint(1) NOT NULL DEFAULT 0,
>   `twitter_id` varchar(255),
>   `github_username` varchar(255),
>   PRIMARY KEY (`id`)
> );
> 
> --
> -- Table structure for table `user_stars`
> --
> CREATE TABLE `user_stars` (
>   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
>   `target_user_id` int(255) UNSIGNED NOT NULL,
>   PRIMARY KEY (`id`)
> );
> ```

- SQLで`COUNT`関数を実行していることからも分かる通り、ユーザーのスター数は`user_stars`テーブルのレコード件数で表現されている。

### `user_stars`テーブルを使用しているコードを探す

IDE（VSCodeなど）で全文検索をし、`user_stars`テーブルを使用しているコードを探す。

- `./domain/model/users/get.go`
- `./domain/model/users/get_list.go`
- `./domain/model/users/add_star.go`
- `./domain/model/eisucon/migrate.sql`
- `./.mysql/init.sql`

これらに含まれている`SELECT`文や`INESRT`文を書き換える想定でテーブルを変更する。

### SQLの変更

#### `users`テーブル (`./.mysql/init.sql`)

`user_start`テーブルを統合する。
- `star_count`を追加する

```diff
  --
  -- Table structure for table `users`
  --
  CREATE TABLE `users` (
    `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL UNIQUE,
    `password` varchar(255) NOT NULL,
    `post_event_availabled` tinyint(1) NOT NULL DEFAULT 0,
    `manage` tinyint(1) NOT NULL DEFAULT 0,
    `admin` tinyint(1) NOT NULL DEFAULT 0,
    `twitter_id` varchar(255),
    `github_username` varchar(255),
+   `star_count` int(255) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
  );
- 
- --
- -- Table structure for table `user_stars`
- --
- CREATE TABLE `user_stars` (
-   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
-   `target_user_id` int(255) UNSIGNED NOT NULL,
-   PRIMARY KEY (`id`)
- );
```

#### ユーザー一覧取得 (`./domain/model/users/get_list.go`)

```diff
- SELECT
- 	u.id, u.name, u.email, u.password,
- 	u.post_event_availabled, u.manage, u.admin,
- 	u.twitter_id, u.github_username,
- 	COUNT(s.target_user_id) AS star_count
- FROM
- 	users u
- LEFT JOIN
- 	user_stars s ON u.id = s.target_user_id
- GROUP BY u.id;```sql
+ SELECT
+ 	id, name, email, password,
+ 	post_event_availabled, manage, admin,
+ 	twitter_id, github_username, star_count
+ FROM users;
```

#### ユーザー取得 (`./domain/model/users/get.go`)

13行目

```diff
- SELECT
- 	u.id, u.name, u.email, u.password,
- 	u.post_event_availabled, u.manage, u.admin,
- 	u.twitter_id, u.github_username,
- 	COUNT(s.target_user_id) AS star_count
- FROM
- 	users u
- LEFT JOIN
- 	user_stars s ON u.id = s.target_user_id
- WHERE u.id = ?
- GROUP BY u.id;
+ SELECT
+ 	id, name, email, password,
+ 	post_event_availabled, manage, admin,
+ 	twitter_id, github_username, star_count
+ FROM users
+ WHERE id = ?;
```

56行目

```diff
- SELECT
- 	u.id, u.name, u.email, u.password,
- 	u.post_event_availabled, u.manage, u.admin,
- 	u.twitter_id, u.github_username,
- 	COUNT(s.target_user_id) AS star_count
- FROM
- 	users u
- LEFT JOIN
- 	user_stars s ON u.id = s.target_user_id
- WHERE u.id = ?
- GROUP BY u.id;
+ SELECT
+ 	id, name, email, password,
+ 	post_event_availabled, manage, admin,
+ 	twitter_id, github_username, star_count
+ FROM users
+ WHERE id = ?;
```


#### スター追加 (`./domain/model/users/add_star.go`)

```diff
- INSERT INTO user_stars (target_user_id) VALUES (?);
+ UPDATE users SET star_count = star_count + 1 WHERE id = ?;
```

#### スター数取得 (`./domain/model/users/add_star.go`)

```diff
- SELECT COUNT(*) FROM user_stars WHERE target_user_id = ?
+ SELECT star_count FROM users WHERE id = ?;
```

### データの変更

`./domain/model/eisucon/migrate.sql` には初期データを登録するクエリが記述されているため、テーブル構造を変更する場合は変更する必要がある。

こういった MySQL の初期データを変更する場合、自分で書き換えると時間がかかるため、MySQL の機能を使ってテーブル構造変更後の`INSERT`文を生成する。

1. `mysql`に接続する。
	```bash
	# [localmachine ~/Downloads/eisucon]
	docker-compose up -d
	docker-compose exec mysql mysql -uroot -psecret prc_hub
	```
1. `ALTER TABLE`文でテーブル構造を変更
	```sql
	ALTER TABLE
		users
	ADD
		star_count int(255) UNSIGNED NOT NULL DEFAULT 0
	AFTER `github_username`;
	```
1. `UPDATE`文で追加したカラムにデータを入れる
	```sql
	UPDATE users
	LEFT JOIN (
		SELECT
			target_user_id, COUNT(target_user_id) AS star_count
		FROM user_stars GROUP BY 1
	) AS counts
	ON users.id = counts.target_user_id
	SET users.star_count = COALESCE(counts.star_count, 0);
	```
1. `sqldump`で`INSERT`文を生成する
	```bash
	# mysqldump で usersテーブルのデータからSQL文を生成し、users.sqlに書き込む
	docker-compose exec mysql mysqldump -psecret prc_hub users > users.sql
	```

1. 生成した`INSERT`文を`./domain/model/eisucon/migrate.sql`に貼り付け、不要な SQL を削除する。

	```diff
	  DELETE FROM `documents`;
	  DELETE FROM `event_datetimes`;
	  DELETE FROM `events`;
	- DELETE FROM `user_stars`;
	  DELETE FROM `users`;

	  ALTER TABLE `documents` auto_increment = 1;
	  ALTER TABLE `event_datetimes` auto_increment = 1;
	  ALTER TABLE `events` auto_increment = 1;
	- ALTER TABLE `user_stars` auto_increment = 1;
	  ALTER TABLE `users` auto_increment = 1;

	+ INSERT INTO `users` VALUES ...;
	- INSERT INTO `users` (name, email, password, post_event_availabled, manage, admin) VALUES
	- ("throbbing-pond", "throbbing-pond@prchub.com", "$2a$10$N6XHQHo9zOIdkv.SiqeGgeRouYBu4VpuCizK25BBmM4.pHyJHwwBi", 1, 1, 1),
	- ("old-wood", "old-wood@prchub.com", "$2a$10$jiRWkv27tJqOPB8RyW/WzuzO1/d05KoqWrYhhJyew.WRInfAW9nq2", 1, 1, 1),
	- ...;
	- 
	- INSERT INTO `user_stars` (target_user_id) VALUES
	- 	(1), (2), (2), (3), (3), (3), (4), (4), (4), (4),
	- 	(5), (5), (5), (5), (5), (6), (6), (6), (6), (6), (6),
	- 	(7), (7), (7), (7), (7), (7), (7),
	- 	(8), (8), (8), (8), (8), (8), (8), (8),
	- 	(9), (9), (9), (9), (9), (9), (9), (9), (9),
	- 	(11), (12), (12), (13), (13), (13), (14), (14), (14), (14),
	- 	(15), (15), (15), (15), (15), (16), (16), (16), (16), (16), (16),
	- 	(17), (17), (17), (17), (17), (17), (17),
	- 	(18), (18), (18), (18), (18), (18), (18), (18),
	- 	(19), (19), (19), (19), (19), (19), (19), (19), (19),
	- 	(21), (22), (22), (23), (23), (23), (24), (24), (24), (24),
	- 	(25), (25), (25), (25), (25), (26), (26), (26), (26), (26), (26),
	- 	(27), (27), (27), (27), (27), (27), (27),
	- 	(28), (28), (28), (28), (28), (28), (28), (28),
	- 	(29), (29), (29), (29), (29), (29), (29), (29), (29),
	- 	(31), (32), (32), (33), (33), (33), (34), (34), (34), (34),
	- 	(35), (35), (35), (35), (35), (36), (36), (36), (36), (36), (36),
	- 	(37), (37), (37), (37), (37), (37), (37),
	- 	(38), (38), (38), (38), (38), (38), (38), (38),
	- 	(39), (39), (39), (39), (39), (39), (39), (39), (39),
	- 	(41), (42), (42), (43), (43), (43), (44), (44), (44), (44),
	- 	(45), (45), (45), (45), (45), (46), (46), (46), (46), (46), (46),
	- 	(47), (47), (47), (47), (47), (47), (47),
	- 	(48), (48), (48), (48), (48), (48), (48), (48),
	- 	(49), (49), (49), (49), (49), (49), (49), (49), (49),
	- 	(51), (52), (52), (53), (53), (53), (54), (54), (54), (54),
	- 	(55), (55), (55), (55), (55), (56), (56), (56), (56), (56), (56),
	- 	(57), (57), (57), (57), (57), (57), (57),
	- 	(58), (58), (58), (58), (58), (58), (58), (58),
	- 	(59), (59), (59), (59), (59), (59), (59), (59), (59),
	- 	(61), (62), (62), (63), (63), (63), (64), (64), (64), (64),
	- 	(65), (65), (65), (65), (65), (66), (66), (66), (66), (66), (66),
	- 	(67), (67), (67), (67), (67), (67), (67),
	- 	(68), (68), (68), (68), (68), (68), (68), (68),
	- 	(69), (69), (69), (69), (69), (69), (69), (69), (69),
	- 	(71), (72), (72), (73), (73), (73), (74), (74), (74), (74),
	- 	(75), (75), (75), (75), (75), (76), (76), (76), (76), (76), (76),
	- 	(77), (77), (77), (77), (77), (77), (77),
	- 	(78), (78), (78), (78), (78), (78), (78), (78),
	- 	(79), (79), (79), (79), (79), (79), (79), (79), (79),
	- 	(81), (82), (82), (83), (83), (83), (84), (84), (84), (84),
	- 	(85), (85), (85), (85), (85), (86), (86), (86), (86), (86), (86),
	- 	(87), (87), (87), (87), (87), (87), (87),
	- 	(88), (88), (88), (88), (88), (88), (88), (88),
	- 	(89), (89), (89), (89), (89), (89), (89), (89), (89),
	- 	(91), (92), (92), (93), (93), (93), (94), (94), (94), (94),
	- 	(95), (95), (95), (95), (95), (96), (96), (96), (96), (96), (96),
	- 	(97), (97), (97), (97), (97), (97), (97),
	- 	(98), (98), (98), (98), (98), (98), (98), (98),
	- 	(99), (99), (99), (99), (99), (99), (99), (99), (99);
	```

### サーバーに適応
  
GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: update query to get users'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: update query and DDL related to user · tingtt/e-isucon@3e179a1 (github.com)](https://github.com/tingtt/e-isucon/commit/3e179a11d096ca7f4fd6ddcdd82a20be00aba746)

## `/events/{id}`

### コードを追っていく。　

> `./presentation/echo/main.go:65`
> ```go
> 	e.GET("/events/:id", server.GetEventsId)
> ```

> `./presentation/echo/handler_event.get.go:11`
> ```go
> // (GET /events/{id})
> func (*Server) GetEventsId(ctx echo.Context) error {
> 	// Get jwt claim
> 	var jwtId *int64
> 	jcc, err := jwt.Check(ctx)
> 	if err == nil {
> 		jwtId = &jcc.Id
> 	}
> 
> 	// Bind id
> 	var id Id
> 	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
> 	if err != nil {
> 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
> 	}
> 
> 	// Bind query
> 	v := ctx.QueryParams()
> 	embed := v["embed"]
> 	query := new(event.GetEventQueryParam)
> 	query.Embed = &embed
> 
> 	// Get event
> 	e, err := event.GetEvent(id, *query, jwtId)
> 	if err != nil {
> 		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
> 	}
> 
> 	return JSONPretty(ctx, http.StatusOK, e)
> }
> ```

- `36`行目あたりが根本処理っぽい

> `./application/event/event.go:50`
> ```go
> func GetEvent(id int64, q GetEventQueryParam, requestUserId *int64) (event.EventEmbed, error) {
> 	u := new(userDomain.User)
> 
> 	if requestUserId != nil {
> 		// リクエスト元のユーザーを取得
> 		var u2 userDomain.User
> 		u2, err := user.Get(*requestUserId)
> 		if err != nil {
> 			return event.EventEmbed{}, err
> 		}
> 		u = &u2
> 	} else if requestUserId == nil {
> 		// リクエストユーザーが指定されていない場合は最小権限のユーザーを仮使用
> 		u = &userDomain.User{
> 			Id:                  0,
> 			PostEventAvailabled: false,
> 			Manage:              false,
> 			Admin:               false,
> 		}
> 	}
> 
> 	return event.GetEvent(id, q, *u)
> }
> ```

- `47`行目で更に関数を呼び出している

> `./domain/model/event/event_get.go:12`
> ```go
> func GetEvent(id int64, q GetEventQueryParam, requestUser user.User) (EventEmbed, error) {
> 	// Get event
> 	// MySQLサーバーに接続
> 	db, err := OpenMysql()
> 	if err != nil {
> 		return EventEmbed{}, err
> 	}
> 	// return時にMySQLサーバーとの接続を閉じる
> 	defer db.Close()
> 
> 	embedUser := false
> 	embedDocuments := false
> 	if q.Embed != nil {
> 		for _, e := range *q.Embed {
> 			if e == "user" {
> 				embedUser = true
> 			}
> 			if e == "documents" {
> 				embedDocuments = true
> 			}
> 		}
> 	}
> 
> 	// `Event`を取得
> 	r1, err := db.Query("SELECT * FROM events WHERE id = ?", id)
> 	if err != nil {
> 		return EventEmbed{}, err
> 	}
> 	defer r1.Close()
> 	if !r1.Next() {
> 		// 1行もレコードが無い場合
> 		// not found
> 		return EventEmbed{}, ErrEventNotFound
> 	}
> 	// 一時変数に割当
> 	var (
> 		id2         int64
> 		name        string
> 		description *string
> 		location    *string
> 		published   bool
> 		completed   bool
> 		userId      int64
> 	)
> 	err = r1.Scan(&id2, &name, &description, &location, &published, &completed, &userId)
> 	if err != nil {
> 		return EventEmbed{}, err
> 	}
> 
> 	// 返り値用変数
> 	event := EventEmbed{
> 		Event: Event{
> 			Id:          id,
> 			Name:        name,
> 			Description: description,
> 			Location:    location,
> 			Datetimes:   []EventDatetime{},
> 			Published:   published,
> 			Completed:   completed,
> 			UserId:      userId,
> 		},
> 	}
> 
> 	// `EventDatetime`を取得
> 	r2, err := db.Query("SELECT * FROM event_datetimes WHERE event_id = ?", id)
> 	if err != nil {
> 		return EventEmbed{}, err
> 	}
> 	defer r2.Close()
> 	for r2.Next() {
> 		var (
> 			eId   string
> 			start *time.Time
> 			end   *time.Time
> 		)
> 		err = r2.Scan(&eId, &start, &end)
> 		if err != nil {
> 			return EventEmbed{}, err
> 		}
> 		// 配列に追加
> 		event.Event.Datetimes = append(event.Event.Datetimes, EventDatetime{*start, *end})
> 	}
> 
> 	if embedUser {
> 		// `User`を取得
> 		u, err := user.Get(event.UserId)
> 		if err != nil {
> 			return EventEmbed{}, err
> 		}
> 		// 変数に追加
> 		event.User = &u
> 	}
> 
> 	if embedDocuments {
> 		// `Documents`を取得
> 		ed, err := GetDocumentList(GetDocumentQueryParam{EventId: &id})
> 		if err != nil {
> 			return EventEmbed{}, err
> 		}
> 		event.Documents = &ed
> 	}
> 
> 	return event, nil
> }
> ```

### クエリ変更

`/events`と同様に`UNION`で対応する。

4 回 → 1 回

```sql
SELECT * FROM (
	WITH params AS ( SELECT <イベントID> as event_id )
		SELECT
			e.id, e.name, e.description, e.location, e.published, e.completed, e.user_id,
			null AS start, null AS end,
			null AS doc_id, null AS doc_name, null AS doc_url,
			u.id AS u_id, u.name AS u_name, u.email AS u_email, u.post_event_availabled, u.manage, u.admin, u.twitter_id, u.github_username
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		WHERE e.id IN (SELECT event_id FROM params)
	UNION ALL
		SELECT
			dt.event_id, null, null, null, null, null, null,
			dt.start, dt.end,
			null, null, null,
			null, null, null, null, null, null, null, null
		FROM event_datetimes dt
		WHERE dt.event_id IN (SELECT event_id FROM params)
	UNION ALL
		SELECT
			doc.event_id, null, null, null, null, null, null,
			null, null,
			doc.id, doc.name, doc.url,
			null, null, null, null, null, null, null, null
		FROM documents doc
		WHERE doc.event_id IN (SELECT event_id FROM params)
	) AS e
ORDER BY id, name IS NULL ASC
```

### `go`のコードを書き換える

省略（下の url から変更例を見れます。）

### サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: update query to get event'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: update query to get event · tingtt/e-isucon@3d3fa15 (github.com)](https://github.com/tingtt/e-isucon/commit/3d3fa15982fb758d0a4398f57fbd9d5bad4d0ed5)

## `/users/sign_in`

### コードを見る

> `/domain/model/users/get.go:46`
> ```go
> func GetByEmail(email string) (User, error) {
> 	// MySQLサーバーに接続
> 	db, err := OpenMysql()
> 	if err != nil {
> 		return User{}, err
> 	}
> 	// return時にMySQLサーバーとの接続を閉じる
> 	defer db.Close()
> 
> 	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
> 	r, err := db.Query(
> 		`SELECT
> 			id, name, email, password,
> 			post_event_availabled, manage, admin,
> 			twitter_id, github_username, star_count
> 		FROM users
> 		WHERE email = ?`,
> 		email,
> 	)
> 	if err != nil {
> 		return User{}, err
> 	}
> 	defer r.Close()
> 	if !r.Next() {
> 		// 1行もレコードが無い場合
> 		// not found
> 		return User{}, ErrUserNotFound
> 	}
> 
> 	// 変数に割り当て
> 	u := User{}
> 	err = r.Scan(
> 		&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
> 		&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
> 	)
> 	if err != nil {
> 		return User{}, err
> 	}
> 
> 	return u, nil
> }
> ```

- `GetByEmail`関数をしているコードを全て見ても、使用されているフィールドが `id`, `email`, `password`, `admin` だけなので無駄なフィールドが読み込まれている。

### コードを変更

```diff
  func GetByEmail(email string) (User, error) {
  	// MySQLサーバーに接続
  	db, err := OpenMysql()
  	if err != nil {
  		return User{}, err
  	}
  	// return時にMySQLサーバーとの接続を閉じる
  	defer db.Close()
  
  	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
  	r, err := db.Query(
- 		`SELECT
- 			id, name, email, password,
- 			post_event_availabled, manage, admin,
- 			twitter_id, github_username, star_count
- 		FROM users
- 		WHERE email = ?`,
+ 		`SELECT id, password, admin FROM users WHERE email = ?`,
  		email,
  	)
  	if err != nil {
  		return User{}, err
  	}
  	defer r.Close()
  	if !r.Next() {
  		// 1行もレコードが無い場合
  		// not found
  		return User{}, ErrUserNotFound
  	}
  
  	// 変数に割り当て
  	u := User{}
- 	err = r.Scan(
- 		&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
- 		&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
- 	)
+ 	err = r.Scan(&u.Id, &u.Password, &u.Admin)
  	if err != nil {
  		return User{}, err
  	}
  
+ 	u.Email = email
  	return u, nil
  }
```

### サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: update query to get user by email'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: update query to get user by email · tingtt/e-isucon@6e0081f (github.com)](https://github.com/tingtt/e-isucon/commit/6e0081f37c7fe6b89b2ed29c23169a6263bed84c)

## `LastInsertId` をやめて`UUID`を使用する。

### テーブルの変更

`./.mysql/init.sql`
```diff
  --
  -- Table structure for table `users`
  --
  CREATE TABLE `users` (
-   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
+   `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL UNIQUE,
    `password` varchar(255) NOT NULL,
    `post_event_availabled` tinyint(1) NOT NULL DEFAULT 0,
    `manage` tinyint(1) NOT NULL DEFAULT 0,
    `admin` tinyint(1) NOT NULL DEFAULT 0,
    `twitter_id` varchar(255),
    `github_username` varchar(255),
    `star_count` int(255) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
  );
  
  --
  -- Table structure for table `events`
  --
  CREATE TABLE `events` (
-   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
+   `id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` varchar(255),
    `location` varchar(255),
    `published` tinyint(1) NOT NULL,
    `completed` tinyint(1) NOT NULL,
-   `user_id` int(255) UNSIGNED NOT NULL,
+   `user_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
  );
  
  --
  -- Table structure for table `event_datetimes`
  --
  CREATE TABLE `event_datetimes` (
-   `event_id` int(255) UNSIGNED NOT NULL,
+   `event_id` varchar(255) NOT NULL,
    `start` datetime NOT NULL,
    `end` datetime NOT NULL,
    FOREIGN KEY (`event_id`) REFERENCES events(`id`) ON DELETE CASCADE
  );
  
  --
  -- Table structure for table `documents`
  --
  CREATE TABLE `documents` (
-   `id` int(255) UNSIGNED AUTO_INCREMENT NOT NULL,
+   `id` varchar(255) NOT NULL,
-   `event_id` int(255) UNSIGNED NOT NULL,
+   `event_id` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `url` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`event_id`) REFERENCES events(`id`) ON DELETE CASCADE
  );
```

### データの変更

1. 仮テーブルとして`id`, `user_id`, `event_id`を`varchar`型に変更したテーブルを作成
	```sql
	CREATE TABLE `users2` (
	  `id` varchar(255) NOT NULL,
	  `name` varchar(255) NOT NULL,
	  `email` varchar(255) NOT NULL UNIQUE,
	  `password` varchar(255) NOT NULL,
	  `post_event_availabled` tinyint(1) NOT NULL DEFAULT 0,
	  `manage` tinyint(1) NOT NULL DEFAULT 0,
	  `admin` tinyint(1) NOT NULL DEFAULT 0,
	  `twitter_id` varchar(255),
	  `github_username` varchar(255),
	  `star_count` int(255) UNSIGNED NOT NULL DEFAULT 0,
	  PRIMARY KEY (`id`)
	);
	
	CREATE TABLE `events2` (
	  `id` varchar(255) NOT NULL,
	  `name` varchar(255) NOT NULL,
	  `description` varchar(255),
	  `location` varchar(255),
	  `published` tinyint(1) NOT NULL,
	  `completed` tinyint(1) NOT NULL,
	  `user_id` varchar(255) NOT NULL,
	  PRIMARY KEY (`id`)
	);
	
	CREATE TABLE `event_datetimes2` (
	  `event_id` varchar(255) NOT NULL,
	  `start` datetime NOT NULL,
	  `end` datetime NOT NULL,
	  FOREIGN KEY (`event_id`) REFERENCES events2(`id`) ON DELETE CASCADE
	);
	
	CREATE TABLE `documents2` (
	  `id` varchar(255) NOT NULL,
	  `event_id` varchar(255) NOT NULL,
	  `name` varchar(255) NOT NULL,
	  `url` varchar(255) NOT NULL,
	  PRIMARY KEY (`id`),
	  FOREIGN KEY (`event_id`) REFERENCES events2(`id`) ON DELETE CASCADE
	);
	```
1. 元データのあるテーブルから変換して`INSERT`

	`int` → `varchar` の型変換は暗黙変換されるため `SELECT`  したものを `INSERT` するだけで良い。
	```sql
	INSERT INTO users2 SELECT * FROM users;
	INSERT INTO events2 SELECT * FROM events;
	INSERT INTO event_datetimes2 SELECT * FROM event_datetimes;
	INSERT INTO documents2 SELECT * FROM documents;
	```
1. `mysqldump` で `INSERT` 文を生成
	```bash
	# mysqldump で usersテーブルのデータからSQL文を生成し、users.sqlに書き込む
	docker-compose exec mysql mysqldump -psecret prc_hub users2 > users.sql
	docker-compose exec mysql mysqldump -psecret prc_hub events2 > events.sql
	docker-compose exec mysql mysqldump -psecret prc_hub event_datetimes2 > event_datetimes.sql
	docker-compose exec mysql mysqldump -psecret prc_hub documents2 > documents.sql
	```
1. 生成した `INSERT` 文で `backend/domain/model/eisucon/migrate.sql` を書き換える。

### 型変更

`domain/model/user/types.go`, `domain/model/event/types.go`, `domain/model/jwt/main.go` の `Id` 関連の定義をすべて `int64` から `string` 型に変更し、**その後エラーになる箇所を潰していく。**

- `domain/model/user/types.go`

	```diff
	  type User struct {
	- 	Id                  int64   `json:"id"`
	+ 	Id                  string  `json:"id"`
	  	Name                string  `json:"name"`
	  	Email               string  `json:"-"`
	  	Password            string  `json:"-"`
	  	PostEventAvailabled bool    `json:"post_event_availabled"`
	  	Manage              bool    `json:"manage"`
	  	Admin               bool    `json:"admin"`
	  	TwitterId           *string `json:"twitter_id,omitempty"`
	  	GithubUsername      *string `json:"github_username,omitempty"`
	  	StarCount           uint64  `json:"star_count"`
	  }
	```

- `domain/model/event/types.go`

	```diff
	  type Event struct {  
	- 	Id          int64           `json:"id"`
	+ 	Id          string          `json:"id"`
	  	Name        string          `json:"name"`
	  	Description *string         `json:"description,omitempty"`
	  	Location    *string         `json:"location,omitempty"`
	  	Datetimes   []EventDatetime `json:"datetimes"`
	  	Published   bool            `json:"published"`
	  	Completed   bool            `json:"completed"`
	- 	UserId      string          `json:"user_id"`
	+ 	UserId      string          `json:"user_id"`
	  }
	  
	  type EventDatetime struct {
	  	Start time.Time `json:"start"`
	  	End   time.Time `json:"end" dh:"end"`
	  }
	  
	  type EventDocument struct {
	- 	EventId int64  `json:"event_id"`
	+ 	EventId string `json:"event_id"`
	- 	Id      int64  `json:"id"`
	+ 	Id      string `json:"id"`
	  	Name    string `json:"name"`
	  	Url     string `json:"url"`
	  }
	```

- `domain/model/jwt/main.go`

	```diff
	  type jwtCustumClaims struct {
	- 	Id    int64  `json:"id"`
	+ 	Id    string  `json:"id"`
	  	Email string `json:"email"`
	  	Admin bool   `json:"admin"`
	  	jwt.StandardClaims
	  }
	  
	  type GenerateTokenParam struct {
	- 	Id    int64
	+ 	Id    string
	  	Email string
	  	Admin bool
	  }
	```

### `UUID`による`ID`生成

`LastInsertId` の代わりに `UUID` を使用する。

- `domain/model/user/create.go`
	```diff
		// `users`テーブルに追加
	+ 	id := uuid.New().String()
		r, err := d.Exec(
	- 		`INSERT INTO users (name, email, password, post_event_availabled, manage, admin, twitter_id, github_username) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
	+ 		`INSERT INTO users (id, name, email, password, post_event_availabled, manage, admin, twitter_id, github_username) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	- 		p.Name, p.Email, string(hashed), false, false, false, p.TwitterId, p.GithubUsername,
	+ 		id, p.Name, p.Email, string(hashed), false, false, false, p.TwitterId, p.GithubUsername,
		)
		if err != nil {
			return UserWithToken{}, err
		}
	- 	id, err := r.LastInsertId()
	- 	if err != nil {
	- 		return UserWithToken{}, err
	- 	}
	```

- `domain/model/event/event_create.go`
	```diff
	  	// `events`テーブルに追加
	+ 	id := uuid.New().String()
	  	r, err := tx.Exec(
	- 		`INSERT INTO events (name, description, location, published, completed, user_id) VALUES (?, ?, ?, ?, ?, ?)`,
	+ 		`INSERT INTO events (id, name, description, location, published, completed, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
	- 		p.Name, p.Description, p.Location, p.Published, p.Completed, requestUser.Id,
	+ 		id, p.Name, p.Description, p.Location, p.Published, p.Completed, requestUser.Id,
	  	)
	  	if err != nil {
	  		return Event{}, err
	  	}
	- 	id, err := r.LastInsertId()
	- 	if err != nil {
	- 		return Event{}, err
	- 	}
	  
	  	// `event_datetimes`テーブルに追加
	  	for _, dt := range datetimes {
	  		_, err = tx.Exec(
	  			"INSERT INTO event_datetimes (event_id, start, end) VALUES (?, ?, ?)",
	  			id, dt.Start, dt.End,
	  		)
	  		if err != nil {
	  			return Event{}, err
	  		}
	  	}
	```

- `domain/model/event/document_create.go`
	```diff
	  	// `documents`テーブルに追加
	+ 	id := uuid.New().String()
	  	r, err := db.Exec(
	- 		`INSERT INTO documents (event_id, name, url) VALUES (?, ?, ?)`,
	+ 		`INSERT INTO documents (id, event_id, name, url) VALUES (?, ?, ?, ?)`,
	- 		p.EventId, p.Name, p.Url,
	+ 		id, p.EventId, p.Name, p.Url,
	  	)
	  	if err != nil {
	  		return EventDocument{}, err
	  	}
	- 	id, err := r.LastInsertId()
	- 	if err != nil {
	- 		return EventDocument{}, err
	- 	}
	```

### サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'perf: change LastInsertId to UUID'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
git pull
make upgrade
```

[perf: change LastInsertId to UUID · tingtt/e-isucon@0733528 (github.com)](https://github.com/tingtt/e-isucon/commit/0733528423c766608ad280a22431879e12005c3b)

# Docker をやめる

docker をやめて仮想化レイヤーを無くすことで高速化を図る。

## MariaDB のインストール

※ `mysql`の代わり
> MariaDB starting with [10.2.1](https://mariadb.com/kb/en/mariadb-1021-release-notes/)
> [WITH - MariaDB Knowledge Base](https://mariadb.com/kb/en/with/)

`With`句は 10.2.1以上のバージョンでしか使用できないため最新版をインストールする。

```bash
# [root@server ~]
cat <<EOF > /etc/yum.repos.d/MariaDB.repo
[mariadb]
name = MariaDB
baseurl = http://yum.mariadb.org/10.5/centos7-amd64
gpgkey=https://yum.mariadb.org/RPM-GPG-KEY-MariaDB
gpgcheck=1
EOF

yum update -y
yum install -y \
	mariadb-server \
	mariadb-client \
	MariaDB-shared \
	MariaDB-devel
yum install -y mariadb-server
systemctl enable --now mariadb
mysql_secure_installation # すべて [Y]、password は `secret`
```

## go のインストール

```bash
# [root@server /usr/local/src]
wget https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile
source ~/.bash_profile

go version
```

## DBへの接続を tcp -> socket に

```diff
  func InitRepository(user string, password string, host string, port uint, db string) {
- 	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, db)
+ 	dsn = fmt.Sprintf("%s:%s@unix(/var/lib/mysql/mysql.sock)/%s?parseTime=true", user, password, db)
  }
```

## インストールスクリプト

`Makefile`

```makefile
.PHONY: build
build:
	go build

.PHONY: install
SEED_FELE	?= /usr/local/lib/eisucon-backend/seed.sql
MODULE	?= prc_hub_back
BIN	?= /usr/local/bin/eisucon-backend
ARG_LOG_LEVEL	?= debug
ARG_MYSQL_HOST	?= localhost
ARG_MYSQL_DB	?= prc_hub
ARG_MYSQL_USER	?= root
ARG_MYSQL_PASSWORD	?= secret
ARGS	?= --log.level=$(ARG_LOG_LEVEL) --jwt.issuer=prc_hub --jwt.secret=e0VzhtkQ --mysql.host=$(ARG_MYSQL_HOST) --mysql.db=$(ARG_MYSQL_DB) --mysql.user=$(ARG_MYSQL_USER) --mysql.password=$(ARG_MYSQL_PASSWORD) --migrate-sql-file=$(SEED_FELE)
define UNITFILE
[Unit]
Description=ECC-ISUCON backend
After=network.target

[Service]
Restart=on-failure
RestartSec=10
ExecStart=$(BIN) $(ARGS)

[Install]
WantedBy=multi-user.target
endef
export UNITFILE
UNITFILE_PATH	?= /usr/local/lib/systemd/system/eisucon-backend.service
install: build
	mkdir -p /usr/local/lib/eisucon-backend/
	mysql -u$(ARG_MYSQL_USER) -p$(ARG_MYSQL_PASSWORD) < ./.mysql/init.sql
	@echo "$$UNITFILE" > ./.tmp.service
	cp -f ./.tmp.service $(UNITFILE_PATH)
	@rm -f ./.tmp.service
	cp -n ./domain/model/eisucon/migrate.sql $(SEED_FELE)
	cp -n ./$(MODULE) $(BIN)
	systemctl daemon-reload
	systemctl enable --now eisucon-backend.service


.PHONY: purge
SEED_FELE	?= /usr/local/lib/eisucon-backend/seed.sql
BIN	?= /usr/local/bin/eisucon-backend
UNITFILE_PATH	?= /usr/local/lib/systemd/system/eisucon-backend.service
ARG_MYSQL_USER	?= root
ARG_MYSQL_PASSWORD	?= secret
ARG_MYSQL_DB	?= prc_hub
define SQL
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).documents;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).event_datetimes;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).events;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).users;
DROP DATABASE IF EXISTS $(ARG_MYSQL_DB);
endef
export SQL
purge:
	systemctl disable --now eisucon-backend.service
	@echo "$$SQL" > ./.tmp.sql
	mysql -u$(ARG_MYSQL_USER) -p$(ARG_MYSQL_PASSWORD) < ./.tmp.sql
	@rm -f ./.tmp.sql
	rm -f $(SEED_FELE) $(BIN) $(UNITFILE_PATH)
	-rmdir /usr/local/lib/eisucon-backend/
```

## サーバーに適応

GitHub に push

```bash
# [localmachine ~/Downloads/eisucon]
git add .
git commit -m 'feat: add install script'
git push origin head
```

GitHub から pull

```bash
# [root@server /usr/local/src/eisucon]
cd /usr/local/src/eisucon
make purge
git pull
make install
```
