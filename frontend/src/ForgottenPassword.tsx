import "./ForgottenPassword.css";
import { createSignal, Show } from "solid-js";
import { type ValueCheck } from "./helpers/helpersTypes.ts";
import Navbar from "./Navbar";
import Errors from "./Errors.tsx";
import TextInput from "./TextInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";

function ForgottenPassword() {
  const [email, setEmail] = createSignal("")
  const [resetAccepted, setResetAccepted] = createSignal(false)
  const [errors, setErrors] = createSignal<string[]>([]);

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

      const response = await fetch("/api/checkPassword", {
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
          <p>{`A reset email has been sent to ${email()}`}</p>
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
            placeholder="Enter email address"/>

          <Errors errors={errors()}/>

          <button 
            class="green-button clickable"
            onClick={(e) => resetPassword(e)}>
              Submit
          </button>
        </form>
      </Show>
    </div>
  )
}

export default ForgottenPassword;