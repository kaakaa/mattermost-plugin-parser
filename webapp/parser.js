const {Parser} = require('acorn');
const parser = Parser.extend(
    require('acorn-jsx')(),
    require('acorn-static-class-features'),
    require('acorn-class-fields'),
);

const parseOption = {
    locations: true,
    allowImportExportEverywhere: true,
}

const parseMattermostPluginRegistryClass = (body) => {
    const parsed = parser.parse(body, parseOption);

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

/**
 * Find usages of plugin methods
 * 
 * @param {String} body source code
 * @param {String} file file path
 */
const parsePluginSource = (body, file) => {
    const parsed = parser.parse(body, parseOption);

    let ret = [];
    let decls = parsed.body.find(statement =>
        statement.type === "ExportDefaultDeclaration"
    )
    if (decls) {
        const funcs = parseExportDefaultDeclarations(decls, file);
        if (funcs) {
            ret = ret.concat(funcs);
        }
    }

    decls = parsed.body.find(statement =>
        statement.type === "ClassDeclaration" &&
        statement.body
    )
    if (decls) {
        const funcs = parseClassDeclarations(decls, file);
        if (funcs) {
            ret = ret.concat(funcs);
        }
    }
    return ret;
}

/**
 * Find usages in ExporDefaultDeclaration
 * 
 * @param {Object} decls ExportDefaultDeclaration object
 * @param {String} file file path
 */
const parseExportDefaultDeclarations = (decls, file) => {
    if (!decls.declaration.body) {
        return;
    }
    return decls.declaration.body.body.filter(statement =>
        statement.type === "MethodDefinition" &&
        statement.key.name === "initialize"
    ).reduce((acc, statement) => acc.concat(statement.value.body.body), []
    ).filter(statement => 
        statement.expression &&
        statement.expression.type === "CallExpression"
    ).map(statement => {
        return statement.expression.callee
    }).map(statement => {
        return {
            file: file,
            identifier: statement.object.name,
            name: statement.property.name,
            loc: statement.loc,
        }
    })
}

/**
 * Find usages in ClassDeclaration
 * 
 * @param {Object} decls ClassDeclaration object
 * @param {String} file file path
 */
const parseClassDeclarations = (decls, file) => {
    return decls.body.body.filter(statement =>
        statement.type === "MethodDefinition" &&
        statement.key.name === "initialize"
    ).find(statement => 
        statement.value.body.body
    ).value.body.body.filter(statement =>
        statement.expression &&
        statement.expression.type === "CallExpression"
    ).map(statement => {
        return statement.expression.callee
    }).map(statement => {
        return {
            file: file,
            identifier: statement.object.name,
            name: statement.property.name,
            loc: statement.loc,
        }
    })
}

module.exports = {
    parseMattermostPluginRegistryClass: parseMattermostPluginRegistryClass,
    parsePluginSource: parsePluginSource,
}