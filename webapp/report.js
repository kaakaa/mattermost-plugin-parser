const fs = require('fs');

const Store = require('./store/store');

const {logger} = require('./logger');

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
        const data = usages.map(u => {
            const i = u.url.lastIndexOf('/');
            const name = u.url.substring(i+1);

            let tmpUrl = u.url;
            if (!u.url.startsWith('http')) {
                tmpUrl = `https://${tmpUrl}`;
                u.url = tmpUrl;
            }

            const refs = u.refs == 'HEAD' ? u.commit_id : u.refs;
            const newUrl = `${tmpUrl}/blob/${refs}/${u.path}#L${u.line}`

            u.loc = `<a href='${newUrl}' target='_blank'>${u.path}#${u.line}</a>`;
            u.name = name;
            return u;
        });
        fs.writeFileSync("../docs/data.json", JSON.stringify(data, null, "  "));

        // Grouping by api type
        const values = data
            .filter((d) => d.refs == "HEAD")
            .reduce((acc, cur) => {
                if (acc[cur.type]) {
                    acc[cur.type].push(cur);
                } else {
                    acc[cur.type] = [cur];
                }
                return acc;
            }, {})
        // Cout usages by API
        const stats = Object.keys(values).map((key) => {
            const s = values[key].reduce((acc, cur) => {
                if (acc[cur.api]) {
                    acc[cur.api]++;
                } else {
                    acc[cur.api] = 1;
                }
                return acc;
            }, {});
            return {
                type: key,
                stats: s,
            };
        })

        fs.writeFileSync("../docs/stats.json", JSON.stringify(stats, null, "  "));

    }).catch(err => logger.error("Failed to output usages: ", err))
    .finally(_ => Store.end());
}

report()
