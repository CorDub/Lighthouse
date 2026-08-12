import type { JSX } from "solid-js";
import { createSignal, Show, onMount, For } from "solid-js";
import "./styles/SocialNetworkLine.css";
import Errors from "./Errors";
import type { ErrorKey } from "./Errors";
import SocialNetworkLineChannel from "./SocialNetworkLineChannel.tsx";
import type { Connection } from "./schemas/connection.ts";

type SocialNetworkLineProps = {
  title: string,
  icon: JSX.Element,
  channels: Connection[],
}

function SocialNetworkLine(props: SocialNetworkLineProps) {
  const [isOpen, setOpen] = createSignal(false)
  const [isOpening, setOpening] = createSignal(false)
  const [isHovered, setHovered] = createSignal(false)
  const [width, setWidth] = createSignal<number>();
  const [height, setHeight] = createSignal<number>();
  let lineRef: HTMLDivElement | undefined
  const [errors, setErrors] = createSignal<ErrorKey[]>([])
  

  onMount(() => {
    if (!lineRef) return;
    const observer = new ResizeObserver(() => {
      setWidth(lineRef!.offsetWidth);
      setHeight(lineRef!.offsetHeight);
      observer.disconnect();
    });
    observer.observe(lineRef);
  })

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

          <Show when={props.channels.length > 0}>
            <Show when={props.title==="YouTube"}>
              <div class="snlo-handle">
                <p class="snlo-handle-title">Please provide your YouTube channel handle</p>
                <p class="snlo-handle-text">For example, '@deDOS' in "youtube.com/@deDOS"</p>
              </div>
            </Show>

            <For each={props.channels}>
              {(channel, index) => (
                <SocialNetworkLineChannel 
                  title={props.title}
                  errors={errors()}
                  errorsSetFn={setErrors}
                  locked={true}/>
              )}
            </For>
          </Show>

          <SocialNetworkLineChannel 
            title={props.title}
            errors={errors()}
            errorsSetFn={setErrors}
            locked={false}
            required={true}/>
        
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