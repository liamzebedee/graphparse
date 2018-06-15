import React from 'react';
import { connect } from 'react-redux'
import classNames from 'classnames';
import { AtlaskitThemeProvider, themed, colors } from '@atlaskit/theme';
import _ from 'lodash'

import {
    changeDepth,
    VIEWS,
    changeView,
    toggleShowDefinitions,
    searchFocusChange,
    loadGraph
} from './actions'

import Search from './ui/search';
import NodeControls from './ui/node-controls';

import './ui.css';
import nodeColor, { getVariantName } from './colours';



class GraphControls extends React.Component {
    state = {
        searchFocused: false
    }
    constructor() {
        super()
    }

    render() {
        let {
            changeDepth, maxDepth,
            toggleShowDefinitions, showDefinitions,

            currentNode
        } = this.props;
        
        return <div styleName="ui-overlay">
            <div styleName='sidebar ui-pane ui-pad ui-rows-pad'>
                <Search/>
                <NodeControls/>
            </div>
        </div>
    }
}


const mapStateToProps = (state) => {
    return {
        nodeTypes: state.graph.nodeTypes,
        clickedNode: state.graph.clickedNode ? state.graph.nodeLookup[state.graph.clickedNode] : null,
        maxDepth: state.graph.maxDepth,
        showDefinitions: state.graph.showDefinitions,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        changeDepth: (depth) => dispatch(changeDepth(depth)),
        toggleShowDefinitions: () => dispatch(toggleShowDefinitions()),
    }
}

const GraphControlsView = connect(
  mapStateToProps,
  mapDispatchToProps
)(GraphControls)
 
export default GraphControlsView;