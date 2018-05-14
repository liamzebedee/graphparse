// @flow
import Viz from 'viz.js';
import _ from 'underscore';
import lodash from 'lodash';
import graphJSON from '../../graph.json';
import {
    removeDuplicates
} from '../util'

type id = number;
type nodeid = id;
type edgeid = id;

const UseEdge = 0;
const DefEdge = 1;


type edge = {|
    source: nodeid,
    target: nodeid,
    variant: 0 | 1,
    shown: boolean,
    id: edgeid,
|};

type nodeVariant = number;

type node = {
    rank: number,
    label: string,
    id: nodeid,
    variant: nodeVariant,
    pos: string,
    debugInfo: string,

    shown: boolean,
    selected: boolean,

    outs: edge[],
};

type nodeLayout = {|
    layout: {
        cx: number,
        cy: number,
        rx: number,
        ry: number,
    },
    id: nodeid,
|};

type edgeLayout = {|
    layout: {
        points: number[],
        arrowPts: number[],
    },
    id: edgeid,
|};

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

export class GraphLogic {
    nodesLayout: nodeLayout[] = [];
    edgesLayout: edgeLayout[] = [];
    selection: nodeid[] = [];
    maxDepth: number = 1;
    graphDOT: string = "";
    
    nodes: node[] = [];

    getNodeById (id : nodeid) {
        let node = _.findWhere(this.nodes, { id, })
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
        return this.nodes.filter(node => node.shown);
    }

    get shownEdges () : edge[] {
        return this.shownNodes.map(node => {
            return node.outs.filter(e => {
                return this.getNodeById(e.target).shown;
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

    refresh(nodes: node[], edges: edge[], currentNode: ?nodeid, selection: nodeid[], maxDepth: number = 3) {
        const nodeExists = (id) => _.findWhere(nodes, { id, }) != null;
        let edgesThatExist = edges.filter(edge => {
            return nodeExists(edge.source) && nodeExists(edge.target);
        })

        this.nodes = nodes.map(node => {
            return {
                ...node,

                shown: false,
                selected: false,
                outs: _.where(edgesThatExist, { source: node.id })
            }
        });

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
            return { id, shown: true }
        })
        this.nodes = mergeByKey('id', this.nodes, tree)
    }

    postFilterNodesAndEdges() {
        return
    }

    // Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
    getSpanningTree() : nodeid[] {
        type traversedNode = {
            fromDef: ?boolean
        } & node;

        let nodesToTraverse: Array<traversedNode> = [];
        let depth = -1;
        let visited = new Set();

        if(this.currentNode == null) {
            return Array.from(visited);
        }
        nodesToTraverse.push(this.getNodeById(this.currentNode))

        const traverse = (parent: traversedNode) : Array<traversedNode> => {
            let parentFromDef = parent.fromDef || true;

            let outs = parent.outs.map(out => {
                let child: traversedNode = this.getNodeById(out.target);
                child.fromDef = (out.variant == DefEdge);
                return child;
            }).filter(child => {
                if(parentFromDef) {
                    // show defs,uses
                    return true;
                } else {
                    // show uses
                    if(child.fromDef) return false;
                    return true;
                }
            })

            return outs
        }

        do {
            // visit
            depth++;

            // two cases
            // 1. parent is only a def chain -> show defs,uses
            // 2. parent is a use chain -> show only more uses

            // if(edge.variant == DefEdge && parent.fromDefChain) 
            // if(parent.fromDefChain == null) fromDefChain = true

            nodesToTraverse = nodesToTraverse.map(traverse)
            .reduce((prev, curr) => prev.concat(curr), [])
            .filter(node => {
                if(visited.has(node.id)) return false;
                else {
                    visited.add(node.id);
                    return true;
                }
            });

        } while(depth < this.maxDepth && nodesToTraverse.length);

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
        let graphDOT = this.generateGraphDOT(this.shownNodes, this.shownEdges)
        let graphvizData = JSON.parse(Viz(graphDOT, { format: 'json' }));
    
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

export const edgeRelationId = (edge: edge) => `${edge.source}${edge.target}`;

