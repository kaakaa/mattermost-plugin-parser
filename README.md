## SetUp

Run MySQL service.

### Docker

This repository has `docker-compose.yml`. If you want to run the analysis locally, you run `docker-compose up -d`.
`docker-compose` command start mysql and metabase services, and setup database by [pre-defined sql](./initdb/).

### Heroku

I use [Metabase on Heroku](https://metabase.com/start/heroku.html). But Metabase on Heroku uses `PostgreSQL` for DB, but this repository uses `MySQL` for DB.
So after setting up Metabase on Heroku, I add [ClearDB add-on](https://metabase.com/start/heroku.html) for that heroku app. And I set Metabase service to refer ClearDB service.

## Parse Mattemrost plugin repositories

### 1. Environement variables

At first, you must set evironment variables for connection MySQL.

If you set up mysql by `docker-compose`...
```
$ export MYSQL_HOST=localhost
$ export MYSQL_PORT=13306
$ export MYSQL_USER=mmuser
$ export MYSQL_PASSWORD=mostest
$ export MYSQL_DATABASE=mmplugin_parser
```

If you set up mysql on Heroku, you can get information for connecting by [`heroku-cli`](https://devcenter.heroku.com/articles/heroku-cli).

```
$ heroku config | grep CLEARDB_DATABASE_URL
CLEARDB_DATABASE_URL => mysql://adffdadf2341:adf4234@us-cdbr-east.cleardb.com/heroku_db?reconnect=true

$ export MYSQL_HOST=us-cdbr-east.cleardb.com
$ export MYSQL_PORT=3306
$ export MYSQL_USER=adffdadf2341
$ export MYSQL_PASSWORD=adf4234
$ export MYSQL_DATABASE=heroku_db
```
* https://devcenter.heroku.com/articles/cleardb

## 2. Setup

* Resolve dependencies for `server`

```
$ cd server
$ dep ensure
```

* Resolve dependencies for `webapp`

```
$ cd webapp
$ npm install
```

## 3. Run

Run script.

```
$ ./run.sh
```

If you add repositories of Mattermost plugin development project, edit [repositories.txt](./repositories.txt)
