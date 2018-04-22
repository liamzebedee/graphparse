import graphJSON from '../../graph.json';
import { combineReducers } from 'redux';
import { searchNodes } from './actions';


import matchSorter from 'match-sorter';

const initialState = {
    interested: [],
    clickedNode: null,
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
        case "CLICK_NODE":
            return Object.assign({}, state, {
                clickedNode: action.id
            })
        case "HOVER_NODE":
            let interested = getSubPaths(state.adjList, action.id)
            return Object.assign({}, state, {
                interested,
            })
        case "SEARCH_NODES":
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

        case "SELECT_NODE_BY_LABEL": {
            let matches = matchSorter(state.nodes, action.label, { keys: ['label'] })
            let match = matches[0];

            return Object.assign({}, state, {
                search: {
                    q: action.label,
                    matches,
                },
                interested: getSubPaths(state.adjList, match.id),
            })
        }
        default:
            return state
    }
}


export default combineReducers({
    graph,
})