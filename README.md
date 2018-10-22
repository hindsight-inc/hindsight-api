# Project hindsight-api

RESTful server for the Hindsight project.

# Setup

## Go

- `brew install go`

## Workspace

[gostart](https://github.com/alco/gostart#faq0)

- Add `export GOPATH=$HOME/go` in your `.zshrc`
- `cd $GOPATH && mkdir src && cd src`
- `git clone git@github.com:hindsight-inc/hindsight-api.git hindsight && cd hindsight-api`
- `go get ./...`: this puts dependencies in `$GOPATH/src`
- `./run.sh`

## MySQL

- Recommended client: `TablePlus` (`Sequel Pro` doesn't support newest encryption protocols)

### Local test environment: macOS

- `brew install mysql`
- Starting service
  - Background: `brew services start mysql`
  - Current session: `mysql.server start`
- `mysql -uroot`
- `CREATE USER 'golang'@'localhost' IDENTIFIED BY 'password';`
- `GRANT ALL PRIVILEGES ON * . * TO 'golang'@'localhost';`
- `create database golang;`

## Tools

- `run.sh`: build, test, and run
- `edit.sh`: edit all source files like a pro
