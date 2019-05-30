//@flow
import css from "@emotion/css";
import emotionNormalize from "emotion-normalize";

export default css`
  ${emotionNormalize}
  @import url("https://fonts.googleapis.com/css?family=Montserrat:400,500,700");
  @import url("https://fonts.googleapis.com/icon?family=Material+Icons");

  html, body, #__next {
    height: 100%;
    min-height: 100%;
    font-family: sans-serif;
  }  
`;