// Design:
// Most popular in order:
// - methods
// - structs (stand out)
// - fields/vars

// grouping
// * files/rootpackage
// * imports
// * structs
// * methods/funcs
// * vars/fields


// logic and data
// logic = methods/funcs
// data = structs/vars/fields


// https://github.com/kevinsawicki/monokai/blob/master/styles/colors.less

// @ghost-white: #F8F8F0;
// @light-ghost-white: #F8F8F2;
// @light-gray: #CCC;
// @gray: #888;
// @brown-gray: #49483E;
// @dark-gray: #282828;

// @yellow: #E6DB74;
// @blue: #66D9EF;
// @pink: #F92672;
// @purple: #AE81FF;
// @brown: #75715E;
// @orange: #FD971F;
// @light-orange: #FFD569;
// @green: #A6E22E;
// @sea-green: #529B2F;


import graphJSON from '../../graph.json';
import _ from 'underscore';
import * as d3 from 'd3v4';

const config = `
Struct              #aec7e8
Method              #ff7f0e
Func                #ffbb78
Field               #70fafd
RootPackage         #2ca02c
File                #52d352
ImportedPackage     white
ImportedFunc        white
FuncCall            white
`.split('\n').filter(Boolean).map(line => {
    let [ nodeType, colour ] = line.split(/\s+/)
    return {
        nodeType,
        colour
    }
})


// TODO MVP hacking.
// match to Go enum (authoritative)
let domain = graphJSON.nodeTypes.map((nodeType, i) => {
    return _.findIndex(config, { nodeType })
})

export function getVariantName(variant) {
    return graphJSON.nodeTypes[variant]
}

// nodeTypesMapping.map(x => x.nodeType)
let range = config.map(x => x.colour)

const nodeColor = d3
    .scaleOrdinal(range)
    .domain(domain)

export default nodeColor;
// export default d3.scaleOrdinal(d3.schemeCategory20);