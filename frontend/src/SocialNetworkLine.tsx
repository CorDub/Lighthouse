import type { JSX } from "solid-js";
import { createSignal, Show, onMount } from "solid-js";
import "./styles/SocialNetworkLine.css";
import Errors from "./Errors";
import type { ErrorKey } from "./Errors";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import LockableTextInput from "./LockableTextInput";
import { BASE_URL } from "./helpers/config";

type SocialNetworkLineProps = {
  title: string,
  icon: JSX.Element;
}

function SocialNetworkLine(props: SocialNetworkLineProps) {
  const [isOpen, setOpen] = createSignal(false)
  const [isOpening, setOpening] = createSignal(false)
  const [isHovered, setHovered] = createSignal(false)
  const [width, setWidth] = createSignal<number>();
  const [height, setHeight] = createSignal<number>();
  let lineRef: HTMLDivElement | undefined
  const [channel, setChannel] = createSignal("");
  const [errors, setErrors] = createSignal<ErrorKey[]>([])
  const [isConnecting, setConnecting] = createSignal(false)

  onMount(() => {
    if (!lineRef) return;
    const observer = new ResizeObserver(() => {
      setWidth(lineRef!.offsetWidth);
      setHeight(lineRef!.offsetHeight);
      observer.disconnect();
    });
    observer.observe(lineRef);
  });

  function openingLine() {
    setOpening(true)
    setTimeout(() => {
      setOpening(false)
      setOpen(true)
    }, 400)
  }

  function closingLine() {
    setOpening(true)
    setTimeout(() => {
      setOpening(false)
      setOpen(false)
    }, 400)
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

  function prependAtOnInput(value:string) {
    if (value.length > 0 && !value.startsWith("@")) {
      value = "@" +value;
    }
    setChannel(value)
  }

  return (
    <div class="social-network-line">
      <Show when={!isOpen()}>
        <div class="sn-line clickable"
          ref={lineRef}
          onClick={openingLine}
          onMouseEnter={() => setHovered(true)}
          onMouseLeave={() => setHovered(false)}
          style={{
            "--resting-width": `${width()}px`,
            "--resting-height": `${height()}px`
          }}
          classList={{ 
            "snl-opening" : isOpening(),
            "sn-line-hover" : isHovered()
          }}>
          <div class="sm-icon">
            {props.icon}
          </div>
          <div class="snl-title">
            {props.title}
          </div>
        </div>
      </Show>
      <Show when={isOpen()}>
        <div class="snl-open clickable">
          <div class="snlo-title">
            <div class="sm-icon">
              {props.icon}
            </div>
            <div class="snl-title">
              {props.title}
            </div>
          </div>

          <Show when={props.title==="YouTube"}>
            <div class="snlo-handle">
              <p class="snlo-handle-title">Please provide your YouTube channel handle</p>
              <p class="snlo-handle-text">The handle starts at @ in the url of your channel</p>
              <p class="snlo-handle-text">For example, @deDOS in "youtube.com/@deDOS"</p>
            </div>
          </Show>

          <form class="snlo-channels"
            onSubmit={connectChannel}>
            <div class="snlo-channel-input">
              <LockableTextInput 
                value={channel()}
                valueSetFn={props.title === "YouTube" ? prependAtOnInput : setChannel}
                errors={errors()}
                errorsSetFn={setErrors}
                locked={false}
                required={true}
                valueCheck={["channelName", channel()]}/>
            </div>
            <div class="snlo-button">
              <button class="green-button"
                type="button"
                onClick={connectChannel}>
                  Connect to Lighthouse
              </button>
            </div>
            
          </form>
        
          <Show when={errors().length > 0}>
            <Errors 
              errors={errors()}
              margin={{
                "marginTop": 0,
                "marginBottom": 0
              }}/>
          </Show>
        </div>
      </Show>
    </div>
  )
}

export default SocialNetworkLine