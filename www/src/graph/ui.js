import React from 'react';
import D3Graph from './graph';

import graphJSON from '../../graph.json';
import { createStore } from 'satcheljs';
import { observer } from 'mobx-react';
import { observable } from 'mobx'

import graphDOT from 'raw-loader!../../graph.dot';

export let getStore = createStore(
    'graphStore',
    observable({ 
        graphDOT,
        interested: [],
        ...graphJSON
    })
);




@observer
export default class GraphUI extends React.Component {
    render() {
        let store = getStore()

        return <div>
            <D3Graph 
                graphDOT={store.graphDOT}
                nodeLookup={store.nodeLookup}
                interested={store.interested}
            />
        </div>
    }
}