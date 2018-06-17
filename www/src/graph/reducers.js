// @flow
import { combineReducers } from 'redux';
import {Enum} from 'enumify';
import matchSorter from 'match-sorter';

import { searchNodes } from './actions';
import {
    toggleInArray
} from '../util'

// $FlowFixMe
import { BASE } from '/Users/liamz/Documents/open-source/proxy-object-defaults';

import type {
    nodeid,
    node,
    nodeLayout,
    edge,
    edgeLayout,
    nodeSel
} from 'graphparse';

import { getNodeSelection } from './selectors';


export class ClickActions extends Enum {}
ClickActions.initEnum([
    'select',
    'relationships',
    'visibility'
]);

type graphState = {|
    grabbing: boolean,

    currentNode: ?nodeid,

    search: {|
        q: string,
        matches: node[],
        state: string,
    |},

    nodes: node[],
    edges: edge[],
    clickAction: ClickActions
|};

const initialState: graphState = {
    grabbing: false,

    currentNode: null,

    search: {
        q: "",
        matches: [],
        state: "blurred"
    },

    nodes: [],
    edges: [],

    clickAction: ClickActions.select
}


function updateNode(nodes, id, cb) {
    return nodes.map(node => {
        if(node.id !== id) return node;
        else return cb(node);
    })
}

function graph(state: graphState = initialState, action: any) {
    switch(action.type) {
        case "LOAD_GRAPH":
            return {
                ...state,
                nodes: action.nodes.map(node => {
                    return {
                        ...node,
                        selection: {
                            ins: {},
                            outs: {}
                        }
                    }
                }),
                edges: action.edges,
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
        
        case "CLICK_NODE":
        case "SELECT_CLICK_ACTION":
            return ui(state, action);
        
        default:
            return state;
    }
}





// function updateNodeFilters(state = initialFilter, action) {
//     switch(action.type) {
//         case "TOGGLE_NODE_TYPE_FILTER":
//             return {
//                 ...state,
//                 shownNodeTypes: toggleInArray(state.shownNodeTypes, action.nodeTypeFilterIdx)
//             }

//         case "CLICK_NODE":
//             return {
//                 ...state,
//                 shown: !state.shown
//             }

//         default:
//             return state;
//     }
// }



function ui(state, action) {
    switch(action.type) {
        case "SELECT_CLICK_ACTION":
            return {
                ...state,
                clickAction: action.action
            }

        case "CLICK_NODE":
            return clickNode(state, action)
        
        // case "TOGGLE_FILTER":
        //     return {
        //         ...state,
        //         nodes: updateNodes(state.nodes, {
        //             ...action,
        //             id: state.selectedNode
        //         })
        //     }
        default: 
            return state;
    }
}

function clickNode(state, action) {
    switch(state.clickAction) {
        case ClickActions.select:
            return {
                ...state,
                nodes: state.nodes.map(node => {
                    return {
                        ...node,
                        selected: node.id === action.id
                    }
                })
            }
        
        case ClickActions.relationships:
            let shiftKeyDown = false;

            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);

                    if(shiftKeyDown) {
                        sel.outs.shown = !sel.outs.shown;
                    } else {
                        sel.ins.shown = !sel.ins.shown;
                    }
                    
                    return {
                        ...node,
                        selection: sel[BASE]
                    }
                })
            }
        
        case ClickActions.visibility:
            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);
                    sel.shown = !sel.shown;

                    return {
                        ...node,
                        selection: sel[BASE]
                    }
                })
            }
        
        default:
            return state;
    }
}

export default graph;