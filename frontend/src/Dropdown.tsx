import { Show, children, type JSXElement } from "solid-js";

type DropdownProps = {
  openStatus: boolean,
  setOpenStatus: (status: boolean) => void,
  anchor: {
    top?: number,
    bottom?: number,
    left?: number,
    right?: number
  },
  children: JSXElement
}

function Dropdown(props: DropdownProps) {
  const resolvedChildren = children(() => props.children)

  return (
    <Show when={props.openStatus}>
      <div class="dd-outclickcatcher"
        onClick={() => props.setOpenStatus(false)}>
      </div>
      <div class="dropdown"
        style={{
          top: props.anchor.top != null ? `${props.anchor.top}px` : undefined,
          bottom: props.anchor.bottom != null ? `${window.innerHeight - props.anchor.bottom}px` : undefined,
          left: props.anchor.left != null ? `${props.anchor.left}px` : undefined,
          right: props.anchor.right != null ? `${window.innerWidth - props.anchor.right}px` : undefined,
        }}>
        {resolvedChildren()}
      </div>
    </Show>
  )
}

export default Dropdown