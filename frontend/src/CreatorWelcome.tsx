import "./styles/CreatorWelcome.css";
import { createSignal, Show } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import Navbar from "./Navbar.tsx";
import Errors from "./Errors.tsx";
import TextInput from "./TextInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { type ValueCheck } from "./types/helpersTypes.ts";
import { type ErrorKey } from "./Errors.tsx";
import creaWelText from "./translations/CreatorWelcome.json";
import Text from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import { getText } from "./helpers/getText.ts";

function CreatorWelcome() {
  const [searchParams] = useSearchParams<{ token: string, name: string }>()
  const [name, setName] = createSignal(searchParams.name ?? "")
  const [confirmed, setConfirmed] = createSignal(false)
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const { defaults } = useDefaults();

  function confirmName(e: Event) {
    e.preventDefault();

    const checks: ValueCheck[] = [
      ["name", name()],
    ]
    const checkResults = checkForErrors(...checks)

    if (checkResults.length > 0) {
      setErrors(checkResults)
      return
    }

    setConfirmed(true)
  }

  return (
    <div class="creator-welcome">
      <Navbar />
      <Show when={!confirmed()}>
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
            <TextInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={name()}
              valueSetFn={setName}
              placeholder={getText(creaWelText.nameInput, defaults().lang)}
              autofocus={true}/>
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
      <Show when={confirmed()}>
        <div class="cw-confirmation">
          <p>
            <Text
              value={creaWelText.nameConfirmed}
              lang={defaults().lang} />
          </p>
        </div>
      </Show>
    </div>
  )
}

export default CreatorWelcome;
