import FieldRange, { FieldRangeStateless } from '@atlaskit/field-range';
import FieldText, { FieldTextStateless } from '@atlaskit/field-text';

import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './depth.css';

const Depth = ({ depth, changeDepth }) => {
    return <div styleName='depth'>
        <FieldTextStateless label="Depth" type="Number" 
            value={depth}
            min={0}
            max={10}
            step={1}
            onChange={(ev) => changeDepth(parseInt(ev.target.value))}/>
    </div>
}

const mapStateToProps = state => {
    return {
    }
}

const mapDispatchToProps = dispatch => {
    return {
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Depth);