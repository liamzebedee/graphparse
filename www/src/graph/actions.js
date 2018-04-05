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
    let nodes = getSubPaths(getStore().adjList, actionMessage.id)
    getStore().interested = nodes;
})


export function getSubPaths(adjList, fromNodeId) {
    let currentNodes = [
        fromNodeId
    ];

    let interested = new Set();
    let visited = new Set();

    do {
        let next = [];

        currentNodes.map(node => {
            if(visited.has(node)) return;
            else visited.add(node)

            interested.add(node)

            let outs = adjList[""+node]
            if(outs) next = next.concat(outs)
        })

        currentNodes = next;

    } while(currentNodes.length)

    return Array.from(interested);
}