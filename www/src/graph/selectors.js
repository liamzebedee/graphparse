// @flow
import _ from 'lodash';

// $FlowFixMe
import proxyDefaults from '/Users/liamz/Documents/open-source/proxy-object-defaults';

import {
    getNodeTypes,
} from './colours';

import type {
    relationshipsSel,
    nodeSel,
    node
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
    maxDepth: 2,
    showDefs: true,
    showUses: true,
    shown: false,
};

const DEFAULT_SELECTION: nodeSel = {
    shown: true,
    ins: DEFAULT_RELATIONSHIPS_SELECTION,
    outs: {
        ...DEFAULT_RELATIONSHIPS_SELECTION,
        shown: true
    },
};

const nodeSelectionLookup = {};

export function getNodeSelection(node: node) : nodeSel {
    return _.defaultsDeep(DEFAULT_SELECTION, node.selection);
}

export function mergeNodeSelection(sel: nodeSel, newVals: any) {
    return _.merge(sel, newVals);
}