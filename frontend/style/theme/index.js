import futaba from './futaba';

let setState = null;

export const themes = {
    futaba,
};

declare var process: {
    browser: boolean,
    env: *,
};

let currentThemeValue = process.browser && document.cookie
        .split(/;\s*/)
        .map(x => x.split('=',2))
        .filter(([name]) => name === 'uitheme')
        .map(([, value]) => value)
        .find((x) => x in themes) || Object.keys(themes)[0];

export const provideSetter = (setter) => {
    setState = setter;
    setTheme(currentThemeValue);
};

export const currentThemeName = () => currentThemeValue;
export const currentTheme = () => themes[currentThemeValue];

export const setTheme = (theme) => {
    if (!setState || !themes[theme]) {
        return;
    }
    currentThemeValue = theme;
    setState({theme: themes[theme]});
    if (Object.keys(themes).length > 1 && process.browser) {
        document.cookie = 'uitheme=' + theme;
    }
};