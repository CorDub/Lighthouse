import "./styles/SocialNetworkLineChannel.css"
import LockableTextInput from "./LockableTextInput"
import { createSignal, Show } from "solid-js";
import { checkForErrors } from "./helpers/checkForErrors";
import Errors from "./Errors.tsx";
import { BASE_URL } from "./helpers/config";
import type { ErrorKey } from "./Errors";
import { onMount, onCleanup } from "solid-js";
import type { Connection } from "./schemas/connection.ts";
import Tooltip from "./Tooltip.tsx";

type SocialNetworkLineChannelProps = {
  title: string,
  errors: string[],
  errorsSetFn: (value: string[]) => void,
  locked?: boolean,
  required?: boolean,
  connected?: boolean,
  channel?: Connection,
  addingNewChannelFn?: (value: boolean) => void
}

function SocialNetworkLineChannel(props: SocialNetworkLineChannelProps) {
  const [channel, setChannel] = createSignal(props.channel?.channelHandle || "")
  const [isConnecting, setConnecting] = createSignal(false)
  const [isConnected, setConnected] = createSignal(props.connected || false)
  const [confirmationShown, setConfirmationShown] = createSignal(false)
  const [errors, setErrors] = createSignal<ErrorKey[]>([])
  const [deactivated, setDeactivated] = createSignal(!props.channel?.active || false)
  const [deactivateButtonRect, setDeactivateButtonRect] = createSignal<DOMRect>()
  const frontendOrigin = import.meta.env.VITE_FRONTEND_ORIGIN
  //has to be top level for later cleanup
  let watcherIntervalId: ReturnType<typeof setInterval> | undefined;
  let deactivateButtonRef: HTMLDivElement | undefined;
  const [tooltipDeactivateButtonOpen, setTooltipDeactivateButtonOpen] = createSignal(false);

  onMount(() => {
    // get a listener for the pop up closing event
    window.addEventListener("message", handlePostMessage)

    //get deactivate button measurement for the tooltip
    setDeactivateButtonRect(deactivateButtonRef?.getBoundingClientRect())
  });

  onCleanup(() => {
    window.removeEventListener("message", handlePostMessage)

    //cleaning up the setInterval for checking the popup is closed or open
    if (watcherIntervalId !== undefined) {
      clearInterval(watcherIntervalId);
    }  
  })

  function handlePostMessage(event: MessageEvent) {
    if (event.origin !== frontendOrigin) return;
    if (event.data?.source !== "lighthouse-oauth") return;

    if (event.data.status === "success") {
      setConnecting(false);
      setConnected(true);
      setConfirmationShown(true);
      setTimeout(() => {
        setConfirmationShown(false);
      }, 1000)
    } else if (event.data.status === "error") {
      setConnecting(false)
      setErrors(["oauthFailed"]);
    } else if (event.data.status === "cancelled") {
      setConnecting(false)
      setErrors(["oauthDenied"]);
    }
  }

  function prependAtOnInput(value:string) {
    if (value.length > 0 && !value.startsWith("@")) {
      value = "@" +value;
    }
    setChannel(value)
  }

  async function connectChannel(e: Event) {
    e.preventDefault();

    // if already connecting returns
    if (isConnecting()) { return };

    // checks 
    const newErrorList = checkForErrors(["channelName", channel()])
    setErrors(newErrorList)
    if (errors().length > 0) { return }

    //open blank pop-up straightaway before any awaits takes too long and gets blocked by popup blockers
    const popup = window.open("", "youtube-oauth", "width")
    if (!popup) {
      setErrors(["popupBlocked"]);
      return;
    }

    setConnecting(true);
    setErrors([]);

    //set an interval to check whether the pop up is still active
    watcherIntervalId = setInterval(() => {
      if (popup.closed) {
        clearInterval(watcherIntervalId);
        watcherIntervalId = undefined;
        if (isConnecting()) {
          setConnecting(false)
        }
      }
    }, 250);

    try {
      const response = await fetch(`${BASE_URL}/api/connectChannel/${props.title.toLowerCase()}`, {
        method: "POST",
        credentials: "include",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          channelName: channel(),
        })
      })

      // error handling
      if (response.status === 404) {
        setErrors(["channelNotFound"])
        popup.close()
        setConnecting(false)
        return
      }
      if (!response.ok) {
        setErrors(["unexpectedError"])
        popup.close()
        setConnecting(false)
        return
      }

      if (response.ok) {
        const data = await response.json()
        popup.location.href = data.authUrl;
      }
    } catch(error) {
      console.log("Error connecting the channel to YT", error)
      setErrors(["unexpectedError"])
      popup.close()
      setConnecting(false); 
    }
  }

  async function toggleConnection() {
    try {
      const response = await fetch(`${BASE_URL}/api/connections/${props.channel?.id}`, {
        method: "PATCH",
        credentials: "include",
        headers: {"Content-Type": "application/json"}
      })

      if (response.ok) {
        const data = await response.json()
        setDeactivated(!data.active)
      }

    } catch(error) {
      console.log("Error deactivating the connection", error)
    }
  }

  return(
    <>
      <form class="snlc-channels"
        onSubmit={connectChannel}>
        <div class="snlc-channel-input">
          <LockableTextInput 
            value={channel()}
            valueSetFn={props.title === "YouTube" ? prependAtOnInput : setChannel}
            errors={props.errors}
            errorsSetFn={props.errorsSetFn}
            locked={props.locked}
            required={props.required}
            valueCheck={["channelName", channel()]}
            openable={!props.locked}
            deactivated={deactivated()}/>
          <Show when={isConnected()}>
            <div class="sm-icon clickable snlc-channel-link"
              ref={deactivateButtonRef}
              classList={{
                "snlc-activated": !deactivated(),
                "lti-deactivated": deactivated()
              }}
              onMouseEnter={() => {setTooltipDeactivateButtonOpen(true)}}
              onMouseLeave={() => {setTooltipDeactivateButtonOpen(false)}}
              onClick={toggleConnection}>
              <Show when={tooltipDeactivateButtonOpen()}>
                <Tooltip 
                  position={"over"}
                  topSender={deactivateButtonRect()?.top || 0}
                  leftSender={deactivateButtonRect()?.left || 0}
                  widthSender={deactivateButtonRect()?.width || 0}
                  heightSender={deactivateButtonRect()?.height || 0}
                  text={deactivated() ? "Reactivate channel" : "Deactivate channel"}/>
              </Show>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640" fill="currentColor"><path d="M451.5 160C434.9 160 418.8 164.5 404.7 172.7C388.9 156.7 370.5 143.3 350.2 133.2C378.4 109.2 414.3 96 451.5 96C537.9 96 608 166 608 252.5C608 294 591.5 333.8 562.2 363.1L491.1 434.2C461.8 463.5 422 480 380.5 480C294.1 480 224 410 224 323.5C224 322 224 320.5 224.1 319C224.6 301.3 239.3 287.4 257 287.9C274.7 288.4 288.6 303.1 288.1 320.8C288.1 321.7 288.1 322.6 288.1 323.4C288.1 374.5 329.5 415.9 380.6 415.9C405.1 415.9 428.6 406.2 446 388.8L517.1 317.7C534.4 300.4 544.2 276.8 544.2 252.3C544.2 201.2 502.8 159.8 451.7 159.8zM307.2 237.3C305.3 236.5 303.4 235.4 301.7 234.2C289.1 227.7 274.7 224 259.6 224C235.1 224 211.6 233.7 194.2 251.1L123.1 322.2C105.8 339.5 96 363.1 96 387.6C96 438.7 137.4 480.1 188.5 480.1C205 480.1 221.1 475.7 235.2 467.5C251 483.5 269.4 496.9 289.8 507C261.6 530.9 225.8 544.2 188.5 544.2C102.1 544.2 32 474.2 32 387.7C32 346.2 48.5 306.4 77.8 277.1L148.9 206C178.2 176.7 218 160.2 259.5 160.2C346.1 160.2 416 230.8 416 317.1C416 318.4 416 319.7 416 321C415.6 338.7 400.9 352.6 383.2 352.2C365.5 351.8 351.6 337.1 352 319.4C352 318.6 352 317.9 352 317.1C352 283.4 334 253.8 307.2 237.5z"/></svg>
            </div>
          </Show>
        </div>

        <div class="snlc-button">
          <Show when={!isConnected()}>
            <Show when={!isConnecting()}>
              <button class="green-button clickable"
                type="button"
                onClick={connectChannel}>
                  Connect to Lighthouse
              </button>
              <div class="snlc-cancel clickable icon"
                onClick={() => props.addingNewChannelFn?.(false)}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640" fill="currentColor"><path d="M504.6 148.5C515.9 134.9 514.1 114.7 500.5 103.4C486.9 92.1 466.7 93.9 455.4 107.5L320 270L184.6 107.5C173.3 93.9 153.1 92.1 139.5 103.4C125.9 114.7 124.1 134.9 135.4 148.5L278.3 320L135.4 491.5C124.1 505.1 125.9 525.3 139.5 536.6C153.1 547.9 173.3 546.1 184.6 532.5L320 370L455.4 532.5C466.7 546.1 486.9 547.9 500.5 536.6C514.1 525.3 515.9 505.1 504.6 491.5L361.7 320L504.6 148.5z"/></svg>
              </div>
            </Show>

            <Show when={isConnecting()}>
              <div class="green-button snlc-connecting">
                Connecting...
              </div>
            </Show>
          </Show>

          <Show when={isConnected()}>
            <Show when={confirmationShown()}>
            <div class="green-button snlc-connecting">
              <p>Connected to Lighthouse</p>
              <div class="small-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640" fill="currentColor"><path d="M530.8 134.1C545.1 144.5 548.3 164.5 537.9 178.8L281.9 530.8C276.4 538.4 267.9 543.1 258.5 543.9C249.1 544.7 240 541.2 233.4 534.6L105.4 406.6C92.9 394.1 92.9 373.8 105.4 361.3C117.9 348.8 138.2 348.8 150.7 361.3L252.2 462.8L486.2 141.1C496.6 126.8 516.6 123.6 530.9 134z"/></svg>
              </div>
            </div>
            </Show>
          </Show>
        </div>
        
        <Show when={errors().length > 0}>
          <Errors 
            errors={errors()}
            margin={{
              "marginTop": 0,
              "marginBottom": 0
            }}/>
        </Show>
      </form>
    </>
  )
}

export default SocialNetworkLineChannel