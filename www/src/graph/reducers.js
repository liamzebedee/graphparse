import graphJSON from '../../graph.json';
import { combineReducers } from 'redux';
import { searchNodes } from './actions';


import matchSorter from 'match-sorter';

const initialState = {
    interested: [],
    search: {
        q: "",
        matches: []
    },
    nodes: graphJSON.nodes,
    edges: graphJSON.edges,
    nodeLookup: graphJSON.nodeLookup,
    adjList: graphJSON.adjList,
    nodeTypes: graphJSON.nodeTypes
}

export function getSubPaths(adjList, fromNodeId) {
    let currentNodes = [
        fromNodeId
    ];

    let interested = new Set();
    let visited = new Set();

    do {
        let next = [];

        currentNodes.map(node => {
            if(visited.has(node)) return;
            else visited.add(node)

            interested.add(node)

            let outs = adjList[""+node]
            if(outs) next = next.concat(outs)
        })

        currentNodes = next;

    } while(currentNodes.length)

    return Array.from(interested);
}

function graph(state = initialState, action) {
    switch(action.type) {
        case "HOVER_NODE":
            let interested = getSubPaths(state.adjList, action.id)
            return Object.assign({}, state, {
                interested,
            })
        case "SEARCH_NODES":
            // let searchList = state.nodes.map(node => node.label)
            let matches = matchSorter(state.nodes, action.q, { keys: ['label'] })
            return Object.assign({}, state, {
                search: {
                    q: action.q,
                    matches,
                }
            })
        case "SELECT_NODE_FROM_SEARCH":
            return Object.assign({}, state, {
                interested: getSubPaths(state.adjList, action.id)
            })
        default:
            return state
    }
}


export default combineReducers({
    graph,
})