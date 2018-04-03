import React from 'react';
import D3Graph from './graph';

import graphJSON from '../../graph.json';
import { createStore } from 'satcheljs';
import { observer } from 'mobx-react';

export const getStore = createStore(
    'graphStore',
    graphJSON
);

@observer
export default class GraphUI extends React.Component {
    render() {
        return <div>
            <D3Graph store={getStore()}/>
        </div>
    }
}