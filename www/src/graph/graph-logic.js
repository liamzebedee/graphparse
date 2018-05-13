// @flow
import Viz from 'viz.js';
import _ from 'underscore';
import lodash from 'lodash';
import graphJSON from '../../graph.json';
import mergeByKey from "array-merge-by-key";

type id = number;
type nodeid = id;
type edgeid = id;
type edgeVariant = number;

type edge = {
    source: nodeid,
    target: nodeid,
    variant: edgeVariant,
    id: edgeid,
};

type nodeVariant = number;

type node = {
    rank: number,
    label: string,
    id: nodeid,
    variant: nodeVariant,
    pos: string,
    debugInfo: string,
};

type nodeLayout = {
    layout: {
        cx: number,
        cy: number,
        rx: number,
        ry: number,
    },
    id: nodeid,
};

type edgeLayout = {
    layout: {
        points: number[],
        arrowPts: number[],
    },
    id: edgeid,
};

function mergeByKey2(key, coll1, coll2) {
    let els = [];
    for(let a of coll1) {
        let b = _.findWhere(coll2, { id: a.id });
        if(!b) throw new Error();
        els.push({
            ...a,
            ...b,
        })
    }
    return els;
}

export class GraphLogic {
    nodesLayout: nodeLayout[] = [];
    edgesLayout: edgeLayout[] = [];
    currentNode: ?nodeid;
    selection: nodeid[] = [];
    maxDepth: number = 1;
    graphDOT: string = "";
    edges: edge[] = [];
    _nodes: node[] = [];
    nodesLookup: node[] = [];
    edgesLookup: edge[] = [];

    get nodes () {
        return this._nodes;
    }
    set nodes (nodes: node[]) {
        this.edges = this.edges.filter(edge => {
            return lodash.find(nodes, { id: edge.source })
        })
        this._nodes = nodes;
    }

    constructor() {
    }
    
    refresh(nodes: node[], edges: edge[], currentNode: ?nodeid, selection: nodeid[], maxDepth: number) {
        this.nodesLookup = nodes;
        this.edgesLookup = edges;

        this.nodes = nodes;
        this.edges = edges;
        this.currentNode = currentNode;
        this.selection = selection;
        this.maxDepth = maxDepth;

        // Merge nodes with selection information
        // nodes = mergeByKey(nodes, selection);

        this.preFilterNodesAndEdges();
        this.generateLayout();
        this.postFilterNodesAndEdges();
    }

    getLayout() {
        return {
            nodes: mergeByKey2('id', this.nodesLayout, this.nodesLookup),
            edges: mergeByKey2('id', this.edgesLayout, this.edges),
            graphDOT: this.graphDOT,
        }
    }

    preFilterNodesAndEdges() {
        this.nodes = this.getSpanningTree().map(id => {
            let node = _.findWhere(this.nodesLookup, { id, });
            return node;
        })
        // this.edges = this.getEdgesForNodes()
    
        // const showOnlyInterestedNodes = (edge) => {
        //     return _.contains(interested, edge.source)
        // }
    
        // edges = edges
        // .filter(showOnlyInterestedNodes)
        // // .filter(edge => edge.variant == 0)
    
        // let nodesToInclude = edges.map(e => [e.source, e.target]).reduce((a, b) => a.concat(b), []);
        // nodes = nodes.filter(node => {
        //     return _.contains(nodesToInclude, node.id)
        // })
    }

    postFilterNodesAndEdges() {
        return
    }

    // Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
    getSpanningTree(nodes: node[] = this.nodesLookup, edges: edge[] = this.edges) : nodeid[] {
        let nodesToTraverse: nodeid[] = [];
        let adjList = constructAdjList(nodes, edges);
        let depth = -1;
        let visited = new Set();

        if(this.currentNode == null) {
            return nodesToTraverse;
        }
        nodesToTraverse.push(this.currentNode)

        do {
            // visit
            depth++;

            nodesToTraverse = nodesToTraverse.map(id => {
                if(visited.has(id)) return [];
                else {
                    visited.add(id)
                    let outs = adjList[id];
                    return outs;
                }
            }).reduce((prev, curr) => prev.concat(curr), [])

        } while(depth < this.maxDepth && nodesToTraverse.length);

        return Array.from(visited);
    }

    generateLayout() {
        if(this.nodes.length < 1 || this.edges.length < 1) {
            this.nodesLayout = [];
            this.edgesLayout = [];
            return;
        }
    
        // Generate dot layout
        this.graphDOT = this.generateGraphDOT()
        let graphvizData = JSON.parse(Viz(this.graphDOT, { format: 'json' }));
    
        this.nodesLayout = graphvizData.objects.map(obj => {    
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
        });
    
        this.edgesLayout = graphvizData.edges.map((edge, i) => {
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
        });
    }

    generateGraphDOT() {
        const nodes = this.nodes;
        const edges = this.edges;

        let edgeWeights = {};

        edges.map(edge => {
            let id = edgeRelationId(edge);
            let weight = edgeWeights[id] || 0;
            edgeWeights[id] = weight+1;
        })

        let seenEdges = [];
        function removeDuplicates(edge) {
            let id = edgeRelationId(edge); 
            if(_.contains(seenEdges, id)) return false;
            seenEdges.push(id)
            return true
        }

        let weightedEdges = edges.filter(removeDuplicates).map(edge => {
            return { 
                ...edge,
                weight: edgeWeights[edgeRelationId(edge)],
            }
        })

        return `
            digraph graphname {
                ${nodes.map(({ id, rank, label }) => {
                    rank = 1;
                    return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"];`
                }).join('\n')}

                ${weightedEdges.map(({ source, target, id, weight }) => `"${target}" -> "${source}" [id=${id}] [weight=${weight}];`).join('\n')}
            }
        `
    }
}

const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

const toSvgPointSpace = point => [ point[0], point[1] ];

const edgeRelationId = (edge) => `${edge.source}${edge.target}`;



type AdjList = {
    [nodeid]: nodeid[]
};

export function constructAdjList(nodes: node[], edges: edge[]) : AdjList {
    let adjList: AdjList = {};

    edges.map(edge => {
        adjList[edge.source] = [];
        adjList[edge.target] = [];
    })
    nodes.map(node => {
        adjList[node.id] = [];
    })
    
    edges.map(({ source, target }) => {
        adjList[source].push(target);
    })
    return adjList;
}