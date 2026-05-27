import "./styles/Errors.css";
import { For, mergeProps } from "solid-js";
import { useDefaults } from "./DefaultsContext.tsx";
import Text from "./Text.tsx";
import errorsText from "./translations/Errors.json";

export type ErrorKey = keyof typeof errorsText;

type MarginProp = {
  marginTop?: number,
  marginBottom?: number,
}

type ErrorsProps = {
  errors: ErrorKey[],
  margin?: MarginProp
}

function Errors(props: ErrorsProps) {
  const errorDefaults = mergeProps({
    margin: {
      marginBottom: 0.5,
      marginTop: 0.5
    }
  }, props)
  const { defaults } = useDefaults();

  return (
    <div class="error"
      style={{
        "margin-top": `${errorDefaults.margin.marginTop}rem`,
        "margin-bottom": `${errorDefaults.margin.marginBottom}rem`
      }}>
      <For each={props.errors}>
        {(error, _) => 
          <p class="error-text"><Text value={errorsText[error]} lang={defaults().lang}/></p>
        }
      </For>
    </div>
  )
}

export default Errors