<h1 align="center"><b>GO project with mysql and nginx reverse proxy on Docker.</b></h1>

## ent
Laravelみたいにmodelとrepositoryが１つになる感じに近い？  
infrastructure/repository部分がentに依存する形になってしまう。  
domain/modelはentに依存しない作りにするために定義すべきだが、下記問題が発生する。
1. entとdomain/modelでの２重管理になる
2. domain/modelで表したものがrepositoryに表現されない可能性がある

つまりなんちゃってonion architectureしか作れない。  
ただし、model部分を自動生成してくれるので実装はかなり早くなる。

#### generate schema files
`go run -mod=mod entgo.io/ent/cmd/ent new User Car`
#### generate entity, migrate & client...etc
`go generate ./ent`
#### migrate
開発環境ではAutomatic migrations（ent）を行う。  
本番環境では開発環境で自動作成され、/ent/migrate/migrationsに格納しておいたmigrationファイルを使って、Version migrations（atlas）を行う。  
- development (ent migration)  
`go run main.go <migration file name>`
- production or staging (atlas migration)
    ```shell
    atlas migrate apply \
    --dir "file://ent/migrate/migrations" \
    --url mysql://sample_user:sample_password@localhost:3306/sample_db
    ```

<br>

## ozzo-validation
~~本来であればdomain/modelにvalidationを書くべきだが、  
entのgenerateで自動生成されるmodelにはvalidationが書くことはできない。  
そのため、clean architectureライクにusecase/inputを作成し、そこにvalidationを実装。  
（entにもvalidationは書けるがあくまでもDB用のvalidation処理）~~
domain/modelに作成

<br>

## aws-sdk-go
infrastructure層でS3のバケットに対する各種操作を行う。
アップロードでは大容量ファイルを想定したgoroutineを用いた並列処理を行う。

<br>

## Docker
#### start
`docker-compose up -d`
#### start go project
```
docker-compose exec -it go sh
go run main.go
```
#### install GO module
```shell
docker-compose exec -it go sh
go get xxxxxxxxxx
```
#### get out of go command
`exit`
#### stop
`docker-compose stop`
#### delete
`docker-compose down`
#### check docker status
`docker-compose ps`
