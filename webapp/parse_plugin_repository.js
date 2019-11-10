const fsp = require('fs-promise');
const glob = require('glob');

const {logger} = require('./logger');

const {parsePluginSource} = require('./parser');

/**
 * Find plugin js/jsx files recursively in `path`, and parse it.
 * Files in node_modules will not be parsed.
 * 
 * @param {String} path root directory 
 */
const findRecursive = async (path) => {
    let ret = [];
    ret = ret.concat(await parseFiles(glob.sync(path + "/**/*.jsx", {ignore: path + "/node_modules/**"}), true));
    ret = ret.concat(await parseFiles(glob.sync(path + "/**/*.js", {ignore: path + "/node_modules/**"}), false));
    return ret
}

const parseFiles = async (files, isJsx) => {
    let ret = [];
    for (idx in files) {
        try {
            const file = files[idx];
            logger.info("Parse %s", file);
            const usages = await parseFile(file, isJsx);
            if (usages.length > 1) {
                logger.info("  %d functions are detected (%s)", usages.length, usages);
            }
            ret = ret.concat(usages)
        } catch (err) {
            logger.error('Failed to parse: %s', err.message)
        }
    }
    return ret;
}

/**
 * Parse plugin source code
 * @param {String} file file path for parsing
 * @param {boolean} isJsx true if file is .jsx
 */
const parseFile = async (file, isJsx) => {
    const contents = fsp.readFileSync(file, 'utf-8');
    return parsePluginSource(contents, file)
}

module.exports = {
    findRecursive: findRecursive,
}

