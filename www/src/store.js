
import { createStore, combineReducers, applyMiddleware, compose } from 'redux'
import thunk from 'redux-thunk';
import {  
    routerMiddleware, 
} from 'react-router-redux'
import rootReducer from './reducers/index';

const enhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;


export default function configureStore(preloadedState = {}) {
    const store = createStore(
        rootReducer,
        preloadedState,
        enhancers(
            applyMiddleware(
                routerMiddleware(history),
                thunk,
            )
        ),
    )

    if(module.hot) {
        module.hot.accept('./reducers/index', () =>{
            const newRootReducer = require('./reducers/index').default;
            store.replaceReducer(newRootReducer)
        });
    }

    return store;
}

