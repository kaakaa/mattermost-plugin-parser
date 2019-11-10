const Store = require('./store/store');

const logger = require('./logger');

const report = async () => {
    await Store.init();

    const usages = await Store.listUsages();
    let cached = {};
    const usagePromises = usages.map(async (u) => {
        let repo = !cached[u.commit_id];
        if (repo) {
            const ret = await Store.selectRepository(u.commit_id)
            if (!ret) {
                return
            }
            repo = ret[0];
            cached[u.commit_id] = repo;
        }
        return data = {
            ...repo,
            ...u,
        }
    })

    Promise.all(usagePromises)
    .then(usages => {
        const usageMap = usages.reduce((pv, cv, index, array) => {
            let api = pv[cv.api];
            if (api) {
                api.push(cv);
            } else {
                api = [cv];
            }
            pv[cv.api] = api;
            return pv
        }, {})

        const ret = [];
        for (key in usageMap) {
            ret.push(toMd(key, usageMap[key]))
        }
        logger.info(ret.join("\n\n"));
    }).catch(err => logger.error("Failed to output usages: ", err))
    .finally(_ => Store.end());
}

const toMd = (api, usages) => {
    const ret = [];
    ret.push(`## ${api}`);
    ret.push(`| Plugin | Path |`);
    ret.push(`|:-------|:-----|`);

    usages.forEach(u => {
        const i = u.url.lastIndexOf('/');
        const name = u.url.substring(i+1);
        ret.push(`| [${name}](${u.url}) | [${u.path}#${u.line}](${u.url}/blob/${u.commit_id}/${u.path}#L${u.line}) |`)
    })
    return ret.join("\n");
}

report()
