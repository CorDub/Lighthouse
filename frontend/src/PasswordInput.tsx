import "./styles/PasswordInput.css";
import { createEffect, mergeProps, createSignal } from "solid-js";
import { getText } from "./helpers/getText";
import pasinText from "./translations/PasswordInput.json";
import { useDefaults } from "./DefaultsContext";

type PasswordInputProps = {
  errors: string[],
  errorsSetFn: (errors: string[]) => void,
  value: string,
  valueSetFn: (value: string) => void,
  placeholder?: string
  autofocus?: boolean
}

function PasswordInput(props: PasswordInputProps) {
  const [visible, setVisible] = createSignal(false)
  const { defaults } = useDefaults()
  const pasinDefaults = mergeProps({
    placeholder: getText(pasinText.placeholder, defaults().lang),
    autofocus: false
  }, props)

  //autofocus
  let inputRef: HTMLInputElement | undefined
  createEffect(() => {
    if (props.autofocus) {
      inputRef?.focus()
    }
  })

  return(
    <div class="password-input">
      <input 
        ref={inputRef}
        type={visible() ? "text" : "password"}
        class="form-input clickable"
        classList={{"form-input-error": props.errors.length > 0}}
        onFocus={() => props.errorsSetFn([])}
        onInput={(e) => props.valueSetFn(e.target.value)}
        value={props.value}
        placeholder={pasinDefaults.placeholder}/>
      <div 
        class="icon password-icon clickable"
        onClick={() => setVisible(!visible())}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M320 96C239.2 96 174.5 132.8 127.4 176.6C80.6 220.1 49.3 272 34.4 307.7C31.1 315.6 31.1 324.4 34.4 332.3C49.3 368 80.6 420 127.4 463.4C174.5 507.1 239.2 544 320 544C400.8 544 465.5 507.2 512.6 463.4C559.4 419.9 590.7 368 605.6 332.3C608.9 324.4 608.9 315.6 605.6 307.7C590.7 272 559.4 220 512.6 176.6C465.5 132.9 400.8 96 320 96zM176 320C176 240.5 240.5 176 320 176C399.5 176 464 240.5 464 320C464 399.5 399.5 464 320 464C240.5 464 176 399.5 176 320zM320 256C320 291.3 291.3 320 256 320C244.5 320 233.7 317 224.3 311.6C223.3 322.5 224.2 333.7 227.2 344.8C240.9 396 293.6 426.4 344.8 412.7C396 399 426.4 346.3 412.7 295.1C400.5 249.4 357.2 220.3 311.6 224.3C316.9 233.6 320 244.4 320 256z"/></svg>
      </div>
    </div>
  )
}

export default PasswordInput;