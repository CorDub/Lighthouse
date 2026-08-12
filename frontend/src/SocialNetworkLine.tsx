import type { JSX } from "solid-js";
import { createSignal, Show, onMount, For } from "solid-js";
import "./styles/SocialNetworkLine.css";
import Errors from "./Errors";
import type { ErrorKey } from "./Errors";
import SocialNetworkLineChannel from "./SocialNetworkLineChannel.tsx";
import type { Connection } from "./schemas/connection.ts";
import AddButton from "./AddButton.tsx";

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
  const [isAddingNewChannel, setAddingNewChannel] = createSignal(false)

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
        <div class="snl-open">
          <div class="snlo-title">
            <div class="snlo-title-left">
              <div class="sm-icon">
                {props.icon}
              </div>
              <div class="snl-title">
                {props.title}
              </div>
            </div>
            
            <div class="snl-close-button">
              <div class="snlc-cancel clickable icon"
                onClick={closingLine}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640" fill="currentColor"><path d="M504.6 148.5C515.9 134.9 514.1 114.7 500.5 103.4C486.9 92.1 466.7 93.9 455.4 107.5L320 270L184.6 107.5C173.3 93.9 153.1 92.1 139.5 103.4C125.9 114.7 124.1 134.9 135.4 148.5L278.3 320L135.4 491.5C124.1 505.1 125.9 525.3 139.5 536.6C153.1 547.9 173.3 546.1 184.6 532.5L320 370L455.4 532.5C466.7 546.1 486.9 547.9 500.5 536.6C514.1 525.3 515.9 505.1 504.6 491.5L361.7 320L504.6 148.5z"/></svg>
              </div>
            </div>
          </div>

          <Show when={props.channels.length > 0}>
            <For each={props.channels}>
              {(channel, _) => (
                <SocialNetworkLineChannel 
                  title={props.title}
                  errors={errors()}
                  errorsSetFn={setErrors}
                  locked={true}
                  connected={true}
                  channel={channel}/>
              )}
            </For>

            <Show when={!isAddingNewChannel()}>
              <div class="snlc-add-button-line">
                <div class="snlc-add-button">
                  <AddButton 
                    clickFn={() => setAddingNewChannel(true)}/>
                </div>
              </div>
            </Show>

            <Show when={isAddingNewChannel()}>
              <SocialNetworkLineChannel 
                title={props.title}
                errors={errors()}
                errorsSetFn={setErrors}
                locked={false}
                required={true}
                addingNewChannelFn={setAddingNewChannel}/>
            </Show>

          </Show>

          <Show when={props.channels.length === 0}>
            <Show when={props.title==="YouTube"}>
              <div class="snlo-handle">
                <p class="snlo-handle-title">Please provide your YouTube channel handle</p>
                <p class="snlo-handle-text">For example, '@deDOS' in "youtube.com/@deDOS"</p>
              </div>
            </Show>
            <SocialNetworkLineChannel 
              title={props.title}
              errors={errors()}
              errorsSetFn={setErrors}
              locked={false}
              required={true}/>
          </Show>
        
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