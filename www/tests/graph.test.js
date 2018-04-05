import {
    getSubPaths
} from '../src/graph/actions';

var assert = require('assert');


describe('find correct path', () => {
    let adjList = {
        0: [1,2,3],
        1: [7, 0],
        2: [4, 5],
        4: [6],
        5: [10],
        7: [8,9,10],
    }

    assert.deepStrictEqual(getSubPaths(adjList, 2), [ 2, 4, 5, 6, 10 ])
})