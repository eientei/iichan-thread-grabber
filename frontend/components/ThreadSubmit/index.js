import {withNamespaces} from "../../lib/locale";
import * as React from "react";
import styled from "@emotion/styled";
import {connect} from "react-redux";
import {threadActions} from "../../lib/redux/store";
import {WideColumnContainer} from "../WideColumnContainer";
import {CenteringRowContainer} from "../CenteringRowContainer";
import {ThreadSubmitProgress} from "../ThreadSubmitProgress";
import {ErrorText} from "../ErrorText";

const TextInput = styled.input`
  width: 50%;
`;

const initialState = {
    threadUrl: null,
    last: 0
};

export const ThreadSubmit = connect(({thread}) => thread)(withNamespaces(["thread_submit", "error"])(class extends React.Component<*> {
    constructor(props) {
        super(props);
        this.state = initialState;
    };

    componentDidMount() {
        this.setState(initialState);
    }

    submitThread() {
        const {dispatch} = this.props;
        const {threadUrl} = this.state;
        this.setState({last: 0, submitted: true});
        dispatch(threadActions.submit(threadUrl));
    }

    render() {
        const {t, submitting, except} = this.props;
        return (
                submitting ? (
                        <ThreadSubmitProgress redirect={true}/>
                ) : (
                        <WideColumnContainer>
                            <CenteringRowContainer>
                                <label htmlFor="thread_url">{t("input_thread_url")}</label>
                            </CenteringRowContainer>
                            <CenteringRowContainer>
                                <TextInput type="text" name="thread_url" id="thread_url"
                                           onChange={(e) => this.setState({threadUrl: e.target.value})}
                                           onKeyDown={(e) => e.key === 'Enter' && this.submitThread()}/>
                                <button onClick={() => this.submitThread()}>{t("common:submit")}</button>
                            </CenteringRowContainer>
                            {except != null && (
                                    <CenteringRowContainer>
                                        <ErrorText>{t("thread_submitting_error", {except: t("error:" + except)})}</ErrorText>
                                    </CenteringRowContainer>
                            )}
                        </WideColumnContainer>
                )
        );
    }
}));
