import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './controls.css';

import {
    selectClickAction
} from './actions';

const ClickControl = ({ clickAction, keybinding, children, selectClickAction, currentAction }) => {
    let active = currentAction == clickAction;
    return <div styleName='click-control' onClick={() => selectClickAction(clickAction) }>
        { active ? <b>{children}</b> : children }
    </div>;
}

const mapStateToProps = state => {
    return {
        currentAction: state.graph.clickAction
    }
}

const mapDispatchToProps = dispatch => {
    return {
        selectClickAction: (action) => dispatch(selectClickAction(action))
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(ClickControl);