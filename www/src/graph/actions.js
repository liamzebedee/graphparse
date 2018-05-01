import graphJSON from '../../graph.json';

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

export function searchNodes(q) {
    return {
        type: "SEARCH_NODES",
        q
    }
}

export function loadInitialFileForTesting() {
    return (dispatch, getState) => {
        dispatch(searchNodes("Server"))

        let top = getState().graph.search.matches[0];
        if(top == null) { throw new Error() }
        dispatch(selectNodeFromSearch(top.id))
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

export function uiChangeView(currentView) {
    return {
        type: "UI_CHANGE_VIEW",
        view: ("show types" == currentView ? "show relationships" : "show types")
    }
}