// Design:
// Most popular in order:
// - methods
// - structs (stand out)



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

const variantColor = d3
    .scaleOrdinal()
    // .domain(graphJSON.nodeTypes)
    .domain([
        "Struct",
        "Method",
        "Func",
        "Field",
        "RootPackage",
        "File",
        "ImportedPackage",
        "ImportedFunc",
        "FuncCall"
    ])
    .range([
    ]);

// export default variantColor;
export default d3.scaleOrdinal(d3.schemeCategory20);