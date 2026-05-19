import { createEffect } from "solid-js";

type TextInputProps = {
  errors: string[],
  errorsSetFn: (errors: string[]) => void,
  value: string,
  valueSetFn: (value: string) => void,
  placeholder?: string
  autofocus?: boolean
}

function TextInput(props: TextInputProps) {
  //autofocus
  let inputRef: HTMLInputElement | undefined
  createEffect(() => {
    if (props.autofocus) {
      inputRef?.focus()
    }
  })

  return (
    <input 
      ref={inputRef}
      type="text"
      class="form-input clickable"
      classList={{"form-input-error": props.errors.length > 0}}
      onFocus={() => props.errorsSetFn([])}
      onInput={(e) => props.valueSetFn(e.target.value)}
      value={props.value}
      placeholder={props.placeholder}/>
  )
}

export default TextInput;