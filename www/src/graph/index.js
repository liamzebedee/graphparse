import React from 'react';
import graphDOT from 'raw-loader!../../graph.dot';

import D3GraphCtn from './graph';
import GraphControlsView from './ui';

import {
    createStore,
    combineReducers,
    applyMiddleware,
    compose
} from 'redux';
import { Provider } from 'react-redux'
import thunk from 'redux-thunk';
import rootReducer from './reducers';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const store = createStore(
    rootReducer,
    composeEnhancers(
        applyMiddleware(thunk)
    )
)


export default class GraphUI extends React.Component {
    render() {
        return <Provider store={store}>
            <div class='ctn'>
                <GraphControlsView/>
                <D3GraphCtn/>
            </div>
        </Provider>
    }
}