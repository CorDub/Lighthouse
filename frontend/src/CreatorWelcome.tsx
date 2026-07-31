import "./styles/CreatorWelcome.css";
import { createSignal, Show } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import Navbar from "./Navbar.tsx";
import Errors from "./Errors.tsx";
import TextInput from "./TextInput.tsx";
import PasswordInput from "./PasswordInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { type ValueCheck } from "./types/helpersTypes.ts";
import { type ErrorKey } from "./Errors.tsx";
import creaWelText from "./translations/CreatorWelcome.json";
import Text from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import { getText } from "./helpers/getText.ts";
import { type LanguageCode } from "./types/langTypes.ts";
import LockableTextInput from "./LockableTextInput.tsx";

function CreatorWelcome() {
  const [searchParams] = useSearchParams<{ token: string, name: string }>()
  const [name, setName] = createSignal(searchParams.name ?? "")
  const [email, setEmail] = createSignal("")
  const [password, setPassword] = createSignal("")
  const [language, setLanguage] = createSignal<LanguageCode>("en")
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const [step, setStep] = createSignal(0);
  const { defaults } = useDefaults();

  async function confirmName(e: Event) {
    e.preventDefault();

    const checks: ValueCheck[] = [
      ["name", name()],
    ]
    const checkResults = checkForErrors(...checks)

    if (checkResults.length > 0) {
      setErrors(checkResults)
      return
    }

    setStep(1)
    ///////
  }

  return (
    <div class="creator-welcome">
      <Navbar />
      <Show when={step() === 0}>
        <form
          class="cw-form"
          onSubmit={(e) => confirmName(e)}>
          <h1 class="cw-heading">
            <Text
              value={creaWelText.welcomeHeading}
              lang={defaults().lang} />
          </h1>
          <p class="cw-message">
            <Text
              value={creaWelText.welcomeMessage}
              lang={defaults().lang} />
          </p>
          <div class="cw-input">
            {/* <TextInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={name()}
              valueSetFn={setName}
              placeholder={getText(creaWelText.nameInput, defaults().lang)}
              autofocus={true}
              name={getText(creaWelText.nameInputTitle, defaults().lang)}
              nameOnTop={true}
              bgColor={"var(--white)"}/> */}
            <LockableTextInput 
              errors={errors()}
              errorsSetFn={setErrors}
              value={name()}
              valueSetFn={setName}/>
          </div>
          <div class="cw-input">
            <TextInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={email()}
              valueSetFn={setEmail}
              // placeholder={getText(creaWelText.emailInput, defaults().lang)}
              name={getText(creaWelText.emailInput, defaults().lang)}
              bgColor={"var(--white)"}/>
          </div>
          <div class="cw-input">
            <PasswordInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={password()}
              valueSetFn={setPassword}
              placeholder={getText(creaWelText.passwordInput, defaults().lang)}/>
          </div>
        
          <Show when={errors().length > 0}>
            <Errors 
              errors={errors()}
              margin={{
                marginTop: 0,
                marginBottom: 2
              }}/>
          </Show>
          
          <button class="green-button clickable" onClick={(e) => confirmName(e)}>
            <Text
              value={creaWelText.confirmButton}
              lang={defaults().lang} />
          </button>
        </form>
      </Show>
      
      {/* <Show when={step() === 1}>
        <div class="cw-step2">
          <div class="cw-message">
            <Text 
              />
          </div>
        </div>
      </Show>  */}
    </div>
  )
}

export default CreatorWelcome;
