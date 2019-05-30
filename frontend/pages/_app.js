//@flow

import * as React from "react";
import { Container } from "next/app";
import { Global } from "@emotion/core";
import { ThemeProvider } from 'emotion-theming'
import { Provider } from "react-redux";
import Head from "next/head";
import globalStyles from "../style/global";
import withRedux from "../lib/redux";
import {appWithTranslation, withNamespaces} from "../lib/locale";
import {provideSetter} from "../style/theme";
import HTML5Backend from 'react-dnd-html5-backend';
import { DragDropContext } from 'react-dnd';

class App extends React.Component<*,*> {
    static async getInitialProps({ Component, ctx }) {
        const pageProps = Component.getInitialProps  ? await Component.getInitialProps(ctx) : {};
        return {
            pageProps,
            namespacesRequired: ['common'],
        };
    }
    constructor(props) {
        super(props);
        this.state = {
            theme: null,
        };
    }

    componentDidMount() {
        provideSetter((state) => this.setState(state));
    }

    render() {
        const {
            Component,
            pageProps,
            store,
            t
        } = this.props;

        const {
            theme,
        } = this.state;

        if (theme == null) {
            return null;
        }

        return (
                <Provider store={store}>
                    <Container>
                        <Global styles={globalStyles} />
                        <Head>
                            <title>{t("title")}</title>
                        </Head>
                        <ThemeProvider theme={theme}>
                            <Component {...pageProps} />
                        </ThemeProvider>
                    </Container>
                </Provider>
        );
    }
}

export default withRedux(appWithTranslation(withNamespaces("common")(DragDropContext(HTML5Backend)(App))));
