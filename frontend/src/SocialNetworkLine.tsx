import type { JSX } from "solid-js";
import { createSignal, Show, onMount } from "solid-js";
import "./styles/SocialNetworkLine.css";
import TextInput from "./TextInput.tsx";
import Errors from "./Errors";
import type { ErrorKey } from "./Errors";
import LockableTextInput from "./LockableTextInput";

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
          <div class="snlo-channels">
            <div class="snlo-channel-input">
              {/* <TextInput 
                value={channel()}
                valueSetFn={setChannel}
                errors={errors()}
                errorsSetFn={setErrors}
                placeholder={"Channel Name"}/> */}
              <LockableTextInput 
                value={channel()}
                valueSetFn={setChannel}
                errors={errors()}
                errorsSetFn={setErrors}
                locked={false}/>
            </div>
          </div>
          
        </div>
      </Show>
    </div>
  )
}

export default SocialNetworkLine