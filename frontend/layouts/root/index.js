//@flow

import * as React from "react";
import styled from "@emotion/styled";
import {Header} from "../../components/Header";
import {ColumnContainer} from "../../components/ColumnContainer";

const Container = styled(ColumnContainer)`
  min-height: 100%;
  justify-content: space-between;
`;

const Content = styled(ColumnContainer)`
  background: ${props => props.theme.content.bg};
  flex-grow: 1;
  color: ${props => props.theme.content.fg};
`;

export const Root = ({ children }: {children?: React.Node}) => {
    return (
            <Container>
                <Header/>
                <Content>
                    {children}
                </Content>
            </Container>
    );
};