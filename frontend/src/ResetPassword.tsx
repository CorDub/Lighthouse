import "./styles/ResetPassword.css";
import { createSignal, Show } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import Navbar from "./Navbar.tsx";
import Errors from "./Errors.tsx";
import PasswordInput from "./PasswordInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { type ValueCheck } from "./helpers/helpersTypes.ts"; 
import { BASE_URL } from "./helpers/config.ts";
import { type ErrorKey } from "./helpers/checkForErrors.ts";
import respasText from "./translations/ResetPassword.json";
import Text from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import SubmitButton from "./SubmitButton.tsx";
import { getText } from "./helpers/getText.ts";

function ResetPassword() {
  const [searchParams] = useSearchParams()
  const paramsToken = searchParams.token
  const [newPassword, setNewPassword] = createSignal("")
  const [confirmed, setConfirmed] = createSignal(false)
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const { defaults } = useDefaults();

  async function changePassword(e: Event) {
    e.preventDefault();

    // do the checks
    const checks: ValueCheck[] = [
      ["password", newPassword()],
    ]
    const checkResults = checkForErrors(...checks)

    if (checkResults.length > 0) {
      setErrors(checkResults)
      return
    }

    // send the new password
    const response = await fetch(`${BASE_URL}/api/changePassword`, {
      method: "POST",
      credentials: "include",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        password: newPassword(),
        token: paramsToken,
      })
    })

    if (response.ok) {
      setConfirmed(true)
    }
  }

  return (
    <div class="reset-password">
      <Navbar />
      <Show when={!confirmed()}>
        <form 
          class="confirm-email-form"
          onSubmit={(e) => changePassword(e)}>
          <PasswordInput 
            errors={errors()}
            errorsSetFn={setErrors}
            value={newPassword()}
            valueSetFn={setNewPassword}
            placeholder={getText(respasText.passwordInput, defaults().lang)}
            autofocus={true}/>

          <Errors errors={errors()} />
          
          <SubmitButton clickFn={(e) => changePassword(e)}/>
        </form>
      </Show>
      <Show when={confirmed()}>
        <div class="respas-confirmation">
          <p style={{ "margin-bottom": "1rem" }}>
            <Text 
              value={respasText.passwordSaved}
              lang={defaults().lang} />
          </p>
          <button class="green-button clickable">
            <Text 
              value={respasText.returnLogin}
              lang={defaults().lang} />
          </button>
        </div>
      </Show>
    </div>
  )
}

export default ResetPassword;