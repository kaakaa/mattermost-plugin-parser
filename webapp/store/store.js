var mysql = require('mysql2');

let connection;

const init = async () => {
    if (!connection) {
        const pool = mysql.createPool({
            host     : process.env.MYSQL_HOST, // 'localhost',
            port     : process.env.MYSQL_PORT || 3306, //13306,
            user     : process.env.MYSQL_USER,
            password : process.env.MYSQL_PASSWORD,
            database : process.env.MYSQL_DATABASE, // 'mmplugin_parser'
        });
        connection = pool.promise();
    }
}

const listRepositories = async () => {
    const [rows, fields] = await connection.query('SELECT * FROM repositories');
    return rows
}

const selectRepository = async (data) => {
    const [rows, fields] = await connection.query('SELECT * FROM repositories WHERE commit_id = ?', data)
    return rows
}
 
const insertRepository = async (data) => {
    const [rows, fields] = await connection.query('INSERT IGNORE INTO repositories SET url = ?, commit_id = ?, created_at = ?', data)
    return rows
}

const listUsages = async () => {
    const [rows, fields] = await connection.query('SELECT * FROM usages');
    return rows;
}

const insertUsage = async (data) => {
    const [rows, fields] = await connection.query('INSERT IGNORE INTO usages SET commit_id = ?, api = ?, path = ?, line = ?, type = "webapp.registry"', data)
    return rows
}

const end = async () => {
    if (connection) {
        await connection.end();
    }
}

module.exports = {
    init: init,
    listRepositories: listRepositories,
    selectRepository: selectRepository,
    insertRepository: insertRepository,
    listUsages: listUsages,
    insertUsage: insertUsage,
    end: end,
}