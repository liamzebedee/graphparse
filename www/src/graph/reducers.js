// @flow
import { combineReducers } from 'redux';
import {Enum} from 'enumify';
import matchSorter from 'match-sorter';
import reduceReducers from 'reduce-reducers';
import _ from 'underscore';

import { searchNodes } from './actions';
import {
    toggleInArray
} from '../util'

import type {
    nodeid,
    node,
    nodeLayout,
    edge,
    edgeLayout,
    nodeSel
} from 'graphparse';

import { getNodeSelection, mergeNodeSelection } from './selectors';

import { graphLogic } from './graph-logic.worker';


export class ClickActions extends Enum {}
ClickActions.initEnum([
    'select',
    'relationships',
    'visibility'
]);

export type graphState = {|
    grabbing: boolean,

    currentNode: ?nodeid,

    search: {|
        q: string,
        matches: node[],
        state: string,
    |},

    nodes: node[],
    edges: edge[],
    clickAction: string,

    generating: boolean,
    spanningTree: nodeid[],
    layout: {|
        nodes: nodeLayout[],
        edges: edgeLayout[],
    |},
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

    clickAction: ClickActions.select.name,

    generating: false,
    spanningTree: [],
    layout: {
        nodes: [],
        edges: []
    }
}


function updateNode(nodes: node[], id: nodeid, cb: (x: node) => node) : node[] {
    return nodes.map(node => {
        if(node.id !== id) return node;
        else return cb(node);
    })
}

const processNode = (nodes, edges) => {
    const nodeExists = (id) => _.findWhere(nodes, { id, }) != null;

    let edgesThatExist = edges.filter(edge => {
        return nodeExists(edge.source) && nodeExists(edge.target);
    })

    return (node) => {
        return {
            // todo reordered this, maybe it's the src of a bug?
            ...node,
            selection: {
                ins: {},
                outs: {},
            },
            outs: _.where(edgesThatExist, { source: node.id }),
            ins: _.where(edgesThatExist, { target: node.id }),
        }
    };
}

function graph(state: graphState = initialState, action: any) {
    switch(action.type) {
        case "LOAD_GRAPH":
            let { nodes, edges } = action;
            return {
                ...state,
                nodes: nodes.map(processNode(nodes, edges)),
                edges,
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
        
        case "CLICK_NODE":
        case "SELECT_CLICK_ACTION":
        case "CHANGE_DEPTH":
            return ui(state, action);
        
        case "SEARCH_NODES":
            let matches = matchSorter(state.nodes, action.q, { keys: ['label'] })
            return Object.assign({}, state, {
                search: {
                    q: action.q,
                    matches,
                }
            })
        
        case "GENERATING":
            return {
                ...state,
                generating: true
            }
        
        case "GENERATE_COMPLETE":
            if(action.payload === {}) return {
                generating: false
            };
            else return {
                ...state,
                ...action.payload,
                generating: false,
            }
        
        default:
            return state;
    }
}

function ui(state, action) {
    switch(action.type) {
        case "SELECT_CLICK_ACTION":
            return {
                ...state,
                clickAction: action.action.name
            }
        
        case "SELECT_NODE":
            return {
                ...state,
                // $FlowFixMe
                nodes: state.nodes.map((node) => {
                    return {
                        ...node,
                        selected: node.id === action.id
                    };
                })
            }
        
        case "TOGGLE_NODE_RELATIONSHIPS":
            let shiftKeyDown = false;

            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);
                    let newSel = {
                        outs: {},
                        ins: {}
                    };

                    if(shiftKeyDown) {
                        sel.outs.shown = !sel.outs.shown;
                    } else {
                        sel.ins.shown = !sel.ins.shown;
                    }

                    return {
                        ...node,
                        selection: mergeNodeSelection(node.selection, newSel)
                    }
                })
            }
        
        case "TOGGLE_NODE_VISIBILITY":
            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);

                    return {
                        ...node,
                        selection: mergeNodeSelection(node.selection, {
                            shown: !sel.shown
                        })
                    }
                })
            }
        
        case "CHANGE_DEPTH":
            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);

                    return {
                        ...node,
                        selection: mergeNodeSelection(node.selection, {
                            [action.relationships]: {
                                maxDepth: action.depth
                            }
                        })
                    }
                })
            }
        
        case "TOGGLE_FILTER":
            break;
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

export default graph;