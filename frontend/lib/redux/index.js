//@flow
import { applyMiddleware, createStore } from "redux";
import { composeWithDevTools } from "redux-devtools-extension";
import { createLogger } from "redux-logger";
import withRedux from "next-redux-wrapper";
import createSagaMiddleware from 'redux-saga';

import rootReducer from "./store";
import rootSaga from './sagas';

const isProduction = process.env.NODE_ENV === "production";

const composeEnhancers = composeWithDevTools({});

const logger = createLogger({
    collapsed: true,
    predicate: () => !isProduction
});

export default withRedux(
        (initialState /*options*/) => {
            const middleware = [logger];
            const delay = [];
            const sagaMiddleware = createSagaMiddleware();
            middleware.push(sagaMiddleware);
            delay.push(() => sagaMiddleware.run(rootSaga));
            const store = createStore(
                    rootReducer,
                    initialState,
                    composeEnhancers(applyMiddleware(...middleware))
            );
            delay.forEach(x => x());
            return store;
        },
        { debug: false }
);