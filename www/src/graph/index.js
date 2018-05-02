import React from 'react';
import graphDOT from 'raw-loader!../../graph.dot';

import D3GraphCtn from './graph';
import GraphControlsView from './ui';



export default class GraphUI extends React.Component {
    render() {
        return <div class='ctn'>
            <GraphControlsView/>
            <D3GraphCtn/>
        </div>
    }
}