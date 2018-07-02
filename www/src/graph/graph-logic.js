// @flow
import Viz from 'viz.js';
import _ from 'underscore';
import lodash from 'lodash';
import graphJSON from '../../graph.json';
import {
    removeDuplicates
} from '../util'

import type {
    nodeid,
    node,
    nodeLayout,
    edge,
    edgeid,
    edgeLayout,
    nodeSel,
    relationshipsSel,
    graphState
} from 'graphparse';

import {
    getNodeSelection,
    getEdges,
    getNodeById,
    getCurrentNodeSelection
} from './selectors';


const UseEdge = 0;
const DefEdge = 1;

type visitedNode = {
    fromDef: ?boolean,
    depth: number,
} & node;

type nodeVisitor = (parent: visitedNode) => Array<visitedNode>;



// if the edges are shown, show them to max depth


// The visitor returns a list of nodes to continue visiting. 
// The philosophy is that we consider edges and we return nodes.
// We work on the parent and not any 'children' (in or out edges).
type depthMap = { [nodeid]: number};
// , depthMap: depthMap

function makeVisitor(nodes: node[], baseDepth: number): nodeVisitor {
    const relsShownFilter = (relsel: relationshipsSel) => {
        return () => relsel.shown;
    }
    const shownFilter = (node: node) => node.selection.shown;
    const depthFilter = (relsel: relationshipsSel, depth: number) => {
        return (node) => depth < relsel.maxDepth;
    }
    const usesFilter = (relsel, useEdge) => useEdge ? relsel.showUses : true;
    const defsFilter = (relsel, defEdge) => defEdge ? relsel.showDefs : true;
    const nodeTypesFilter = (shownNodeTypes, child) => {
        return _.contains(shownNodeTypes, child.variant)
    }

    return (parent: visitedNode) : Array<visitedNode> => {
        // if(!shownFilter(parent)) return [];

        if(
            [parent]
            .filter(shownFilter)
            .length == 0
        ) return [];

        const edgeContext = (getNode: (edge) => nodeid) => {
            return (edge: edge) => {
                let node = getNodeById(nodes, getNode(edge));
                let useEdge = edge.variant == UseEdge;
                let defEdge = edge.variant == DefEdge;
                return { node, useEdge, defEdge }
            }
        }

        return [].concat(
            parent.ins
            .map(edgeContext(edge => edge.source))
            // .filter(depthFilter(parent.selection.ins, baseDepth))
            .filter(relsShownFilter(parent.selection.ins))
            .filter(({ node })    => shownFilter(node))
            .filter(({ useEdge }) => usesFilter(parent.selection.ins, useEdge))
            .filter(({ defEdge }) => defsFilter(parent.selection.ins, defEdge))
            .filter(({ node })    => nodeTypesFilter(parent.selection.ins.shownNodeTypes, node))
            .map(({ node })       => node),

            parent.outs
            .map(edgeContext(edge => edge.target))
            // .filter(depthFilter(parent.selection.outs, baseDepth))
            .filter(relsShownFilter(parent.selection.outs))
            .filter(({ node })    => shownFilter(node))
            .filter(({ useEdge }) => usesFilter(parent.selection.outs, useEdge))
            .filter(({ defEdge }) => defsFilter(parent.selection.outs, defEdge))
            .filter(({ node })    => nodeTypesFilter(parent.selection.outs.shownNodeTypes, node))
            .map(({ node })       => node),
        );
    }
}

// Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
function generateSpanningTree(state: graphState) : nodeid[] {
    let fromNode: ?nodeid = state.currentNode;

    let nodes = state.nodes;

    let tree: nodeid[] = [];
    let visited: Set<nodeid> = new Set();
    let toVisit: Array<visitedNode> = [];

    if(fromNode == null) {
        return Array.from(visited);
    }
    let depth = 0;
    let current = getNodeById(nodes, fromNode)
    current = {
        ...current,
        selection: getCurrentNodeSelection(current),
        depth,
    }
    toVisit.push(current)
    visited.add(current.id)

    do {
        // visit
        depth++;
        // let depthMap: depthMap = {
        //     fromNode: 0,
        // };



        toVisit = toVisit
        .map(makeVisitor(nodes, depth))
        .reduce((prev, curr) => prev.concat(curr), [])
        .filter(node => {
            if(visited.has(node.id)) return false;
            else {
                visited.add(node.id);
                return true;
            }
        })
    
    } while(toVisit.length > 0);
    
    return Array.from(visited);
}



const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

const toSvgPointSpace = point => [ point[0], point[1] ];

export const edgeRelationId = (edge: edge) => `${edge.source}${edge.target}`;




function generateGraphDOT(nodes: node[], edges: edge[]) {
    let edgeWeights = {};

    edges.map(edge => {
        let id = edgeRelationId(edge);
        let weight = edgeWeights[id] || 0;
        edgeWeights[id] = weight+1;
    })

    let weightedEdges = edges.filter(removeDuplicates(edgeRelationId)).map(edge => {
        return { 
            ...edge,
            weight: edgeWeights[edgeRelationId(edge)],
        }
    })

    return `
        digraph graphname {
            graph [ordering=in];
            ${nodes.map(({ id, rank, label, shown }) => {
                // rank = 1;

                let fixedPos = "";
                // if(shown) {
                //     let prevlayout: ?nodeLayout = _.findWhere(this.nodesLayout, { id, });
                //     if(!prevlayout) throw new Error("Not found");

                //     let { cx, cy } = prevlayout.layout;
                //     fixedPos = `[pos="${cx},${cy}!"]`;
                // }

                return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"] ${fixedPos};`
            }).join('\n')}

            ${weightedEdges.map(({ source, target, id, weight }) => `"${target}" -> "${source}" [id=${id}] [weight=${weight}];`).join('\n')}
        }
    `
};

function generateElk(nodes: node[], edges: edge[]) {
    let edgeWeights = {};

    edges.map(edge => {
        let id = edgeRelationId(edge);
        let weight = edgeWeights[id] || 0;
        edgeWeights[id] = weight+1;
    })

    let weightedEdges = edges.filter(removeDuplicates(edgeRelationId)).map(edge => {
        return { 
            ...edge,
            weight: edgeWeights[edgeRelationId(edge)],
        }
    })
    return `
algorithm: layered

${nodes.map(({ id, rank, label, shown }) => {
    return `
    node id${id} {
        nodeLabelPlacement: "INSIDE H_CENTER V_CENTER"
        label "${label}"
    }
    `
}).join('\n')}



${weightedEdges.map(({ source, target, id, weight }) => {
    return `edge id${target} -> id${source}`
}).join('\n')}
}
    `
}


/*

digraph {
    node parse.go
    node generatecodegraph

    parse.go -> generatecodegraph
    generatecodegraph -> generatecodegraphfromprog
    generatecodegraphfromprog -> generatecodegraphfromprog_body
 
    subgraph generategraphfromprog {
        node newgraph
        node pkginfo
        node visit

        newgraph -> pkginfo
        pkginfo -> visit
        visit -> generategraphfromprog_body
    }

    subgraph generategraphfromprog_body {
        parseimportspec
        parse2
        parse3
        parse4
        parse5
    }
}

*/



type layout = {|
    nodes: nodeLayout[],
    edges: edgeLayout[],
|};

function buildLayout(state: graphState, nodes: node[]) : layout {
    let layout: layout = {
        nodes: [],
        edges: [],
    };

    if(nodes.length == 0) {
        return layout;
    }

    let edges = getEdges(nodes);

    // Generate dot layout
    let graphDOT = generateGraphDOT(nodes, edges)
    // console.log(generateElk(nodes, edges))
    let graphvizData = JSON.parse(Viz(graphDOT, {
        format: 'json',
        // engine: 'neato',
        engine: 'dot'
    }));
    
    layout.nodes = nodes.length > 0 ? graphvizData.objects.map(obj => {    
        let pos = obj.pos.split(',').map(Number);
        let id = parseInt(obj.name)

        return {
            layout: {
                cx: pos[0],
                cy: pos[1],
                rx: obj._draw_[1].rect[2],
                ry: obj._draw_[1].rect[3],
            },
            id,
        };
    }) : [];

    layout.edges = edges.length > 1 ? graphvizData.edges.map((edge, i) => {
        let points = edge._draw_[1].points.map(toSvgPointSpace);
        let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace);
        let id = parseInt(edge.id);

        return {
            layout: {
                points,
                arrowPts,
            },
            id,
        };
    }) : [];

    // layout.graphDOT = graphDOT;
    
    return layout;
}


const arrayEqual = (a, b) => {
    return a.filter(e => !b.includes(e)).length > 0;
}

export function graphLogic(state: graphState) {
    state.nodes = state.nodes.map(node => {
        return {
            ...node,
            selection: getNodeSelection(node),
        }
    })

    let spanningTree = generateSpanningTree(state)
    let layout = buildLayout(state, spanningTree.map(id => getNodeById(state.nodes, id)));

    return {
        spanningTree,
        layout,
    }
}


// $FlowFixMe
if (module.hot) {
    module.hot.accept()
}