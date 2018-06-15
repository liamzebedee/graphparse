import { combineReducers } from 'redux';
import { searchNodes } from './actions';
import {
    toggleInArray
} from '../util'

import matchSorter from 'match-sorter';

import graphJSON from '../../graph.json';

const initialState = {
    firstLoad: true,
    grabbing: false,

    currentNode: null,
    selectedNode: null,
    maxDepth: 1,

    clickedNode: null,
    search: {
        q: "",
        matches: [],
        state: "blurred"
    },

    showDefinitions: false,

    nodes: [],
    edges: [],
    nodeTypes: graphJSON.nodeTypes,
}


function graph(state = initialState, action) {
    switch(action.type) {
        case "LOAD_GRAPH":
            return {
                ...state,
                nodes: updateNodes(action.nodes, {
                    id: 'all' 
                }),
                edges: action.edges
            }

        case "TOGGLE_NODE_TYPE_FILTER":
            return {
                ...state,
                nodes: updateNodes(state.nodes, {
                    ...action,
                    id: state.selectedNode
                })
            }
        case "CLICK_NODE":
            return {
                ...state,
                selectedNode: action.id,
                nodes: updateNodes(state.nodes, action),
            };
        
        case "SELECT_NODE_FROM_SEARCH":
            return Object.assign({}, state, {
                currentNode: action.id
            })

        case "SELECT_NODE_BY_LABEL": {
            let matches = matchSorter(state.nodes, action.label, { keys: ['label'] })
            let match = matches[0];

            return Object.assign({}, state, {
                search: {
                    q: action.label,
                    matches,
                },
                currentNode: match.id,
            })
        }

        case "SEARCH_NODES":
            let matches = matchSorter(state.nodes, action.q, { keys: ['label'] })
            return Object.assign({}, state, {
                search: {
                    q: action.q,
                    matches,
                }
            })

        case "CHANGE_DEPTH":
            return {
                ...state,
                maxDepth: action.depth,
            }

        case "GRABBING_CHANGE":
            return {
                ...state,
                grabbing: action.grabbing,
            }
            
        case "toggleShowDefinitions":
            return {
                ...state,
                showDefinitions: !state.showDefinitions
            }

        case "FIRST_LOAD":
            return {
                ...state,
                firstLoad: false,
            }
            
        default:
            return state
    }
}


const initialFilter = {
    shownNodeTypes: graphJSON.nodeTypes.map((a,i) => i),
    selected: false,
    shown: false,
};

function updateNodes(nodes = [], action) {
    return nodes.map(node => {
        if(action.id == 'all' || node.id === action.id) {
            return {
                ...node,
                filters: updateNodeFilters(node.filters, action)
            }
        }
        return node;
    })
}

function updateNodeFilters(state = initialFilter, action) {
    switch(action.type) {
        case "TOGGLE_NODE_TYPE_FILTER":
            return {
                ...state,
                shownNodeTypes: toggleInArray(state.shownNodeTypes, action.nodeTypeFilterIdx)
            }

        case "CLICK_NODE":
            return {
                ...state,
                shown: !state.shown
            }

        default:
            return state;
    }
}


export default graph;