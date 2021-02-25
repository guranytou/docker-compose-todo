# docker-composeチュートリアルメモ

## はじめに
本チュートリアルでは、ローカルPC上で、docker-composeを利用してToDoアプリを構築し、docker-composeの使い方やdocker-compose.ymlの書き方を学習します。

## この記事でできるようになること
　複数のコンテナを一括で構築することができるようになる

## 想定読者
 - docker-composeを使って業務はしているが、docker-composeが何か分からない人
 - docker-composeの事はなんとなく分かるが、docker-compose.yml書き方が分からない人

## ToDoアプリの構成技術
以下の言語/FWで構成されています。なお、これらのソースコード/sqlファイルについては今回触れません。 
 - frontend : Vue.js
 - backend: Go
 - database: Mariadb

## docker-composeとは
docker-composeとは、複数のコンテナを定義し実行する Docker アプリケーションのためのツールです。 
YAML ファイルを使ってアプリケーションの設定を行い、コマンドを１つ実行するだけで、設定内容に基づいたアプリケーションの生成/起動を行うことが出来ます。

## 今回使う構成について
今回の構成は以下の通りです。　
![docker-compose-tutorial](https://user-images.githubusercontent.com/42028429/109179903-b692b100-77cd-11eb-91c8-aee2e4414dc5.png)

## 起動方法
上記リポジトリをローカルへclone後、`docker-compose up`すると立ち上がります。
その後、ブラウザで `http://localhost`へアクセスするとToDoアプリにアクセスします。

## よく利用するdocker-composeコマンド
 `docker-compose up`  
docker-compose.ymlに記載されている内容を元にコンテナを起動します。  
`-d`オプションでバックグラウンド処理ができます。

 `docker-compose down`  
起動したコンテナやコンテナネットワークを停止した後に削除します。

　`docker-compose stop`  
起動したコンテナを停止する（削除はしたくない）場合はこのコマンドを利用します。

　`docker-compose run`  
`docker-compose stop`で停止したコンテナを起動する場合はこのコマンドを利用します。

## 実際に触ってみよう
### frontのポート番号を変更してみる
docker-compose.ymlのポート番号を変更してみましょう
```yml
front:
 image: nginx:1.19.7
 container_name: tutorial-front
 hostname: tutorial-front
 ports:
  - 18080:80 // 80から18080に変更
```
変更後に `docker-compose up`を実行し、 `http://localhost:18080`でToDoアプリにアクセスできれば成功です。

### 後片付け
`docker-compose down`を行うと、先ほど `docker-compose up`で作成したコンテナ・ネットワークが自動的に削除されます。


## docker-compose.ymlの解説
- `version: `  
利用するdocker-composeのバージョンを指定します。  
特段の理由がない場合は、最新版（2021/2月時点での最新はv3.8）を利用しましょう。  
※参考：https://matsuand.github.io/docs.docker.jp.onthefly/compose/compose-file/

- `services:`  
services配下に構築したいサービス名を記載します。その下にそれぞれのサービス毎の設定を記載していきます。  
今回の場合は、フロントエンド・バックエンド・データベースの3サービスを構築したいので、それぞれfront/api/dbと記載します。  
```yml
services:
 front:
 	（以下略）
 api:
 	（以下略）
 db:
 	（以下略）
```

- `image:`  
利用したいDockerイメージ名を書きます。  
今回フロントではnginxを、データベースではMariaDBを利用したいので、以下のように記載しています。  
※後述するbuild:を利用したapiでもimageを利用していますが、それはbuild:の項で説明します。
```yml
front:
 	image: nginx:1.19.7
db:
  image: mariadb:10.5.8
```

- `build:`
Dockerfileを利用してサービスを構築したい場合に利用します。  
今回バックエンドではGoのビルドやホットリロードを行うために、Dockerfileを使って構築しています。  
```yml
api:
 build:
   context: .  // Docker buildを実行する際のディレクトリ。ここではルートディレクトリを設定しています
   dockerfile: ./backend/Dockerfile // どのファイルをDockerfileとするかを設定します
```
また、ここで作成するイメージ名を設定する時にはimageを利用します。
```yml
api:
 image: docker-compose-tutorial-api // この名前でイメージが作成されます
```

- `container_name:`
コンテナの名前を設定できます。 

- `hostname:`
Dockerネットワーク内で利用するホスト名を設定できます。  

- `ports:`
コンテナが公開するポートを設定することが出来ます。  
ホスト側で利用するポート：コンテナ側で開いているポート、と記載します。  
```yml
front:
 ports:
  - 80:80 
api:
 ports:
  - 8080:8080 
db:
 ports:
  - 3306:3306 
```

- `networks:`
docker-compose.ymlの下部で独自のDockerネットワークを作成し、そのネットワークにコンテナを参加させることができます。
```yml
 // 各コンテナを作成した'todo_app'という名前のDockerネットワークに参加させている
service:
 front:
  networks:
   - todo_app
 api:
  networks:
   - todo_app
 db:
  networks:
   - todo_app
            
 // ここで'todo_app'というネットワークを作成している
networks:
 todo_app:
  name: todo_app
```

- `environment:`
環境変数の設定を行うことができます。  
ここではdbコンテナのルートパスワード設定およびToDoアプリで利用するデータベース作成を行っています。  
```yml
db:
 environment:
  MYSQL_ROOT_PASSWORD: admin //ルートパスワードを'admin'に設定
  MYSQL_DATABASE: todo // todoという名前のデータベースを作成
```

- `volumes:`
ホストにあるディレクトリやファイルを、コンテナ内のディレクトリにマウントすることができます。
``` yml
// nginxの公開用ファイルを置くディレクトリに、index.htmlとindex.jsがあるfrontendフォルダをマウントしている
front:
 volumes:
  - ./frontend:/usr/share/nginx/html
 
// Goのソースを置くディレクトリにファイルをマウントしている
api:
 volumes:
  - ./backend/app:/go/src/github.com/guranytou/docker-compose-todo
  
//　docker-entrypoint-initdb.dディレクトリにSQLファイルを格納すると、コンテナ起動時に自動実行するので、テーブル作成とテストデータ投入用のSQLファイルをマウントしている
db:
 - ./db/:/docker-entrypoint-initdb.d/
```

- `command:`
コンテナ起動時に実行したいコマンドを記載します。  
今回はDBコンテナに文字コードを設定するために利用しています。
```yml
db:
 command: mysqld --character-set-server=utf8 --collation-server=utf8_unicode_ci
```

## 参考資料
　docker-compose.ymlのレファレンス
http://docs.docker.jp/compose/compose-file.html
