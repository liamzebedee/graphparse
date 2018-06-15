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
    toggleNodeTypeFilter
} from '../actions';
import './filters.css';

const SHOWN_FILTERS = [
    'Struct',
    'Method',
    'Func',
    'Field',
    'File'
]

const Filters = ({ currentNode, toggleFilter }) => {
    return <div styleName='filters'>
        { getNodeTypes().map((variantName, variant) => {
            if(!_.contains(SHOWN_FILTERS, variantName)) return;
            let checked = false;
            if(currentNode) {
                checked = _.contains(currentNode.filters.shownNodeTypes, variant);
            }

            return <CheckboxStateless
            isChecked={checked}
            onChange={() => toggleFilter(variant)}
            isActive={false}
            label={variantName}
            isFullWidth={false}
            />
        }) } 
    </div>
}

const mapStateToProps = state => {
    let { nodes, selectedNode } = state.graph;
    let currentNode = _.findWhere(nodes, { id: selectedNode });

    return {
        currentNode,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        toggleFilter: (variant) => dispatch(toggleNodeTypeFilter(variant))
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Filters);