const request = require('request');

const {parseMattermostPluginRegistryClass} = require('./parser')

/**
 * URL to source code of Mattermost PluginRegistry class
 */
const pluginRegistryClassSourceURL = 'https://raw.githubusercontent.com/mattermost/mattermost-webapp/master/plugins/registry.js';

/**
 * Do http request
 * @param {*} options http request
 */
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

/**
 * Collect method names except for constructor in PluginRegistry class 
 */
const parsePluginResistryClassMethod = async () => {
    const body = await doRequest({
        url: pluginRegistryClassSourceURL,
        method: "GET",
    })

    return parseMattermostPluginRegistryClass(body);
};

module.exports = {
    parsePluginResistryClassMethod: parsePluginResistryClassMethod,
}