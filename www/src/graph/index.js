import React from 'react';
import { connect } from 'react-redux'

import Graph from './graph';
import UI from './ui';

import './index.css';

import { load } from './actions';

import queryString from 'query-string';

class GraphUI extends React.Component {
    state = {
        firstLoad: true
    }

    componentDidMount() {
        let { name } = queryString.parse(this.props.location.search);
        this.props.load(name, this.state.firstLoad)
        this.setState({ firstLoad: false })
    }

    render() {
        return <div styleName='ctn'>
            <Graph/>
            <UI/>
        </div>
    }
}



export default connect(null, dispatch => ({
    load: (codebaseId, firstLoad) => {
        dispatch(load(codebaseId, firstLoad));
    }
}))(GraphUI)