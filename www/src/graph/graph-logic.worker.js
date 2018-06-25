// // @flow
// import Viz from 'viz.js';
// import _ from 'underscore';
// import lodash from 'lodash';
// import graphJSON from '../../graph.json';
// import {
//     removeDuplicates
// } from '../util'

// import type {
//     nodeid,
//     node,
//     nodeLayout,
//     edge,
//     edgeid,
//     edgeLayout,
//     nodeSel,
//     relationshipsSel
// } from 'graphparse';

// import { getNodeSelection } from './selectors';



// const UseEdge = 0;
// const DefEdge = 1;

// type visitedNode = {
//     fromDef: ?boolean
// } & node;

// type nodeVisitor = (parent: visitedNode) => Array<visitedNode>;

// function getNodeById(nodes: node[], id: nodeid) {
//     let node: node = _.findWhere(nodes, { id, })
//     if(!node) {
//         throw new Error(`node not found: ${id}`)
//     }
//     return node
// }



// function getEdges(nodes: node[]) : edge[] {
//     return nodes.map(node => {
//         return [].concat(node.ins, node.outs).filter((edge: edge) => {
//             return _.findWhere(nodes, { id: edge.target }) && _.findWhere(nodes, { id: edge.source })
//         })
//     }).reduce((prev, curr) => prev.concat(curr), []);
// }


// // getLayout() {
// //     return {
// //         nodes: mergeByKey('id', this.shownNodes, this.nodesLayout),

// //         // TODO
// //         // duplicate edges are filtered out in generateGraphDOT
// //         // this means that merging here will error unless we relax
// //         // since there are edges that don't have a layout due to being removed as duplicates
// //         edges: mergeByKey('id', this.edgesLayout, this.shownEdges, false),
// //         graphDOT: this.graphDOT,
// //     }
// // }

// // preFilterNodesAndEdges() {
// //     let tree = this.getSpanningTree().map(id => {
// //         return { id, inTree: true }
// //     })
// //     this.nodes = mergeByKey('id', this.nodes, tree)
// // }

// // postFilterNodesAndEdges() {
// //     // this.nodes = this.shownNodes;
// //     this.nodes = this.nodes.map(node => {
// //         return {
// //             ...node,
// //             selection: node.selection[FULL]
// //         }
// //     })
// // }

// function makeVisitor(nodes: node[], depth: number): nodeVisitor {
//     const relsShownFilter = (relsel: relationshipsSel) => () => relsel.shown;
//     const shownFilter = (node: node) => node.selection.shown;
//     const depthFilter = (relsel: relationshipsSel) => {
//         return (node) => depth < relsel.maxDepth;
//     }
//     const usesFilter = (relsel, useEdge) => useEdge ? relsel.showUses : true;
//     const defsFilter = (relsel, defEdge) => defEdge ? relsel.showDefs : true;
//     const nodeTypesFilter = (shownNodeTypes, child) => {
//         return _.contains(shownNodeTypes, child.variant)
//     }

//     return (parent: visitedNode) : Array<visitedNode> => {
//         // if(!shownFilter(parent)) return [];
//         if(
//             [parent]
//             .filter(shownFilter)
//             .length == 0
//         ) return [];

//         const edgeContext = (getNode: (edge) => nodeid) => {
//             return (edge: edge) => {
//                 let node = getNodeById(nodes, getNode(edge));
//                 let useEdge = edge.variant == UseEdge;
//                 let defEdge = edge.variant == DefEdge;
//                 return { node, useEdge, defEdge }
//             }
//         }

//         return [].concat(
//             parent.ins
//             .map(edgeContext(edge => edge.source))
//             .filter(depthFilter(parent.selection.ins))
//             .filter(relsShownFilter(parent.selection.ins))
//             .filter(({ node })    => shownFilter(node))
//             .filter(({ useEdge }) => usesFilter(parent.selection.ins, useEdge))
//             .filter(({ defEdge }) => defsFilter(parent.selection.ins, defEdge))
//             .filter(({ node })    => nodeTypesFilter(parent.selection.ins.shownNodeTypes, node))
//             .map(({ node })       => node),

//             parent.outs
//             .map(edgeContext(edge => edge.target))
//             .filter(depthFilter(parent.selection.outs))
//             .filter(relsShownFilter(parent.selection.outs))
//             .filter(({ node })    => shownFilter(node))
//             .filter(({ useEdge }) => usesFilter(parent.selection.outs, useEdge))
//             .filter(({ defEdge }) => defsFilter(parent.selection.outs, defEdge))
//             .filter(({ node })    => nodeTypesFilter(parent.selection.outs.shownNodeTypes, node))
//             .map(({ node })       => node),
//         );
//     }
// }

// // Returns array of node id's that represent the spanning tree of from currentNode down, until maxdepth
// function generateSpanningTree(state: graphState) : nodeid[] {
//     let fromNode: ?nodeid = state.currentNode;

//     let nodes = state.nodes;

//     let tree: nodeid[] = [];
//     let visited: Set<nodeid> = new Set();
//     let depth = 0;
//     let toVisit: Array<visitedNode> = [];

//     if(fromNode == null) {
//         return Array.from(visited);
//     }
//     let current = getNodeById(nodes, fromNode);
//     current.selection.shown = true; // TODO hack
//     toVisit.push(current)
//     visited.add(current.id)

//     do {
//         // visit
//         depth++;
//         // debugger
    
//         toVisit = toVisit
//         .map(makeVisitor(nodes, depth))
//         .reduce((prev, curr) => prev.concat(curr), [])
//         .filter(node => {
//             if(visited.has(node.id)) return false;
//             else {
//                 visited.add(node.id);
//                 return true;
//             }
//         })
    
//     } while(toVisit.length > 0);
    
//     return Array.from(visited);
// }



// const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

// const toSvgPointSpace = point => [ point[0], point[1] ];

// export const edgeRelationId = (edge: edge) => `${edge.source}${edge.target}`;




// function generateGraphDOT(nodes: node[], edges: edge[]) {
//     let edgeWeights = {};

//     edges.map(edge => {
//         let id = edgeRelationId(edge);
//         let weight = edgeWeights[id] || 0;
//         edgeWeights[id] = weight+1;
//     })

//     let weightedEdges = edges.filter(removeDuplicates(edgeRelationId)).map(edge => {
//         return { 
//             ...edge,
//             weight: edgeWeights[edgeRelationId(edge)],
//         }
//     })

//     return `
//         digraph graphname {
//             graph [ordering=in];
//             ${nodes.map(({ id, rank, label, shown }) => {
//                 // rank = 1;

//                 let fixedPos = "";
//                 // if(shown) {
//                 //     let prevlayout: ?nodeLayout = _.findWhere(this.nodesLayout, { id, });
//                 //     if(!prevlayout) throw new Error("Not found");

//                 //     let { cx, cy } = prevlayout.layout;
//                 //     fixedPos = `[pos="${cx},${cy}!"]`;
//                 // }

//                 return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"] ${fixedPos};`
//             }).join('\n')}

//             ${weightedEdges.map(({ source, target, id, weight }) => `"${target}" -> "${source}" [id=${id}] [weight=${weight}];`).join('\n')}
//         }
//     `
// };

// /*

// digraph {
//     node parse.go
//     node generatecodegraph

//     parse.go -> generatecodegraph
//     generatecodegraph -> generatecodegraphfromprog
//     generatecodegraphfromprog -> generatecodegraphfromprog_body
 
//     subgraph generategraphfromprog {
//         node newgraph
//         node pkginfo
//         node visit

//         newgraph -> pkginfo
//         pkginfo -> visit
//         visit -> generategraphfromprog_body
//     }

//     subgraph generategraphfromprog_body {
//         parseimportspec
//         parse2
//         parse3
//         parse4
//         parse5
//     }
// }

// */
// /*


// type traversedNode = {
//     fromDef: ?boolean
// } & node;

// let nodesToTraverse: Array<traversedNode> = [];
// let depth = 0;
// let visited: Set<nodeid> = new Set();

// if(fromNode == null) {
//     return Array.from(visited);
// }
// let current = this.getNodeById(fromNode);
// nodesToTraverse.push(current)
// visited.add(current.id)

// const traverse = (parent: traversedNode) : Array<traversedNode> => {
//     let parentFromDef = parent.fromDef || true;
//     let outs = parent.outs

//     .map(out => {
//         let child: traversedNode = this.getNodeById(out.target);
//         child.fromDef = (out.variant == DefEdge);
//         return child;
//     })
//     .filter(child => {
//         if(parentFromDef) {
//             // show defs,uses
//             return true;
//         } else {
//             // show uses
//             if(child.fromDef) return false;
//             return true;
//         }
//     })
//     .filter(child => {
//         // return _.contains(parent.filters.shownNodeTypes, child.variant)
//     })

//     return outs
// }

// do {
//     // visit
//     depth++;

//     nodesToTraverse = nodesToTraverse
//     .map(traverse)
//     .reduce((prev, curr) => prev.concat(curr), [])
//     .filter(node => {
//         if(visited.has(node.id)) return false;
//         else {
//             visited.add(node.id);
//             return true;
//         }
//     })
//     .filter(node => {
//         if(node.filters.shown) return true;
//         return depth < this.maxDepth;
//     })

// } while(nodesToTraverse.length);

// return Array.from(visited);






// let parentFromDef = parent.fromDef || true;
//             let outs = parent.outs
        
//             // .filter(node => {
//             //     if(node.filters.shown) return true;
//             //     return depth < this.maxDepth;
//             // })

//             .map(out => {
//                 let child: visitedNode = this.getNodeById(out.target);
//                 child.fromDef = (out.variant == DefEdge);
//                 return child;
//             })
//             .filter(child => {
//                 if(parentFromDef) {
//                     // show defs,uses
//                     return true;
//                 } else {
//                     // show uses
//                     if(child.fromDef) return false;
//                     return true;
//                 }
//             })
//             .filter(child => {
//                 // return _.contains(parent.filters.shownNodeTypes, child.variant)
//             })




//             // // let ins = parent.ins
//             // .map(edge => {
//             //     let node = this.getNodeById(edge.source);
//             //     let useEdge = edge.variant == UseEdge;
//             //     let defEdge = edge.variant == DefEdge;
//             //     return { node, useEdge, defEdge }
//             // })
            

//             // let outs = parent.outs
//             // .map(edge => {
//             //     let node = this.getNodeById(edge.target);
//             //     let useEdge = edge.variant == UseEdge;
//             //     let defEdge = edge.variant == DefEdge;
//             //     return { node, useEdge, defEdge }
//             // })
//             // .filter(({ node }) => shownFilter(node))
//             // .filter(({ useEdge }) => usesFilter(parent.selection.ins, useEdge)
//             // .filter(({ defEdge }) => defsFilter(parent.selection.ins, defEdge)
//             // .filter(({ node }) => nodeTypesFilter(parent.selection.ins.shownNodeTypes, node))
//             // .map(({ node }) => node)
//             */



// type layout = {|
//     nodes: nodeLayout[],
//     edges: edgeLayout[],
// |};

// function buildLayout(state: graphState, nodes: node[]) : layout {
//     let layout: layout = {
//         nodes: [],
//         edges: []
//     };

//     if(nodes.length < 1) {
//         return layout;
//     }

//     let edges = getEdges(nodes);

//     // Generate dot layout
//     let graphDOT = generateGraphDOT(nodes, edges)
//     let graphvizData = JSON.parse(Viz(graphDOT, {
//         format: 'json',
//         // engine: 'neato',
//         engine: 'dot'
//     }));
    
//     layout.nodes = graphvizData.objects.map(obj => {    
//         let pos = obj.pos.split(',').map(Number);
//         let id = parseInt(obj.name)

//         return {
//             layout: {
//                 cx: pos[0],
//                 cy: pos[1],
//                 rx: obj._draw_[1].rect[2],
//                 ry: obj._draw_[1].rect[3],
//             },
//             id,
//         };
//     });

//     layout.edges = edges.length > 1 ? graphvizData.edges.map((edge, i) => {
//         let points = edge._draw_[1].points.map(toSvgPointSpace);
//         let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace);
//         let id = parseInt(edge.id);

//         return {
//             layout: {
//                 points,
//                 arrowPts,
//             },
//             id,
//         };
//     }) : [];

//     // layout.graphDOT = graphDOT;
    
//     return layout;
// }

// import type {
//     graphState,
// } from './reducers';

// const arrayEqual = (a, b) => {
//     return a.filter(e => !b.includes(e)).length > 0;
// }

// export function graphLogic(state: graphState) {
//     state.nodes = state.nodes.map(node => {
//         return {
//             ...node,
//             selection: getNodeSelection(node),
//         }
//     })
//     let spanningTree = generateSpanningTree(state)
//     if(arrayEqual(state.spanningTree, spanningTree)) {
//         return {};
//     }

//     let layout = buildLayout(state, spanningTree.map(id => getNodeById(state.nodes, id)));

//     return {
//         spanningTree,
//         layout,
//     }
// }



// // yield instructions
// // produce effects

// // instruction: load graph
// // effects:
// //  - process ins, outs and selection
// //  - generate spanning tree
// //  - generate a layout
// //  - set this data.





//     // switch(action.type) {
//     //     case "RENDER_GRAPH":
//     //         return {
//     //             nodes: {},
//     //             edges: {}
//     //         }
//     //     case "RENDER_COMPLETE":

//     //     default:
//     //         return state;
//     // }

//     // pass to a web worker for layout
//     // action: render graph
//     // action: begin graph render
//     // 


//     // generate layout for nodes, edges
//         // nodes = set ins and outs from edges, selection is full selection object

//         // generate spanning tree from (nodes, edges) data

//         // generate graphdot from spanning tree
//         // return if graphdot the same

//         // else generate layout from graphdot


// // type refreshMsg = {
// //     state: graphState
// // };

// self.addEventListener('message', (ev) => {
//     let msg: graphState = ev.data;
//     self.postMessage(graphLogic(msg));
// })

// // $FlowFixMe
// if (module.hot) {
//     module.hot.accept()
// }