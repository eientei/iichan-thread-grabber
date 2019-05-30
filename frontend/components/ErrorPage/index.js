//@flow

import * as React from "react";
import styled from "@emotion/styled";
import {withNamespaces} from "../../lib/locale";

const ErrorWrapper = styled.h1`
  padding: 1em;
`;

type GetInitialPropsArgs = {
    err?: any,
    pathname: string,
    query: any,
    req?: any,
    res?: any,
    xhr: any
};

export const ErrorPage = withNamespaces("error")(class extends React.Component<{ t : Function, statusCode: number }> {
    static getInitialProps({ res, xhr } : GetInitialPropsArgs) {
        const statusCode = res ? res.statusCode : xhr ? xhr.status : null;
        return {
            statusCode,
            namespacesRequired: ['common'],
        };
    }

    render() {
        const {statusCode, t} = this.props;
        return (
                <ErrorWrapper>
                    {statusCode ? t(statusCode + "_code") : t("404_code")}
                </ErrorWrapper>
        );
    }
});
