![gh-action](https://github.com/kaakaa/mattermost-plugin-parser/workflows/Run/badge.svg)

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

### 2. Setup

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

### 3. Run

Run script.

```
$ ./run.sh
```

If you add repositories of Mattermost plugin development project, edit [repositories.txt](./repositories.txt)

### 4. Report

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

## Tips

### Remove invalid data from MySQL database

Since tables in `mmplugin_parser` database has some constrains, if you want to remove a row, it's needed to remove data with correct order.

```shell
export ID=${COMMIT_ID_YOU_WANNA_DELETE}
mysql -ummuser -pmostest mmplugin_parser -e "DELETE FROM manifest WHERE commit_id = '$ID'; DELETE FROM settings_schema WHERE commit_id = '$ID'; DELETE FROM plugin_settings WHERE  commit_id = '$ID'; DELETE FROM props WHERE commit_id = '$ID'; DELETE FROM usages WHERE commit_id = '$ID'; DELETE FROM repositories WHERE commit_id = '$ID';"
```