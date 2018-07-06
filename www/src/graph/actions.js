// @flow
import underscore from 'underscore';
import _ from 'lodash';

import { axios } from '../util';
import { getNodeById } from './selectors';
import { Base64 } from 'js-base64';

import { compose } from 'redux'

import type {
    nodeid,
    nodeVariant,
    graphState,
    node,
    edge
} from 'graphparse';

type relationships = 'ins' | 'outs';

// export const hoverNode = compose(regenerateGraph, (id) => {
//     return {
//         type: "HOVER_NODE",
//         id,
//     }
// })

export const hoverNode = (id: nodeid) => {
    return {
        type: "HOVER_NODE",
        id,
    }
}

export const blurHover = () => {
    return {
        type: "BLUR_HOVER",
    }
}

// $FlowFixMe
export const clickNode = compose(regenerateGraph, (id: nodeid, shiftKey: boolean) => {
    return {
        type: "CLICK_NODE",
        id,
        shiftKey,
    };
})

export const clearSelection = compose(regenerateGraph, function(x: any) {
    return {
        type: "BLUR_SELECTED_NODE"
    }
})

// $FlowFixMe
export const toggleFilter = compose(regenerateGraph, function(id: nodeid, relationships: relationships, variant: nodeVariant) {
    return {
        type: "TOGGLE_FILTER",
        id,
        relationships,
        variant,
    }
})

export const searchNodes = function(q: string) {
    return {
        type: "SEARCH_NODES",
        q
    }
}


export const selectNodeFromSearch = compose(regenerateGraph, function(node: node) {
    return {
        type: "SELECT_NODE_FROM_SEARCH",
        node,
    }
})

export const selectNodeByLabel = compose(regenerateGraph, function(label: string) {
    return {
        type: "SELECT_NODE_BY_LABEL",
        label
    }
})

// $FlowFixMe
export const changeDepth = compose(regenerateGraph, function(id: nodeid, relationships: relationships, depth: number) {
    return {
        type: "CHANGE_DEPTH",
        id,
        relationships,
        depth,
    }
})

// $FlowFixMe
type graphResponse = {|
    nodes: node[],
    edges: edge[],
    rootNode: nodeid
|};

// $FlowFixMe
export const load = compose(regenerateGraph, function(graphID: string, firstLoad: boolean) {
    return (dispatch, getState) => {
        axios.get(`/graph/public/${Base64.encode(graphID)}`)
        .then(res => {
            let graph: graphResponse = res.data;
            dispatch(loadGraph(
                graph.nodes,
                graph.edges,
            ));

            let rootNode = getNodeById(graph.nodes, graph.rootNode);
            if(firstLoad) dispatch(selectNodeFromSearch(rootNode));
        })
        .catch(err => {
            if(err.response.status === 404) {
                dispatch({
                    type: "ERROR",
                    message: "Codebase was not found"
                })
            } else throw err;
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

function regenerateGraph(action: any) {
    return (dispatch, getState) => {
        dispatch({
            type: "GENERATING"
        });

        dispatch(action);

        let state: graphState = _.cloneDeep(getState().graph);

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



