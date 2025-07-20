import { render } from "preact";
import App from "./App";
import "./index.css";
import "./app.css";

render(<App />, document.getElementById("app") as HTMLElement);
