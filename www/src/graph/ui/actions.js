// @flow
import type {
    nodeid,
    node,
    nodeLayout,
    edge,
    edgeLayout,
    nodeVariant,
    nodeSel
} from 'graphparse';

type relationships = 'ins' | 'outs';

export function selectClickAction(action: any) {
    return {
        type: "SELECT_CLICK_ACTION",
        action,
    }
}

export function changeDepth(node: nodeid, relationships: relationships, depth: number) {
    return {
        type: "CHANGE_DEPTH",
        node,
        relationships,
        depth,
    }
}

export function toggleFilter(node: nodeid, relationships: relationships, variant: nodeVariant) {
    return {
        type: "TOGGLE_FILTER",
        node,
        relationships,
        variant,
    }
}