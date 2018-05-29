import FieldRange, { FieldRangeStateless } from '@atlaskit/field-range';

import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './depth.css';
import {
    changeDepth
} from '../actions';

const Depth = ({ maxDepth, changeDepth }) => {
    return <div styleName='depth'>
        <FieldRange
            value={maxDepth}
            min={0}
            max={10}
            step={1}
            onChange={(x) => changeDepth(x)}
          />
    </div>
}

const mapStateToProps = state => {
    return {
        maxDepth: state.maxDepth
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