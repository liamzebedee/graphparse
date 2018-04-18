import React from 'react';
import graphDOT from 'raw-loader!../../graph.dot';

import D3GraphCtn from './graph';
import GraphControlsView from './ui';
import './style.css';

import {
    createStore,
    combineReducers,
    applyMiddleware,
    compose
} from 'redux';
import { Provider } from 'react-redux'
import rootReducer from './reducers';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const middleware = [() => {}];
const store = createStore(
    rootReducer,
    composeEnhancers(
        // applyMiddleware(middleware)
    )
)


export default class GraphUI extends React.Component {
    render() {
        return <Provider store={store}>
            <div>
                <GraphControlsView/>
                <D3GraphCtn/>
            </div>
        </Provider>
    }
}