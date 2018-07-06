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
    getCurrentNodeSelection,
    getEdgeById
} from './selectors';

import {
    compose
} from 'redux';


const UseEdge = 0;
const DefEdge = 1;

type visitedNode = {
    node: node,
    edge: edge,
};
type nodeVisitor = (parent: node) => Array<visitedNode>;


// The visitor returns a list of nodes to continue visiting. 
// The philosophy is that we consider edges and we return nodes.
// We work on the parent and not any 'children' (in or out edges).
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

    return (parent: node) : Array<visitedNode> => {
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
                return { 
                    node, 
                    edge,
                    useEdge, 
                    defEdge
                }
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
            .map(({ node, edge }) => {
                return { node, edge };
            }),

            parent.outs
            .map(edgeContext(edge => edge.target))
            // .filter(depthFilter(parent.selection.outs, baseDepth))
            .filter(relsShownFilter(parent.selection.outs))
            .filter(({ node })    => shownFilter(node))
            .filter(({ useEdge }) => usesFilter(parent.selection.outs, useEdge))
            .filter(({ defEdge }) => defsFilter(parent.selection.outs, defEdge))
            .filter(({ node })    => nodeTypesFilter(parent.selection.outs.shownNodeTypes, node))
            .map(({ node, edge }) => {
                return { node, edge };
            }),
        );
    }
}

const flattenArray = (prev, curr) => prev.concat(curr);


// Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
function generateSpanningTree(state: graphState) : graphState {
    let fromNode: ?nodeid = state.currentNode;

    let nodes = state.nodes;

    let tree: nodeid[] = [];
    let visited: Set<nodeid> = new Set();
    let edges: Set<edgeid> = new Set();
    let toVisit: Array<node> = [];

    if(fromNode == null) {
        return state;
    }
    let depth = 0;
    let current = getNodeById(nodes, fromNode)
    current = {
        ...current,
        selection: getCurrentNodeSelection(current),
    }
    toVisit.push(current)
    visited.add(current.id)

    do {
        // visit
        depth++;

        toVisit = toVisit
        .map(makeVisitor(nodes, depth))
        .reduce(flattenArray, [])
        .map(({ node, edge } : visitedNode) => {
            edges.add(edge.id);
            return node;
        })
        .filter(node => {
            if(visited.has(node.id)) return false;
            else {
                visited.add(node.id);
                return true;
            }
        })
    
    } while(toVisit.length > 0);

    let g = new Graph(
        state.nodes,
        state.edges,
    );

    g.nodes = Array.from(visited).map(id => getNodeById(state.nodes, id));
    g.edges = Array.from(edges).map(id => getEdgeById(state.edges, id));
    
    return {
        ...state,
        ...g.pack(),
    }
}




class Graph {
    _nodes: node[];
    _edges: edge[];

    constructor(nodes, edges) {
        this._nodes = nodes;
        this._edges = edges;
    }

    set edges(edges) {
        // update nodes
        this._nodes = this._nodes.map(node => {
            return {
                ...node,
                ins: _.where(edges,  { source: node.id }),
                outs: _.where(edges, { target: node.id }),
            }
        });
        this._edges = edges.map(edge => {
            return {
                ...edge
            }
        });
    }

    set nodes(nodes) {
        // update edges
        this._edges = getEdges(nodes);
        this._nodes = nodes.map(node => {
            return {
                ...node
            }
        })
    }

    get edges() {
        return this._edges;
    }

    get nodes() {
        return this._nodes;
    }

    pack() {
        let nodes = this._nodes;
        let edges = this._edges;
        return {
            nodes,
            edges,
        }
    }
}

const getOuts = (edges: edge[], id: nodeid) => {
    return _.where(edges, { source: id })
}

const getIns = (edges: edge[], id: nodeid) => {
    return _.where(edges, { target: id })
}


function postFilterGraph(state : graphState) : graphState {
    // Hide DEF edges if the node is already linked to.
    // let { nodes, edges } = graph(state.nodes, state.edges);
    let g = new Graph(state.nodes, state.edges);

    g.nodes = g.nodes.map(node => {
        return {
            ...node,
            ins: node.ins.length > 1 
                  ? node.ins.filter(edge => edge.variant != DefEdge) 
                  : node.ins
        }
    });

    return {
        ...state,
        ...g.pack()
    }
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

    // ${weightedEdges.map(({ source, target, id, weight }) => `"${target}" -> "${source}" [id=${id}] [weight=${weight}];`).join('\n')}

    return `
        digraph graphname {
            graph [ordering=in];
            rankdir=LR;
            ${nodes.map(({ id, rank, label, shown }) => {
                // rank = 1;

                let fixedPos = "";

                return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"] ${fixedPos};`
            }).join('\n')}

            ${weightedEdges.map(({ source, target, id, weight }) => `"${source}" -> "${target}" [id=${id}] [weight=${weight}];`).join('\n')}
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




type layout = {|
    nodes: nodeLayout[],
    edges: edgeLayout[],
|};

function buildLayout(state: graphState) : layout {
    let { nodes, edges } = state;

    let layout: layout = {
        nodes: [],
        edges: [],
    };

    if(nodes.length == 0) {
        return layout;
    }

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

    let layout = compose(buildLayout, postFilterGraph, generateSpanningTree)(state);
    
    return {
        spanningTree: null,
        layout,
    }
}
