// @flow
import { combineReducers } from 'redux';

import matchSorter from 'match-sorter';
import reduceReducers from 'reduce-reducers';
import underscore from 'underscore';
import _ from 'lodash';

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
    nodeSel,
    graphState
} from 'graphparse';

import { getNodeSelection, mergeNodeSelection } from './selectors';
import { ClickActions } from './types';




const initialState: graphState = {
    grabbing: false,

    currentNode: null,
    hoveredNode: null,

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


function updateNode(nodes: node[], id: nodeid, cb: (x: node) => any) : node[] {
    return nodes.map(node => {
        if(node.id !== id) return node;
        else return _.merge(
            _.cloneDeep(node),
            cb(node)
        );
    })
}

const processNode = (nodes: node[], edges: edge[]) => {
    const nodeExists = (id) => underscore.findWhere(nodes, { id, }) != null;

    let edgesThatExist = edges.filter(edge => {
        return nodeExists(edge.source) && nodeExists(edge.target);
    })

    return (node: node) => {
        return {
            ...node,
            selection: {
                ins: {},
                outs: {},
            },
            outs: underscore.where(edgesThatExist, { source: node.id }),
            ins: underscore.where(edgesThatExist, { target: node.id }),
        }
    };
}

function graph(state: graphState = initialState, action: any) {
    switch(action.type) {
        case "ERROR":
            return {
                ...state,
                error: action.message
            }

        case "LOAD_GRAPH":
            let { nodes, edges } = action;
            return {
                ...state,
                nodes: nodes.map(processNode(nodes, edges)),
                edges,
            }
        
        case "SELECT_NODE_FROM_SEARCH":
            return {
                ...state,
                currentNode: action.node.id,
                search: {
                    q: action.node.label
                }
            }

        case "SELECT_NODE_BY_LABEL": {
            let matches = matchSorter(state.nodes, action.label, { keys: ['label'] })
            let match = matches[0];

            return {
                ...state,
                search: {
                    q: action.label,
                    matches,
                },
                currentNode: match.id,
            }
        }
        
        case "CLICK_NODE":
        case "SELECT_CLICK_ACTION":
        case "CHANGE_DEPTH":
        case "TOGGLE_FILTER":
        case "HOVER_NODE":
        case "BLUR_HOVER":
            return ui(state, action);
        
        case "SEARCH_NODES":
            let matches = matchSorter(state.nodes, action.q, { keys: ['label'] })
            return {
                ...state,
                search: {
                    q: action.q,
                    matches,
                }
            }
        
        case "GENERATING":
            return {
                ...state,
                generating: true
            }
        
        case "GENERATE_COMPLETE":
            return {
                ...state,
                layout: action.payload.layout,
                spanningTree: action.payload.spanningTree,
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
        
        case "HOVER_NODE":
            return {
                ...state,
                hoveredNode: action.id
            }
        case "BLUR_HOVER":
            return {
                ...state,
                hoveredNode: null
            }

        case "CLICK_NODE":
            return clickNode(state, action)
        
        case "CHANGE_DEPTH":
            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);

                    return {
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

function clickNode(state, action) {
    switch(state.clickAction) {
        case ClickActions.select.name:
            return {
                ...state,
                nodes: state.nodes.map(node => {
                    return {
                        ...node,
                        selected: node.id === action.id
                    }
                })
            }
        
        case ClickActions.relationships.name:
            let shiftKeyDown = false;
            
            return {
                ...state,
                nodes: updateNode(state.nodes, action.id, (node) => {
                    let sel = getNodeSelection(node);
                    let newSel = {
                        outs: {},
                        ins: {}
                    };

                    if(!action.shiftKey) {
                        newSel.outs.shown = !sel.outs.shown;
                    } else {
                        newSel.ins.shown = !sel.ins.shown;
                    }

                    return {
                        ...node,
                        selection: mergeNodeSelection(node.selection, newSel)
                    }
                })
            }
        
        case ClickActions.visibility.name:
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
        
        default:
            return state;
    }
}

export default graph;