import { createEffect, Show, createSignal } from "solid-js";
import "./styles/TextInput.css";

type TextInputProps = {
  errors: string[],
  errorsSetFn: (errors: string[]) => void,
  value: string,
  valueSetFn: (value: string) => void,
  placeholder?: string,
  autofocus?: boolean,
  name?: string
  bgColor?: string,
}

function TextInput(props: TextInputProps) {
  //autofocus
  let inputRef: HTMLInputElement | undefined
  createEffect(() => {
    if (props.autofocus) {
      inputRef?.focus()
    }
  })
  const [focused, setFocused] = createSignal(false)
  const [titleOnTop, setTitleOnTop] = createSignal(false)

  function nameFocus() {
    console.log("focus", focused())
    if (focused()) {
      return
    }

    console.log("props.value", props.value)
    if (!props.value) {
      setTitleOnTop(false)
    } 
  }

  return (
    <div class="text-input">
      <Show when={props.name}>
        <div class={`text-input-name ${titleOnTop() ? "title-on-top" : ""} clickable`}
          style={props.bgColor ? {"background-color": props.bgColor} : ""}
          onClick={nameFocus}>
          <p>{props.name}</p>
        </div>
      </Show>
      <input 
        ref={inputRef}
        type="text"
        class="form-input"
        style={props.bgColor ? {"background-color": props.bgColor} : ""}
        classList={{"form-input-error": props.errors.length > 0}}
        onFocus={() => props.errorsSetFn([])}
        onClick={() => setTitleOnTop(true)}
        onBlur={() => nameFocus()}
        onInput={(e) => props.valueSetFn(e.target.value)}
        value={props.value}
        placeholder={props.placeholder}/>
    </div>
  )
}

export default TextInput;