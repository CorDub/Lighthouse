import "./Errors.css"
import { For, mergeProps } from "solid-js"

type MarginProp = {
  marginTop?: number,
  marginBottom?: number,
}

type ErrorsProps = {
  errors: string[],
  margin?: MarginProp
}

function Errors(props: ErrorsProps) {
  const defaults = mergeProps({
    margin: {
      marginBottom: 0.5,
      marginTop: 0.5
    }
  }, props)

  return (
    <div class="error"
      style={{
        "margin-top": `${defaults.margin.marginTop}rem`,
        "margin-bottom": `${defaults.margin.marginBottom}rem`
      }}>
      <For each={props.errors}>
        {(error, _) => 
          <p class="error-text">{error}</p>
        }
      </For>
    </div>
  )
}

export default Errors