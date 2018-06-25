// @flow
import _ from 'lodash';
import underscore from 'underscore';

import {
    getNodeTypes,
} from './colours';

import type {
    relationshipsSel,
    nodeSel,
    node,
    edge,
    nodeid,
} from 'graphparse';

export function getSelectedNode(state: any) {
    let selectedNodes = state.nodes
    .filter(node => node.selected)
    .map(node => {
        return {
            ...node,
            selection: getNodeSelection(node.selection)
        }
    });

    return selectedNodes.length ? selectedNodes[0] : null;
}

const DEFAULT_RELATIONSHIPS_SELECTION: relationshipsSel = {
    shownNodeTypes: getNodeTypes().map((a,i) => i),
    maxDepth: 0,
    showDefs: true,
    showUses: true,
    shown: false,
};

const DEFAULT_SELECTION: nodeSel = {
    shown: true,
    ins: DEFAULT_RELATIONSHIPS_SELECTION,
    outs: DEFAULT_RELATIONSHIPS_SELECTION,
};

const nodeSelectionLookup = {};

export function getNodeSelection(node: node) : nodeSel {
    return _.defaultsDeep(
        {},
        _.cloneDeep(node.selection),
        DEFAULT_SELECTION
    );
}

export function mergeNodeSelection(sel: nodeSel, newVals: any) {
    return _.merge(_.cloneDeep(sel), _.cloneDeep(newVals));
}

export function getNodeById(nodes: node[], id: nodeid) {
    let node: node = underscore.findWhere(nodes, { id, })
    if(!node) {
        throw new Error(`node not found: ${id}`)
    }
    return node
}

export function getEdges(nodes: node[]) : edge[] {
    return nodes.map(node => {
        return [].concat(node.ins, node.outs).filter((edge: edge) => {
            return underscore.findWhere(nodes, { id: edge.target }) && underscore.findWhere(nodes, { id: edge.source })
        })
    }).reduce((prev, curr) => prev.concat(curr), []);
}