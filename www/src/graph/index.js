import React from 'react';
import { connect } from 'react-redux'

import Graph from './graph';
import UI from './ui';

import './index.css';

import { load } from './actions';

import queryString from 'query-string';

class GraphUI extends React.Component {
    componentDidMount() {
        let { name } = queryString.parse(this.props.location.search);
        this.props.load(name)
    }

    render() {
        return <div styleName='ctn'>
            <UI/>
            <Graph/>
        </div>
    }
}



export default connect(null, dispatch => ({
    load: (codebaseId) => {
        dispatch(load(codebaseId));
        // dispatch(loadGraph())
        // dispatch(loadInitialFileForTesting())
    }
}))(GraphUI)