const request = require('request');
const esprima = require('esprima');

const URL = 'https://raw.githubusercontent.com/mattermost/mattermost-webapp/master/plugins/registry.js';

function doRequest(options) {
    return new Promise(function(resolve, reject) {
        request(options, function(err, resp, body) {
            if (!err && resp.statusCode == 200) {
                resolve(body);
            } else {
                reject(err);
            }
        })
    })
}

const parsePluginResistryClassMethod = async () => {
    const body = await doRequest({
        url: URL,
        method: "GET",
    })
    const parsed = esprima.parseModule(body, { comment: false, attachComment: false });
    return parsed.body.find(statement =>
        statement.type === "ExportDefaultDeclaration" &&
        statement.declaration.id.name === "PluginRegistry"
    ).declaration.body.body.filter(statement =>
        statement.type === "MethodDefinition" &&
        statement.key.name !== "constructor"
    ).map(statement => ({
        name: statement.key.name,
    }))
}

module.exports = {
    parsePluginResistryClassMethod: parsePluginResistryClassMethod,
}