import React from 'react';
import graphDOT from 'raw-loader!../../graph.dot';

import D3GraphCtn from './graph';
import GraphControlsView from './ui';

import './index.css';


export default class GraphUI extends React.Component {
    render() {
        return <div styleName='ctn'>
            <GraphControlsView/>
            <D3GraphCtn/>
        </div>
    }
}