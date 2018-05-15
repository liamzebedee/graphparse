import graphJSON from '../../graph.json';

import _ from 'underscore';

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

export function loadInitialFileForTesting() {
    return (dispatch, getState) => {
        // dispatch(searchNodes("Server"))
        dispatch(searchNodes("parse.go"))

        let topMatch = getState().graph.search.matches[0];
        if(topMatch == null) { return }
        dispatch(selectNodeFromSearch(topMatch.id))
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


export const VIEWS = [
    "show relationships",
    "show types"
];

export function changeView(uiView) {
    if(!_.contains(VIEWS, uiView)) throw new Error();
    return {
        type: "UI_CHANGE_VIEW",
        uiView
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