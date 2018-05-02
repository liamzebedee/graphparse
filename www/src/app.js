import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, combineReducers, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'
import createHistory from 'history/createBrowserHistory'
import { Route } from 'react-router'
import thunk from 'redux-thunk';
import { 
    ConnectedRouter, 
    routerReducer, 
    routerMiddleware, 
} from 'react-router-redux'
import { hot } from 'react-hot-loader'

import Home from './home/'
import GraphUI from './graph';
import ASTView from './ast';
import graph from './graph/reducers';

const history = createHistory()
const enhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
    combineReducers({
        router: routerReducer,
        graph,
    }),
    enhancers(
        applyMiddleware(
            routerMiddleware(history),
            thunk,
        )
    ),
)

const App = () => <Provider store={store}>
    <ConnectedRouter history={history}>
        <div>
            <Route exact path="/" component={Home}/>
            <Route path="/repo/:id" component={GraphUI}/>
            <Route path="/ast" component={ASTView}/>
        </div>
    </ConnectedRouter>
</Provider>;

export default App;