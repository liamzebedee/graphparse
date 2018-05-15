// @flow
var assert = require('assert');
import {describe} from 'mocha'

import {
    GraphLogic,
    mergeByKey,
    edgeRelationId
} from '../src/graph/graph-logic';
import {
    removeDuplicates
} from '../src/util';
const graphJSON = require('../../test/graph.json');
import lodash from 'lodash';
import _ from 'underscore';
var expect = require('chai').expect




function setup() {
    let logic = new GraphLogic();

    let {
        nodes, edges
    } = graphJSON;
    let currentNode = null;
    let selection = {};
    let maxDepth = 2;

    return {
        logic,
        nodes,
        edges,
        currentNode,
        selection,
        maxDepth
    }
}

// function bareData() {
//     let adjList = {
//         0: [1,2,3],
//         1: [7, 0],
//         2: [4, 5],
//         3: [],
//         4: [6],
//         5: [10],
//         6: [],
//         7: [8,9,10],
//         8: [],
//         9: [],
//         10: [],
//     }

//     let nodes = [0,1,2,3,4,5,6,7,8,9,10].map((id, i) => ({
//         id,
//         variant: -1,
//         label: `node #${i}`
//     }));
//     let edgeI = 0;
//     let edges = Object.keys(adjList).map(source => {
//         return adjList[source].map(target => {
//             return { 
//                 source: parseInt(source), 
//                 target: parseInt(target),
//                 variant: -1,
//                 id: edgeI++
//             }
//         })
//     }).reduce((prev, curr) => prev.concat(curr))

//     return { adjList, nodes, edges };
// }

function assertEmptyArray(val) {
    assert(Array.isArray(val));
    assert.equal(val.length, 0);
}

describe("mergeByKey should throw on merging key that doesn't exist in base collection", () => {
    let dataColl = [
        {
            id: 0,
            name: "Test 0"
        },
        {
            id: 1,
            name: "Test 1"
        },
        {
            id: 2,
            name: "Test 2"
        }
    ]

    let layoutColl = [
        {
            id: 0,
            layout: {
                x: 1, y: 2
            }
        },
        {
            id: 1,
            layout: {
                x: 1, y: 2
            }
        },
        {
            id: 2000,
            layout: {
                x: 1, y: 2
            }
        }
    ]

    let merged = () => mergeByKey('id', dataColl, layoutColl)
    expect(merged).to.throw();
})

describe('constructs graph correctly when currentNode is null', () => {
    let { logic } = setup();
    assert.equal(logic.currentNode, null);
    assertEmptyArray(logic.getSpanningTree());
    
    logic.preFilterNodesAndEdges();
    assertEmptyArray(logic.nodes)
    assertEmptyArray(logic.edges)

    logic.generateLayout();
    assertEmptyArray(logic.nodesLayout)

    logic.postFilterNodesAndEdges();
    assertEmptyArray(logic.nodes)
    assertEmptyArray(logic.edges)
})

// describe('constructs spanning tree correctly when currentNode is set', () => {
//     let { logic } = setup();
//     let { nodes, edges } = bareData();

//     logic.currentNode = 2;
//     logic.maxDepth = 2;
//     let tree = logic.getSpanningTree(nodes, edges);
//     assert.deepStrictEqual(tree, [ 2, 4, 5, 6, 10 ])

//     logic.maxDepth = 1;
//     tree = logic.getSpanningTree(nodes, edges);
//     assert.deepStrictEqual(tree, [ 2, 4, 5 ], "respects max depth param")
// })

// describe("preFilterNodesAndEdges sets the nodes and edges", () => {
//     let { logic } = setup();
//     let { nodes, edges } = bareData();

//     logic.nodes = nodes;
//     logic.edges = edges;
//     logic.currentNode = 2;
//     logic.nodesLookup = nodes;
//     logic.preFilterNodesAndEdges();
//     assert.notEqual(logic.nodes.length, 0)
// })

// describe('edges is always synced with this.nodes', () => {
//     let { logic } = setup();
//     assert.equal(logic.currentNode, null)
//     assert.deepStrictEqual(logic.edges, edges)
//     assert.deepStrictEqual(logic.edges, newEdges);
// })

describe('layout handles case where no nodes or edges', () => {
    let { logic } = setup();

    logic.refresh([], [], null, {});

    assertEmptyArray(logic.nodesLayout)
    assertEmptyArray(logic.edgesLayout)
})

describe('layout is generated and merged', () => {
    let { logic, nodes, edges, currentNode, selection, maxDepth } = setup();

    currentNode = _.findWhere(nodes, { label: "main.go" }).id;
    assert.notEqual(currentNode, null, "unit testing check");

    logic.maxDepth = 100;
    logic.refresh(nodes, edges, currentNode, selection, maxDepth);
    assert.equal(logic.nodes.length, nodes.length)

    let tree = logic.getSpanningTree();
    assert.notEqual(tree.length, 0)

    assert.notEqual(logic.shownNodes.length, 0);
    assert.notEqual(logic.shownEdges.length, 0);

    assert.notDeepEqual(logic.nodesLayout, []);
    assert.notDeepEqual(logic.edgesLayout, []);

    // Layout should only contain shown nodes and edges.
    assert.equal(logic.nodesLayout.length, logic.shownNodes.length)
    assert.deepStrictEqual(
        logic.edgesLayout.map(e => e.id).sort(), 
        logic.shownEdges.filter(removeDuplicates(edgeRelationId)).map(e => e.id).sort(),
    )

    logic.nodesLayout.map(x => {
        assert.notEqual(x.id, null);
        assert.notEqual(x.id, NaN);
    }) 
    logic.edgesLayout.map(x => {
        assert.notEqual(x.id, null);
        assert.notStrictEqual(x.id, NaN);
    })

    let layout = logic.getLayout();
    assert.equal(layout.nodes.length, 7);
    assert.equal(layout.edges.length, 7);

    layout.nodes.map(node => {
        // console.log(node)
        assert.notEqual(node.variant, null)
        assert.notEqual(node.layout, null)
        assert.notEqual(node.layout, {})
    }) 
    layout.edges.map(edge => { 
        console.log(edge)
        assert.notEqual(edge.variant, null)
        assert.notEqual(edge.layout, null)
        assert.notEqual(edge.layout, {})
    })
})