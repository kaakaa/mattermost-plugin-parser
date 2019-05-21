const Store = require('./store/store');
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
        console.log(ret.join("\n\n"));
    }).catch(err => console.log("ERROR: ", err))
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

/*
REPOSITORY_NAME
* [${usage.api}](${repository.url}/blob/${repository.commit_id}/${usage.path}#l${usage.line})
*/

/*
## #{usage.api}
| repo | file |
|:-----|:-----|
| [${REPONAME}](${repository.url}) | [${usage.path}](${repository.url}/blob/master/${usage.path}#l${usage.line}) |
*/