import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';
import classNames from 'classnames';

import './search.css';

import {
    searchNodes,
    selectNodeFromSearch
} from '../actions';

import nodeColor, { getVariantName } from '../colours';


class Search extends React.Component {
    state = {
        searchFocused: false,
    }
    
    render() {
        let {
            q, 
            matches,

            selectNode,
            searchNodes
        } = this.props;

        return <div>
            <div styleName='search'>
                <input type='text' className="" placeholder="Search types, files" onChange={(ev) => searchNodes(ev.target.value)} value={q} 
                onFocus={() => this.setState({ searchFocused: true })} 
                onBlur={() => this.setState({ searchFocused: false })} />
            </div>

            <div styleName={classNames('results', { 'active': this.state.searchFocused })}>
                { matches && matches.length > 0 ? matches.map((node, i) => {
                    return <NodeMatch key={i} onClick={() => selectNode(node)} {...node}/>
                }) : 'none' }
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
    }
}

const mapDispatchToProps = dispatch => {
    return {
        searchNodes: (q) => dispatch(searchNodes(q)),
        selectNode:  (node) => dispatch(selectNodeFromSearch(node)),
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Search);