const parseMattermost = require('./parse_mattermost');
const parsePluginRepository = require('./parse_plugin_repository');

const Store = require('./store/store');
const {logger} = require('./logger');

/**
 * Find usages of plugin method
 * 
 * @param {String} path path of root directory of plugin
 */
const parse = async (path) => {
    // Parse Matermost PluginRegistry class
    logger.info('Fetch and parse Mattermost PluginRegistry class')
    const pluginRegistryClassMethod = await parseMattermost.parsePluginResistryClassMethod();
    logger.info('Complete: %d functions (%s)',
        pluginRegistryClassMethod.length,
        JSON.stringify(pluginRegistryClassMethod.map(m => m.name))
    );

    // Parse plugin source code
    logger.info('Parse plugin source code')
    const calledFuncs = await parsePluginRepository.findRecursive(path);
    logger.info("Commplete: %d functions (%s)", calledFuncs.length, JSON.stringify(calledFuncs));

    return calledFuncs
        .filter(f => pluginRegistryClassMethod.some(m => m.name == f.name))
}

const save = async (f, basePath) => {
    logger.info("Save usages of webapp api: %s: %d", f.file, f.loc.start.line);

    await Store.init()
    let rows = await Store.selectRepository([f.commit_id])
    if (rows.length === 0) {
        rows = await Store.insertRepository([f.repository, f.commit_id, f.commited_at])
    }

    const regexp = new RegExp(basePath + '/', 'g');
    const relPath = f.file.replace(regexp, '');
    return await Store.insertUsage([f.commit_id, f.name, relPath, f.loc.start.line])
}

if (process.argv.length !== 6) {
    logger.error("Must be 4 arguments. [repository_url] [commit_id] [repo_dir] [commited_at]");
    process.exit(9)
}

// Commmend line options
const repository = process.argv[2]
const commitId = process.argv[3]
const basePath = process.argv[4]
const commitedAt = process.argv[5]

logger.info("START webapp parser. ", process.argv);

// Main logic
parse(basePath)
    .then(funcs => {
        const ret = funcs.map(f => {
            return {
                repository: repository,
                commit_id: commitId,
                commited_at: commitedAt,
                ...f
            }
        }).map(f => save(f, basePath))
        Promise.all(ret)
            .catch(err => logger.error("Error occurs: %s", err))
            .finally(_ => Store.end())
    })
    .catch(err => logger.info("Error occurs: %s", err));
