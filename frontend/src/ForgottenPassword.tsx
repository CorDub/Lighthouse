import "./styles/ForgottenPassword.css";
import { createSignal, Show } from "solid-js";
import { type ValueCheck } from "./types/helpersTypes.ts";
import Navbar from "./Navbar";
import Errors from "./Errors.tsx";
import TextInput from "./TextInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { BASE_URL } from "./helpers/config.ts";
import { type ErrorKey } from "./Errors.tsx";
import Text from "./Text.tsx";
import forpasText from "./translations/ForgottenPassword.json";
import { useDefaults } from "./DefaultsContext.tsx";
import { getText } from "./helpers/getText.ts";

function ForgottenPassword() {
  const [email, setEmail] = createSignal("")
  const [resetAccepted, setResetAccepted] = createSignal(false)
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const { defaults } = useDefaults()

  async function resetPassword(e: Event) {
    try {
      e.preventDefault()

      //check if the email sent is a valid password
      const checks: ValueCheck[] = [
        ["email", email()]
      ]
      const checksResults = checkForErrors(...checks)

      if (checksResults.length > 0) {
        setErrors(checksResults)
        return
      }

      const response = await fetch(`${BASE_URL}/api/checkPassword`, {
        method: "POST",
        credentials: "include",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          email: email()
        })
      })

      if (response.ok) {
        setResetAccepted(true)
      }
    } catch(error) {
      console.error(error)
    }
  }

  return (
    <div class="forgotten-password">
      <Navbar/>
      <Show when={resetAccepted()}>
        <div>
          <p>
            <Text 
              value={forpasText.resetEmailSent} 
              lang={defaults().lang}
              var={[email()]}/>
          </p>
        </div>
      </Show>
      <Show when={!resetAccepted()}>
        <form 
          class="confirm-email-form"
          onSubmit={(e) => resetPassword(e)}>
          <TextInput 
            errors={errors()}
            errorsSetFn={setErrors}
            value={email()}
            valueSetFn={setEmail}
            placeholder={getText(forpasText.passwordInput, defaults().lang)}/>

          <Errors errors={errors()}/>

          <button 
            class="green-button clickable"
            onClick={(e) => resetPassword(e)}>
              <Text 
                value={forpasText.submit}
                lang={defaults().lang}/>
          </button>
        </form>
      </Show>
    </div>
  )
}

export default ForgottenPassword;