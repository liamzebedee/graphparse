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

export function selectNodeFromSearch(id) {
    return {
        type: "SELECT_NODE_FROM_SEARCH",
        id
    }
}

export function selectNodeByLabel(label) {
    return {
        type: "SELECT_NODE_BY_LABEL",
        label
    }
}