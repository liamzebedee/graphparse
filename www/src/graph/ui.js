import React from 'react';

import {
    searchNodes,
    loadInitialFileForTesting,
    selectNodeByLabel,
    selectNodeFromSearch,
    changeDepth,
    VIEWS,
    changeView,
    toggleShowDefinitions,
    searchFocusChange
} from './actions'

import { connect } from 'react-redux'

import classNames from 'classnames';

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
        
        return <div className="infoView">
            <div>
                {/* <div className="btn-group btn-group-sm" role="group">
                { VIEWS.map(view => 
                    <button 
                        className={classNames("btn btn-secondary", { "active": uiView == view })} 
                        onClick={() => changeView(view)}>
                        {view}
                    </button>
                )}
                </div> */}

                {/* <input className="depth" type="number" value={maxDepth} onChange={(ev) => changeDepth(ev.target.value)}/> */}

                {/* <div className="form-check">
                    <input className="form-check-input" type="checkbox" 
                        checked={showDefinitions}
                        onChange={toggleShowDefinitions}/>
                    <label className="form-check-label" for="defaultCheck1">
                        Definitions
                    </label>
                </div> */}

                <div className='search'>
                    <input type='text' className="form-control" placeholder="Search types, files" onChange={(ev) => searchNodes(ev.target.value)} value={q} 
                    onFocus={() => this.setState({ searchFocused: true })} 
                    onBlur={() => this.setState({ searchFocused: false })} />
                </div>

                
                <div className={classNames('results', { 'active': searchFocused })}>
                    { matches.length > 0 ? matches.map((node, i) => {
                        return <NodeMatch key={i} onClick={() => selectNode(node.id)} {...node}/>
                    }) : 'none' }
                </div>

                <footer>Basemap</footer>
            </div>

            {/* <div className='debug'>
                <pre>{ clickedNode ? clickedNode.debugInfo : null }</pre>
            </div> */}

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