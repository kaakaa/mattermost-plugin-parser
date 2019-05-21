
const parseMattermost = require('./parse_mattermost');
const parsePluginRepository = require('./parse_plugin_repository');

const Store = require('./store/store');

/*
const repository = "https://github.com/mattermost/mattermost-plugin-demo";
const commitId = "ffcf4e3e3841e55a801c7e497e3f24f740c79e1d";
const basePath = "./mattermost-plugin-demo";
*/

const parse = async (path) => {
    const pluginResistryClassMethod = await parseMattermost.parsePluginResistryClassMethod();
    console.log("  * Complete to parse mattermost webapp plugin resistry. ");
    const calledFuncs = await parsePluginRepository.findRecursive(path);
    console.log("  * Commplete to parse plugin webapp.")

    return calledFuncs
        .filter(f => pluginResistryClassMethod.some(m => m.name == f.name))
}

const save = async (f, basePath) => {
    console.log("    * Save usages of webapp api: ", f.file, ":", f.loc.start.line)
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
    console.log("Must be 4 arguments. [repository_url] [commit_id] [repo_dir] [commited_at]")
    process.exit(9)
}

const repository = process.argv[2]
const commitId = process.argv[3]
const basePath = process.argv[4]
const commitedAt = process.argv[5]

console.log("Run webapp parser. ", process.argv)

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
            .catch(err => console.log("ERROR: ", err))
            .finally(_ => Store.end())
    })
    .catch(err => console.log("ERROR: ", err))
