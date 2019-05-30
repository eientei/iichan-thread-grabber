import * as React from "react";
import {Root} from "../../layouts/root";
import {ThreadEdit} from "../../components/ThreadEdit";

export default class extends React.Component<*> {
    static async getInitialProps() {
        return {
            namespacesRequired: ['common'],
        }
    }

    render() {
        return (
                <Root>
                    <ThreadEdit/>
                </Root>
        );
    }
}