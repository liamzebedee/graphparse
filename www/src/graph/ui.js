import React from 'react';
import { connect } from 'react-redux'
import classNames from 'classnames';
import { AtlaskitThemeProvider, themed, colors } from '@atlaskit/theme';


import {
    searchNodes,
    loadInitialFileForTesting,
    selectNodeByLabel,
    selectNodeFromSearch,
    changeDepth,
    VIEWS,
    changeView,
    toggleShowDefinitions,
    searchFocusChange,
    loadGraph
} from './actions'


import Filters from './ui/filters';
import Depth from './ui/depth';

import './ui.css';
import nodeColor, { getVariantName } from './colours';

class GraphControls extends React.Component {
    state = {
        searchFocused: false
    }
    constructor() {
        super()
    }

    componentDidMount() {
        this.props.firstLoad()
    }

    render() {
        let {
            nodeTypes, q, matches, searchNodes, selectNode, clickedNode, 
            uiView, changeView,
            changeDepth, maxDepth,
            toggleShowDefinitions, showDefinitions,

            focusSearch, blurSearch
        } = this.props;

        let { searchFocused } = this.state;
        
        return <div styleName="ui-overlay">
            <div styleName='ui-floating'>
                <div styleName='search'>
                    <input type='text' className="form-control" placeholder="Search types, files" onChange={(ev) => searchNodes(ev.target.value)} value={q} 
                    onFocus={() => this.setState({ searchFocused: true })} 
                    onBlur={() => this.setState({ searchFocused: false })} />
                </div>

                <div styleName={classNames('results', { 'active': searchFocused })}>
                    { matches && matches.length > 0 ? matches.map((node, i) => {
                        return <NodeMatch key={i} onClick={() => selectNode(node.id)} {...node}/>
                    }) : 'none' }
                </div>
            </div>
            
            <div styleName='ui-floating options'>
                <Filters/>
                <Depth/>
            </div>
        </div>
    }
}

const NodeMatch = ({ onClick, label, variant }) => {
    return <div onClick={onClick}>
        {label}
        <span className="badge badge-light" style={{
            backgroundColor: nodeColor(variant),
            float: 'right'
        }}>{getVariantName(variant)}</span>
    </div>
}


const mapStateToProps = state => {
    return {
        ...state.graph.search,
        nodeTypes: state.graph.nodeTypes,
        clickedNode: state.graph.clickedNode ? state.graph.nodeLookup[state.graph.clickedNode] : null,
        uiView: state.graph.uiView,
        maxDepth: state.graph.maxDepth,
        showDefinitions: state.graph.showDefinitions,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        firstLoad:   () => {
            dispatch(loadGraph())
            dispatch(loadInitialFileForTesting())
        },
        changeDepth: (depth) => dispatch(changeDepth(depth)),
        searchNodes: (q) => dispatch(searchNodes(q)),
        selectNode:  (id) => dispatch(selectNodeFromSearch(id)),
        changeView: (view) => dispatch(changeView(view)),
        toggleShowDefinitions: () => dispatch(toggleShowDefinitions()),
    }
}

const GraphControlsView = connect(
  mapStateToProps,
  mapDispatchToProps
)(GraphControls)
 
export default GraphControlsView;