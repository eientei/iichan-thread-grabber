//@flow
import * as React from "react";
import styled from "@emotion/styled";
import {withNamespaces} from "../../lib/locale";
import Link from 'next/link';
import {RowContainer} from "../RowContainer";

const Left = styled.div`
  max-width: 10em;
  text-align: left;
`;

const Center = styled.h3`
  text-align: center;
`;

const Right = styled.div`
  max-width: 10em;
  text-align: right;
`;

const Container = styled(RowContainer)`
  height: 4em;
  justify-content: space-between;
  align-items: center;
  background: ${props => props.theme.panel.bg};
  color: ${props => props.theme.panel.fg};
  padding: 0 1em;
`;

export const Header = withNamespaces("header")(({t}) => (
        <Container>
            <Left>
                <Link href="https://iibooru.org"><a>{t("back_to_iibooru")}</a></Link>
            </Left>
            <Center>
            </Center>
            <Right>
            </Right>
        </Container>
));