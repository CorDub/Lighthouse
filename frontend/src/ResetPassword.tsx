import "./ResetPassword.css";
import { createSignal, Show } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import Navbar from "./Navbar.tsx";
import Errors from "./Errors.tsx";
import PasswordInput from "./PasswordInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { type ValueCheck } from "./helpers/helpersTypes.ts"; 
import { BASE_URL } from "./helpers/config.ts";

function ResetPassword() {
  const [searchParams] = useSearchParams()
  const paramsToken = searchParams.token
  const [newPassword, setNewPassword] = createSignal("")
  const [confirmed, setConfirmed] = createSignal(false)
  const [errors, setErrors] = createSignal<string[]>([]);

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
            placeholder={"Enter your new password"}
            autofocus={true}/>

          <Errors errors={errors()} />
          
          <button
            class="green-button clickable"
            onClick={(e) => changePassword(e)}>
              Submit
          </button>
        </form>
      </Show>
      <Show when={confirmed()}>
        <div class="respas-confirmation">
          <p style={{ "margin-bottom": "1rem" }}>Your new password has been saved</p>
          <button class="green-button clickable">Return to login</button>
        </div>
      </Show>
    </div>
  )
}

export default ResetPassword;