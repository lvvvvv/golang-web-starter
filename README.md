＃ Golang Web starter
the minimal starter to build web/webapi projects

## What included
 - generic folder structure
 - toml config file, could override by argument flags
 - general iris middlewares setup
 - support i18n
 - included db and redis cache setup
 - jwt for api auth

## How to start
  - clone the project and rename as yours
  - modify main.go import path by yours package name

  - create database (in example, the db name is 'starter')
  - edit dbconf.yml with your db setting
  - download & install goose (bitbucket.org/liamstask/goose/cmd/goose)
  - exec mirgation (cmd: goose up)
  - download & install glide (https://github.com/Masterminds/glide)
  - exec glide install
  - exec go run main.go
