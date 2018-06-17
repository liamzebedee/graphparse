import React from 'react';
import { connect } from 'react-redux'

import Checkbox, {
    CheckboxStateless,
    CheckboxGroup
} from '@atlaskit/checkbox';

import _ from 'underscore';

import {
    getNodeTypes,
    getVariantName
} from '../colours';

import {
    toggleFilter
} from '../actions';

import {
    getSelectedNode
} from '../selectors';

import './filters.css';

const SHOWN_FILTERS = [
    'Struct',
    'Method',
    'Func',
    'Field',
    'File'
]

const Filters = ({ shownNodeTypes, node, toggleFilter }) => {
    return <div styleName='filters'>
        <CheckboxGroup>
        { getNodeTypes().map((variantName, variant) => {
            if(!_.contains(SHOWN_FILTERS, variantName)) return;
            
            let checked = false;
            if(node) {
                checked = _.contains(shownNodeTypes, variant);
            }

            return <CheckboxStateless
                isChecked={checked}
                onChange={toggleFilter}
                isActive={false}
                label={variantName}
                isFullWidth={false}
            />
        }) }
        </CheckboxGroup>
    </div>
}

const mapStateToProps = state => {
    return {
        node: getSelectedNode(state.graph)
    }
}

const mapDispatchToProps = dispatch => {
    return {
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Filters);