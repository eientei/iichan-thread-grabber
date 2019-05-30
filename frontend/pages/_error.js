//@flow

import * as React from "react";
import {Root} from "../layouts/root";
import {ErrorPage} from "../components/ErrorPage";

export default class extends React.Component<*> {
    static async getInitialProps() {
        return {
            namespacesRequired: ['common'],
        }
    }

    render() {
        return (
                <Root>
                    <ErrorPage statusCode={this.props.statusCode}/>
                </Root>
        );
    }
};