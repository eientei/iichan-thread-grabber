import { select, put, takeLatest, all } from 'redux-saga/effects'
import {threadActions} from "../store";

function* initSaga() {
    yield null
}

function* threadSubmit({payload: {thread}}) {
    const params = new URLSearchParams;
    params.append('thread_url', thread);
    try {
        const resp = yield fetch(process.env.API_HOST + '/', {
            method: 'POST',
            mode: 'cors',
            headers: {
                'Content-type': 'application/x-www-form-urlencoded',
            },
            body: params.toString(),
            credentials: 'same-origin',
        });
        switch (resp.status / 100) {
            case 4:
                throw new Error("invalid_request");
            case 5:
                throw new Error("server_error");
        }

        const reader = resp.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const res = yield reader.read();
            const chunk = decoder.decode(res.value || new Uint8Array());
            const parts = chunk.split("\n").map(x => {
                try {
                    return JSON.parse(x);
                } catch {
                    return null;
                }
            }).filter(x => x != null);
            yield all(parts.map(p => put(threadActions.progress(p))));
            if (res.done) {
                yield put(threadActions.complete(null));
                break;
            }
        }
    } catch (e) {
        yield put(threadActions.complete(e.message));
    }
}

function* threadGet({payload: {path}}) {
    try {
        const resp = yield fetch(process.env.API_HOST + path, {
            method: 'GET',
            mode: 'cors',
            credentials: 'same-origin',
        });
        switch (resp.status / 100) {
            case 4:
                throw new Error("invalid_request");
            case 5:
                throw new Error("server_error");
        }

        const reader = resp.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const res = yield reader.read();
            const chunk = decoder.decode(res.value || new Uint8Array());
            const parts = chunk.split("\n").map(x => {
                try {
                    return JSON.parse(x);
                } catch {
                    return null;
                }
            }).filter(x => x != null);
            yield all(parts.map(p => put(threadActions.result(p && p.data, null))));
            if (res.done) {
                break;
            }
        }
    } catch (e) {
        yield put(threadActions.result(null, e.message));
    }
}

function* threadSave({payload: {path}}) {
    const {thread: {result}} = yield select();
    const json = JSON.stringify(result);
    try {
        const resp = yield fetch(process.env.API_HOST + path, {
            method: 'POST',
            mode: 'cors',
            headers: {
                'Content-type': 'application/json',
            },
            body: json,
            credentials: 'same-origin',
        });
        switch (resp.status / 100) {
            case 4:
                throw new Error("invalid_request");
            case 5:
                throw new Error("server_error");
        }

        const reader = resp.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const res = yield reader.read();
            const chunk = decoder.decode(res.value || new Uint8Array());
            const parts = chunk.split("\n").map(x => {
                try {
                    return JSON.parse(x);
                } catch {
                    return null;
                }
            }).filter(x => x != null);
            yield all(parts.map(p => put(threadActions.progress(p))));
            if (res.done) {
                yield put(threadActions.saved(null));
                break;
            }
        }
    } catch (e) {
        yield put(threadActions.saved(e.message));
    }
}

function* dispatchRedux() {
    yield all([
        takeLatest(threadActions.submit().type, threadSubmit),
        takeLatest(threadActions.get().type, threadGet),
        takeLatest(threadActions.save().type, threadSave),
    ]);
}

export default function* rootSaga() {
    yield all([
        initSaga(),
        dispatchRedux(),
    ]);
}