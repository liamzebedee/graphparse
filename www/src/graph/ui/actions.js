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