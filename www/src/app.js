import React from 'react'

import ReactDOM from 'react-dom'


import { createStore, combineReducers, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'

import createHistory from 'history/createBrowserHistory'
import { Route } from 'react-router'

import { 
    ConnectedRouter, 
    routerReducer, 
    routerMiddleware, 
} from 'react-router-redux'


import Home from './home/'
import GraphUI from './graph';
// import reducers from './reducers'

const history = createHistory()
const middleware = routerMiddleware(history)

export const store = createStore(
  combineReducers({
    router: routerReducer
  }),
  applyMiddleware(middleware)
)
    

export const App = () => <Provider store={store}>
    <ConnectedRouter history={history}>
        <div>
            <Route exact path="/" component={Home}/>
            <Route path="/repo/:id" component={GraphUI}/>
            {/* <Route path="/ast" component={GraphUI}/> */}
        </div>
    </ConnectedRouter>
</Provider>;