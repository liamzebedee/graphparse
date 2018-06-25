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

export function selectClickAction(action: any) {
    return {
        type: "SELECT_CLICK_ACTION",
        action,
    }
}