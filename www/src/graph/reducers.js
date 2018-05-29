import graphJSON from '../../graph.json';
import { combineReducers } from 'redux';
import { searchNodes } from './actions';
import {
    toggleInArray
} from '../util'

import matchSorter from 'match-sorter';


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
    nodes: graphJSON.nodes,
    edges: graphJSON.edges,
    nodeLookup: graphJSON.nodeLookup,
    adjList: graphJSON.adjList,
    nodeTypes: graphJSON.nodeTypes,

    uiView: "show relationships",
    showDefinitions: false,
}


function graph(state = initialState, action) {
    switch(action.type) {
        case "CLICK_NODE":
            return {
                ...state,
                selectedNode: action.id,
                nodes: state.nodes.map(node => {
                    if(node.id == action.id) {
                        return {
                            ...node,
                            shown: !node.shown
                        };
                    } else return node;
                })
            };
        
        case "BLUR_SELECTED_NODE":
            return {
                ...state,
                selectedNode: null
            }
        
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
        
        case "UI_CHANGE_VIEW":
            return {
                ...state,
                uiView: action.uiView,
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


export default graph;