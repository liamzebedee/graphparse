// import graphJSON from '../../graph.json';
import _ from 'underscore';
import { axios } from '../util';
import { Base64 } from 'js-base64';

export function hoverNode(id) {
    return {
        type: "HOVER_NODE",
        id,
    }
}

export function clickNode(id) {
    return {
        type: "CLICK_NODE",
        id,
    }
}

export function clearSelection() {
    return {
        type: "BLUR_SELECTED_NODE"
    }
}

export function toggleNodeTypeFilter(nodeTypeFilterIdx) {
    return {
        type: "TOGGLE_NODE_TYPE_FILTER",
        nodeTypeFilterIdx,
    }
}

export function searchNodes(q) {
    return {
        type: "SEARCH_NODES",
        q
    }
}

export function selectNodeFromSearch(id) {
    return {
        type: "SELECT_NODE_FROM_SEARCH",
        id,
    }
}

export function selectNodeByLabel(label) {
    return {
        type: "SELECT_NODE_BY_LABEL",
        label
    }
}

export function setGrabbing(grabbing) {
    return {
        type: "GRABBING_CHANGE",
        grabbing,
    }
}


export function changeDepth(depth) {
    return {
        type: "CHANGE_DEPTH",
        depth
    }
}

export function toggleShowDefinitions() {
    return {
        type: "toggleShowDefinitions"
    }
}

export function searchFocusChange(state) {
    return {
        type: "searchFocusChange",
        state
    }
}

export function load(graphID) {
    return (dispatch, getState) => {
        axios.get(`/graph/public/${Base64.encode(graphID)}`)
        .then(res => {
            let graph = res.data;

            dispatch({
                type: "LOAD_GRAPH",
                nodes: graph.nodes,
                edges: graph.edges,
            });
            dispatch(selectNodeFromSearch(graph.rootNode));
        })
        .catch(err => {
            throw err;
        })
    }
}