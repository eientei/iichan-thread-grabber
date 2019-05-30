//@flow

import * as React from "react";
import Document, { Head, Main, NextScript } from "next/document";

export default class extends Document {
    render() {
        return (
                <html lang="en">
                <Head>
                    <link href="/static/favicons/favicon.ico" rel="shortcut icon" />
                </Head>
                <body>
                <Main />
                <NextScript />
                </body>
                </html>
        );
    }
}
