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
    if (pluginRegistryClassMethod.length == 0) {
        logger.error("No plugin methods are found.");
        process.exit(9);
    }
    logger.info('Complete: %d functions are detected.', pluginRegistryClassMethod.length);
    pluginRegistryClassMethod.forEach(m => logger.info('  %s', m.name));

    // Parse plugin source code
    logger.info('Parse plugin source code')
    const calledFuncs = await parsePluginRepository.findRecursive(path);
    if (calledFuncs.length == 0) {
        console.log('No functions are detected.')
        process.exit(0);
    }
    logger.info("Complete: %d functions are detected.", calledFuncs.length)
    calledFuncs.forEach(f => logger.info('  %s.%s:%d (%s)', f.identifier, f.name, f.loc.start.line, f.file));

    return calledFuncs
        .filter(f => pluginRegistryClassMethod.some(m => m.name == f.name))
}

const save = async (f, basePath) => {
    logger.info("Save usages of webapp api: %s: %d", f.file, f.loc.start.line);

    await Store.init()
    let rows = await Store.selectRepository([f.commit_id])
    if (rows.length === 0) {
        rows = await Store.insertRepository([f.repository, f.commit_id, f.commited_at, f.commit_refs])
    }

    const regexp = new RegExp(basePath + '/', 'g');
    const relPath = f.file.replace(regexp, '');
    return await Store.insertUsage([f.commit_id, f.name, relPath, f.loc.start.line])
}

if (process.argv.length !== 7) {
    logger.error("Must be 4 arguments. [repository_url] [commit_id] [repo_dir] [commited_at]");
    process.exit(9)
}

// Commmend line options
const repository = process.argv[2]
const commitId = process.argv[3]
const basePath = process.argv[4]
const commitedAt = process.argv[5]
const commitRefs = process.argv[6]

logger.info("START webapp parser. ", process.argv);

// Main logic
parse(basePath)
    .then(funcs => {
        const ret = funcs.map(f => {
            return {
                repository: repository,
                commit_id: commitId,
                commited_at: commitedAt,
                commit_refs: commitRefs,
                ...f
            }
        }).map(f => save(f, basePath))
        Promise.all(ret)
            .catch(err => logger.error("Error occurs: %s", err))
            .finally(_ => Store.end())
    })
    .catch(err => logger.info("Error occurs: %s", err));
