## mattermost-plugin-parser

mattermost-plugin-parser parses project of Mattermost Plugin, and make reports about usages of Mattermsot Plugin API.

## SetUp

mattermost-plugin-parser send data to MySQL server.

### Docker

This repository has `docker-compose.yml`. If you want to run the analysis locally, you should run `docker-compose up -d`.
`docker-compose` command start mysql service, and setup database by [pre-defined sql](./initdb/).

## Parse Mattemrost plugin repositories

### 1. Environement variables

At first, you must set evironment variables for connection MySQL.

If you set up mysql by `docker-compose`...
```
export MYSQL_HOST=localhost \
&& export MYSQL_PORT=13306 \
&& export MYSQL_USER=mmuser \
&& export MYSQL_PASSWORD=mostest \
&& export MYSQL_DATABASE=mmplugin_parser
```

## 2. Setup

* Resolve dependencies for `server`

```
$ cd server
$ go mod tidy
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

## 4. Report

Run `report.js`

```
cd webapp/
node report.js
```

`report.js` gets data from MySQL, and generate JSON file to `docs/data.json`.
You can see, filter, sort this data by [Tabulator](http://tabulator.info/) UI. If you want to access this, run the following command and access http://localhost:8000.

```
cd docs
python -m http.server        # Python v3
python -m SimpleHTTPServer   # Python v2
```