// import graphJSON from '../../graph.json';
import _ from 'underscore';
import { axios } from '../util';
import { Base64 } from 'js-base64';

import { compose } from 'redux'
import { ClickActions } from './reducers';

export const hoverNode = compose(regenerateGraph, (id) => {
    return {
        type: "HOVER_NODE",
        id,
    }
})

export const clickNode = (id) => {
    return (dispatch, getState) => {
        switch(getState().graph.clickAction) {
            case ClickActions.select.name:
                return dispatch(selectNode(id))
            case ClickActions.relationships.name:
                return dispatch(toggleNodeRelationships(id))
            case ClickActions.visibility.name:
                return dispatch(toggleNodeVisibility(id))
        }
    }
}

export const selectNode = compose(regenerateGraph, (id) => {
    return {
        type: "SELECT_NODE",
        id,
    };
})

export const toggleNodeRelationships = compose(regenerateGraph, (id) => {
    return {
        type: "TOGGLE_NODE_RELATIONSHIPS",
        id,
    };
})

export const toggleNodeVisibility = compose(regenerateGraph, (id) => {
    return {
        type: "TOGGLE_NODE_VISIBILITY",
        id,
    };
})

export const clearSelection = compose(regenerateGraph, function() {
    return {
        type: "BLUR_SELECTED_NODE"
    }
})

export const toggleFilter = compose(regenerateGraph, function(id: nodeid, relationships: relationships, variant: nodeVariant) {
    return {
        type: "TOGGLE_FILTER",
        id,
        relationships,
        variant,
    }
})

export const searchNodes = function(q) {
    return {
        type: "SEARCH_NODES",
        q
    }
}

export const selectNodeFromSearch = compose(regenerateGraph, function(id) {
    return {
        type: "SELECT_NODE_FROM_SEARCH",
        id,
    }
})

export const selectNodeByLabel = compose(regenerateGraph, function(label) {
    return {
        type: "SELECT_NODE_BY_LABEL",
        label
    }
})


export const changeDepth = compose(regenerateGraph, function(id: nodeid, relationships: relationships, depth: number) {
    return {
        type: "CHANGE_DEPTH",
        id,
        relationships,
        depth,
    }
})


export const load = compose(regenerateGraph, function(graphID, firstLoad) {
    return (dispatch, getState) => {
        axios.get(`/graph/public/${Base64.encode(graphID)}`)
        .then(res => {
            let graph = res.data;
            dispatch(loadGraph(
                graph.nodes,
                graph.edges,
            ))
            if(firstLoad) dispatch(selectNodeFromSearch(graph.rootNode));
        })
        .catch(err => {
            throw err;
        })
    }
})



function generateSpanningTree() {
    //  - generate a layout
    //  - set this data.
}

function loadGraph(nodes, edges) {
    return {
        type: "LOAD_GRAPH",
        nodes,
        edges,
    }
}







// Internal
// --------

import { graphLogic } from './graph-logic';

function regenerateGraph(action) {
    return (dispatch, getState) => {
        dispatch({
            type: "GENERATING"
        });

        dispatch(action);

        let state: graphState = getState().graph;

        new Promise((resolve, reject) => {
            // worker.addEventListener("message", (ev) => {
            //     let msg = ev.data;
            //     if(msg === {}) resolve();
            //     else resolve(msg)
            // }, { once: true });
            
            // worker.postMessage(state);
            resolve(graphLogic(state))

        }).then((payload) => {
            dispatch({
                type: "GENERATE_COMPLETE",
                payload,
            })
        })
    }
}



// import Worker from './graph-logic.worker';
// const worker = new Worker();
import type { graphState } from './reducers';



