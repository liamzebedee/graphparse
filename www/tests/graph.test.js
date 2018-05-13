
var assert = require('assert');

import {
    GraphLogic,
    constructAdjList,
    mergeByKey
} from '../src/graph/graph-logic';
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
    let selection = [];
    let maxDepth = 2;

    // logic.refresh({ nodes, edges, currentNode, selection, maxDepth });
    logic.nodes = nodes;
    logic.edges = edges;
    logic.currentNode = currentNode;
    logic.selection = selection;
    logic.maxDepth = maxDepth;

    return {
        logic,
        nodes,
        edges,
        currentNode,
        selection,
        maxDepth
    }
}

function bareData() {
    let adjList = {
        0: [1,2,3],
        1: [7, 0],
        2: [4, 5],
        3: [],
        4: [6],
        5: [10],
        6: [],
        7: [8,9,10],
        8: [],
        9: [],
        10: [],
    }

    let nodes = [0,1,2,3,4,5,6,7,8,9,10].map((id, i) => ({
        id,
        variant: -1,
        label: `node #${i}`
    }));
    let edgeI = 0;
    let edges = Object.keys(adjList).map(source => {
        return adjList[source].map(target => {
            return { 
                source: parseInt(source), 
                target: parseInt(target),
                variant: -1,
                id: edgeI++
            }
        })
    }).reduce((prev, curr) => prev.concat(curr))

    return { adjList, nodes, edges };
}

function assertEmptyArray(val) {
    assert(Array.isArray(val), true);
    assert.equal(val.length, 0);
}


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

describe('constructs spanning tree correctly when currentNode is set', () => {
    let { logic } = setup();
    let { nodes, edges } = bareData();

    logic.currentNode = 2;
    logic.maxDepth = 2;
    let tree = logic.getSpanningTree(nodes, edges);
    assert.deepStrictEqual(tree, [ 2, 4, 5, 6, 10 ])

    logic.maxDepth = 1;
    tree = logic.getSpanningTree(nodes, edges);
    assert.deepStrictEqual(tree, [ 2, 4, 5 ], "respects max depth param")
})

describe("constructAdjList works", () => {
    let { logic } = setup();
    let { nodes, edges, adjList } = bareData();
    assert.deepEqual(constructAdjList(nodes, edges), adjList);
})

describe("preFilterNodesAndEdges sets the nodes and edges", () => {
    let { logic } = setup();
    let { nodes, edges } = bareData();

    logic.nodes = nodes;
    logic.edges = edges;
    logic.currentNode = 2;
    logic.nodesLookup = nodes;
    logic.preFilterNodesAndEdges();
    assert.notEqual(logic.nodes.length, 0)
})

describe('edges is always synced with this.nodes', () => {
    let { logic } = setup();
    assert.equal(logic.currentNode, null)

    let { nodes, edges, adjList } = bareData();
    logic._nodes = nodes;
    logic.edges = edges;
    assert.deepStrictEqual(logic.edges, edges)

    logic.nodes = [{ id: 2 }];
    let newEdges = adjList[2].map((target, i) => {
        return { source: 2, target, variant: -1, id: 5+i }
    })
    assert.deepStrictEqual(logic.edges, newEdges);
})

describe('layout handles case where no nodes or edges', () => {
    let { logic } = setup();
    logic._nodes = [];
    logic.edges = [];

    logic.generateLayout();
    assertEmptyArray(logic.nodesLayout)
    assertEmptyArray(logic.edgesLayout)
})

describe('layout is generated and merged', () => {
    let { logic, nodes, edges, currentNode, selection, maxDepth } = setup();

    currentNode = _.findWhere(nodes, { label: "main.go" }).id;
    logic.maxDepth = 100;
    logic.refresh(nodes, edges, currentNode, selection, maxDepth);    

    let layout = logic.getLayout();
    
    logic.nodesLayout.map(x => {
        assert.notEqual(x.id, null);
        assert.notEqual(x.id, NaN);
    }) 
    logic.edgesLayout.map(x => {
        assert.notEqual(x.id, null);
        assert.notStrictEqual(x.id, NaN);
    })

    assert.notDeepEqual(logic.nodesLayout, []);
    assert.notDeepEqual(logic.edgesLayout, []);

    assert.equal(layout.nodes.length, 8);
    assert.equal(layout.edges.length, 8);

    layout.nodes.map(node => {
        // console.log(node)
        assert.notEqual(node.variant, null)
        assert.notEqual(node.layout, null)
        assert.notEqual(node.layout, {})
    }) 
    layout.edges.map(edge => { 
        // console.log(edge) 
        assert.notEqual(edge.variant, null)
        assert.notEqual(edge.layout, null)
        assert.notEqual(edge.layout, {})
    })
})