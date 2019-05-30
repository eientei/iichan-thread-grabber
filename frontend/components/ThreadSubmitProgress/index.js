import {withNamespaces} from "../../lib/locale";
import * as React from "react";
import {connect} from "react-redux";
import {threadActions} from "../../lib/redux/store";
import {WideColumnContainer} from "../WideColumnContainer";
import {CenteringRowContainer} from "../CenteringRowContainer";
import Router from 'next/router';

const initialState = {
    progress: {status: "waiting"},
};

export const ThreadSubmitProgress = connect(({thread}) => thread)(withNamespaces(["thread_submit", "error", "common"])(class extends React.Component<*> {
    constructor(props) {
        super(props);
        this.state = initialState;
    };

    cancelSubmit() {
        const {dispatch} = this.props;
        dispatch(threadActions.cancel());
    }

    static getDerivedStateFromProps({submitting, progress, last, dispatch, redirect, except}, state) {
        if (progress.length > 0) {
            const last = progress[progress.length-1];
            if (redirect && submitting && last.status === 'complete' && except == null) {
                const base = new URL(process.env.API_HOST).pathname;
                const target = process.env.PUBLIC_PREFIX + new URL(last.data.base).pathname.substring(base.length);
                Router.push("/edit", target);
            }
            return ({...state, progress: last, last: progress.length});
        }
        return state;
    }

    renderProgressWaiting() {
        const {t} = this.props;
        return (
                <>
                    <CenteringRowContainer>{t("waiting")}</CenteringRowContainer>
                </>
        );
    }

    renderProgressQueue({position}) {
        const {t} = this.props;
        return (
                <>
                    <CenteringRowContainer>{t("queue_progress", {position})}</CenteringRowContainer>
                </>
        );
    }

    renderProgressDownload({current_download, total_download}) {
        const {t} = this.props;
        return (
                <>
                    <CenteringRowContainer>{t("download_progress", {current: current_download, total: total_download})}</CenteringRowContainer>
                </>
        );
    }

    renderProgressComplete({base}) {
        const {t} = this.props;
        return (
                <>
                    <CenteringRowContainer>{t("complete")}</CenteringRowContainer>
                    <CenteringRowContainer>{base}</CenteringRowContainer>
                </>
        );
    }

    renderProgressUnknown(unk) {
        console.log("unknown", unk);
        return JSON.stringify(unk);
    }

    renderProgress({status, data}) {
        switch(status) {
            case "waiting":
                return this.renderProgressWaiting();
            case "queue_progress":
                return this.renderProgressQueue(data);
            case "download_progress":
                return this.renderProgressDownload(data);
            case "complete":
                return this.renderProgressComplete(data);
            default:
                return this.renderProgressUnknown({status, data});
        }
    }

    render() {
        const {t, thread} = this.props;
        const {progress} = this.state;
        return (
                <WideColumnContainer>
                    <CenteringRowContainer>
                        {t("thread_submitting", {thread})}
                    </CenteringRowContainer>
                    <CenteringRowContainer>
                        <button onClick={() => this.cancelSubmit()}>{t("common:cancel")}</button>
                    </CenteringRowContainer>
                    {this.renderProgress(progress)}
                </WideColumnContainer>
        );
    }
}));