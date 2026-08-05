import { createSignal, Show } from "solid-js";
import TextInput from "./TextInput";
import PasswordInput from "./PasswordInput";
import LockableTextInput from "./LockableTextInput";

type TextInputModel = "simple" | "password" | "lockable"

type TextInputRouterProps = {
  model: TextInputModel,
  errors: string[],
  errorsSetFn: (errors: string[]) => void,
  value: string,
  valueSetFn: (value: string) => void,
  valueSetFnArg?: string,
  autocomplete?: boolean,
  autocompleteValue?: string,
  setAutocompleteFn?: (value: string) => void,
  searchFn?: () => void,
  enterFn?: (...args: any[]) => any,
  enterFnArg?: any,
  onBlurFn?: (...args: any[]) => void,
  onBlurFnArg?: any,
  placeholder?: string,
  autofocus?: boolean,
  name?: string,
  nameOnTop?: boolean,
  bgColor?: string,
  newClass?: string,
  newClassToggle?: boolean,
  //lockable
  locked?: boolean
}

function TextInputRouter(props: TextInputRouterProps) {
  const [titleOnTop, setTitleOnTop] = createSignal(false)

  //autofocus logic
  function autofocus(ref: HTMLInputElement) {
    if (props.autofocus) {
      ref?.focus()
    }
  }

  function setNameOnTop(ref: HTMLInputElement) {
    if (props.name && props.nameOnTop) {
      setTitleOnTop(true)
      ref?.focus
    }
  }

  function resetNameAndInput() {
    if (props.name && titleOnTop()) {
      props.valueSetFn("")
      setTitleOnTop(false)
      props.searchFn?.()
      if (props.setAutocompleteFn) { props.setAutocompleteFn("") }
    }
  }

  function runOnBlurFn() {
    if (props.onBlurFn) {
      props.onBlurFn(props.onBlurFnArg)
      return
    }

    resetNameAndInput();
   }

   function triggerSearch(value: string) {
    props.valueSetFn(value)
    props.searchFn?.()
   }

  return(
    <>
      <Show when={props.model === "simple"}>
        <TextInput 
          errors={props.errors}
          errorsSetFn={props.errorsSetFn}
          value={props.value}
          valueSetFn={props.valueSetFn}
          valueSetFnArg={props.valueSetFnArg}
          autocomplete={props.autocomplete}
          autocompleteValue={props.autocompleteValue}
          setAutocompleteFn={props.setAutocompleteFn}
          searchFn={props.searchFn}
          enterFn={props.enterFn}
          enterFnArg={props.enterFnArg}
          onBlurFn={props.onBlurFn}
          onBlurFnArg={props.onBlurFnArg}
          placeholder={props.placeholder}
          autofocus={props.autofocus}
          name={props.name}
          nameOnTop={props.nameOnTop}
          bgColor={props.bgColor}
          newClass={props.newClass}
          newClassToggle={props.newClassToggle}/>
      </Show>

      <Show when={props.model === "password"}>
        <PasswordInput 
          errors={props.errors}
          errorsSetFn={props.errorsSetFn}
          value={props.value}
          valueSetFn={props.valueSetFn}/>
      </Show>

      <Show when={props.model === "lockable"}>
        <LockableTextInput 
          errors={props.errors}
          errorsSetFn={props.errorsSetFn}
          value={props.value}
          valueSetFn={props.valueSetFn}/>
      </Show>
    </>
  )
}

export default TextInputRouter;

