import { action } from 'satcheljs';
import { mutator } from 'satcheljs';

import { getStore } from './ui';

export const hoverNode = action(
    'HOVER_NODE',
    (id) => ({ id })
);

export const unhoverNode = action(
    'UNHOVER_NODE',
    (id) => ({ id })
);

export const selectNode = action(
    'SELECT_NODE',
    (id) => ({ id })
);

mutator(hoverNode, (actionMessage) => {
    // highlight paths
    getStore().highlightedPaths = [];
})