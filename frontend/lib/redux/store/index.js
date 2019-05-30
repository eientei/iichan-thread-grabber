//@flow

import {createActions, handleActions} from "redux-actions";
import {combineReducers} from "redux";

export const appActions = createActions({
    APP: {
        ERROR: (except, origin) => ({except, origin}),
    },
}).app;

const app = handleActions({
    APP: {
        ERROR: (state, {payload: {except, origin}}) => ({...state, except, origin}),
    },
}, {
    except: null,
    origin: null,
});

export const threadActions = createActions({
    THREAD: {
        SUBMIT: (thread) => ({thread}),
        PROGRESS: (message) => ({message}),
        CONSUME: (consumed) => ({consumed}),
        COMPLETE: (except) => ({except}),
        CANCEL: () => null,
        GET: (path) => ({path}),
        RESULT: (result) => ({result}),
        DRAG: (insert, {fromgroup, fromimg}, {togroup, toimg}) => ({insert, from: {fromgroup, fromimg}, to: {togroup: togroup, toimg}}),
        SAVE: (path) => ({path}),
        SAVED: (except) => ({except}),
        UPDATE: (updates) => ({updates}),
    },
}).thread;


const thread = handleActions({
    THREAD: {
        SUBMIT: (state, {payload: {thread}}) => ({...state, thread, except: null, submitting: true, result: null, progress: []}),
        PROGRESS: ({progress, ...state}, {payload: {message}}) => state.submitting ? ({...state, progress: [...progress, message]}) : ({...state, progress}),
        CONSUME: (state, {payload: {consumed}}) => ({...state, progress: state.progress.filter((_, i) => consumed.indexOf(i) < 0)}),
        COMPLETE: ({submitting, ...state}, {payload: {except}}) => ({...state, submitting: submitting && except == null, except}),
        CANCEL: (state) => ({...state, except: "cancelled", submitting: false}),
        GET: (state) => ({...state, except: null, result: null}),
        RESULT: (state, {payload: {result, except}}) => ({...state, submitting: false, result, except}),
        DRAG: (state, {payload: {insert, from: {fromgroup, fromimg}, to: {togroup, toimg}}}) => {
            let image = null;
            const groups = state.result.groups.map((g, i) => {
                if (i === fromgroup) {
                    image = g.find((_, n) => n === fromimg);
                    return g.filter((_, n) => n !== fromimg);
                }
                return g;
            });

            if (insert) {
                groups.splice(togroup+1, 0, [image]);
            } else {
                if (groups.length <= togroup) {
                    groups.push([image]);
                } else {
                    groups.find((g, i) => i === togroup).splice(toimg, 0, image);
                }
            }
            return ({...state, result: {...state.result, groups: groups.filter((g,i) => i === 0 || g.length > 0)}});
        },
        SAVE: (state) => ({...state, submitting: true}),
        SAVED: (state, {payload: {except}}) => ({...state, submitting: false, except}),
        UPDATE: (state, {payload: {updates}}) => {
            const groups = state.result.groups.map((g, i) => {
                const gupds = updates.filter(({group}) => group === i);
                if (gupds.length > 0) {
                    return g.map((img, n) => {
                        const upds = gupds.filter(({idx}) => idx === n);
                        if (upds.length > 0) {
                            return upds.reduce((acc, {group, idx, ...props}) => ({...acc, ...props}), img);
                        }
                        return img;
                    });
                }
                return g;
            });
            return ({...state, result: {...state.result, groups}});
        }
    },
}, {
    thread: null,
    submitting: false,
    progress: [],
    except: null,
    result: null,
});

export default combineReducers({
    app,
    thread,
});