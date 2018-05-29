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
    edgeLayout,
    nodeSel
} from 'graphparse';

const UseEdge = 0;
const DefEdge = 1;

// Merge items from coll2 into coll1 by key 'id', overwriting values.
export function mergeByKey(key: string, coll1: any[], coll2: any[], failOnMissing: boolean = true) : any[] {
    let map = {};

    coll1.map(a => {
        map[a.id] = a;
    })
    coll2.map(b => {
       let a = map[b.id];
       if(!a && failOnMissing) throw new Error(`base item not found for key ${b.id}`);
       map[b.id] = {
           ...a,
           ...b
       }
    })

    return coll1.map(a => map[a.id])
}

type refreshArgs = {
    nodes: node[], 
    edges: edge[], 
    currentNode: ?nodeid, 
    maxDepth: number
};

export class GraphLogic {
    nodesLayout: nodeLayout[] = [];
    edgesLayout: edgeLayout[] = [];
    maxDepth: number;
    graphDOT: string = "";
    
    nodes: node[] = [];

    getNodeById (id : nodeid) {
        let node: node = _.findWhere(this.nodes, { id, })
        if(!node) {
            // console.log(this.nodes)
            throw new Error(`node not found: ${id}`)
        }
        return node
    }

    get edges () : edge[] {
        return this.nodes.map(node => {
            return node.outs
        }).reduce((prev, curr) => prev.concat(curr), []);
    }

    get shownNodes () : node[] {
        return this.nodes.filter(node => node.inTree);
    }

    get shownEdges () : edge[] {
        return this.shownNodes.map(node => {
            return node.outs.filter(e => {
                return this.getNodeById(e.target).inTree;
            })
        }).reduce((prev, curr) => prev.concat(curr), []);
    }

    _currentNode: ?nodeid;
    get currentNode () {
        return this._currentNode;
    }
    set currentNode (id: ?nodeid) {
        if(id == null) return;
        this.getNodeById(id);
        this._currentNode = id;
    }

    constructor() {
    }

    refresh({
        nodes,
        edges,
        currentNode,
        maxDepth = 1
    }: refreshArgs) {
        const nodeExists = (id) => _.findWhere(nodes, { id, }) != null;
        let edgesThatExist = edges.filter(edge => {
            return nodeExists(edge.source) && nodeExists(edge.target);
        })

        this.nodes = nodes.map(node => {
            return {
                outs: _.where(edgesThatExist, { source: node.id }),
                ...node,
            }
        });

        this.currentNode = currentNode;
        this.maxDepth = maxDepth;

        this.preFilterNodesAndEdges();
        this.generateLayout();
        this.postFilterNodesAndEdges();
    }

    getLayout() {
        return {
            nodes: mergeByKey('id', this.shownNodes, this.nodesLayout),

            // TODO
            // duplicate edges are filtered out in generateGraphDOT
            // this means that merging here will error unless we relax
            // since there are edges that don't have a layout due to being removed as duplicates
            edges: mergeByKey('id', this.edgesLayout, this.shownEdges, false),
            graphDOT: this.graphDOT,
        }
    }

    preFilterNodesAndEdges() {
        let tree = this.getSpanningTree().map(id => {
            return { id, inTree: true }
        })
        this.nodes = mergeByKey('id', this.nodes, tree)
    }

    postFilterNodesAndEdges() {
        // this.nodes = this.shownNodes;
    }

    // Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
    getSpanningTree() : nodeid[] {
        type traversedNode = {
            fromDef: ?boolean
        } & node;

        let nodesToTraverse: Array<traversedNode> = [];
        let depth = 0;
        let visited: Set<nodeid> = new Set();

        if(this.currentNode == null) {
            return Array.from(visited);
        }
        let current = this.getNodeById(this.currentNode);
        nodesToTraverse.push(current)
        visited.add(current.id)

        const traverse = (parent: traversedNode) : Array<traversedNode> => {
            let parentFromDef = parent.fromDef || true;

            let outs = parent.outs
            .map(out => {
                let child: traversedNode = this.getNodeById(out.target);
                child.fromDef = (out.variant == DefEdge);
                return child;
            })
            .filter(child => {
                if(parentFromDef) {
                    // show defs,uses
                    return true;
                } else {
                    // show uses
                    if(child.fromDef) return false;
                    return true;
                }
            })
            .filter(child => {
                return _.contains(parent.filters.shownNodeTypes, child.variant)
            })

            return outs
        }

        do {
            // visit
            depth++;

            nodesToTraverse = nodesToTraverse
            .map(traverse)
            .reduce((prev, curr) => prev.concat(curr), [])
            .filter(node => {
                if(visited.has(node.id)) return false;
                else {
                    visited.add(node.id);
                    return true;
                }
            })
            .filter(node => {
                if(node.filters.shown) return true;
                return depth < this.maxDepth;
            })

        } while(nodesToTraverse.length);

        return Array.from(visited);
    }

    generateLayout() {
        if(this.shownNodes.length < 1 || this.shownEdges.length < 1) {
            this.nodesLayout = [];
            this.edgesLayout = [];
            this.graphDOT = "";
            return;
        }
    
        // Generate dot layout
        // let graphDOT = this.generateGraphDOT(this.nodes, this.edges)
        let graphDOT = this.generateGraphDOT(this.shownNodes, this.shownEdges)
        let graphvizData = JSON.parse(Viz(graphDOT, {
            format: 'json',
            engine: 'dot',
            // engine: 'dot'
        }));
    
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

        this.graphDOT = graphDOT;
    }

    generateGraphDOT(nodes: node[], edges: edge[]) {
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
    }
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

const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

const toSvgPointSpace = point => [ point[0], point[1] ];

export const edgeRelationId = (edge: edge) => `${edge.source}${edge.target}`;


const logic = new GraphLogic();

type refreshMsg = {
    data: refreshArgs,
    type: 'refresh'
};

self.addEventListener('message', (ev) => {
    let msg: refreshMsg = ev.data;

    switch(msg.type) {
        case 'refresh':
            logic.refresh(msg.data)

            self.postMessage({
                type: 'layout',
                data: logic.getLayout()
            })
    }
})

// $FlowFixMe
if (module.hot) {
    module.hot.accept()
}