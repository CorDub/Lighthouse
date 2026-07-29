import { createEffect, Show, createSignal } from "solid-js";
import "./styles/TextInput.css";

type TextInputProps = {
  errors: string[],
  errorsSetFn: (errors: string[]) => void,
  value?: string,
  valueSetFn: (value: string) => void,
  valueSetFnArg?: string,
  autocomplete?: boolean,
  autocompleteValue?: string,
  setAutocompleteFn?: (value: string) => void,
  searchFn?: () => void,
  enterFn?: (...args: any[]) => any,
  enterFnArg?: any,
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
    if (focused()) {
      return
    }

    if (props.name && titleOnTop()) {
      props.valueSetFn("")
      setTitleOnTop(false)
      props.searchFn?.()
      if (props.setAutocompleteFn) { props.setAutocompleteFn("") }
    }
  }

  function clickOnName() {
    setTitleOnTop(true)
    inputRef?.focus()
  }

  function triggerSearch(value: string) {
    props.valueSetFn(value)
    props.searchFn?.()
  }

  return (
    <div class="text-input">
      <Show when={props.name}>
        <div class={`text-input-name ${titleOnTop() ? "title-on-top" : ""}`}
          style={props.bgColor ? {"background-color": props.bgColor} : ""}
          onClick={() => clickOnName()}>
          <p>{props.name}</p>
        </div>
      </Show>
      <div class="ti-wrapper">
        <input 
          ref={inputRef}
          type="text"
          class="form-input ti-real"
          classList={{"form-input-error": props.errors.length > 0}}
          onFocus={() => props.errorsSetFn([])}
          onClick={() => setTitleOnTop(true)}
          onBlur={() => nameFocus()}
          onInput={(e) => triggerSearch(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && props.enterFn) {
              props.enterFn(props.enterFnArg)
            }
          }}
          value={props.value}
          placeholder={props.placeholder}/>
        <input 
          class="form-input ti-ghost"
          style={props.bgColor ? {"background-color": props.bgColor} : ""}
          type="text"
          disabled
          value={props.autocompleteValue ?? ""}/>
      </div>
    </div>
  )
}

export default TextInput;