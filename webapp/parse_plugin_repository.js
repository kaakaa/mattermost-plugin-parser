const esprima = require('esprima');
const fsp = require('fs-promise');
const glob = require('glob');

const findRecursive = async (path) => {
    let ret = [];
    ret = ret.concat(await parseFiles(glob.sync(path + "/**/*.jsx", {ignore: path + "/node_modules/**"}), true));
    ret = ret.concat(await parseFiles(glob.sync(path + "/**/*.js", {ignore: path + "/node_modules/**"}), false));
    return ret
}

const parseFiles = async (files, isJsx) => {
    let ret = [];
    for (i in files) {
        try {
            ret = ret.concat(await parseFile(files[i], isJsx))
        } catch (e) { console.log("  => ERROR: ", e.message)}
    }
    return ret;
}

const parseFile = async (file, isJsx) => {
    const data = fsp.readFileSync(file, 'utf-8')
    console.log("## Parse file: ", file)
    const parsed = esprima.parseModule(data, {
        jsx: isJsx,
        loc: true,
        comment: false,
        attachComment: false
    });

    let ret = [];
    let decls = parsed.body.find(statement =>
        statement.type === "ExportDefaultDeclaration"
    )
    if (decls) {
        ret = ret.concat(parseExportDefaultDeclarations(decls, file));
    }

    decls = parsed.body.find(statement =>
        statement.type === "ClassDeclaration" &&
        statement.body
    )
    if (decls) {
        ret = ret.concat(parseClassDeclarations(decls, file));
    }
    return ret;
}

const parseExportDefaultDeclarations = (decls, file) => {
    return decls.declaration.body.body.filter(statement =>{
        return statement.type === "MethodDefinition" && statement.key.name === "initialize"
    }).flatMap(statement => {
        return statement.value.body.body
    }).filter(statement => {
        return statement.expression.type === "CallExpression"
    }).map(statement => {
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

const parseClassDeclarations = (decls, file) => {
    return decls.body.body.filter(statement =>
        statement.type === "MethodDefinition" &&
        statement.key.name === "initialize"
    ).find(statement => 
        statement.value.body.body
    ).value.body.body.filter(statement =>
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
    findRecursive: findRecursive,
}

