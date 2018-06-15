import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './node-controls.css';

import Filters from './filters';
import Depth from './depth';

import {

} from '../actions';

const NodeControls = ({ currentNode }) => {
    return <div styleName='ui-pad'>
        <h4>Current node</h4>
        { currentNode == null 
            ? <i>No node selected</i>
            : <div>
                <div styleName='current-node'>
                    <u>{currentNode.label}</u>

                    <div>Shown <input type='checkbox'/></div>

                    <h5>Relations</h5>
                    <h6>In</h6>
                    <RelationControl/>

                    <h6>Out</h6>
                    <RelationControl/>
                </div>
            </div>
        }
    </div>
}

const RelationControl = ({}) => {
    return <div>
        <Filters/>
        <Depth/>
    </div>
}

const mapStateToProps = state => {
    let currentNode = null;
    if(state.graph.currentNode) {
        currentNode = _.find(state.graph.nodes, { id: state.graph.currentNode });
    }
    return {
        currentNode,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeControls);