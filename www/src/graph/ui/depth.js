import FieldRange, { FieldRangeStateless } from '@atlaskit/field-range';
import FieldText, { FieldTextStateless } from '@atlaskit/field-text';

import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './depth.css';
import {
    changeDepth
} from '../actions';

const Depth = ({ maxDepth, changeDepth }) => {
    return <div styleName='depth'>
        <FieldTextStateless type="Number" value={maxDepth}
            min={0}
            max={10}
            step={1}
            onChange={(ev) => changeDepth(ev.target.value)}/>
    </div>
}

const mapStateToProps = state => {
    return {
        maxDepth: state.graph.maxDepth
    }
}

const mapDispatchToProps = dispatch => {
    return {
        changeDepth: (depth) => dispatch(changeDepth(depth))
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Depth);