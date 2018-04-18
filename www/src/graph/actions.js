import graphJSON from '../../graph.json';

export function hoverNode(id) {
    return {
        type: "HOVER_NODE",
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