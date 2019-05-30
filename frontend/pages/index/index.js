//@flow
import * as React from "react";
import {Root} from "../../layouts/root";
import styled from "@emotion/styled";
import {ThreadSubmit} from "../../components/ThreadSubmit";
import {RowContainer} from "../../components/RowContainer";

const Center = styled(RowContainer)`
  justify-content: center;
  padding-top: 10%;
`;

export default class extends React.Component<*> {
    static async getInitialProps() {
        return {
            namespacesRequired: ['common'],
        }
    }

    render() {
        return (
                <Root>
                    <Center>
                        <ThreadSubmit/>
                    </Center>
                </Root>
        );
    }
}